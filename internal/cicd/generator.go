package cicd

import (
	"fmt"
)

type CICDPlatform string

const (
	GitHubActions CICDPlatform = "github"
	GitLabCI      CICDPlatform = "gitlab"
	Jenkins       CICDPlatform = "jenkins"
)

type PipelineConfig struct {
	Project    string
	Platform   CICDPlatform
	Languages  []string
	Steps      []PipelineStep
	Trigger    TriggerConfig
	Secrets    []string
}

type PipelineStep struct {
	Name    string
	Command string
	Env     map[string]string
}

type TriggerConfig struct {
	OnPush    bool
	OnPR      bool
	OnRelease bool
	Schedule  string
}

type PipelineGenerator interface {
	Name() string
	Generate(config *PipelineConfig) (string, error)
}

func GetGenerator(platform CICDPlatform) (PipelineGenerator, error) {
	switch platform {
	case GitHubActions:
		return &GitHubActionsGenerator{}, nil
	case GitLabCI:
		return &GitLabCIGenerator{}, nil
	case Jenkins:
		return &JenkinsGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
