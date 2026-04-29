package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skel-tech/mdpress/internal/config"
	"github.com/spf13/cobra"
)

var projectConfig bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mdpress configuration",
	Long: `Initialize mdpress configuration with default settings.

By default, creates the global config directory structure at ~/.config/mdpress/
including a config file with all available options, a logos directory,
and a default template.

With --project, creates a local mdpress.yml in the current directory
for project-specific overrides.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&projectConfig, "project", false, "create local project config (./mdpress.yml)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	if projectConfig {
		return runProjectInit(cmd)
	}
	return runGlobalInit(cmd)
}

func runProjectInit(cmd *cobra.Command) error {
	configPath := config.ProjectConfigPath()
	out := cmd.OutOrStdout()

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Fprintf(out, "Skipped %s (already exists)\n", configPath)
		return nil
	}

	// Create the project config file
	if err := os.WriteFile(configPath, []byte(projectConfigTemplate), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", configPath, err)
	}

	fmt.Fprintf(out, "Created %s\n", configPath)
	return nil
}

func runGlobalInit(cmd *cobra.Command) error {
	// Get the global config path and derive base directory
	configPath := config.GlobalConfigPath()
	if configPath == "" {
		return fmt.Errorf("could not determine home directory")
	}
	baseDir := filepath.Dir(configPath)

	// Convert to absolute path for display
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	result := &scaffoldResult{}

	// Create base directory
	if err := result.scaffoldDir(absBaseDir); err != nil {
		return err
	}

	// Create mdpress.yml
	configFilePath := filepath.Join(absBaseDir, "mdpress.yml")
	if err := result.scaffoldFile(configFilePath, globalConfigTemplate); err != nil {
		return err
	}

	// Create logos directory
	logosDir := filepath.Join(absBaseDir, "logos")
	if err := result.scaffoldDir(logosDir); err != nil {
		return err
	}

	// Create templates directory
	templatesDir := filepath.Join(absBaseDir, "templates")
	if err := result.scaffoldDir(templatesDir); err != nil {
		return err
	}

	// Create templates/default.yml
	defaultTemplatePath := filepath.Join(templatesDir, "default.yml")
	if err := result.scaffoldFile(defaultTemplatePath, defaultTemplateContent); err != nil {
		return err
	}

	// Print results
	out := cmd.OutOrStdout()

	// If everything was skipped, print summary
	if len(result.created) == 0 && len(result.skipped) > 0 {
		fmt.Fprintf(out, "Config already initialized at %s\n", absBaseDir)
		return nil
	}

	// Print created items
	for _, path := range result.created {
		fmt.Fprintf(out, "Created %s\n", path)
	}

	// Print skipped items
	for _, path := range result.skipped {
		fmt.Fprintf(out, "Skipped %s (already exists)\n", path)
	}

	return nil
}

// scaffoldResult tracks what was created or skipped during initialization.
type scaffoldResult struct {
	created []string
	skipped []string
}

// scaffoldFile creates a file with the given content if it doesn't exist.
func (r *scaffoldResult) scaffoldFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		r.skipped = append(r.skipped, path)
		return nil
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	r.created = append(r.created, path)
	return nil
}

// scaffoldDir creates a directory if it doesn't exist.
func (r *scaffoldResult) scaffoldDir(path string) error {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		r.skipped = append(r.skipped, path)
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", path, err)
	}

	r.created = append(r.created, path)
	return nil
}

const globalConfigTemplate = `# mdpress global configuration
# For documentation, see: https://github.com/Skel-Tech/mdpress
#
# This is the global config file that provides user-wide defaults.
# Project-specific configs (./mdpress.yml) override these settings.

# Config file version (required)
version: "1"

# Path to logo image (relative to this config file or absolute)
# Place logo files in the logos/ directory and reference them here.
# Supported formats: PNG, JPG, SVG
# logo: "logos/company-logo.png"

# Logo position on the page
# Options: top-left, top-center, top-right, bottom-left, bottom-center, bottom-right
logo_position: "top-right"

# Logo width in millimeters
logo_width: 120

# Default template to use for rendering
# Templates are YAML files that override config values for specific document types.
# Reference templates from the templates/ directory.
# default_template: "templates/invoice.yml"

# Font family for the document
# Options: Helvetica, Times, Courier (built-in), or path to custom font file
font: "Helvetica"

# Accent color for headings and other elements (hex format)
# accent_color: "#336699"

# Page margins in millimeters
margins:
  top: 20
  right: 20
  bottom: 20
  left: 20
`

const defaultTemplateContent = `# mdpress template: default
# This is a template for document-type-specific configuration overrides.
#
# Templates allow you to define presets for different document types
# (e.g., invoices, reports, letters). When rendering, specify a template
# with --template or set default_template in your config.
#
# All fields below are optional. Uncomment and modify only the values
# you want to override from your global or project config.

# Config file version (required if using this template standalone)
# version: "1"

# Path to logo image (relative to this file or absolute)
# logo: "logos/template-logo.png"

# Logo position on the page
# Options: top-left, top-center, top-right, bottom-left, bottom-center, bottom-right
# logo_position: "top-right"

# Logo width in millimeters
# logo_width: 120

# Font family for the document
# Options: Helvetica, Times, Courier (built-in), or path to custom font file
# font: "Helvetica"

# Accent color for headings and other elements (hex format)
# accent_color: "#336699"

# Page margins in millimeters
# margins:
#   top: 20
#   right: 20
#   bottom: 20
#   left: 20
`

const projectConfigTemplate = `# mdpress project configuration
# This file overrides your global config (~/.config/mdpress/mdpress.yml)
# for this project. Only include settings you want to customize.
#
# For documentation, see: https://github.com/Skel-Tech/mdpress

# Config file version (required)
version: "1"

# Path to project logo (relative to this file or absolute)
# logo: "assets/logo.png"

# Accent color for headings and other elements (hex format)
# accent_color: "#336699"

# Default template to use for rendering
# default_template: "templates/project-default.yml"

# Less commonly overridden settings:
# logo_position: "top-right"
# logo_width: 120
# font: "Helvetica"
# margins:
#   top: 20
#   right: 20
#   bottom: 20
#   left: 20
`
