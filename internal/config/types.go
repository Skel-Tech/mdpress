// Package config provides configuration management for mdpress.
package config

import (
	// viper is used for config file loading and binding to these structs
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

// Config holds the complete mdpress configuration.
type Config struct {
	Version         string  `mapstructure:"version"`
	Logo            string  `mapstructure:"logo"`
	LogoPosition    string  `mapstructure:"logo_position"`
	LogoWidth       float64 `mapstructure:"logo_width"`
	DefaultTemplate string  `mapstructure:"default_template"`
	Font            string  `mapstructure:"font"`
	FontSize        float64 `mapstructure:"font_size"`
	AccentColor     string  `mapstructure:"accent_color"`
	Header          bool    `mapstructure:"header"`
	Footer          string  `mapstructure:"footer"`
	Margins         Margins `mapstructure:"margins"`

	// SourcePath stores the path of the config file this was loaded from.
	// Used for resolving relative paths (e.g., logo paths) later.
	SourcePath string `mapstructure:"-"`
}
