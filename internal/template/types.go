// Package template provides template management for mdpress.
package template

import (
	// viper is used for template file loading and binding to these structs
	// via the mapstructure tags defined on the struct fields.
	_ "github.com/spf13/viper"
)

// Margins defines the page margins in millimeters.
type Margins struct {
	Top    float64 `mapstructure:"top"`
	Right  float64 `mapstructure:"right"`
	Bottom float64 `mapstructure:"bottom"`
	Left   float64 `mapstructure:"left"`
}

// Template holds the configuration for an mdpress template.
type Template struct {
	// Required fields
	Name    string `mapstructure:"name"`    // Display name for the template
	Version string `mapstructure:"version"` // Schema version, must be "1"

	// Optional metadata
	Description string `mapstructure:"description"` // Description for display

	// Overridable config fields (all optional in templates)
	Logo         string  `mapstructure:"logo"`
	LogoPosition string  `mapstructure:"logo_position"`
	LogoWidth    float64 `mapstructure:"logo_width"`
	Font         string  `mapstructure:"font"`
	FontSize     float64 `mapstructure:"font_size"`
	AccentColor  string  `mapstructure:"accent_color"`
	Footer       string  `mapstructure:"footer"`
	Margins      Margins `mapstructure:"margins"`

	// Header uses pointer type so we can distinguish "not set" from "explicitly false"
	Header *bool `mapstructure:"header"`

	// SourcePath stores the path of the template file this was loaded from.
	SourcePath string `mapstructure:"-"`
}
