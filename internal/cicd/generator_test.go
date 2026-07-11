package cicd

import (
	"testing"
)

func TestGetGeneratorGitHub(t *testing.T) {
	gen, err := GetGenerator(GitHubActions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gen.Name() != "GitHub Actions" {
		t.Errorf("expected name 'GitHub Actions', got %s", gen.Name())
	}
}

func TestGetGeneratorGitLab(t *testing.T) {
	gen, err := GetGenerator(GitLabCI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gen.Name() != "GitLab CI" {
		t.Errorf("expected name 'GitLab CI', got %s", gen.Name())
	}
}

func TestGetGeneratorJenkins(t *testing.T) {
	gen, err := GetGenerator(Jenkins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gen.Name() != "Jenkins" {
		t.Errorf("expected name 'Jenkins', got %s", gen.Name())
	}
}

func TestGetGeneratorInvalid(t *testing.T) {
	_, err := GetGenerator("invalid")
	if err == nil {
		t.Error("expected error for invalid platform")
	}
}

func TestGitHubActionsGenerate(t *testing.T) {
	gen := &GitHubActionsGenerator{}
	config := &PipelineConfig{
		Project:   "myapp",
		Platform:  GitHubActions,
		Languages: []string{"go"},
		Trigger: TriggerConfig{
			OnPush: true,
			OnPR:   true,
		},
	}

	output, err := gen.Generate(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestGitLabCIGenerate(t *testing.T) {
	gen := &GitLabCIGenerator{}
	config := &PipelineConfig{
		Project:   "myapp",
		Platform:  GitLabCI,
		Languages: []string{"node"},
		Trigger: TriggerConfig{
			OnPush: true,
		},
	}

	output, err := gen.Generate(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestJenkinsGenerate(t *testing.T) {
	gen := &JenkinsGenerator{}
	config := &PipelineConfig{
		Project:   "myapp",
		Platform:  Jenkins,
		Languages: []string{"python"},
		Trigger: TriggerConfig{
			OnPush: true,
		},
		Secrets: []string{"AWS_ACCESS_KEY"},
	}

	output, err := gen.Generate(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == "" {
		t.Error("expected non-empty output")
	}
}
