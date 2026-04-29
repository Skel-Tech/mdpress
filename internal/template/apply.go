package template

import "github.com/skel-tech/mdpress/internal/config"

// Apply overrides config values with non-zero template values.
// Template values take precedence, completely overriding the corresponding config values.
// Zero/empty values in the template do not modify the config.
func Apply(cfg *config.Config, tmpl *Template) {
	if tmpl == nil {
		return
	}

	if tmpl.Logo != "" {
		cfg.Logo = tmpl.Logo
	}
	if tmpl.LogoPosition != "" {
		cfg.LogoPosition = tmpl.LogoPosition
	}
	if tmpl.LogoWidth != 0 {
		cfg.LogoWidth = tmpl.LogoWidth
	}
	if tmpl.Font != "" {
		cfg.Font = tmpl.Font
	}
	if tmpl.FontSize != 0 {
		cfg.FontSize = tmpl.FontSize
	}
	if tmpl.AccentColor != "" {
		cfg.AccentColor = tmpl.AccentColor
	}
	if tmpl.Footer != "" {
		cfg.Footer = tmpl.Footer
	}

	// Header uses *bool so we can distinguish "not set" from "explicitly false"
	if tmpl.Header != nil {
		cfg.Header = *tmpl.Header
	}

	// Apply margins individually - each non-zero margin overrides
	if tmpl.Margins.Top != 0 {
		cfg.Margins.Top = tmpl.Margins.Top
	}
	if tmpl.Margins.Right != 0 {
		cfg.Margins.Right = tmpl.Margins.Right
	}
	if tmpl.Margins.Bottom != 0 {
		cfg.Margins.Bottom = tmpl.Margins.Bottom
	}
	if tmpl.Margins.Left != 0 {
		cfg.Margins.Left = tmpl.Margins.Left
	}
}
