package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromFile_ValidTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "valid.yml")

	content := `
name: "Corporate Template"
version: "1"
description: "A corporate branding template"
logo: "logo.png"
logo_position: "top-right"
logo_width: 80
font: "Helvetica"
font_size: 11
accent_color: "#336699"
header: true
footer: "Page {page} of {pages}"
margins:
  top: 25
  right: 20
  bottom: 25
  left: 20
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	tmpl, err := LoadFromFile(templatePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v, want nil", err)
	}

	if tmpl.Name != "Corporate Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Corporate Template")
	}
	if tmpl.Version != "1" {
		t.Errorf("Version = %q, want %q", tmpl.Version, "1")
	}
	if tmpl.Description != "A corporate branding template" {
		t.Errorf("Description = %q, want %q", tmpl.Description, "A corporate branding template")
	}
	if tmpl.Logo != "logo.png" {
		t.Errorf("Logo = %q, want %q", tmpl.Logo, "logo.png")
	}
	if tmpl.LogoPosition != "top-right" {
		t.Errorf("LogoPosition = %q, want %q", tmpl.LogoPosition, "top-right")
	}
	if tmpl.LogoWidth != 80 {
		t.Errorf("LogoWidth = %v, want %v", tmpl.LogoWidth, 80)
	}
	if tmpl.Font != "Helvetica" {
		t.Errorf("Font = %q, want %q", tmpl.Font, "Helvetica")
	}
	if tmpl.FontSize != 11 {
		t.Errorf("FontSize = %v, want %v", tmpl.FontSize, 11)
	}
	if tmpl.AccentColor != "#336699" {
		t.Errorf("AccentColor = %q, want %q", tmpl.AccentColor, "#336699")
	}
	if tmpl.Header == nil || *tmpl.Header != true {
		t.Errorf("Header = %v, want true", tmpl.Header)
	}
	if tmpl.Footer != "Page {page} of {pages}" {
		t.Errorf("Footer = %q, want %q", tmpl.Footer, "Page {page} of {pages}")
	}
	if tmpl.Margins.Top != 25 {
		t.Errorf("Margins.Top = %v, want %v", tmpl.Margins.Top, 25)
	}
	if tmpl.Margins.Right != 20 {
		t.Errorf("Margins.Right = %v, want %v", tmpl.Margins.Right, 20)
	}
	if tmpl.Margins.Bottom != 25 {
		t.Errorf("Margins.Bottom = %v, want %v", tmpl.Margins.Bottom, 25)
	}
	if tmpl.Margins.Left != 20 {
		t.Errorf("Margins.Left = %v, want %v", tmpl.Margins.Left, 20)
	}
	if tmpl.SourcePath != templatePath {
		t.Errorf("SourcePath = %q, want %q", tmpl.SourcePath, templatePath)
	}
}

func TestLoadFromFile_MinimalTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "minimal.yml")

	// Only required fields
	content := `
name: "Minimal Template"
version: "1"
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	tmpl, err := LoadFromFile(templatePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v, want nil", err)
	}

	if tmpl.Name != "Minimal Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Minimal Template")
	}
	if tmpl.Version != "1" {
		t.Errorf("Version = %q, want %q", tmpl.Version, "1")
	}
	// Optional fields should have zero values
	if tmpl.Description != "" {
		t.Errorf("Description = %q, want empty", tmpl.Description)
	}
	if tmpl.Header != nil {
		t.Errorf("Header = %v, want nil (not set)", tmpl.Header)
	}
}

func TestLoadFromFile_HeaderExplicitlyFalse(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "no-header.yml")

	content := `
name: "No Header Template"
version: "1"
header: false
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	tmpl, err := LoadFromFile(templatePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v, want nil", err)
	}

	if tmpl.Header == nil {
		t.Error("Header should not be nil when explicitly set to false")
	} else if *tmpl.Header != false {
		t.Errorf("Header = %v, want false", *tmpl.Header)
	}
}

func TestLoadFromFile_FileNotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/template.yml")
	if err == nil {
		t.Fatal("LoadFromFile() should return error for non-existent file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention 'not found': %v", err)
	}
}

func TestLoadFromFile_MissingName(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "no-name.yml")

	content := `
version: "1"
description: "A template without a name"
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromFile(templatePath)
	if err == nil {
		t.Fatal("LoadFromFile() should return error for missing name")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("error should mention 'name': %v", err)
	}
}

func TestLoadFromFile_MissingVersion(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "no-version.yml")

	content := `
name: "Template Without Version"
description: "This template is missing a version"
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromFile(templatePath)
	if err == nil {
		t.Fatal("LoadFromFile() should return error for missing version")
	}
	if !strings.Contains(err.Error(), "version") {
		t.Errorf("error should mention 'version': %v", err)
	}
}

func TestLoadFromFile_InvalidVersion(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "bad-version.yml")

	content := `
name: "Bad Version Template"
version: "2"
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromFile(templatePath)
	if err == nil {
		t.Fatal("LoadFromFile() should return error for invalid version")
	}
	if !strings.Contains(err.Error(), "version") {
		t.Errorf("error should mention 'version': %v", err)
	}
	if !strings.Contains(err.Error(), "1") {
		t.Errorf("error should mention expected version '1': %v", err)
	}
}

func TestLoadFromFile_UnknownFields(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "unknown-fields.yml")

	content := `
name: "Template With Unknown Field"
version: "1"
unknown_field: "this should cause an error"
another_bad_field: 42
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromFile(templatePath)
	if err == nil {
		t.Fatal("LoadFromFile() should return error for unknown fields")
	}
	if !strings.Contains(err.Error(), "unknown_field") && !strings.Contains(err.Error(), "another_bad_field") {
		t.Errorf("error should mention the unknown field: %v", err)
	}
}

func TestLoadFromFile_MalformedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "malformed.yml")

	content := `
name: "Bad Template"
version: [invalid
  - yaml
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromFile(templatePath)
	if err == nil {
		t.Fatal("LoadFromFile() should return error for malformed YAML")
	}
}

func TestLoadFromFile_NestedMargins(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "margins.yml")

	content := `
name: "Margins Template"
version: "1"
margins:
  top: 30
  right: 25
  bottom: 30
  left: 25
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	tmpl, err := LoadFromFile(templatePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v, want nil", err)
	}

	if tmpl.Margins.Top != 30 {
		t.Errorf("Margins.Top = %v, want %v", tmpl.Margins.Top, 30)
	}
	if tmpl.Margins.Right != 25 {
		t.Errorf("Margins.Right = %v, want %v", tmpl.Margins.Right, 25)
	}
	if tmpl.Margins.Bottom != 30 {
		t.Errorf("Margins.Bottom = %v, want %v", tmpl.Margins.Bottom, 30)
	}
	if tmpl.Margins.Left != 25 {
		t.Errorf("Margins.Left = %v, want %v", tmpl.Margins.Left, 25)
	}
}

func TestLoadFromFile_PartialMargins(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "partial-margins.yml")

	content := `
name: "Partial Margins Template"
version: "1"
margins:
  top: 40
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	tmpl, err := LoadFromFile(templatePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v, want nil", err)
	}

	// Only top should be set
	if tmpl.Margins.Top != 40 {
		t.Errorf("Margins.Top = %v, want %v", tmpl.Margins.Top, 40)
	}
	// Others should be zero (not set)
	if tmpl.Margins.Right != 0 {
		t.Errorf("Margins.Right = %v, want %v", tmpl.Margins.Right, 0)
	}
	if tmpl.Margins.Bottom != 0 {
		t.Errorf("Margins.Bottom = %v, want %v", tmpl.Margins.Bottom, 0)
	}
	if tmpl.Margins.Left != 0 {
		t.Errorf("Margins.Left = %v, want %v", tmpl.Margins.Left, 0)
	}
}

func TestLoadFromFile_UnknownFieldInMargins(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "bad-margins.yml")

	content := `
name: "Bad Margins Template"
version: "1"
margins:
  top: 20
  invalid_margin: 10
`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromFile(templatePath)
	if err == nil {
		t.Fatal("LoadFromFile() should return error for unknown field in margins")
	}
	if !strings.Contains(err.Error(), "invalid_margin") {
		t.Errorf("error should mention the unknown field 'invalid_margin': %v", err)
	}
}
