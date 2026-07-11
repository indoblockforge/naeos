package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})
	if s == nil {
		t.Fatal("expected server to be created")
	}
	if s.Addr != ":8080" {
		t.Errorf("expected addr ':8080', got %s", s.Addr)
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	s.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}
}

func TestSpecsEndpointGET(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("GET", "/api/v1/specs", nil)
	w := httptest.NewRecorder()

	s.handleSpecs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSpecsEndpointPOST(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	spec := map[string]interface{}{
		"project": "test",
		"modules": []interface{}{},
	}
	body, _ := json.Marshal(spec)

	req := httptest.NewRequest("POST", "/api/v1/specs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSpecs(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestSpecValidateEndpoint(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("POST", "/api/v1/specs/validate", nil)
	w := httptest.NewRecorder()

	s.handleSpecValidate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestPipelineRunEndpoint(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("POST", "/api/v1/pipeline/run", nil)
	w := httptest.NewRecorder()

	s.handlePipelineRun(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("DELETE", "/api/v1/specs", nil)
	w := httptest.NewRecorder()

	s.handleSpecs(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}
