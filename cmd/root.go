package cmd

import (
	"fmt"
	"os"

	"github.com/skel-tech/mdpress/internal/config"
	"github.com/spf13/cobra"
)

var (
	// configDebug enables verbose config loading output
	configDebug bool

	// loadedConfig holds the loaded configuration for use by subcommands
	loadedConfig *config.Config

	// loadedSources holds the source information for config debug output
	loadedSources map[string]config.ConfigSource

	// globalConfigExists tracks whether a global config file exists
	globalConfigExists bool
)

var rootCmd = &cobra.Command{
	Use:     "mdpress",
	Short:   "Convert Markdown to PDF",
	Long:    "mdpress is a CLI tool for converting Markdown documents into styled PDF files.",
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if global config exists before loading
		globalPath := config.GlobalConfigPath()
		if _, err := os.Stat(globalPath); err == nil {
			globalConfigExists = true
		}

		// Load configuration
		result, err := config.LoadWithSources()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		loadedConfig = result.Config
		loadedSources = result.Sources

		// Print config debug info if requested
		if configDebug {
			printConfigDebug()
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&configDebug, "config-debug", false, "show config resolution details")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Show first-run hint after successful execution
	// Only show if no global config exists and we haven't shown it before
	if !globalConfigExists && !configDebug {
		showFirstRunHint()
	}
}

// showFirstRunHint prints a subtle hint about creating a global config.
// Uses a marker file to avoid repeating the hint on every run.
func showFirstRunHint() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	// Check if we've already shown the hint
	markerPath := home + "/.config/mdpress/.hint-shown"
	if _, err := os.Stat(markerPath); err == nil {
		return // Marker exists, don't show hint again
	}

	// Print the hint
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Tip: Run 'mdpress init' to create a global config with your defaults.")

	// Create the marker file to avoid showing hint again
	// Ensure directory exists
	if err := os.MkdirAll(home+"/.config/mdpress", 0755); err != nil {
		return
	}
	_ = os.WriteFile(markerPath, []byte{}, 0644)
}

// printConfigDebug prints detailed config resolution information.
func printConfigDebug() {
	fmt.Fprintln(os.Stderr, "Config Debug:")
	fmt.Fprintln(os.Stderr, "")

	// Show which config files were found
	fmt.Fprintln(os.Stderr, "Config files:")
	globalPath := config.GlobalConfigPath()
	if _, err := os.Stat(globalPath); err == nil {
		fmt.Fprintf(os.Stderr, "  global:  %s (found)\n", globalPath)
	} else {
		fmt.Fprintf(os.Stderr, "  global:  %s (not found)\n", globalPath)
	}

	projectPath := config.ProjectConfigPath()
	if _, err := os.Stat(projectPath); err == nil {
		fmt.Fprintf(os.Stderr, "  project: %s (found)\n", projectPath)
	} else {
		fmt.Fprintf(os.Stderr, "  project: %s (not found)\n", projectPath)
	}
	fmt.Fprintln(os.Stderr, "")

	// Show final merged config values with source attribution
	fmt.Fprintln(os.Stderr, "Resolved config values:")
	printConfigValue("version", loadedConfig.Version, loadedSources["version"])
	printConfigValue("logo", loadedConfig.Logo, loadedSources["logo"])
	printConfigValue("logo_position", loadedConfig.LogoPosition, loadedSources["logo_position"])
	printConfigValue("logo_width", fmt.Sprintf("%.0f", loadedConfig.LogoWidth), loadedSources["logo_width"])
	printConfigValue("default_template", loadedConfig.DefaultTemplate, loadedSources["default_template"])
	printConfigValue("font", loadedConfig.Font, loadedSources["font"])
	printConfigValue("font_size", fmt.Sprintf("%.0f", loadedConfig.FontSize), loadedSources["font_size"])
	printConfigValue("accent_color", loadedConfig.AccentColor, loadedSources["accent_color"])
	printConfigValue("header", fmt.Sprintf("%v", loadedConfig.Header), loadedSources["header"])
	printConfigValue("footer", loadedConfig.Footer, loadedSources["footer"])
	printConfigValue("margins.top", fmt.Sprintf("%.0f", loadedConfig.Margins.Top), loadedSources["margins.top"])
	printConfigValue("margins.right", fmt.Sprintf("%.0f", loadedConfig.Margins.Right), loadedSources["margins.right"])
	printConfigValue("margins.bottom", fmt.Sprintf("%.0f", loadedConfig.Margins.Bottom), loadedSources["margins.bottom"])
	printConfigValue("margins.left", fmt.Sprintf("%.0f", loadedConfig.Margins.Left), loadedSources["margins.left"])
	fmt.Fprintln(os.Stderr, "")
}

// printConfigValue prints a single config value with its source.
func printConfigValue(name, value string, source config.ConfigSource) {
	sourceStr := source.Source
	if source.Path != "" {
		sourceStr = fmt.Sprintf("%s (%s)", source.Source, source.Path)
	}
	if value == "" {
		value = "(empty)"
	}
	fmt.Fprintf(os.Stderr, "  %-20s = %-20s [%s]\n", name, value, sourceStr)
}
