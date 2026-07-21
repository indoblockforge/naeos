package ai

import (
	"strings"
	"testing"
)

func TestBuildEnrichPrompt(t *testing.T) {
	s := &LLMService{}
	prompt := s.buildEnrichPrompt("spec: test")
	if !strings.Contains(prompt, "spec: test") {
		t.Error("expected prompt to contain spec content")
	}
}

func TestTruncatePromptShort(t *testing.T) {
	s := &LLMService{}
	input := "short prompt"
	got := s.truncatePrompt(input)
	if got != input {
		t.Errorf("expected unchanged short prompt")
	}
}

func TestSplitSSELineEdgeCases(t *testing.T) {
	parts := splitSSELine("")
	if parts != nil {
		t.Error("expected nil for empty line")
	}
	parts = splitSSELine("no colon")
	if parts != nil {
		t.Error("expected nil for line without colon")
	}
	parts = splitSSELine("key:val")
	if len(parts) != 2 || parts[0] != "key" || parts[1] != "val" {
		t.Errorf("expected [key val], got %v", parts)
	}
	parts = splitSSELine("data: hello world")
	if len(parts) != 2 || parts[0] != "data" || parts[1] != "hello world" {
		t.Errorf("expected [data hello world], got %v", parts)
	}
}

func TestBuildSuggestionsPrompt(t *testing.T) {
	s := &LLMService{}
	prompt := s.buildSuggestionsPrompt("spec content here")
	if !strings.Contains(prompt, "spec content") {
		t.Error("expected prompt to contain spec content")
	}
}

func TestBuildExplainPrompt(t *testing.T) {
	s := &LLMService{}
	prompt := s.buildExplainPrompt("concept", "arch")
	if !strings.Contains(prompt, "concept") || !strings.Contains(prompt, "arch") {
		t.Error("expected prompt to contain concept and arch")
	}
}

func TestBuildCompilerPrompt(t *testing.T) {
	s := &LLMService{}
	prompt := s.buildCompilerPrompt("target", "spec")
	if !strings.Contains(prompt, "target") || !strings.Contains(prompt, "spec") {
		t.Error("expected prompt to contain target and spec")
	}
}
