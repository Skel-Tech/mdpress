package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGlobalTemplateDir(t *testing.T) {
	dir := GlobalTemplateDir()

	// Should return a path under user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	expected := filepath.Join(home, ".config", "mdpress", "templates")
	if dir != expected {
		t.Errorf("GlobalTemplateDir() = %q, want %q", dir, expected)
	}
}

func TestProjectTemplateDir(t *testing.T) {
	dir := ProjectTemplateDir()
	if dir != "templates" {
		t.Errorf("ProjectTemplateDir() = %q, want %q", dir, "templates")
	}
}

func TestListTemplates_EmptyDirs(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "project")
	globalDir := filepath.Join(tmpDir, "global")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		t.Fatalf("failed to create global dir: %v", err)
	}

	templates, err := listTemplatesFromDir(projectDir, "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 0 {
		t.Errorf("listTemplatesFromDir() returned %d templates, want 0", len(templates))
	}
}

func TestListTemplates_NonexistentDir(t *testing.T) {
	templates, err := listTemplatesFromDir("/nonexistent/path", "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 0 {
		t.Errorf("listTemplatesFromDir() returned %d templates, want 0", len(templates))
	}
}

func TestListTemplates_FindsTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid template
	templateContent := `
name: "Test Template"
version: "1"
description: "A test template for unit testing"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.yml"), []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	templates, err := listTemplatesFromDir(tmpDir, "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 1 {
		t.Fatalf("listTemplatesFromDir() returned %d templates, want 1", len(templates))
	}

	tmpl := templates[0]
	if tmpl.Name != "Test Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Test Template")
	}
	if tmpl.Description != "A test template for unit testing" {
		t.Errorf("Description = %q, want %q", tmpl.Description, "A test template for unit testing")
	}
	if tmpl.Source != "project" {
		t.Errorf("Source = %q, want %q", tmpl.Source, "project")
	}
	if tmpl.Path != filepath.Join(tmpDir, "test.yml") {
		t.Errorf("Path = %q, want %q", tmpl.Path, filepath.Join(tmpDir, "test.yml"))
	}
}

func TestListTemplates_MultipleTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	templates := []struct {
		filename string
		name     string
	}{
		{"corporate.yml", "Corporate Template"},
		{"invoice.yaml", "Invoice Template"},
		{"report.yml", "Report Template"},
	}

	for _, tmpl := range templates {
		content := "name: \"" + tmpl.name + "\"\nversion: \"1\"\n"
		if err := os.WriteFile(filepath.Join(tmpDir, tmpl.filename), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write template %s: %v", tmpl.filename, err)
		}
	}

	result, err := listTemplatesFromDir(tmpDir, "global")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(result) != 3 {
		t.Errorf("listTemplatesFromDir() returned %d templates, want 3", len(result))
	}

	// Verify all templates are returned with correct source
	for _, r := range result {
		if r.Source != "global" {
			t.Errorf("template %q has Source = %q, want %q", r.Name, r.Source, "global")
		}
	}
}

func TestListTemplates_IgnoresNonYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid template
	if err := os.WriteFile(filepath.Join(tmpDir, "valid.yml"), []byte("name: \"Valid\"\nversion: \"1\"\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Create non-YAML files that should be ignored
	if err := os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# README"), 0644); err != nil {
		t.Fatalf("failed to write readme: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "config.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	templates, err := listTemplatesFromDir(tmpDir, "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 1 {
		t.Errorf("listTemplatesFromDir() returned %d templates, want 1", len(templates))
	}
}

func TestListTemplates_SkipsInvalidTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid template
	if err := os.WriteFile(filepath.Join(tmpDir, "valid.yml"), []byte("name: \"Valid\"\nversion: \"1\"\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Create invalid template (missing required fields)
	if err := os.WriteFile(filepath.Join(tmpDir, "invalid.yml"), []byte("description: \"No name or version\"\n"), 0644); err != nil {
		t.Fatalf("failed to write invalid template: %v", err)
	}

	templates, err := listTemplatesFromDir(tmpDir, "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 1 {
		t.Errorf("listTemplatesFromDir() returned %d templates, want 1", len(templates))
	}
	if templates[0].Name != "Valid" {
		t.Errorf("expected the valid template to be returned, got %q", templates[0].Name)
	}
}

func TestListTemplates_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid template
	if err := os.WriteFile(filepath.Join(tmpDir, "valid.yml"), []byte("name: \"Valid\"\nversion: \"1\"\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Create a subdirectory
	if err := os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	templates, err := listTemplatesFromDir(tmpDir, "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 1 {
		t.Errorf("listTemplatesFromDir() returned %d templates, want 1", len(templates))
	}
}

func TestListTemplates_SupportsYAMLExtension(t *testing.T) {
	tmpDir := t.TempDir()

	// Create template with .yaml extension
	if err := os.WriteFile(filepath.Join(tmpDir, "template.yaml"), []byte("name: \"YAML Extension\"\nversion: \"1\"\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	templates, err := listTemplatesFromDir(tmpDir, "project")
	if err != nil {
		t.Fatalf("listTemplatesFromDir() error = %v, want nil", err)
	}
	if len(templates) != 1 {
		t.Fatalf("listTemplatesFromDir() returned %d templates, want 1", len(templates))
	}
	if templates[0].Name != "YAML Extension" {
		t.Errorf("Name = %q, want %q", templates[0].Name, "YAML Extension")
	}
}

func TestSourceConstants(t *testing.T) {
	if SourceProject != "project" {
		t.Errorf("SourceProject = %q, want %q", SourceProject, "project")
	}
	if SourceGlobal != "global" {
		t.Errorf("SourceGlobal = %q, want %q", SourceGlobal, "global")
	}
	if SourceCloud != "cloud" {
		t.Errorf("SourceCloud = %q, want %q", SourceCloud, "cloud")
	}
}

func TestTemplateInfo_FreeField(t *testing.T) {
	// Verify the Free field exists and works correctly
	info := TemplateInfo{
		Name:        "Test",
		Description: "Test template",
		Source:      SourceCloud,
		Path:        "",
		Free:        true,
	}
	if !info.Free {
		t.Error("TemplateInfo.Free should be true")
	}

	info.Free = false
	if info.Free {
		t.Error("TemplateInfo.Free should be false")
	}
}

func TestTemplateExistsLocal_FindsInProject(t *testing.T) {
	// Save original function and restore after test
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create a valid template in project directory
	templateContent := `name: "My Template"
version: "1"
description: "Test template"
`
	if err := os.WriteFile(filepath.Join(projectDir, "test.yml"), []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Change to temp directory so ProjectTemplateDir() finds our test templates
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if !TemplateExistsLocal("My Template") {
		t.Error("TemplateExistsLocal should return true for existing project template")
	}
}

func TestTemplateExistsLocal_NotFound(t *testing.T) {
	// Use a temp directory with no templates
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if TemplateExistsLocal("Nonexistent Template") {
		t.Error("TemplateExistsLocal should return false for nonexistent template")
	}
}

func TestTemplateExistsLocal_EmptyDirectories(t *testing.T) {
	// Use a temp directory with empty templates directories
	tmpDir := t.TempDir()

	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// No templates directory exists - should return false without error
	if TemplateExistsLocal("Any Template") {
		t.Error("TemplateExistsLocal should return false when no templates exist")
	}
}
