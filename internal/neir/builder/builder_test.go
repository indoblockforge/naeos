package builder

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/specification/resolver"
)

func TestBuilderCreatesNEIRFromResolvedSpec(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"project": "acme-api",
		"modules": []map[string]any{{"name": "auth", "path": "./internal/auth"}},
		"services": []map[string]any{{"name": "gateway", "kind": "http", "port": 8080}},
	}}

	builder := NewBuilder()
	neir, err := builder.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Project == nil || neir.Project.Name != "acme-api" {
		t.Fatalf("expected project acme-api, got %v", neir.Project)
	}
	if len(neir.Modules) != 1 {
		t.Fatalf("expected one module, got %d", len(neir.Modules))
	}
	if len(neir.Services) != 1 {
		t.Fatalf("expected one service, got %d", len(neir.Services))
	}
	if neir.Services[0].Name != "gateway" {
		t.Fatalf("expected service gateway, got %s", neir.Services[0].Name)
	}
	if neir.Services[0].Port != 8080 {
		t.Fatalf("expected service port 8080, got %d", neir.Services[0].Port)
	}
	if neir.Architecture != nil {
		t.Fatalf("expected nil architecture, got %v", neir.Architecture)
	}
}

func TestBuilderExtractsArchitecture(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"project": "acme-api",
		"modules": []map[string]any{{"name": "core", "path": "./internal/core"}},
		"architecture": map[string]any{"pattern": "clean", "description": "Clean architecture"},
	}}

	b := NewBuilder()
	neir, err := b.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Architecture == nil {
		t.Fatal("expected architecture to be set")
	}
	if neir.Architecture.Pattern != "clean" {
		t.Fatalf("expected pattern clean, got %s", neir.Architecture.Pattern)
	}
	if neir.Architecture.Description != "Clean architecture" {
		t.Fatalf("expected description, got %s", neir.Architecture.Description)
	}
}

func TestBuilderWithNilInput(t *testing.T) {
	b := NewBuilder()
	_, err := b.Build(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestBuilderExtractsDeployment(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"project": "acme-api",
		"deployment": map[string]any{
			"strategy": "canary",
			"environments": []any{
				map[string]any{"name": "staging"},
				"production",
			},
		},
	}}

	b := NewBuilder()
	neir, err := b.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Deployment == nil {
		t.Fatal("expected deployment to be set")
	}
	if neir.Deployment.Strategy != "canary" {
		t.Fatalf("expected strategy canary, got %s", neir.Deployment.Strategy)
	}
	if len(neir.Deployment.Environments) != 2 {
		t.Fatalf("expected 2 environments, got %d", len(neir.Deployment.Environments))
	}
	if neir.Deployment.Environments[0].Name != "staging" {
		t.Fatalf("expected first env staging, got %s", neir.Deployment.Environments[0].Name)
	}
	if neir.Deployment.Environments[1].Name != "production" {
		t.Fatalf("expected second env production, got %s", neir.Deployment.Environments[1].Name)
	}
}

func TestBuilderExtractsTesting(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"project": "acme-api",
		"testing": map[string]any{
			"strategy": "unit",
			"coverage": "high",
		},
	}}

	b := NewBuilder()
	neir, err := b.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Testing == nil {
		t.Fatal("expected testing to be set")
	}
	if neir.Testing.Strategy != "unit" {
		t.Fatalf("expected strategy unit, got %s", neir.Testing.Strategy)
	}
	if neir.Testing.Coverage == nil {
		t.Fatal("expected coverage to be set")
	}
	if neir.Testing.Coverage.MinPercent != 80.0 {
		t.Fatalf("expected coverage 80.0, got %f", neir.Testing.Coverage.MinPercent)
	}
}

func TestBuilderExtractsTestingMediumCoverage(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"testing": map[string]any{"coverage": "medium"},
	}}

	b := NewBuilder()
	neir, err := b.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Testing.Coverage.MinPercent != 60.0 {
		t.Fatalf("expected coverage 60.0, got %f", neir.Testing.Coverage.MinPercent)
	}
}

func TestBuilderExtractsTestingLowCoverage(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"testing": map[string]any{"coverage": "low"},
	}}

	b := NewBuilder()
	neir, err := b.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Testing.Coverage.MinPercent != 40.0 {
		t.Fatalf("expected coverage 40.0, got %f", neir.Testing.Coverage.MinPercent)
	}
}
