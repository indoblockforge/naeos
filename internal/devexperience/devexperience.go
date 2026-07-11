package devexperience

import (
	"fmt"
	"strings"
)

// VS Code Extension

type VSCodeExtension struct {
	Name        string
	Version     string
	Description string
	Author      string
	Languages   []string
	Features    []string
}

func NewVSCodeExtension(name, version, description, author string, languages []string) *VSCodeExtension {
	return &VSCodeExtension{
		Name:        name,
		Version:     version,
		Description: description,
		Author:      author,
		Languages:   languages,
		Features:    []string{"syntax highlighting", "autocomplete", "linting"},
	}
}

func (e *VSCodeExtension) GeneratePackageJSON() string {
	return fmt.Sprintf(`{
  "name": "%s",
  "displayName": "%s",
  "description": "%s",
  "version": "%s",
  "publisher": "%s",
  "engines": {
    "vscode": "^1.80.0"
  },
  "categories": ["Programming Languages", "Linters"],
  "contributes": {
    "languages": [%s],
    "commands": [
      {
        "command": "naeos.compile",
        "title": "NAEOS: Compile Project"
      },
      {
        "command": "naeos.validate",
        "title": "NAEOS: Validate Spec"
      }
    ]
  }
}`,
		e.Name,
		e.Name,
		e.Description,
		e.Version,
		e.Author,
		e.generateLanguagesJSON(),
	)
}

func (e *VSCodeExtension) generateLanguagesJSON() string {
	var langs []string
	for _, lang := range e.Languages {
		langs = append(langs, fmt.Sprintf(`"%s"`, lang))
	}
	return strings.Join(langs, ",")
}

func (e *VSCodeExtension) GenerateSyntaxJSON() string {
	return `{
  "scopeName": "source.naeos",
  "fileTypes": ["naeos.yaml", "naeos.yml", "naeos.json"],
  "patterns": [
    {
      "match": "^\\s*(name|version|description):",
      "name": "keyword.other.naeos"
    },
    {
      "match": "^\\s*(language|framework|type):",
      "name": "keyword.control.naeos"
    },
    {
      "match": "^\\s*(dependencies|adapters|plugins):",
      "name": "keyword.declaration.naeos"
    }
  ]
}`
}

// CLI Completion

type CompletionEngine struct {
	commands []string
	options  map[string][]string
}

func NewCompletionEngine() *CompletionEngine {
	e := &CompletionEngine{
		commands: []string{"init", "compile", "validate", "watch", "api", "ws", "graphql", "monitor", "auth", "db", "search", "workflow", "gateway", "cloud", "cicd", "pluginsdk"},
		options: map[string][]string{
			"init":      {"--name", "--type", "--language", "--framework"},
			"compile":   {"--input", "--output", "--language"},
			"validate":  {"--input", "--strict"},
			"watch":     {"--input", "--debounce"},
			"api":       {"--port", "--auth", "--secret"},
			"ws":        {"--port"},
			"graphql":   {"--port"},
			"monitor":   {"--port"},
			"auth":      {"login", "logout", "whoami"},
			"db":        {"connect", "disconnect", "migrate"},
			"search":    {"index", "query", "delete"},
			"workflow":  {"create", "list", "approve"},
			"gateway":   {"start", "stop"},
			"cloud":     {"deploy", "plan", "export"},
			"cicd":      {"generate", "list"},
			"pluginsdk": {"list", "info"},
		},
	}
	return e
}

func (e *CompletionEngine) Complete(input string) []string {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return e.commands
	}

	if len(parts) == 1 {
		var matches []string
		for _, cmd := range e.commands {
			if strings.HasPrefix(cmd, parts[0]) {
				matches = append(matches, cmd)
			}
		}
		return matches
	}

	if len(parts) == 2 {
		cmd := parts[0]
		if opts, ok := e.options[cmd]; ok {
			prefix := parts[1]
			var matches []string
			for _, opt := range opts {
				if strings.HasPrefix(opt, prefix) {
					matches = append(matches, opt)
				}
			}
			return matches
		}
	}

	return nil
}

func (e *CompletionEngine) GenerateBashCompletion() string {
	return `#!/bin/bash
_naeos_completions() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    commands="init compile validate watch api ws graphql monitor auth db search workflow gateway cloud cicd pluginsdk"
    
    if [ ${COMP_CWORD} -eq 1 ]; then
        COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
    fi
}
complete -F _naeos_completions naeos`
}

func (e *CompletionEngine) GenerateZshCompletion() string {
	return `#compdef naeos
_naeos() {
    _arguments \
        '1:command:(init compile validate watch api ws graphql monitor auth db search workflow gateway cloud cicd pluginsdk)' \
        '*::arg:->args'
}
_naeos "$@"`
}

func (e *CompletionEngine) GeneratePowerShellCompletion() string {
	return `Register-ArgumentCompleter -Native -CommandName naeos -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commands = @('init', 'compile', 'validate', 'watch', 'api', 'ws', 'graphql', 'monitor', 'auth', 'db', 'search', 'workflow', 'gateway', 'cloud', 'cicd', 'pluginsdk')
    $commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}`
}

// Snippets

type SnippetManager struct {
	snippets map[string]string
}

func NewSnippetManager() *SnippetManager {
	sm := &SnippetManager{
		snippets: make(map[string]string),
	}

	sm.snippets["project"] = `name: my-project
version: 0.1.0
description: A new NAEOS project

language: go
framework: gin

dependencies:
  - github.com/gin-gonic/gin

adapters:
  - name: auth
    type: oauth2
  - name: db
    type: postgresql

plugins:
  - name: logger
    version: 1.0.0
`
	sm.snippets["api-endpoint"] = `func HandleRequest(c *gin.Context) {
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    result, err := Service.DoSomething(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, result)
}
`
	sm.snippets["test"] = `func TestSomething(t *testing.T) {
    t.Run("valid input", func(t *testing.T) {
        result, err := DoSomething("input")
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if result != expected {
            t.Errorf("expected %v, got %v", expected, result)
        }
    })
}
`

	return sm
}

func (sm *SnippetManager) Get(name string) (string, bool) {
	snippet, ok := sm.snippets[name]
	return snippet, ok
}

func (sm *SnippetManager) List() []string {
	names := make([]string, 0, len(sm.snippets))
	for name := range sm.snippets {
		names = append(names, name)
	}
	return names
}

func (sm *SnippetManager) Add(name, snippet string) {
	sm.snippets[name] = snippet
}

// Dev Experience Stack

type Stack struct {
	Extension *VSCodeExtension
	Engine    *CompletionEngine
	Snippets  *SnippetManager
}

func NewStack() *Stack {
	return &Stack{
		Extension: NewVSCodeExtension("naeos", "1.0.0", "NAEOS project support", "NAEOS", []string{"yaml", "json"}),
		Engine:    NewCompletionEngine(),
		Snippets:  NewSnippetManager(),
	}
}
