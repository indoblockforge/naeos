package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/NAEOS-foundation/naeos/pkg/pipeline"
)

var version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	root := newRootCommand()
	root.SetArgs(args)
	return root.Execute()
}

func newRootCommand() *cobra.Command {
	var rootVerbose bool

	root := &cobra.Command{
		Use:           "naeos",
		Short:         "NAEOS CLI",
		Long:          "NAEOS is a declarative engineering runtime for specification-driven project delivery.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().BoolVar(&rootVerbose, "verbose", false, "enable verbose logging")

	root.AddCommand(newInitCommand())
	root.AddCommand(newRunCommand())
	root.AddCommand(newValidateCommand())
	root.AddCommand(newKernelCommand())
	root.AddCommand(newVersionCommand())
	return root
}

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show NAEOS version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "naeos %s\n", version); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	}
	return cmd
}

func newInitCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a default NAEOS config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			content := strings.Join([]string{
				"pipeline:",
				"  name: naeos-dev",
				"  mode: development",
				"  verbose: true",
				"  output_dir: ./out",
			}, "\n") + "\n"

			if err := os.WriteFile(output, []byte(content), 0o644); err != nil {
				return fmt.Errorf("write config: %w", err)
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", output); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "config.example.yaml", "path for the generated config file")
	return cmd
}

func newRunCommand() *cobra.Command {
	var configPath string
	var input string
	var inputFile string
	var outputFormat string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute the NAEOS pipeline",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			inputValue, err := loadInput(input, inputFile)
			if err != nil {
				return err
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}
			result, err := p.Run(inputValue)
			if err != nil {
				return fmt.Errorf("pipeline run failed: %w", err)
			}

			payload := map[string]any{
				"pipeline":   cfg.Name,
				"mode":       cfg.Mode,
				"verbose":    cfg.Verbose,
				"output_dir": cfg.OutputDir,
				"artifacts":  len(result.Artifacts),
				"tasks":      len(result.Tasks),
			}

			var rendered []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				data, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
				rendered = append(data, '\n')
			case "yaml":
				data, err := yaml.Marshal(payload)
				if err != nil {
					return fmt.Errorf("encode yaml output: %w", err)
				}
				rendered = data
			default:
				rendered = []byte(fmt.Sprintf("pipeline=%s mode=%s verbose=%t output_dir=%s\nartifacts=%d tasks=%d\n", result.NEIR.Project, cfg.Mode, cfg.Verbose, cfg.OutputDir, len(result.Artifacts), len(result.Tasks)))
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, rendered, 0o644); err != nil {
					return fmt.Errorf("write output file: %w", err)
				}
				return nil
			}

			if _, err := cmd.OutOrStdout().Write(rendered); err != nil {
				return fmt.Errorf("write output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input to process")
	cmd.Flags().StringVar(&inputFile, "input-file", "", "path to a specification file")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "optional file path to write the formatted output")

	return cmd
}

func newValidateCommand() *cobra.Command {
	var configPath string
	var input string
	var inputFile string
	var outputFormat string
	var outputFile string

	cmd := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v"},
		Short:   "Validate a specification using the NAEOS pipeline",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			inputValue, err := loadInput(input, inputFile)
			if err != nil {
				return err
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}
			result, err := p.Validate(inputValue)
			if err != nil {
				return fmt.Errorf("pipeline validate failed: %w", err)
			}

			payload := map[string]any{
				"pipeline":   cfg.Name,
				"mode":       cfg.Mode,
				"verbose":    cfg.Verbose,
				"output_dir": cfg.OutputDir,
				"status":     "valid",
				"project":    result.NEIR.Project,
				"source_len": len(result.Source),
			}

			var rendered []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				data, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
				rendered = append(data, '\n')
			case "yaml":
				data, err := yaml.Marshal(payload)
				if err != nil {
					return fmt.Errorf("encode yaml output: %w", err)
				}
				rendered = data
			default:
				rendered = []byte(fmt.Sprintf("config=%s mode=%s verbose=%t output_dir=%s\nstatus=valid project=%v source_len=%d\n",
					cfg.Name, cfg.Mode, cfg.Verbose, cfg.OutputDir, result.NEIR.Project, len(result.Source)))
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, rendered, 0o644); err != nil {
					return fmt.Errorf("write output file: %w", err)
				}
				return nil
			}

			if _, err := cmd.OutOrStdout().Write(rendered); err != nil {
				return fmt.Errorf("write output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.Flags().StringVar(&input, "input", "", "specification input to process")
	cmd.Flags().StringVar(&inputFile, "input-file", "", "path to a specification file")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "optional file path to write the formatted output")

	return cmd
}

func newKernelCommand() *cobra.Command {
	var configPath string
	var outputFormat string
	var topic string
	var payload string

	cmd := &cobra.Command{
		Use:   "kernel",
		Short: "Inspect the NAEOS kernel and service registry",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "services",
		Short: "List registered kernel services",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			services := p.RegisteredKernelServices()
			var rendered []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				data, err := json.MarshalIndent(services, "", "  ")
				if err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
				rendered = append(data, '\n')
			case "yaml":
				data, err := yaml.Marshal(services)
				if err != nil {
					return fmt.Errorf("encode yaml output: %w", err)
				}
				rendered = data
			default:
				rendered = []byte(strings.Join(services, "\n") + "\n")
			}

			if _, err := cmd.OutOrStdout().Write(rendered); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "metrics",
		Short: "Show kernel telemetry metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			metrics := p.KernelMetrics()
			var rendered []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				data, err := json.MarshalIndent(metrics, "", "  ")
				if err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
				rendered = append(data, '\n')
			case "yaml":
				data, err := yaml.Marshal(metrics)
				if err != nil {
					return fmt.Errorf("encode yaml output: %w", err)
				}
				rendered = data
			default:
				rendered = []byte(fmt.Sprintf("events=%d\nlast_event=%s\n", metrics.Events, metrics.LastEvent.Name))
			}

			if _, err := cmd.OutOrStdout().Write(rendered); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "events",
		Short: "List active kernel event topics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			topics := p.KernelTopics()
			var rendered []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				data, err := json.MarshalIndent(topics, "", "  ")
				if err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
				rendered = append(data, '\n')
			case "yaml":
				data, err := yaml.Marshal(topics)
				if err != nil {
					return fmt.Errorf("encode yaml output: %w", err)
				}
				rendered = data
			default:
				rendered = []byte(strings.Join(topics, "\n") + "\n")
			}

			if _, err := cmd.OutOrStdout().Write(rendered); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "publish",
		Short: "Publish an event to the kernel event bus",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}
			if topic == "" {
				return fmt.Errorf("missing required --topic")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			if err := p.Publish(topic, payload); err != nil {
				return err
			}
			if _, err := cmd.OutOrStdout().Write([]byte(fmt.Sprintf("published topic=%s payload=%v\n", topic, payload))); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe to a kernel event topic and optionally publish a sample payload",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("missing required --config")
			}
			if topic == "" {
				return fmt.Errorf("missing required --topic")
			}

			cfg, err := pipeline.ConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p, err := pipeline.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to construct pipeline: %w", err)
			}

			var received any
			if err := p.Subscribe(topic, func(payload any) {
				received = payload
			}); err != nil {
				return err
			}

			if payload != "" {
				if err := p.Publish(topic, payload); err != nil {
					return err
				}
			}

			var rendered []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				data, err := json.MarshalIndent(map[string]any{"topic": topic, "received": received}, "", "  ")
				if err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
				rendered = append(data, '\n')
			case "yaml":
				data, err := yaml.Marshal(map[string]any{"topic": topic, "received": received})
				if err != nil {
					return fmt.Errorf("encode yaml output: %w", err)
				}
				rendered = data
			default:
				rendered = []byte(fmt.Sprintf("topic=%s received=%v\n", topic, received))
			}

			if _, err := cmd.OutOrStdout().Write(rendered); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	})

	cmd.PersistentFlags().StringVar(&configPath, "config", "", "path to JSON or YAML config file")
	cmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "output format: text, json, or yaml")
	cmd.PersistentFlags().StringVar(&topic, "topic", "", "kernel event topic")
	cmd.PersistentFlags().StringVar(&payload, "payload", "", "event payload to publish")
	return cmd
}

func loadInput(input, inputFile string) (string, error) {
	if input == "" && inputFile == "" {
		return "", fmt.Errorf("missing required --input or --input-file")
	}
	if input != "" && inputFile != "" {
		return "", fmt.Errorf("cannot use both --input and --input-file")
	}
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("read input file: %w", err)
		}
		return string(data), nil
	}
	return input, nil
}
