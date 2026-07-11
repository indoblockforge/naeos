package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	conn    net.Conn
	server  *Server
	send    chan []byte
	id      string
	created time.Time
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	Time    time.Time   `json:"time"`
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			s.sendSystemMessage("client connected")

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()

		case message := <-s.broadcast:
			s.mu.RLock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if !isWebSocketUpgrade(r) {
		http.Error(w, "not a websocket request", http.StatusBadRequest)
		return
	}

	conn, err := upgradeConnection(w, r)
	if err != nil {
		http.Error(w, "upgrade failed", http.StatusInternalServerError)
		return
	}

	client := &Client{
		conn:    conn,
		server:  s,
		send:    make(chan []byte, 256),
		id:      generateID(),
		created: time.Now(),
	}

	s.register <- client
	go client.writePump()
	go client.readPump()
}

func (s *Server) Broadcast(eventType string, payload interface{}) {
	msg := Message{
		Type:    eventType,
		Payload: payload,
		Time:    time.Now(),
	}
	data, _ := json.Marshal(msg)
	s.broadcast <- data
}

func (s *Server) sendSystemMessage(text string) {
	s.Broadcast("system", map[string]string{"message": text})
}

func (s *Server) ClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

func (c *Client) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	for {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		msg, err := readFrame(c.conn)
		if err != nil {
			return
		}

		var incoming Message
		if err := json.Unmarshal(msg, &incoming); err != nil {
			continue
		}

		switch incoming.Type {
		case "ping":
			pong, _ := json.Marshal(Message{Type: "pong", Time: time.Now()})
			writeFrame(c.conn, pong)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				writeCloseFrame(c.conn)
				return
			}
			writeFrame(c.conn, message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			ping, _ := json.Marshal(Message{Type: "ping", Time: time.Now()})
			if err := writeFrame(c.conn, ping); err != nil {
				return
			}
		}
	}
}

// Raw WebSocket framing
const (
	webSocketGUID        = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	wsOpText             = 1
	wsOpClose            = 8
	wsOpPing             = 9
	wsOpPong             = 10
	wsFinalBit           = 0x80
	wsMaskBit            = 0x80
	wsMaxFrameSize       = 65536
)

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Connection"), "upgrade") &&
		strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return nil, fmt.Errorf("missing Sec-WebSocket-Key")
	}

	h := sha1.New()
	h.Write([]byte(key + webSocketGUID))
	acceptKey := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, fmt.Errorf("server does not support hijacking")
	}

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)
	w.WriteHeader(http.StatusSwitchingProtocols)

	conn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func readFrame(conn net.Conn) ([]byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	opcode := header[0] & 0x0F
	masked := (header[1] & wsMaskBit) != 0
	length := uint64(header[1] & 0x7F)

	switch length {
	case 126:
		ext := make([]byte, 2)
		if _, err := io.ReadFull(conn, ext); err != nil {
			return nil, err
		}
		length = uint64(binary.BigEndian.Uint16(ext))
	case 127:
		ext := make([]byte, 8)
		if _, err := io.ReadFull(conn, ext); err != nil {
			return nil, err
		}
		length = binary.BigEndian.Uint64(ext)
	}

	if length > wsMaxFrameSize {
		return nil, fmt.Errorf("frame too large: %d", length)
	}

	var mask [4]byte
	if masked {
		if _, err := io.ReadFull(conn, mask[:]); err != nil {
			return nil, err
		}
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}

	if masked {
		for i := uint64(0); i < length; i++ {
			payload[i] ^= mask[i%4]
		}
	}

	// Handle control frames
	switch opcode {
	case wsOpClose:
		return nil, fmt.Errorf("connection closed")
	case wsOpPing:
		pong, _ := json.Marshal(Message{Type: "pong", Time: time.Now()})
		writeFrame(conn, pong)
		return readFrame(conn)
	}

	return payload, nil
}

func writeFrame(conn net.Conn, data []byte) error {
	length := len(data)

	var header []byte
	if length < 126 {
		header = []byte{wsFinalBit | wsOpText, byte(length)}
	} else if length < 65536 {
		header = make([]byte, 4)
		header[0] = wsFinalBit | wsOpText
		header[1] = 126
		binary.BigEndian.PutUint16(header[2:], uint16(length))
	} else {
		header = make([]byte, 10)
		header[0] = wsFinalBit | wsOpText
		header[1] = 127
		binary.BigEndian.PutUint64(header[2:], uint64(length))
	}

	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

func writeCloseFrame(conn net.Conn) error {
	header := []byte{wsFinalBit | wsOpClose, 0}
	_, err := conn.Write(header)
	return err
}

func generateID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}

// EventBroadcaster sends events to all connected clients
type EventBroadcaster struct {
	server *Server
}

func NewEventBroadcaster(server *Server) *EventBroadcaster {
	return &EventBroadcaster{server: server}
}

func (b *EventBroadcaster) PipelineStarted(pipelineID string) {
	b.server.Broadcast("pipeline.started", map[string]string{"pipeline_id": pipelineID})
}

func (b *EventBroadcaster) PipelineCompleted(pipelineID string, duration string) {
	b.server.Broadcast("pipeline.completed", map[string]string{"pipeline_id": pipelineID, "duration": duration})
}

func (b *EventBroadcaster) PipelineFailed(pipelineID string, errMsg string) {
	b.server.Broadcast("pipeline.failed", map[string]string{"pipeline_id": pipelineID, "error": errMsg})
}

func (b *EventBroadcaster) SpecValidated(valid bool, errors []string) {
	b.server.Broadcast("spec.validated", map[string]interface{}{"valid": valid, "errors": errors})
}

func (b *EventBroadcaster) ArtifactGenerated(name string, path string) {
	b.server.Broadcast("artifact.generated", map[string]string{"name": name, "path": path})
}

func (b *EventBroadcaster) LogMessage(level string, message string) {
	b.server.Broadcast("log", map[string]string{"level": level, "message": message})
}
