package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newTestInitCmd creates a fresh init command for testing.
// This avoids state pollution between tests.
func newTestInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize mdpress configuration",
		RunE:  runInit,
	}
	cmd.Flags().BoolVar(&projectConfig, "project", false, "create local project config (./mdpress.yml)")
	return cmd
}

func TestInit_FreshGlobal(t *testing.T) {
	// Create isolated temp directory for HOME
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Reset the projectConfig flag
	projectConfig = false

	// Create a fresh command for testing
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Verify mdpress.yml exists and contains expected content
	configPath := filepath.Join(tmpDir, ".config", "mdpress", "mdpress.yml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml: %v", err)
	}
	if !strings.Contains(string(content), "version: \"1\"") {
		t.Error("mdpress.yml should contain version: \"1\"")
	}
	if !strings.Contains(string(content), "logo_position:") {
		t.Error("mdpress.yml should contain logo_position setting")
	}

	// Verify logos/ directory exists
	logosDir := filepath.Join(tmpDir, ".config", "mdpress", "logos")
	info, err := os.Stat(logosDir)
	if err != nil {
		t.Fatalf("logos directory should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("logos should be a directory")
	}

	// Verify templates/default.yml exists and contains expected content
	templatePath := filepath.Join(tmpDir, ".config", "mdpress", "templates", "default.yml")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read templates/default.yml: %v", err)
	}
	if !strings.Contains(string(templateContent), "mdpress template: default") {
		t.Error("templates/default.yml should contain template header")
	}
}

func TestInit_IdempotentGlobal(t *testing.T) {
	// Create isolated temp directory for HOME
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Reset the projectConfig flag
	projectConfig = false

	// First run - create files
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("first init command failed: %v", err)
	}

	// Record file contents after first run
	configPath := filepath.Join(tmpDir, ".config", "mdpress", "mdpress.yml")
	originalContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml: %v", err)
	}

	// Second run - should skip all files
	buf.Reset()
	projectConfig = false
	cmd2 := newTestInitCmd()
	cmd2.SetOut(buf)
	cmd2.SetArgs([]string{})
	err = cmd2.Execute()
	if err != nil {
		t.Fatalf("second init command failed: %v", err)
	}

	// Verify file contents are not modified
	newContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml after second run: %v", err)
	}
	if string(originalContent) != string(newContent) {
		t.Error("mdpress.yml should not be modified on second run")
	}

	// Verify output indicates skipping
	output := buf.String()
	if !strings.Contains(output, "already") {
		t.Error("second run output should indicate config already initialized")
	}
}

func TestInit_PartialGlobal(t *testing.T) {
	// Create isolated temp directory for HOME
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create only mdpress.yml manually (without logos/ and templates/)
	configDir := filepath.Join(tmpDir, ".config", "mdpress")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "mdpress.yml")
	originalContent := []byte("# Pre-existing config\nversion: \"1\"\n")
	err = os.WriteFile(configPath, originalContent, 0644)
	if err != nil {
		t.Fatalf("failed to write pre-existing config: %v", err)
	}

	// Reset the projectConfig flag
	projectConfig = false

	// Run init
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Verify logos/ directory is created
	logosDir := filepath.Join(configDir, "logos")
	info, err := os.Stat(logosDir)
	if err != nil {
		t.Fatalf("logos directory should be created: %v", err)
	}
	if !info.IsDir() {
		t.Error("logos should be a directory")
	}

	// Verify templates/default.yml is created
	templatePath := filepath.Join(configDir, "templates", "default.yml")
	_, err = os.Stat(templatePath)
	if err != nil {
		t.Fatalf("templates/default.yml should be created: %v", err)
	}

	// Verify mdpress.yml is skipped (original content preserved)
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml: %v", err)
	}
	if string(content) != string(originalContent) {
		t.Error("mdpress.yml should be skipped (original content should be preserved)")
	}

	// Verify output shows skipped message for config file
	output := buf.String()
	if !strings.Contains(output, "Skipped") || !strings.Contains(output, "mdpress.yml") {
		t.Error("output should indicate mdpress.yml was skipped")
	}
}

func TestInit_FreshProject(t *testing.T) {
	// Create isolated temp directory for the project
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Also set HOME to avoid side effects
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Set the projectConfig flag
	projectConfig = true

	// Create and execute command
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--project"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("init --project command failed: %v", err)
	}

	// Verify ./mdpress.yml is created
	configPath := filepath.Join(tmpDir, "mdpress.yml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml: %v", err)
	}

	// Verify content is project-specific template
	if !strings.Contains(string(content), "version: \"1\"") {
		t.Error("project mdpress.yml should contain version: \"1\"")
	}
	if !strings.Contains(string(content), "mdpress project configuration") {
		t.Error("project mdpress.yml should contain project-specific header")
	}
	if !strings.Contains(string(content), "overrides your global config") {
		t.Error("project mdpress.yml should mention it overrides global config")
	}
}

func TestInit_IdempotentProject(t *testing.T) {
	// Create isolated temp directory for the project
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Also set HOME to avoid side effects
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Set the projectConfig flag
	projectConfig = true

	// First run - create file
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--project"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("first init --project command failed: %v", err)
	}

	// Record file contents after first run
	configPath := filepath.Join(tmpDir, "mdpress.yml")
	originalContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml: %v", err)
	}

	// Second run - should skip without error
	buf.Reset()
	projectConfig = true
	cmd2 := newTestInitCmd()
	cmd2.SetOut(buf)
	cmd2.SetArgs([]string{"--project"})
	err = cmd2.Execute()
	if err != nil {
		t.Fatalf("second init --project command failed: %v", err)
	}

	// Verify file contents are not modified
	newContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read mdpress.yml after second run: %v", err)
	}
	if string(originalContent) != string(newContent) {
		t.Error("mdpress.yml should not be modified on second run")
	}

	// Verify output indicates skipping
	output := buf.String()
	if !strings.Contains(output, "Skipped") {
		t.Error("second run output should indicate file was skipped")
	}
}

func TestInit_FeedbackCreated(t *testing.T) {
	// Create isolated temp directory for HOME
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Reset the projectConfig flag
	projectConfig = false

	// Create and execute command
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Verify "Created" messages appear in output
	output := buf.String()
	if !strings.Contains(output, "Created") {
		t.Error("fresh init should output 'Created' messages")
	}
	if !strings.Contains(output, "mdpress.yml") {
		t.Error("output should mention mdpress.yml")
	}
	if !strings.Contains(output, "logos") {
		t.Error("output should mention logos directory")
	}
	if !strings.Contains(output, "templates") {
		t.Error("output should mention templates directory")
	}
}

func TestInit_FeedbackSkipped(t *testing.T) {
	// Create isolated temp directory for HOME
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Reset the projectConfig flag
	projectConfig = false

	// First run - create files
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("first init command failed: %v", err)
	}

	// Second run - should show skipped messages
	buf.Reset()
	projectConfig = false
	cmd2 := newTestInitCmd()
	cmd2.SetOut(buf)
	cmd2.SetArgs([]string{})
	err = cmd2.Execute()
	if err != nil {
		t.Fatalf("second init command failed: %v", err)
	}

	// Verify output shows summary message when all files exist
	output := buf.String()
	if !strings.Contains(output, "already initialized") {
		t.Errorf("re-init should output 'already initialized' message, got: %s", output)
	}
}

func TestInit_FeedbackProject(t *testing.T) {
	// Create isolated temp directory for the project
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Also set HOME to avoid side effects
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Set the projectConfig flag
	projectConfig = true

	// First run - fresh init
	cmd := newTestInitCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--project"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("init --project command failed: %v", err)
	}

	// Verify "Created" message on fresh init
	output := buf.String()
	if !strings.Contains(output, "Created") {
		t.Error("fresh project init should output 'Created' message")
	}
	if !strings.Contains(output, "mdpress.yml") {
		t.Error("output should mention mdpress.yml")
	}

	// Second run - re-init
	buf.Reset()
	projectConfig = true
	cmd2 := newTestInitCmd()
	cmd2.SetOut(buf)
	cmd2.SetArgs([]string{"--project"})
	err = cmd2.Execute()
	if err != nil {
		t.Fatalf("second init --project command failed: %v", err)
	}

	// Verify "Skipped" message on re-init
	output = buf.String()
	if !strings.Contains(output, "Skipped") {
		t.Error("re-init should output 'Skipped' message")
	}
	if !strings.Contains(output, "already exists") {
		t.Error("re-init should mention file already exists")
	}
}
