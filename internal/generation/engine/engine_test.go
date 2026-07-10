package engine

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
)

func TestGeneratorCreatesArtifactsFromNEIR(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "acme-api"},
		Modules: []module.Module{{Name: "auth", Path: "./internal/auth"}},
	}

	engine := NewEngine()
	artifacts, err := engine.Generate(neir)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if len(artifacts) < 2 {
		t.Fatalf("expected at least two artifacts, got %d", len(artifacts))
	}
	if artifacts[0].Path != "README.md" {
		t.Fatalf("expected README artifact first, got %s", artifacts[0].Path)
	}
}

func TestGenerateForLanguageGo(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "acme-api"},
		Modules: []module.Module{{Name: "auth", Path: "./internal/auth"}},
	}

	engine := NewEngine()
	artifacts, err := engine.GenerateForLanguage(neir, language.LanguageGo)
	if err != nil {
		t.Fatalf("GenerateForLanguage returned error: %v", err)
	}

	if len(artifacts) < 3 {
		t.Fatalf("expected at least 3 artifacts (go.mod, main.go, Dockerfile), got %d", len(artifacts))
	}

	foundGoMod := false
	foundMain := false
	foundDockerfile := false
	for _, a := range artifacts {
		if a.Path == "go.mod" {
			foundGoMod = true
			if !contains(a.Content, "module github.com/example/acme-api") {
				t.Errorf("go.mod should contain module path")
			}
		}
		if a.Path == "src/main.go" {
			foundMain = true
			if !contains(a.Content, "hello from acme-api") {
				t.Errorf("main.go should contain project name")
			}
		}
		if a.Path == "Dockerfile" {
			foundDockerfile = true
			if !contains(a.Content, "golang:1.22") {
				t.Errorf("Dockerfile should use Go image")
			}
		}
	}

	if !foundGoMod {
		t.Error("expected go.mod artifact")
	}
	if !foundMain {
		t.Error("expected src/main.go artifact")
	}
	if !foundDockerfile {
		t.Error("expected Dockerfile artifact")
	}
}

func TestGenerateForLanguageTypeScript(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "web-app"},
	}

	engine := NewEngine()
	artifacts, err := engine.GenerateForLanguage(neir, language.LanguageTypeScript)
	if err != nil {
		t.Fatalf("GenerateForLanguage returned error: %v", err)
	}

	foundPackageJson := false
	foundMain := false
	for _, a := range artifacts {
		if a.Path == "package.json" {
			foundPackageJson = true
			if !contains(a.Content, "web-app") {
				t.Errorf("package.json should contain project name")
			}
		}
		if a.Path == "src/main.ts" {
			foundMain = true
			if !contains(a.Content, "hello from web-app") {
				t.Errorf("main.ts should contain project name")
			}
		}
	}

	if !foundPackageJson {
		t.Error("expected package.json artifact")
	}
	if !foundMain {
		t.Error("expected src/main.ts artifact")
	}
}

func TestGenerateForLanguagePython(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "ml-service"},
	}

	engine := NewEngine()
	artifacts, err := engine.GenerateForLanguage(neir, language.LanguagePython)
	if err != nil {
		t.Fatalf("GenerateForLanguage returned error: %v", err)
	}

	foundPyproject := false
	for _, a := range artifacts {
		if a.Path == "pyproject.toml" {
			foundPyproject = true
			if !contains(a.Content, "ml-service") {
				t.Errorf("pyproject.toml should contain project name")
			}
		}
	}

	if !foundPyproject {
		t.Error("expected pyproject.toml artifact")
	}
}

func TestGenerateForLanguageWithModules(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "fullstack"},
		Modules: []module.Module{
			{Name: "auth", Path: "./internal/auth"},
			{Name: "api", Path: "./internal/api"},
		},
	}

	engine := NewEngine()
	artifacts, err := engine.GenerateForLanguage(neir, language.LanguageGo)
	if err != nil {
		t.Fatalf("GenerateForLanguage returned error: %v", err)
	}

	authFiles := 0
	apiFiles := 0
	for _, a := range artifacts {
		if contains([]byte(a.Path), "auth") {
			authFiles++
		}
		if contains([]byte(a.Path), "api") {
			apiFiles++
		}
	}

	if authFiles == 0 {
		t.Error("expected auth module files")
	}
	if apiFiles == 0 {
		t.Error("expected api module files")
	}
}

func TestGenerateForLanguageNilNEIR(t *testing.T) {
	engine := NewEngine()
	_, err := engine.GenerateForLanguage(nil, language.LanguageGo)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestGenerateForLanguageInvalidLanguage(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "test"},
	}

	engine := NewEngine()
	_, err := engine.GenerateForLanguage(neir, "invalid-lang")
	if err == nil {
		t.Fatal("expected error for invalid language")
	}
}

func contains(haystack []byte, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 && 
		string(haystack) != "" && 
		(len(haystack) >= len(needle)) && 
		searchBytes(haystack, needle)
}

func searchBytes(data []byte, s string) bool {
	for i := 0; i <= len(data)-len(s); i++ {
		if string(data[i:i+len(s)]) == s {
			return true
		}
	}
	return false
}
