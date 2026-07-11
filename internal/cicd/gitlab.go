package cicd

import (
	"fmt"
	"strings"
)

type GitLabCIGenerator struct{}

func (g *GitLabCIGenerator) Name() string {
	return "GitLab CI"
}

func (g *GitLabCIGenerator) Generate(config *PipelineConfig) (string, error) {
	var sb strings.Builder

	sb.WriteString("stages:\n")
	sb.WriteString("  - build\n")
	sb.WriteString("  - test\n")
	sb.WriteString("  - deploy\n\n")

	// Default image
	if len(config.Languages) > 0 {
		switch config.Languages[0] {
		case "go":
			sb.WriteString("image: golang:1.22\n\n")
		case "node", "typescript":
			sb.WriteString("image: node:20\n\n")
		case "python":
			sb.WriteString("image: python:3.12\n\n")
		case "java":
			sb.WriteString("image: maven:3.9-eclipse-temurin-21\n\n")
		case "rust":
			sb.WriteString("image: rust:latest\n\n")
		}
	}

	// Build job
	sb.WriteString("build:\n")
	sb.WriteString("  stage: build\n")
	sb.WriteString("  script:\n")
	for _, lang := range config.Languages {
		switch lang {
		case "go":
			sb.WriteString("    - go build ./...\n")
		case "node", "typescript":
			sb.WriteString("    - npm ci\n")
			sb.WriteString("    - npm run build\n")
		case "python":
			sb.WriteString("    - pip install -r requirements.txt\n")
		case "java":
			sb.WriteString("    - mvn clean compile\n")
		case "rust":
			sb.WriteString("    - cargo build --release\n")
		}
	}
	sb.WriteString("\n")

	// Test job
	sb.WriteString("test:\n")
	sb.WriteString("  stage: test\n")
	sb.WriteString("  script:\n")
	for _, lang := range config.Languages {
		switch lang {
		case "go":
			sb.WriteString("    - go test ./...\n")
		case "node", "typescript":
			sb.WriteString("    - npm test\n")
		case "python":
			sb.WriteString("    - pytest\n")
		case "java":
			sb.WriteString("    - mvn test\n")
		case "rust":
			sb.WriteString("    - cargo test\n")
		}
	}
	sb.WriteString("\n")

	// Deploy job
	sb.WriteString("deploy:\n")
	sb.WriteString("  stage: deploy\n")
	sb.WriteString("  only:\n")
	sb.WriteString("    - main\n")
	sb.WriteString("  script:\n")
	sb.WriteString("    - echo 'Deploying...'\n")

	// Custom steps
	for _, step := range config.Steps {
		sb.WriteString(fmt.Sprintf("    - %s\n", step.Command))
	}
	sb.WriteString("\n")

	return sb.String(), nil
}
