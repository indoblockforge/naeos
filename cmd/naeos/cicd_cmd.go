package main

import (
	"fmt"
	"os"
	"github.com/NAEOS-foundation/naeos/internal/cicd"
	"github.com/spf13/cobra"
)

var (
	cicdPlatform string
	cicdOutput   string
)

func newCICDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cicd",
		Short: "Generate CI/CD pipeline configuration",
		Long:  `Generate CI/CD pipeline configuration for GitHub Actions, GitLab CI, or Jenkins.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			platform := cicd.CICDPlatform(cicdPlatform)
			gen, err := cicd.GetGenerator(platform)
			if err != nil {
				return err
			}

			config := &cicd.PipelineConfig{
				Project:  "myapp",
				Platform: platform,
				Languages: []string{"go"},
				Trigger: cicd.TriggerConfig{
					OnPush: true,
					OnPR:   true,
				},
			}

			output, err := gen.Generate(config)
			if err != nil {
				return err
			}

			if cicdOutput != "" {
				return os.WriteFile(cicdOutput, []byte(output), 0644)
			}

			fmt.Println(output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&cicdPlatform, "platform", "p", "github", "CI/CD platform (github, gitlab, jenkins)")
	cmd.Flags().StringVarP(&cicdOutput, "output", "o", "", "Output file path")

	return cmd
}
