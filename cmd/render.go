package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	mdpress "github.com/skel-tech/mdpress-core"
	"github.com/skel-tech/mdpress/internal/auth"
	"github.com/skel-tech/mdpress/internal/data"
	"github.com/skel-tech/mdpress/internal/template"
	"github.com/spf13/cobra"
)

var (
	output       string
	templateFlag string
	dataFile     string
	logo         string
	logoPos      string
	logoWidth    float64
	font         string
	accentColor  string
	marginTop    float64
	marginRight  float64
	marginBot    float64
	marginLeft   float64
)

var renderCmd = &cobra.Command{
	Use:   "render [file]",
	Short: "Render a Markdown file to PDF",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		// Determine the input source (file or interpolated content)
		var inputReader io.Reader

		if dataFile != "" {
			// Data-driven rendering requires Pro authentication
			if err := auth.RequirePro("data"); err != nil {
				return err
			}

			// Data-driven rendering: load data, read markdown, interpolate
			dataValues, err := data.LoadFile(dataFile)
			if err != nil {
				return fmt.Errorf("loading data file: %w", err)
			}

			mdContent, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading markdown file: %w", err)
			}

			interpolated, err := data.Interpolate(string(mdContent), dataValues)
			if err != nil {
				return err
			}

			inputReader = strings.NewReader(interpolated)
		} else {
			// Standard rendering: use file directly
			f, err := os.Open(input)
			if err != nil {
				return fmt.Errorf("opening input: %w", err)
			}
			defer f.Close()
			inputReader = f
		}

		if output == "" {
			output = stripExt(input) + ".pdf"
		}

		out, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating output: %w", err)
		}
		defer out.Close()

		opts := mdpress.DefaultOptions()

		// Apply config values (if config was loaded)
		if loadedConfig != nil {
			applyConfigToOptions(&opts)
		}

		// Resolve and apply template
		if err := applyTemplate(cmd); err != nil {
			return err
		}

		// Apply CLI flag overrides (CLI flags take precedence over config)
		applyFlagOverrides(cmd, &opts)

		if err := mdpress.Render(inputReader, out, opts); err != nil {
			return fmt.Errorf("rendering: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Written to %s\n", output)
		return nil
	},
}

func init() {
	// Existing flags
	renderCmd.Flags().StringVarP(&output, "output", "o", "", "output PDF path (default: input with .pdf extension)")
	renderCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "template name or path to apply (e.g. invoice or ./custom.yml)")
	renderCmd.Flags().StringVarP(&dataFile, "data", "d", "", "data file (JSON or YAML) for template interpolation")

	// New config-related flags
	renderCmd.Flags().StringVar(&logo, "logo", "", "path to logo image")
	renderCmd.Flags().StringVar(&logoPos, "logo-position", "", "logo position (top-left, top-right, bottom-left, bottom-right)")
	renderCmd.Flags().Float64Var(&logoWidth, "logo-width", 0, "logo width in millimeters")
	renderCmd.Flags().StringVar(&font, "font", "", "font family (Helvetica, Times, Courier)")
	renderCmd.Flags().StringVar(&accentColor, "accent-color", "", "accent color in hex format (#RGB or #RRGGBB)")
	renderCmd.Flags().Float64Var(&marginTop, "margin-top", 0, "top margin in millimeters")
	renderCmd.Flags().Float64Var(&marginRight, "margin-right", 0, "right margin in millimeters")
	renderCmd.Flags().Float64Var(&marginBot, "margin-bottom", 0, "bottom margin in millimeters")
	renderCmd.Flags().Float64Var(&marginLeft, "margin-left", 0, "left margin in millimeters")

	rootCmd.AddCommand(renderCmd)
}

// applyConfigToOptions applies the loaded config values to render options.
func applyConfigToOptions(opts *mdpress.Options) {
	// Apply font
	if loadedConfig.Font != "" {
		opts.DefaultFamily = loadedConfig.Font
	}

	// Apply margins
	if loadedConfig.Margins.Top > 0 {
		opts.MarginTop = loadedConfig.Margins.Top
	}
	if loadedConfig.Margins.Right > 0 {
		opts.MarginRight = loadedConfig.Margins.Right
	}
	if loadedConfig.Margins.Bottom > 0 {
		opts.MarginBottom = loadedConfig.Margins.Bottom
	}
	if loadedConfig.Margins.Left > 0 {
		opts.MarginLeft = loadedConfig.Margins.Left
	}

	// Logo, logo position, logo width, and accent color are tracked in config
	// but will be applied when the core library supports these features.
	// For now, we store them for future use:
	if logo == "" && loadedConfig.Logo != "" {
		logo = loadedConfig.Logo
	}
	if logoPos == "" && loadedConfig.LogoPosition != "" {
		logoPos = loadedConfig.LogoPosition
	}
	if logoWidth == 0 && loadedConfig.LogoWidth > 0 {
		logoWidth = loadedConfig.LogoWidth
	}
	if accentColor == "" && loadedConfig.AccentColor != "" {
		accentColor = loadedConfig.AccentColor
	}

	// TODO: Apply font_size when mdpress-core supports it
	// if loadedConfig.FontSize > 0 { opts.FontSize = loadedConfig.FontSize }

	// TODO: Apply header when mdpress-core supports it
	// if loadedConfig.Header { opts.Header = loadedConfig.Header }

	// TODO: Apply footer when mdpress-core supports it
	// if loadedConfig.Footer != "" { opts.Footer = loadedConfig.Footer }
}

// applyFlagOverrides applies CLI flag values that override config values.
func applyFlagOverrides(cmd *cobra.Command, opts *mdpress.Options) {
	// Only apply flag values if they were explicitly set on the command line
	if cmd.Flags().Changed("font") {
		opts.DefaultFamily = font
	}
	if cmd.Flags().Changed("margin-top") {
		opts.MarginTop = marginTop
	}
	if cmd.Flags().Changed("margin-right") {
		opts.MarginRight = marginRight
	}
	if cmd.Flags().Changed("margin-bottom") {
		opts.MarginBottom = marginBot
	}
	if cmd.Flags().Changed("margin-left") {
		opts.MarginLeft = marginLeft
	}

	// TODO: Apply logo, logo-position, logo-width, and accent-color
	// when the mdpress-core library supports these features.
	_ = logo
	_ = logoPos
	_ = logoWidth
	_ = accentColor
}

// applyTemplate resolves and applies the template to the loaded config.
// Resolution order:
// 1. CLI --template flag (highest priority)
// 2. Config default_template field
// 3. No template (use config values directly)
func applyTemplate(cmd *cobra.Command) error {
	// Determine which template to use
	templateRef := ""
	if cmd.Flags().Changed("template") {
		// CLI flag takes highest priority
		templateRef = templateFlag
	} else if loadedConfig != nil && loadedConfig.DefaultTemplate != "" {
		// Fall back to config default_template
		templateRef = loadedConfig.DefaultTemplate
	}

	// No template specified
	if templateRef == "" {
		return nil
	}

	// Resolve the template (handles both names and paths)
	tmpl, err := template.Resolve(templateRef)
	if err != nil {
		return fmt.Errorf("template: %w", err)
	}

	// Apply template values to config
	if loadedConfig != nil {
		template.Apply(loadedConfig, tmpl)
	}

	return nil
}

func stripExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[:i]
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return path
}
