package cicd

import (
	"fmt"
	"strings"
)

type JenkinsGenerator struct{}

func (g *JenkinsGenerator) Name() string {
	return "Jenkins"
}

func (g *JenkinsGenerator) Generate(config *PipelineConfig) (string, error) {
	var sb strings.Builder

	sb.WriteString("pipeline {\n")
	sb.WriteString("    agent any\n\n")

	// Environment
	if len(config.Secrets) > 0 {
		sb.WriteString("    environment {\n")
		for _, secret := range config.Secrets {
			sb.WriteString(fmt.Sprintf("        %s = credentials('%s')\n", strings.ToUpper(secret), secret))
		}
		sb.WriteString("    }\n\n")
	}

	// Stages
	sb.WriteString("    stages {\n")

	// Build stage
	sb.WriteString("        stage('Build') {\n")
	sb.WriteString("            steps {\n")
	for _, lang := range config.Languages {
		switch lang {
		case "go":
			sb.WriteString("                sh 'go build ./...'\n")
		case "node", "typescript":
			sb.WriteString("                sh 'npm ci'\n")
			sb.WriteString("                sh 'npm run build'\n")
		case "python":
			sb.WriteString("                sh 'pip install -r requirements.txt'\n")
		case "java":
			sb.WriteString("                sh 'mvn clean compile'\n")
		case "rust":
			sb.WriteString("                sh 'cargo build --release'\n")
		}
	}
	sb.WriteString("            }\n")
	sb.WriteString("        }\n\n")

	// Test stage
	sb.WriteString("        stage('Test') {\n")
	sb.WriteString("            steps {\n")
	for _, lang := range config.Languages {
		switch lang {
		case "go":
			sb.WriteString("                sh 'go test ./...'\n")
		case "node", "typescript":
			sb.WriteString("                sh 'npm test'\n")
		case "python":
			sb.WriteString("                sh 'pytest'\n")
		case "java":
			sb.WriteString("                sh 'mvn test'\n")
		case "rust":
			sb.WriteString("                sh 'cargo test'\n")
		}
	}
	sb.WriteString("            }\n")
	sb.WriteString("        }\n\n")

	// Deploy stage
	sb.WriteString("        stage('Deploy') {\n")
	sb.WriteString("            when {\n")
	sb.WriteString("                branch 'main'\n")
	sb.WriteString("            }\n")
	sb.WriteString("            steps {\n")
	sb.WriteString("                echo 'Deploying...'\n")
	for _, step := range config.Steps {
		sb.WriteString(fmt.Sprintf("                sh '%s'\n", step.Command))
	}
	sb.WriteString("            }\n")
	sb.WriteString("        }\n")

	sb.WriteString("    }\n\n")

	// Post
	sb.WriteString("    post {\n")
	sb.WriteString("        always {\n")
	sb.WriteString("            cleanWs()\n")
	sb.WriteString("        }\n")
	sb.WriteString("        success {\n")
	sb.WriteString("            echo 'Pipeline succeeded!'\n")
	sb.WriteString("        }\n")
	sb.WriteString("        failure {\n")
	sb.WriteString("            echo 'Pipeline failed!'\n")
	sb.WriteString("        }\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}
