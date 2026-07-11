package websocket

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s == nil {
		t.Fatal("expected server to be created")
	}
}

func TestIsWebSocketUpgrade(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected bool
	}{
		{
			name: "valid websocket upgrade",
			headers: map[string]string{
				"Connection": "upgrade",
				"Upgrade":    "websocket",
			},
			expected: true,
		},
		{
			name: "not websocket",
			headers: map[string]string{
				"Connection": "keep-alive",
			},
			expected: false,
		},
		{
			name: "case insensitive",
			headers: map[string]string{
				"Connection": "Upgrade",
				"Upgrade":    "WebSocket",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/ws", nil)
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			if got := isWebSocketUpgrade(r); got != tt.expected {
				t.Errorf("isWebSocketUpgrade() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMessageSerialization(t *testing.T) {
	msg := Message{
		Type:    "test",
		Payload: map[string]string{"key": "value"},
		Time:    time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	if decoded.Type != "test" {
		t.Errorf("expected type 'test', got %s", decoded.Type)
	}
}

func TestBroadcast(t *testing.T) {
	s := NewServer()
	go s.Run()

	s.Broadcast("test", map[string]string{"message": "hello"})

	if s.ClientCount() != 0 {
		t.Errorf("expected 0 clients, got %d", s.ClientCount())
	}
}

func TestEventBroadcaster(t *testing.T) {
	s := NewServer()
	broadcaster := NewEventBroadcaster(s)
	go s.Run()

	broadcaster.PipelineStarted("pipeline-123")
	broadcaster.PipelineCompleted("pipeline-123", "10s")
	broadcaster.PipelineFailed("pipeline-123", "error")
	broadcaster.SpecValidated(true, []string{})
	broadcaster.ArtifactGenerated("main.go", "./cmd/main.go")
	broadcaster.LogMessage("info", "test message")
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == id2 {
		t.Error("expected unique IDs")
	}

	if !strings.HasPrefix(id1, "client-") {
		t.Errorf("expected ID to start with 'client-', got %s", id1)
	}
}
