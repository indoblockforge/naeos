package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitCreatesConfigFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "config.yaml")

	err := run([]string{"init", "--output", output})
	if err != nil {
		t.Fatalf("run init returned error: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read generated config: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected generated config file to contain content")
	}
}

func TestValidateUsesConfigFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	err := run([]string{"validate", "--config", configPath, "--input", "sample specification"})
	if err != nil {
		t.Fatalf("run validate returned error: %v", err)
	}
}

func TestRunSupportsJSONOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	err := run([]string{"run", "--config", configPath, "--input", "sample specification", "--output", "json"})
	if err != nil {
		t.Fatalf("run run returned error: %v", err)
	}
}

func TestRunSupportsYAMLOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	err := run([]string{"run", "--config", configPath, "--input", "sample specification", "--output", "yaml"})
	if err != nil {
		t.Fatalf("run run returned error: %v", err)
	}
}

func TestRunWritesOutputToFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	outputPath := filepath.Join(dir, "result.json")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	err := run([]string{"run", "--config", configPath, "--input", "sample specification", "--output", "json", "--output-file", outputPath})
	if err != nil {
		t.Fatalf("run run returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected output file to contain content")
	}
}

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	root.SilenceErrors = true
	root.SilenceUsage = true
	_, err := root.ExecuteC()
	return buf.String(), err
}

func TestValidateCobraJSONOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "validate", "--config", configPath, "--input", "sample specification", "--output", "json")
	if err != nil {
		t.Fatalf("execute validate failed: %v", err)
	}

	if !strings.Contains(output, `"status": "valid"`) {
		t.Fatalf("expected json output, got %q", output)
	}
}

func TestValidateCobraYAMLOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "validate", "--config", configPath, "--input", "sample specification", "--output", "yaml")
	if err != nil {
		t.Fatalf("execute validate failed: %v", err)
	}

	if !strings.Contains(output, "status: valid") {
		t.Fatalf("expected yaml output, got %q", output)
	}
}

func TestValidateCobraAliasV(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	if _, err := executeCommand(root, "v", "--config", configPath, "--input", "sample specification"); err != nil {
		t.Fatalf("execute v failed: %v", err)
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	os.Stdout = oldStdout
	return string(out)
}

func TestValidateSupportsJSONOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	output := captureOutput(t, func() {
		if err := run([]string{"validate", "--config", configPath, "--input", "sample specification", "--output", "json"}); err != nil {
			t.Fatalf("run validate returned error: %v", err)
		}
	})

	if !strings.Contains(output, `"status": "valid"`) {
		t.Fatalf("expected JSON output, got %q", output)
	}
}

func TestValidateSupportsYAMLOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	output := captureOutput(t, func() {
		if err := run([]string{"validate", "--config", configPath, "--input", "sample specification", "--output", "yaml"}); err != nil {
			t.Fatalf("run validate returned error: %v", err)
		}
	})

	if !strings.Contains(output, "status: valid") {
		t.Fatalf("expected YAML output, got %q", output)
	}
}

func TestValidateAliasV(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if err := run([]string{"v", "--config", configPath, "--input", "sample specification"}); err != nil {
		t.Fatalf("run v returned error: %v", err)
	}
}

func TestVersionCommand(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "version")
	if err != nil {
		t.Fatalf("execute version failed: %v", err)
	}
	if !strings.Contains(output, "naeos ") {
		t.Fatalf("expected version output, got %q", output)
	}
}

func TestKernelServicesCommand(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "kernel", "services", "--config", configPath, "--output", "text")
	if err != nil {
		t.Fatalf("execute kernel services failed: %v", err)
	}
	if !strings.Contains(output, "pipeline") || !strings.Contains(output, "parser") {
		t.Fatalf("expected kernel service list, got %q", output)
	}
}

func TestKernelMetricsCommand(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "kernel", "metrics", "--config", configPath, "--output", "text")
	if err != nil {
		t.Fatalf("execute kernel metrics failed: %v", err)
	}
	if !strings.Contains(output, "events=") {
		t.Fatalf("expected kernel metrics output, got %q", output)
	}
}

func TestKernelEventsCommand(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "kernel", "events", "--config", configPath, "--output", "text")
	if err != nil {
		t.Fatalf("execute kernel events failed: %v", err)
	}
	if strings.TrimSpace(output) != "" {
		t.Fatalf("expected no events output when no topics are registered, got %q", output)
	}
}

func TestKernelPublishSubscribeCommand(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	publishOutput, err := executeCommand(root, "kernel", "publish", "--config", configPath, "--topic", "test", "--payload", "hello", "--output", "text")
	if err != nil {
		t.Fatalf("execute kernel publish failed: %v", err)
	}
	if !strings.Contains(publishOutput, "published topic=test payload=hello") {
		t.Fatalf("expected publish output, got %q", publishOutput)
	}

	subscribeOutput, err := executeCommand(root, "kernel", "subscribe", "--config", configPath, "--topic", "test", "--payload", "hello", "--output", "text")
	if err != nil {
		t.Fatalf("execute kernel subscribe failed: %v", err)
	}
	if !strings.Contains(subscribeOutput, "topic=test") || !strings.Contains(subscribeOutput, "received=hello") {
		t.Fatalf("expected subscribe output, got %q", subscribeOutput)
	}
}

func TestRootVerboseFlag(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "--verbose", "validate", "--config", filepath.Join(t.TempDir(), "config.yaml"), "--input", "sample specification")
	if err == nil {
		t.Fatalf("expected validate to fail with missing config file, got output: %q", output)
	}
}
