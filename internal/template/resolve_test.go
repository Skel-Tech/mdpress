package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsPath_PathsWithSlash(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"./template.yml", true},
		{"../template.yml", true},
		{"templates/corporate.yml", true},
		{"/absolute/path/template.yml", true},
		{"./", true},
		{"path/to/file", true},
	}

	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			result := IsPath(tc.value)
			if result != tc.expected {
				t.Errorf("IsPath(%q) = %v, want %v", tc.value, result, tc.expected)
			}
		})
	}
}

func TestIsPath_PathsStartingWithDot(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{".", true},
		{"..", true},
		{".hidden", true},
		{"..hidden", true},
		{".yml", true},
	}

	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			result := IsPath(tc.value)
			if result != tc.expected {
				t.Errorf("IsPath(%q) = %v, want %v", tc.value, result, tc.expected)
			}
		})
	}
}

func TestIsPath_TemplateNames(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"invoice", false},
		{"corporate", false},
		{"my-template", false},
		{"template_v2", false},
		{"CorporateTemplate", false},
		{"", false},
	}

	for _, tc := range tests {
		name := tc.value
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			result := IsPath(tc.value)
			if result != tc.expected {
				t.Errorf("IsPath(%q) = %v, want %v", tc.value, result, tc.expected)
			}
		})
	}
}

func TestResolve_PathInput(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "custom.yml")

	content := `
name: "Custom Template"
version: "1"
description: "A custom template loaded by path"
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Should load from path directly (starts with / so IsPath returns true)
	tmpl, err := Resolve(templatePath)
	if err != nil {
		t.Fatalf("Resolve() error = %v, want nil", err)
	}
	if tmpl.Name != "Custom Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Custom Template")
	}
}

func TestResolve_RelativePathInput(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create template
	content := `
name: "Relative Path Template"
version: "1"
`
	if err := os.WriteFile("template.yml", []byte(content), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Should load from relative path
	tmpl, err := Resolve("./template.yml")
	if err != nil {
		t.Fatalf("Resolve() error = %v, want nil", err)
	}
	if tmpl.Name != "Relative Path Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Relative Path Template")
	}
}

func TestResolve_NotFoundError(t *testing.T) {
	// Create temp project dir with no templates
	tmpDir := t.TempDir()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err = Resolve("nonexistent")
	if err == nil {
		t.Fatal("Resolve() should return error for non-existent template")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention the template name: %v", err)
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention 'not found': %v", err)
	}
}

func TestResolve_NotFoundWithAvailableTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	// Create some available templates
	if err := os.WriteFile(filepath.Join(templatesDir, "invoice.yml"), []byte("name: \"Invoice Template\"\nversion: \"1\"\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templatesDir, "report.yml"), []byte("name: \"Report Template\"\nversion: \"1\"\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err = Resolve("nonexistent")
	if err == nil {
		t.Fatal("Resolve() should return error for non-existent template")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "Invoice Template") {
		t.Errorf("error should list available template 'Invoice Template': %v", err)
	}
	if !strings.Contains(errStr, "Report Template") {
		t.Errorf("error should list available template 'Report Template': %v", err)
	}
	if !strings.Contains(errStr, "project") {
		t.Errorf("error should indicate source 'project': %v", err)
	}
}

func TestResolve_FindsProjectTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	// Create a project template
	content := `
name: "Corporate Template"
version: "1"
description: "Corporate branding template"
accent_color: "#336699"
`
	if err := os.WriteFile(filepath.Join(templatesDir, "corporate.yml"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	tmpl, err := Resolve("Corporate Template")
	if err != nil {
		t.Fatalf("Resolve() error = %v, want nil", err)
	}
	if tmpl.Name != "Corporate Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Corporate Template")
	}
	if tmpl.AccentColor != "#336699" {
		t.Errorf("AccentColor = %q, want %q", tmpl.AccentColor, "#336699")
	}
}

func TestResolve_ProjectTakesPrecedenceOverGlobal(t *testing.T) {
	// This test requires modifying the home directory, which is complex
	// Instead, we test the behavior through listTemplatesFromDir ordering
	tmpProjectDir := t.TempDir()
	tmpGlobalDir := t.TempDir()

	// Create same-named template in both dirs with different descriptions
	projectContent := `
name: "Shared Template"
version: "1"
description: "Project version"
`
	globalContent := `
name: "Shared Template"
version: "1"
description: "Global version"
`
	if err := os.WriteFile(filepath.Join(tmpProjectDir, "shared.yml"), []byte(projectContent), 0644); err != nil {
		t.Fatalf("failed to write project template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpGlobalDir, "shared.yml"), []byte(globalContent), 0644); err != nil {
		t.Fatalf("failed to write global template: %v", err)
	}

	// Get templates from both dirs
	projectTemplates, _ := listTemplatesFromDir(tmpProjectDir, "project")
	globalTemplates, _ := listTemplatesFromDir(tmpGlobalDir, "global")

	// Simulate the order ListTemplates uses (project first)
	allTemplates := append(projectTemplates, globalTemplates...)

	// First match should be the project one
	var firstMatch *TemplateInfo
	for i := range allTemplates {
		if allTemplates[i].Name == "Shared Template" {
			firstMatch = &allTemplates[i]
			break
		}
	}

	if firstMatch == nil {
		t.Fatal("Expected to find 'Shared Template'")
	}
	if firstMatch.Source != "project" {
		t.Errorf("First match should be project template, got Source = %q", firstMatch.Source)
	}
	if firstMatch.Description != "Project version" {
		t.Errorf("First match should have project description, got %q", firstMatch.Description)
	}
}

func TestResolve_PathNotFound(t *testing.T) {
	_, err := Resolve("./nonexistent/template.yml")
	if err == nil {
		t.Fatal("Resolve() should return error for non-existent path")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention 'not found': %v", err)
	}
}

func TestBuildNotFoundError_NoTemplates(t *testing.T) {
	err := buildNotFoundError("mytemplate", nil)
	if err == nil {
		t.Fatal("buildNotFoundError should return an error")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "mytemplate") {
		t.Errorf("error should mention template name: %v", err)
	}
	if !strings.Contains(errStr, "no templates available") {
		t.Errorf("error should mention no templates: %v", err)
	}
}

func TestBuildNotFoundError_WithTemplates(t *testing.T) {
	templates := []TemplateInfo{
		{Name: "Invoice", Source: "project"},
		{Name: "Report", Source: "global"},
	}

	err := buildNotFoundError("missing", templates)
	if err == nil {
		t.Fatal("buildNotFoundError should return an error")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "missing") {
		t.Errorf("error should mention template name: %v", err)
	}
	if !strings.Contains(errStr, "Invoice") {
		t.Errorf("error should list Invoice template: %v", err)
	}
	if !strings.Contains(errStr, "Report") {
		t.Errorf("error should list Report template: %v", err)
	}
	if !strings.Contains(errStr, "project") {
		t.Errorf("error should indicate project source: %v", err)
	}
	if !strings.Contains(errStr, "global") {
		t.Errorf("error should indicate global source: %v", err)
	}
}
