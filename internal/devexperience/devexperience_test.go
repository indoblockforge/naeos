package devexperience

import (
	"strings"
	"testing"
)

func TestVSCodeExtension(t *testing.T) {
	ext := NewVSCodeExtension("naeos", "1.0.0", "Test extension", "Test Author", []string{"yaml", "json"})

	if ext.Name != "naeos" {
		t.Error("expected name 'naeos'")
	}

	pkg := ext.GeneratePackageJSON()
	if !strings.Contains(pkg, "naeos") {
		t.Error("expected package JSON to contain name")
	}
	if !strings.Contains(pkg, "yaml") {
		t.Error("expected package JSON to contain yaml")
	}

	syntax := ext.GenerateSyntaxJSON()
	if !strings.Contains(syntax, "naeos.yaml") {
		t.Error("expected syntax to contain naeos.yaml")
	}
}

func TestCompletionEngine(t *testing.T) {
	e := NewCompletionEngine()

	// Empty input
	completions := e.Complete("")
	if len(completions) != len(e.commands) {
		t.Errorf("expected %d completions, got %d", len(e.commands), len(completions))
	}

	// Partial match
	completions = e.Complete("co")
	if len(completions) != 1 {
		t.Errorf("expected 1 completion, got %d", len(completions))
	}

	// Partial match with more results
	completions = e.Complete("c")
	if len(completions) != 3 {
		t.Errorf("expected 3 completions, got %d", len(completions))
	}

	// Full command
	completions = e.Complete("compile")
	if len(completions) != 1 {
		t.Errorf("expected 1 completion, got %d", len(completions))
	}
}

func TestCompletionEngineOptions(t *testing.T) {
	e := NewCompletionEngine()

	completions := e.Complete("compile --")
	if len(completions) != 3 {
		t.Errorf("expected 3 completions, got %d", len(completions))
	}

	completions = e.Complete("compile --in")
	if len(completions) != 1 {
		t.Errorf("expected 1 completion, got %d", len(completions))
	}
}

func TestCompletionEngineShellScripts(t *testing.T) {
	e := NewCompletionEngine()

	bash := e.GenerateBashCompletion()
	if !strings.Contains(bash, "_naeos_completions") {
		t.Error("expected bash completion function")
	}

	zsh := e.GenerateZshCompletion()
	if !strings.Contains(zsh, "_naeos") {
		t.Error("expected zsh completion function")
	}

	ps := e.GeneratePowerShellCompletion()
	if !strings.Contains(ps, "Register-ArgumentCompleter") {
		t.Error("expected PowerShell completer")
	}
}

func TestSnippetManager(t *testing.T) {
	sm := NewSnippetManager()

	snippet, ok := sm.Get("project")
	if !ok {
		t.Error("expected snippet to exist")
	}
	if !strings.Contains(snippet, "name:") {
		t.Error("expected snippet to contain name")
	}

	snippets := sm.List()
	if len(snippets) != 3 {
		t.Errorf("expected 3 snippets, got %d", len(snippets))
	}
}

func TestSnippetManagerAdd(t *testing.T) {
	sm := NewSnippetManager()

	sm.Add("custom", "custom snippet")
	snippet, ok := sm.Get("custom")
	if !ok || snippet != "custom snippet" {
		t.Error("expected custom snippet")
	}
}

func TestDevExperienceStack(t *testing.T) {
	stack := NewStack()

	if stack.Extension == nil {
		t.Error("expected extension")
	}
	if stack.Engine == nil {
		t.Error("expected engine")
	}
	if stack.Snippets == nil {
		t.Error("expected snippets")
	}
}
