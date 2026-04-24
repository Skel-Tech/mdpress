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

	"github.com/skel-tech/mdpress/internal/auth"
	"github.com/skel-tech/mdpress/internal/cloud"
	"github.com/spf13/cobra"
)

// newTestPullCmd creates a fresh pull command for testing.
func newTestPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull <name>",
		Short: "Pull a template from the cloud",
		Args:  cobra.ExactArgs(1),
		RunE:  pullCmd.RunE,
	}
	cmd.Flags().Bool("force", false, "Overwrite existing template if it exists")
	return cmd
}

func TestTemplatesPull_SuccessfulPullFreeTemplate(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	templateContent := "name: basic\nversion: \"1\"\ndescription: Basic template\n"

	// Setup mock API server
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "basic", Description: "Basic template", Free: true},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/v1/templates/basic" {
			resp := cloud.FetchTemplateResponse{
				Name:    "basic",
				Content: templateContent,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"basic"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("pull command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "basic") {
		t.Error("output should mention template name")
	}
	if !strings.Contains(output, "saved") {
		t.Error("output should indicate template was saved")
	}

	// Verify template was written to global templates directory
	expectedPath := filepath.Join(tmpDir, ".config", "mdpress", "templates", "basic.yml")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("template file should exist at %s: %v", expectedPath, err)
	}
	if string(content) != templateContent {
		t.Errorf("template content mismatch: got %q, want %q", string(content), templateContent)
	}
}

func TestTemplatesPull_ProTemplateWhenAuthenticated(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create auth file
	createAuthFile(t, tmpDir, "mdp_valid_key")

	// Set up auth.ValidateKey to accept our test key
	oldValidateKey := auth.ValidateKey
	auth.ValidateKey = func(key string) (*auth.LicenseInfo, error) {
		if key == "mdp_valid_key" {
			return &auth.LicenseInfo{}, nil
		}
		return nil, &auth.FeatureGatedError{Feature: "templates pull"}
	}
	defer func() { auth.ValidateKey = oldValidateKey }()

	templateContent := "name: pro-template\nversion: \"1\"\ndescription: Pro template\n"

	// Setup mock API server
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "pro-template", Description: "Pro template", Free: false},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/v1/templates/pro-template" {
			resp := cloud.FetchTemplateResponse{
				Name:    "pro-template",
				Content: templateContent,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"pro-template"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("pull command should succeed for authenticated user: %v", err)
	}

	// Verify template was saved
	expectedPath := filepath.Join(tmpDir, ".config", "mdpress", "templates", "pro-template.yml")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("template file should exist at %s", expectedPath)
	}
}

func TestTemplatesPull_ProTemplateWhenUnauthenticated(t *testing.T) {
	_, cleanup := setupTempHome(t)
	defer cleanup()

	// No auth file = unauthenticated

	// Ensure ValidateKey is nil (fail-closed)
	oldValidateKey := auth.ValidateKey
	auth.ValidateKey = nil
	defer func() { auth.ValidateKey = oldValidateKey }()

	// Setup mock API server
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "pro-only", Description: "Pro template", Free: false},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"pro-only"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("pull command should fail for unauthenticated user on pro template")
	}

	// Should return a FeatureGatedError or similar
	errStr := err.Error()
	if !strings.Contains(errStr, "Pro") && !strings.Contains(errStr, "login") && !strings.Contains(errStr, "account") {
		t.Errorf("error should mention Pro account requirement, got: %v", err)
	}
}

func TestTemplatesPull_ForceSkipsConfirmation(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create an existing template
	globalDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	createLocalTemplate(t, globalDir, "existing", "Old content")

	newContent := "name: existing\nversion: \"1\"\ndescription: New content from cloud\n"

	// Setup mock API server
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "existing", Description: "Cloud version", Free: true},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/v1/templates/existing" {
			resp := cloud.FetchTemplateResponse{
				Name:    "existing",
				Content: newContent,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--force", "existing"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("pull --force should succeed: %v", err)
	}

	// Verify template was overwritten
	expectedPath := filepath.Join(tmpDir, ".config", "mdpress", "templates", "existing.yml")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("template file should exist: %v", err)
	}
	if string(content) != newContent {
		t.Errorf("template should be overwritten with new content")
	}
}

func TestTemplatesPull_OverwritePromptAborted(t *testing.T) {
	tmpDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create an existing template
	globalDir := filepath.Join(tmpDir, ".config", "mdpress", "templates")
	originalContent := "name: existing\nversion: \"1\"\ndescription: Original\n"
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, "existing.yml"), []byte(originalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Setup mock API server
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "existing", Description: "Cloud version", Free: true},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	// Simulate user input "n" (no)
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		w.WriteString("n\n")
		w.Close()
	}()

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"existing"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("pull should fail when user aborts overwrite")
	}
	if !strings.Contains(err.Error(), "Aborted") {
		t.Errorf("error should indicate abort, got: %v", err)
	}

	// Verify original content is preserved
	content, _ := os.ReadFile(filepath.Join(globalDir, "existing.yml"))
	if string(content) != originalContent {
		t.Error("original template should not be modified when user aborts")
	}
}

func TestTemplatesPull_TemplateNotFound(t *testing.T) {
	_, cleanup := setupTempHome(t)
	defer cleanup()

	// Setup mock API server with no matching template
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "other-template", Description: "Some other template", Free: true},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("pull should fail for nonexistent template")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should indicate template not found, got: %v", err)
	}
}

func TestTemplatesPull_NetworkFailure(t *testing.T) {
	_, cleanup := setupTempHome(t)
	defer cleanup()

	// Setup a server and immediately close it
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

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

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"any-template"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("pull should fail on network error")
	}
	if !strings.Contains(err.Error(), "network") {
		t.Errorf("error should indicate network failure, got: %v", err)
	}
}

func TestTemplatesPull_MissingArgument(t *testing.T) {
	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("pull should fail without template name argument")
	}
}

func TestTemplatesPull_OverwriteConfirmedWithYes(t *testing.T) {
	homeDir, cleanup := setupTempHome(t)
	defer cleanup()

	// Create an existing template
	globalDir := filepath.Join(homeDir, ".config", "mdpress", "templates")
	originalContent := "name: existing\nversion: \"1\"\ndescription: Original\n"
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, "existing.yml"), []byte(originalContent), 0644); err != nil {
		t.Fatal(err)
	}

	newContent := "name: existing\nversion: \"1\"\ndescription: Updated\n"

	// Setup mock API server
	server, serverCleanup := setupMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/templates" {
			resp := cloud.ListTemplatesResponse{
				Templates: []cloud.CloudTemplate{
					{Name: "existing", Description: "Cloud version", Free: true},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/v1/templates/existing" {
			resp := cloud.FetchTemplateResponse{
				Name:    "existing",
				Content: newContent,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	defer serverCleanup()
	_ = server

	// Change to temp dir
	projectDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(projectDir)

	// Simulate user input "y" (yes)
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		w.WriteString("y\n")
		w.Close()
	}()

	cmd := newTestPullCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"existing"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("pull should succeed when user confirms overwrite: %v", err)
	}

	// Verify template was overwritten
	content, _ := os.ReadFile(filepath.Join(globalDir, "existing.yml"))
	if string(content) != newContent {
		t.Error("template should be overwritten with new content after user confirms")
	}
}
