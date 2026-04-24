package config

import (
	"fmt"
	"regexp"
)

// validLogoPositions contains the allowed values for logo position.
var validLogoPositions = map[string]bool{
	"top-left":     true,
	"top-right":    true,
	"bottom-left":  true,
	"bottom-right": true,
}

// hexColorRegex matches #RGB or #RRGGBB hex color formats.
var hexColorRegex = regexp.MustCompile(`^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$`)

// Validate checks all config values and returns an error if any are invalid.
func Validate(cfg *Config) error {
	if err := ValidateVersion(cfg.Version); err != nil {
		return err
	}
	if err := ValidateLogoPosition(cfg.LogoPosition); err != nil {
		return err
	}
	if err := ValidateLogoWidth(cfg.LogoWidth); err != nil {
		return err
	}
	if err := ValidateFontSize(cfg.FontSize); err != nil {
		return err
	}
	if err := ValidateAccentColor(cfg.AccentColor); err != nil {
		return err
	}
	if err := ValidateMargins(cfg.Margins); err != nil {
		return err
	}
	return nil
}

// ValidateVersion checks that the version is a supported value.
// Currently only "1" is valid.
func ValidateVersion(version string) error {
	if version != "1" {
		return fmt.Errorf("version: must be \"1\", got %q", version)
	}
	return nil
}

// ValidateLogoPosition checks that the logo position is one of the allowed values.
// Valid positions: top-left, top-right, bottom-left, bottom-right.
func ValidateLogoPosition(position string) error {
	if !validLogoPositions[position] {
		return fmt.Errorf("logo_position: must be one of top-left, top-right, bottom-left, bottom-right; got %q", position)
	}
	return nil
}

// ValidateLogoWidth checks that the logo width is a positive integer.
func ValidateLogoWidth(width float64) error {
	if width <= 0 {
		return fmt.Errorf("logo_width: must be a positive number, got %v", width)
	}
	if width != float64(int(width)) {
		return fmt.Errorf("logo_width: must be an integer, got %v", width)
	}
	return nil
}

// ValidateFontSize checks that the font size is a positive number.
func ValidateFontSize(size float64) error {
	if size <= 0 {
		return fmt.Errorf("font_size: must be a positive number, got %v", size)
	}
	return nil
}

// ValidateAccentColor checks that the accent color is a valid hex color.
// Valid formats: #RGB or #RRGGBB. Empty string is allowed (no accent color).
func ValidateAccentColor(color string) error {
	if color == "" {
		return nil
	}
	if !hexColorRegex.MatchString(color) {
		return fmt.Errorf("accent_color: must be a hex color in #RGB or #RRGGBB format, got %q", color)
	}
	return nil
}

// ValidateMargins checks that all margin values are non-negative.
func ValidateMargins(margins Margins) error {
	if margins.Top < 0 {
		return fmt.Errorf("margins.top: must be non-negative, got %v", margins.Top)
	}
	if margins.Right < 0 {
		return fmt.Errorf("margins.right: must be non-negative, got %v", margins.Right)
	}
	if margins.Bottom < 0 {
		return fmt.Errorf("margins.bottom: must be non-negative, got %v", margins.Bottom)
	}
	if margins.Left < 0 {
		return fmt.Errorf("margins.left: must be non-negative, got %v", margins.Left)
	}
	return nil
}
