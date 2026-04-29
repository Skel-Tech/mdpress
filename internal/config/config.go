package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// GlobalConfigPath returns the path to the global config file.
func GlobalConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "mdpress", "mdpress.yml")
}

// ProjectConfigPath returns the path to the project config file in the current directory.
func ProjectConfigPath() string {
	return "mdpress.yml"
}

// ConfigSource tracks the source of each config value for debugging.
type ConfigSource struct {
	Source string // "default", "global", or "project"
	Path   string // file path if loaded from file
}

// LoadResult contains the loaded config and source information.
type LoadResult struct {
	Config  *Config
	Sources map[string]ConfigSource // field name -> source
}

// Load loads and merges configuration from default, global, and project sources.
// It returns the merged config or an error if YAML is malformed or contains unknown fields.
func Load() (*Config, error) {
	result, err := LoadWithSources()
	if err != nil {
		return nil, err
	}
	return result.Config, nil
}

// LoadWithSources loads configuration and returns source tracking information.
func LoadWithSources() (*LoadResult, error) {
	cfg := DefaultConfig()
	sources := make(map[string]ConfigSource)

	// Initialize all fields as coming from defaults
	initDefaultSources(sources)

	globalPath := GlobalConfigPath()
	projectPath := ProjectConfigPath()

	// Load global config if it exists
	globalCfg, globalLoaded, err := loadConfigFile(globalPath)
	if err != nil {
		return nil, fmt.Errorf("global config: %w", err)
	}
	if globalLoaded {
		mergeConfig(&cfg, globalCfg, sources, "global", globalPath)
	}

	// Load project config if it exists
	projectCfg, projectLoaded, err := loadConfigFile(projectPath)
	if err != nil {
		return nil, fmt.Errorf("project config: %w", err)
	}
	if projectLoaded {
		mergeConfig(&cfg, projectCfg, sources, "project", projectPath)
	}

	// Resolve relative paths based on their source
	if err := resolveRelativePaths(&cfg, sources); err != nil {
		return nil, err
	}

	// Validate the final merged config
	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &LoadResult{
		Config:  &cfg,
		Sources: sources,
	}, nil
}

// loadConfigFile loads a config from a YAML file.
// Returns (config, loaded, error) where loaded indicates if the file existed.
func loadConfigFile(path string) (*Config, bool, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, false, nil
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, false, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var cfg Config
	decoderConfig := mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &cfg,
		TagName:     "mapstructure",
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(v.AllSettings()); err != nil {
		return nil, false, fmt.Errorf("invalid config in %s: %w", path, err)
	}

	cfg.SourcePath = path
	return &cfg, true, nil
}

// mergeConfig merges values from src into dst, tracking sources.
// Only non-zero values from src override dst.
func mergeConfig(dst *Config, src *Config, sources map[string]ConfigSource, sourceName, sourcePath string) {
	if src.Version != "" {
		dst.Version = src.Version
		sources["version"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Logo != "" {
		dst.Logo = src.Logo
		sources["logo"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.LogoPosition != "" {
		dst.LogoPosition = src.LogoPosition
		sources["logo_position"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.LogoWidth != 0 {
		dst.LogoWidth = src.LogoWidth
		sources["logo_width"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.DefaultTemplate != "" {
		dst.DefaultTemplate = src.DefaultTemplate
		sources["default_template"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Font != "" {
		dst.Font = src.Font
		sources["font"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.FontSize != 0 {
		dst.FontSize = src.FontSize
		sources["font_size"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.AccentColor != "" {
		dst.AccentColor = src.AccentColor
		sources["accent_color"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Header {
		dst.Header = src.Header
		sources["header"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Footer != "" {
		dst.Footer = src.Footer
		sources["footer"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}

	// Merge margins - check each field individually
	if src.Margins.Top != 0 {
		dst.Margins.Top = src.Margins.Top
		sources["margins.top"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Margins.Right != 0 {
		dst.Margins.Right = src.Margins.Right
		sources["margins.right"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Margins.Bottom != 0 {
		dst.Margins.Bottom = src.Margins.Bottom
		sources["margins.bottom"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}
	if src.Margins.Left != 0 {
		dst.Margins.Left = src.Margins.Left
		sources["margins.left"] = ConfigSource{Source: sourceName, Path: sourcePath}
	}

	// Track the most recent source path
	dst.SourcePath = sourcePath
}

// initDefaultSources initializes all sources as coming from defaults.
func initDefaultSources(sources map[string]ConfigSource) {
	defaultSource := ConfigSource{Source: "default", Path: ""}
	sources["version"] = defaultSource
	sources["logo"] = defaultSource
	sources["logo_position"] = defaultSource
	sources["logo_width"] = defaultSource
	sources["default_template"] = defaultSource
	sources["font"] = defaultSource
	sources["font_size"] = defaultSource
	sources["accent_color"] = defaultSource
	sources["header"] = defaultSource
	sources["footer"] = defaultSource
	sources["margins.top"] = defaultSource
	sources["margins.right"] = defaultSource
	sources["margins.bottom"] = defaultSource
	sources["margins.left"] = defaultSource
}

// resolveRelativePaths resolves relative paths in the config based on their source.
// Paths from global config are resolved relative to ~/.config/mdpress/
// Paths from project config are resolved relative to the project root (cwd)
func resolveRelativePaths(cfg *Config, sources map[string]ConfigSource) error {
	// Resolve logo path
	if cfg.Logo != "" {
		resolved, err := resolvePath(cfg.Logo, sources["logo"])
		if err != nil {
			return fmt.Errorf("failed to resolve logo path: %w", err)
		}
		cfg.Logo = resolved
	}

	// Resolve default_template path
	if cfg.DefaultTemplate != "" {
		resolved, err := resolvePath(cfg.DefaultTemplate, sources["default_template"])
		if err != nil {
			return fmt.Errorf("failed to resolve default_template path: %w", err)
		}
		cfg.DefaultTemplate = resolved
	}

	return nil
}

// resolvePath resolves a path relative to its source location.
func resolvePath(path string, source ConfigSource) (string, error) {
	// Absolute paths are returned unchanged
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Default source paths stay as-is (no file to resolve from)
	if source.Source == "default" || source.Path == "" {
		return path, nil
	}

	// Get the directory containing the config file
	configDir := filepath.Dir(source.Path)

	// For project config, resolve relative to cwd (where mdpress.yml is)
	// For global config, resolve relative to ~/.config/mdpress/
	resolved := filepath.Join(configDir, path)

	// Clean the path to resolve any ../ or ./ components
	return filepath.Clean(resolved), nil
}
