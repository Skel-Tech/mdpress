package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skel-tech/mdpress/internal/cloud"
	"github.com/spf13/cobra"
)

// newTestListCmd creates a fresh list command for testing.
func newTestListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		RunE:  listCmd.RunE,
	}
	cmd.Flags().Bool("local-only", false, "Only show local templates, skip cloud fetch")
	return cmd
}

func newTestTemplatesCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "templates",
		RunE: templatesCmd.RunE,
	}
}

// setupTempHome creates a temp HOME directory and returns cleanup function.
func setupTempHome(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	return tmpDir, func() {
		os.Setenv("HOME", oldHome)
	}
}

// setupMockAPIServer creates a mock API server and returns cleanup function.
func setupMockAPIServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	oldURL := os.Getenv("MDPRESS_API_URL")
	os.Setenv("MDPRESS_API_URL", server.URL)
	return server, func() {
		server.Close()
		if oldURL != "" {
			os.Setenv("MDPRESS_API_URL", oldURL)
		} else {
			os.Unsetenv("MDPRESS_API_URL")
		}
	}
}

// createLocalTemplate creates a template YAML file in the given directory.
func createLocalTemplate(t *testing.T, dir, name, description string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create template dir: %v", err)
	}
	content := "name: " + name + "\nversion: \"1\"\ndescription: " + description + "\n"
	path := filepath.Join(dir, name+".yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}
}

// createAuthFile creates an auth.yml file with the given API key.
func createAuthFile(t *testing.T, homeDir, apiKey string) {
	t.Helper()
	authDir := filepath.Join(homeDir, ".config", "mdpress")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		t.Fatalf("failed to create auth dir: %v", err)
	}
	content := "api_key: " + apiKey + "\n"
	if err := os.WriteFile(filepath.Join(authDir, "auth.yml"), []byte(content), 0600); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}
}

func TestTemplatesList_LocalOnly(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create a local global template
	globalDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	createLocalTemplate(t, globalDir, "mytemplate", "My local template")

	// Create project template in current directory
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)
	createLocalTemplate(t, "templates", "projecttmpl", "Project template")

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--local-only"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "mytemplate") {
		t.Error("output should contain local global template 'mytemplate'")
	}
	if !strings.Contains(output, "projecttmpl") {
		t.Error("output should contain local project template 'projecttmpl'")
	}
	if !strings.Contains(output, "global") {
		t.Error("output should show 'global' source")
	}
	if !strings.Contains(output, "project") {
		t.Error("output should show 'project' source")
	}
}

func TestTemplatesList_CloudTemplatesAppearAfterLocal(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create a local template
	globalDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	createLocalTemplate(t, globalDir, "local-first", "Local template")

	// Setup mock API server returning cloud templates
	cloudTemplates := []cloud.CloudTemplate{
		{Name: "cloud-second", Description: "Cloud template", Free: true},
	}
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := cloud.ListTemplatesResponse{Templates: cloudTemplates}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir to avoid project templates
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	output := buf.String()
	localIdx := strings.Index(output, "local-first")
	cloudIdx := strings.Index(output, "cloud-second")

	if localIdx == -1 {
		t.Error("output should contain local template 'local-first'")
	}
	if cloudIdx == -1 {
		t.Error("output should contain cloud template 'cloud-second'")
	}
	if localIdx > cloudIdx {
		t.Error("local templates should appear before cloud templates")
	}
}

func TestTemplatesList_LocalOnlySkipsCloudFetch(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create a local template
	globalDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	createLocalTemplate(t, globalDir, "localonly", "Local template")

	// Setup mock API server that fails if called
	apiCalled := false
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		apiCalled = true
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--local-only"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	if apiCalled {
		t.Error("--local-only should prevent cloud API fetch")
	}

	output := buf.String()
	if !strings.Contains(output, "localonly") {
		t.Error("output should contain local template")
	}
}

func TestTemplatesList_NetworkFailureShowsWarningButListsLocal(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create a local template
	globalDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	createLocalTemplate(t, globalDir, "localfallback", "Local template")

	// Setup mock API server that fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close() // Close immediately to simulate network failure

	oldURL := os.Getenv("MDPRESS_API_URL")
	os.Setenv("MDPRESS_API_URL", serverURL)
	defer func() {
		if oldURL != "" {
			os.Setenv("MDPRESS_API_URL", oldURL)
		} else {
			os.Unsetenv("MDPRESS_API_URL")
		}
	}()

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command should not fail on network error: %v", err)
	}

	output := buf.String()

	// Local templates should still be listed
	if !strings.Contains(output, "localfallback") {
		t.Error("output should contain local template despite network failure")
	}

	// Warning should be shown (OutOrStderr goes to output when output is set)
	if !strings.Contains(output, "Warning") {
		t.Error("output should contain warning about cloud fetch failure")
	}
}

func TestTemplatesList_NonFreeTemplatesShowRequiresLoginWhenUnauthenticated(t *testing.T) {
	_, cleanup := setupTempHome(t)
	defer cleanup()

	// No auth file = unauthenticated

	// Setup mock API server returning a pro template
	cloudTemplates := []cloud.CloudTemplate{
		{Name: "pro-template", Description: "Pro features", Free: false},
		{Name: "free-template", Description: "Free features", Free: true},
	}
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := cloud.ListTemplatesResponse{Templates: cloudTemplates}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	output := buf.String()
	// Pro template should show "requires login"
	if !strings.Contains(output, "pro-template") {
		t.Error("output should contain pro-template")
	}
	if !strings.Contains(output, "requires login") {
		t.Error("pro template should show 'requires login' when unauthenticated")
	}
	// Free template should NOT show "requires login"
	freeIdx := strings.Index(output, "free-template")
	loginIdx := strings.Index(output, "requires login")
	if freeIdx > loginIdx {
		// If "requires login" appears before free-template, it's for pro-template (ok)
		// But we should verify free-template line doesn't have it
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "free-template") && strings.Contains(line, "requires login") {
				t.Error("free template should not show 'requires login'")
			}
		}
	}
}

func TestTemplatesList_AuthenticatedUserDoesNotSeeRequiresLogin(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create auth file with API key
	createAuthFile(t, tmpDir, "mdp_test_key")

	// Setup mock API server returning a pro template
	cloudTemplates := []cloud.CloudTemplate{
		{Name: "pro-template", Description: "Pro features", Free: false},
	}
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := cloud.ListTemplatesResponse{Templates: cloudTemplates}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "pro-template") {
		t.Error("output should contain pro-template")
	}
	if strings.Contains(output, "requires login") {
		t.Error("authenticated user should not see 'requires login' for pro templates")
	}
}

func TestTemplatesList_NoTemplatesFound(t *testing.T) {
	_, cleanup := setupTempHome(t)
	defer cleanup()

	// No local templates, empty cloud response
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := cloud.ListTemplatesResponse{Templates: []cloud.CloudTemplate{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir with no templates
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No templates found") {
		t.Error("output should indicate no templates found")
	}
}

func TestTemplates_NoTemplatesFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Change to temp dir so no project templates are found
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	buf := new(bytes.Buffer)
	cmd := newTestTemplatesCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No templates found") {
		t.Errorf("expected 'No templates found', got %q", out)
	}
	if !strings.Contains(out, "project:") {
		t.Error("should show project template directory hint")
	}
	if !strings.Contains(out, "global:") {
		t.Error("should show global template directory hint")
	}
}

func TestTemplates_WithTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create a global template
	globalTemplateDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	os.MkdirAll(globalTemplateDir, 0755)
	templateContent := `name: "invoice"
version: "1"
description: "Invoice template"
font: "Times"
`
	os.WriteFile(filepath.Join(globalTemplateDir, "invoice.yml"), []byte(templateContent), 0644)

	// Change to temp dir with no project templates
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	buf := new(bytes.Buffer)
	cmd := newTestTemplatesCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Available templates") {
		t.Errorf("expected 'Available templates' header, got %q", out)
	}
	if !strings.Contains(out, "invoice") {
		t.Errorf("expected 'invoice' template in listing, got %q", out)
	}
}

func TestTemplates_ProjectTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create a project template directory
	projectTemplateDir := filepath.Join(tmpDir, "templates")
	os.MkdirAll(projectTemplateDir, 0755)
	templateContent := `name: "report"
version: "1"
description: "Report template"
font: "Courier"
`
	os.WriteFile(filepath.Join(projectTemplateDir, "report.yml"), []byte(templateContent), 0644)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	buf := new(bytes.Buffer)
	cmd := newTestTemplatesCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "report") {
		t.Errorf("expected 'report' template in listing, got %q", out)
	}
}
