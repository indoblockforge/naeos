package parser

import "testing"

func TestNewParserParsesJSON(t *testing.T) {
	p := NewParser()
	input := `{"project":{"name":"demo"},"version":1}`

	doc, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if doc.Data == nil {
		t.Fatal("expected parsed data")
	}
}

func TestNewParserParsesYAML(t *testing.T) {
	p := NewParser()
	input := "project:\n  name: demo\nversion: 1\n"

	doc, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if doc.Data == nil {
		t.Fatal("expected parsed data")
	}
}

func TestNewParserRejectsInvalidInput(t *testing.T) {
	p := NewParser()
	input := "{unclosed"

	doc, err := p.Parse(input)
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
	if doc != nil {
		t.Fatal("expected nil document on invalid input")
	}
}
