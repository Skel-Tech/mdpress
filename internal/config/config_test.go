package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlobalConfigPath(t *testing.T) {
	path := GlobalConfigPath()
	if path == "" {
		t.Fatal("GlobalConfigPath() returned empty string")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	expected := filepath.Join(home, ".config", "mdpress", "mdpress.yml")
	if path != expected {
		t.Errorf("GlobalConfigPath() = %q, want %q", path, expected)
	}
}

func TestProjectConfigPath(t *testing.T) {
	path := ProjectConfigPath()
	if path != "mdpress.yml" {
		t.Errorf("ProjectConfigPath() = %q, want %q", path, "mdpress.yml")
	}
}

func TestLoad_NoConfigFiles(t *testing.T) {
	// Create a temp directory and change to it
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Also ensure no global config exists by using a non-existent home
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Should return defaults
	defaults := DefaultConfig()
	if cfg.Version != defaults.Version {
		t.Errorf("Version = %q, want %q", cfg.Version, defaults.Version)
	}
	if cfg.LogoPosition != defaults.LogoPosition {
		t.Errorf("LogoPosition = %q, want %q", cfg.LogoPosition, defaults.LogoPosition)
	}
	if cfg.LogoWidth != defaults.LogoWidth {
		t.Errorf("LogoWidth = %v, want %v", cfg.LogoWidth, defaults.LogoWidth)
	}
	if cfg.Font != defaults.Font {
		t.Errorf("Font = %q, want %q", cfg.Font, defaults.Font)
	}
}

func TestLoad_GlobalConfigOnly(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)

	// Create a working directory without project config
	workDir := filepath.Join(tmpDir, "work")
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create global config directory and file
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	globalConfig := `
version: "1"
logo_position: "bottom-left"
font: "Arial"
logo_width: 80
`
	os.WriteFile(filepath.Join(globalDir, "mdpress.yml"), []byte(globalConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Should have values from global config
	if cfg.LogoPosition != "bottom-left" {
		t.Errorf("LogoPosition = %q, want %q", cfg.LogoPosition, "bottom-left")
	}
	if cfg.Font != "Arial" {
		t.Errorf("Font = %q, want %q", cfg.Font, "Arial")
	}
	if cfg.LogoWidth != 80 {
		t.Errorf("LogoWidth = %v, want %v", cfg.LogoWidth, 80)
	}
}

func TestLoad_ProjectConfigOnly(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory (no global config there)
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create project config
	projectConfig := `
version: "1"
logo_position: "top-left"
accent_color: "#FF0000"
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(projectConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Should have values from project config
	if cfg.LogoPosition != "top-left" {
		t.Errorf("LogoPosition = %q, want %q", cfg.LogoPosition, "top-left")
	}
	if cfg.AccentColor != "#FF0000" {
		t.Errorf("AccentColor = %q, want %q", cfg.AccentColor, "#FF0000")
	}

	// Should have defaults for unspecified fields
	defaults := DefaultConfig()
	if cfg.Font != defaults.Font {
		t.Errorf("Font = %q, want default %q", cfg.Font, defaults.Font)
	}
}

func TestLoad_ProjectOverridesGlobal(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create global config
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	globalConfig := `
version: "1"
logo_position: "bottom-left"
font: "Arial"
accent_color: "#00FF00"
`
	os.WriteFile(filepath.Join(globalDir, "mdpress.yml"), []byte(globalConfig), 0644)

	// Create project config that overrides some values
	projectConfig := `
version: "1"
logo_position: "top-right"
accent_color: "#FF0000"
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(projectConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Values from project config should override global
	if cfg.LogoPosition != "top-right" {
		t.Errorf("LogoPosition = %q, want %q (from project)", cfg.LogoPosition, "top-right")
	}
	if cfg.AccentColor != "#FF0000" {
		t.Errorf("AccentColor = %q, want %q (from project)", cfg.AccentColor, "#FF0000")
	}

	// Values only in global config should be preserved
	if cfg.Font != "Arial" {
		t.Errorf("Font = %q, want %q (from global)", cfg.Font, "Arial")
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create malformed project config
	malformedConfig := `
version: "1"
logo_position: [invalid
  - yaml
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(malformedConfig), 0644)

	_, err = Load()
	if err == nil {
		t.Fatal("Load() should return error for malformed YAML")
	}
	if !strings.Contains(err.Error(), "project config") {
		t.Errorf("error should mention 'project config': %v", err)
	}
}

func TestLoad_MalformedGlobalYAML(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)

	// Create a working directory without project config
	workDir := filepath.Join(tmpDir, "work")
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create malformed global config
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	malformedConfig := `
version: "1"
  bad indentation
    here
`
	os.WriteFile(filepath.Join(globalDir, "mdpress.yml"), []byte(malformedConfig), 0644)

	_, err = Load()
	if err == nil {
		t.Fatal("Load() should return error for malformed YAML")
	}
	if !strings.Contains(err.Error(), "global config") {
		t.Errorf("error should mention 'global config': %v", err)
	}
}

func TestLoad_UnknownFields(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create project config with unknown field
	configWithUnknown := `
version: "1"
logo_position: "top-right"
unknown_field: "should cause error"
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(configWithUnknown), 0644)

	_, err = Load()
	if err == nil {
		t.Fatal("Load() should return error for unknown fields")
	}
	if !strings.Contains(err.Error(), "unknown_field") {
		t.Errorf("error should mention the unknown field: %v", err)
	}
}

func TestLoad_UnknownFieldsInGlobal(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)

	// Create a working directory without project config
	workDir := filepath.Join(tmpDir, "work")
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create global config with unknown field
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	configWithUnknown := `
version: "1"
not_a_real_field: true
`
	os.WriteFile(filepath.Join(globalDir, "mdpress.yml"), []byte(configWithUnknown), 0644)

	_, err = Load()
	if err == nil {
		t.Fatal("Load() should return error for unknown fields in global config")
	}
	if !strings.Contains(err.Error(), "global config") {
		t.Errorf("error should mention 'global config': %v", err)
	}
}

func TestLoad_RelativePathFromProjectConfig(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory (no global config)
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create project config with relative paths
	projectConfig := `
version: "1"
logo: "assets/logo.png"
default_template: "./templates/default.md"
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(projectConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Relative paths from project config should resolve relative to project root (cwd)
	// Since we're in tmpDir, the paths should be relative to tmpDir
	expectedLogo := filepath.Clean("assets/logo.png")
	if cfg.Logo != expectedLogo {
		t.Errorf("Logo = %q, want %q", cfg.Logo, expectedLogo)
	}

	expectedTemplate := filepath.Clean("templates/default.md")
	if cfg.DefaultTemplate != expectedTemplate {
		t.Errorf("DefaultTemplate = %q, want %q", cfg.DefaultTemplate, expectedTemplate)
	}
}

func TestLoad_RelativePathFromGlobalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)

	// Create a working directory without project config
	workDir := filepath.Join(tmpDir, "work")
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create global config with relative paths
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	globalConfig := `
version: "1"
logo: "logos/company.png"
`
	os.WriteFile(filepath.Join(globalDir, "mdpress.yml"), []byte(globalConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Relative paths from global config should resolve relative to ~/.config/mdpress/
	expectedLogo := filepath.Join(globalDir, "logos", "company.png")
	if cfg.Logo != expectedLogo {
		t.Errorf("Logo = %q, want %q", cfg.Logo, expectedLogo)
	}
}

func TestLoad_AbsolutePathUnchanged(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create project config with absolute path
	absPath := "/absolute/path/to/logo.png"
	projectConfig := `
version: "1"
logo: "` + absPath + `"
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(projectConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Absolute paths should remain unchanged
	if cfg.Logo != absPath {
		t.Errorf("Logo = %q, want %q (unchanged absolute path)", cfg.Logo, absPath)
	}
}

func TestLoadWithSources(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create global config
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "mdpress.yml")
	globalConfig := `
version: "1"
font: "Arial"
`
	os.WriteFile(globalPath, []byte(globalConfig), 0644)

	// Create project config
	projectPath := filepath.Join(tmpDir, "mdpress.yml")
	projectConfig := `
version: "1"
accent_color: "#FF0000"
`
	os.WriteFile(projectPath, []byte(projectConfig), 0644)

	result, err := LoadWithSources()
	if err != nil {
		t.Fatalf("LoadWithSources() error = %v, want nil", err)
	}

	// Check font source (from global)
	if src, ok := result.Sources["font"]; !ok {
		t.Error("Sources should contain 'font'")
	} else if src.Source != "global" {
		t.Errorf("font source = %q, want %q", src.Source, "global")
	}

	// Check accent_color source (from project)
	if src, ok := result.Sources["accent_color"]; !ok {
		t.Error("Sources should contain 'accent_color'")
	} else if src.Source != "project" {
		t.Errorf("accent_color source = %q, want %q", src.Source, "project")
	}

	// Check logo_position source (from default, since not overridden)
	if src, ok := result.Sources["logo_position"]; !ok {
		t.Error("Sources should contain 'logo_position'")
	} else if src.Source != "default" {
		t.Errorf("logo_position source = %q, want %q", src.Source, "default")
	}
}

func TestLoad_ValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create project config with invalid value
	projectConfig := `
version: "1"
logo_position: "invalid-position"
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(projectConfig), 0644)

	_, err = Load()
	if err == nil {
		t.Fatal("Load() should return validation error")
	}
	if !strings.Contains(err.Error(), "logo_position") {
		t.Errorf("error should mention 'logo_position': %v", err)
	}
}

func TestLoad_NestedMargins(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create global config with some margins
	globalDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(globalDir, 0755)
	globalConfig := `
version: "1"
margins:
  top: 30
  left: 25
`
	os.WriteFile(filepath.Join(globalDir, "mdpress.yml"), []byte(globalConfig), 0644)

	// Create project config that overrides only some margins
	projectConfig := `
version: "1"
margins:
  top: 40
`
	os.WriteFile(filepath.Join(tmpDir, "mdpress.yml"), []byte(projectConfig), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Top should be from project
	if cfg.Margins.Top != 40 {
		t.Errorf("Margins.Top = %v, want %v (from project)", cfg.Margins.Top, 40)
	}

	// Left should be from global
	if cfg.Margins.Left != 25 {
		t.Errorf("Margins.Left = %v, want %v (from global)", cfg.Margins.Left, 25)
	}

	// Right and Bottom should be from defaults
	defaults := DefaultConfig()
	if cfg.Margins.Right != defaults.Margins.Right {
		t.Errorf("Margins.Right = %v, want %v (from default)", cfg.Margins.Right, defaults.Margins.Right)
	}
	if cfg.Margins.Bottom != defaults.Margins.Bottom {
		t.Errorf("Margins.Bottom = %v, want %v (from default)", cfg.Margins.Bottom, defaults.Margins.Bottom)
	}
}

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		source   ConfigSource
		wantFunc func(string) string
	}{
		{
			name:   "absolute path unchanged",
			path:   "/absolute/path/file.png",
			source: ConfigSource{Source: "project", Path: "mdpress.yml"},
			wantFunc: func(_ string) string {
				return "/absolute/path/file.png"
			},
		},
		{
			name:   "relative from project",
			path:   "assets/logo.png",
			source: ConfigSource{Source: "project", Path: "mdpress.yml"},
			wantFunc: func(_ string) string {
				return "assets/logo.png"
			},
		},
		{
			name:   "relative from global",
			path:   "logos/company.png",
			source: ConfigSource{Source: "global", Path: "/home/user/.config/mdpress/mdpress.yml"},
			wantFunc: func(_ string) string {
				return "/home/user/.config/mdpress/logos/company.png"
			},
		},
		{
			name:   "default source unchanged",
			path:   "some/path.png",
			source: ConfigSource{Source: "default", Path: ""},
			wantFunc: func(_ string) string {
				return "some/path.png"
			},
		},
		{
			name:   "dot-slash relative path",
			path:   "./templates/doc.md",
			source: ConfigSource{Source: "project", Path: "mdpress.yml"},
			wantFunc: func(_ string) string {
				return "templates/doc.md"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePath(tt.path, tt.source)
			if err != nil {
				t.Fatalf("resolvePath() error = %v", err)
			}
			want := tt.wantFunc(got)
			if got != want {
				t.Errorf("resolvePath() = %q, want %q", got, want)
			}
		})
	}
}
