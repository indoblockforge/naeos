package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Addr    string
	Router  *http.ServeMux
	server  *http.Server
	Auth    *AuthConfig
	Limiter *RateLimiter
}

type AuthConfig struct {
	JWTSecret string
	Enabled   bool
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewServer(addr string, auth *AuthConfig) *Server {
	s := &Server{
		Addr:   addr,
		Router: http.NewServeMux(),
		Auth:   auth,
		Limiter: NewRateLimiter(100, time.Minute),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health
	s.Router.HandleFunc("/api/v1/health", s.handleHealth)

	// Spec endpoints
	s.Router.HandleFunc("/api/v1/specs", s.handleSpecs)
	s.Router.HandleFunc("/api/v1/specs/validate", s.handleSpecValidate)
	s.Router.HandleFunc("/api/v1/specs/compile", s.handleSpecCompile)

	// Pipeline endpoints
	s.Router.HandleFunc("/api/v1/pipeline/run", s.handlePipelineRun)
	s.Router.HandleFunc("/api/v1/pipeline/status", s.handlePipelineStatus)

	// Artifact endpoints
	s.Router.HandleFunc("/api/v1/artifacts", s.handleArtifacts)

	// Context endpoints
	s.Router.HandleFunc("/api/v1/context/generate", s.handleContextGenerate)

	// MCP endpoints
	s.Router.HandleFunc("/api/v1/mcp/message", s.handleMCPMessage)
}

func (s *Server) handlerWithMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Rate limit
		if !s.Limiter.Allow() {
			s.writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Auth
		if s.Auth.Enabled && r.URL.Path != "/api/v1/health" {
			token := r.Header.Get("Authorization")
			if token == "" {
				s.writeError(w, http.StatusUnauthorized, "authorization required")
				return
			}
			// TODO: Validate JWT
		}

		handler(w, r)
	}
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	})
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"version": "0.4.0",
		"uptime":  time.Since(startTime).String(),
	})
}

func (s *Server) handleSpecs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"specs": []string{"spec.yaml", "spec.json"},
		})
	case "POST":
		var spec map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
			s.writeError(w, http.StatusBadRequest, "invalid spec")
			return
		}
		s.writeJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "spec received",
			"modules": len(spec["modules"].([]interface{})),
		})
	default:
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleSpecValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"valid":   true,
		"errors":  []string{},
		"warnings": []string{},
	})
}

func (s *Server) handleSpecCompile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"compiled": true,
		"targets":  []string{"copilot", "claude", "cursor", "gemini", "codex", "opencode"},
	})
}

func (s *Server) handlePipelineRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"pipeline_id": fmt.Sprintf("pipeline-%d", time.Now().UnixNano()),
		"status":      "running",
	})
}

func (s *Server) handlePipelineStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "idle",
		"last_run": nil,
	})
}

func (s *Server) handleArtifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"artifacts": []string{},
		})
	case "POST":
		s.writeJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "artifact stored",
		})
	default:
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleContextGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"context": "generated",
		"format":  "markdown",
	})
}

func (s *Server) handleMCPMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"jsonrpc": "2.0",
		"result":  "mcp response",
	})
}

var startTime = time.Now()

func (s *Server) Start() error {
	wrappedHandler := s.handlerWithMiddleware(s.Router.ServeHTTP)

	s.server = &http.Server{
		Addr:         s.Addr,
		Handler:      wrappedHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}()

	log.Printf("Starting NAEOS API server on %s", s.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}
