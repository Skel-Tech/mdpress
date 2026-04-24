package template

import (
	"testing"

	"github.com/skel-tech/mdpress/internal/config"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestApply(t *testing.T) {
	t.Run("nil template does nothing", func(t *testing.T) {
		cfg := &config.Config{
			Logo: "original.png",
			Font: "Arial",
		}
		Apply(cfg, nil)

		if cfg.Logo != "original.png" {
			t.Errorf("Logo changed unexpectedly: got %q", cfg.Logo)
		}
		if cfg.Font != "Arial" {
			t.Errorf("Font changed unexpectedly: got %q", cfg.Font)
		}
	})

	t.Run("empty template does not override values", func(t *testing.T) {
		cfg := &config.Config{
			Logo:         "original.png",
			LogoPosition: "top-left",
			LogoWidth:    100.0,
			Font:         "Arial",
			FontSize:     12.0,
			AccentColor:  "#FF0000",
			Header:       true,
			Footer:       "Page {page}",
			Margins: config.Margins{
				Top:    20.0,
				Right:  15.0,
				Bottom: 20.0,
				Left:   15.0,
			},
		}
		tmpl := &Template{}
		Apply(cfg, tmpl)

		if cfg.Logo != "original.png" {
			t.Errorf("Logo changed: got %q", cfg.Logo)
		}
		if cfg.LogoPosition != "top-left" {
			t.Errorf("LogoPosition changed: got %q", cfg.LogoPosition)
		}
		if cfg.LogoWidth != 100.0 {
			t.Errorf("LogoWidth changed: got %v", cfg.LogoWidth)
		}
		if cfg.Font != "Arial" {
			t.Errorf("Font changed: got %q", cfg.Font)
		}
		if cfg.FontSize != 12.0 {
			t.Errorf("FontSize changed: got %v", cfg.FontSize)
		}
		if cfg.AccentColor != "#FF0000" {
			t.Errorf("AccentColor changed: got %q", cfg.AccentColor)
		}
		if cfg.Header != true {
			t.Errorf("Header changed: got %v", cfg.Header)
		}
		if cfg.Footer != "Page {page}" {
			t.Errorf("Footer changed: got %q", cfg.Footer)
		}
		if cfg.Margins.Top != 20.0 {
			t.Errorf("Margins.Top changed: got %v", cfg.Margins.Top)
		}
		if cfg.Margins.Right != 15.0 {
			t.Errorf("Margins.Right changed: got %v", cfg.Margins.Right)
		}
		if cfg.Margins.Bottom != 20.0 {
			t.Errorf("Margins.Bottom changed: got %v", cfg.Margins.Bottom)
		}
		if cfg.Margins.Left != 15.0 {
			t.Errorf("Margins.Left changed: got %v", cfg.Margins.Left)
		}
	})

	t.Run("template overrides all config fields", func(t *testing.T) {
		cfg := &config.Config{
			Logo:         "original.png",
			LogoPosition: "top-left",
			LogoWidth:    100.0,
			Font:         "Arial",
			FontSize:     12.0,
			AccentColor:  "#FF0000",
			Header:       false,
			Footer:       "Page {page}",
			Margins: config.Margins{
				Top:    20.0,
				Right:  15.0,
				Bottom: 20.0,
				Left:   15.0,
			},
		}
		tmpl := &Template{
			Logo:         "template.png",
			LogoPosition: "top-right",
			LogoWidth:    200.0,
			Font:         "Helvetica",
			FontSize:     14.0,
			AccentColor:  "#00FF00",
			Header:       boolPtr(true),
			Footer:       "Template Footer",
			Margins: Margins{
				Top:    30.0,
				Right:  25.0,
				Bottom: 30.0,
				Left:   25.0,
			},
		}
		Apply(cfg, tmpl)

		if cfg.Logo != "template.png" {
			t.Errorf("Logo not overridden: got %q, want %q", cfg.Logo, "template.png")
		}
		if cfg.LogoPosition != "top-right" {
			t.Errorf("LogoPosition not overridden: got %q, want %q", cfg.LogoPosition, "top-right")
		}
		if cfg.LogoWidth != 200.0 {
			t.Errorf("LogoWidth not overridden: got %v, want %v", cfg.LogoWidth, 200.0)
		}
		if cfg.Font != "Helvetica" {
			t.Errorf("Font not overridden: got %q, want %q", cfg.Font, "Helvetica")
		}
		if cfg.FontSize != 14.0 {
			t.Errorf("FontSize not overridden: got %v, want %v", cfg.FontSize, 14.0)
		}
		if cfg.AccentColor != "#00FF00" {
			t.Errorf("AccentColor not overridden: got %q, want %q", cfg.AccentColor, "#00FF00")
		}
		if cfg.Header != true {
			t.Errorf("Header not overridden: got %v, want %v", cfg.Header, true)
		}
		if cfg.Footer != "Template Footer" {
			t.Errorf("Footer not overridden: got %q, want %q", cfg.Footer, "Template Footer")
		}
		if cfg.Margins.Top != 30.0 {
			t.Errorf("Margins.Top not overridden: got %v, want %v", cfg.Margins.Top, 30.0)
		}
		if cfg.Margins.Right != 25.0 {
			t.Errorf("Margins.Right not overridden: got %v, want %v", cfg.Margins.Right, 25.0)
		}
		if cfg.Margins.Bottom != 30.0 {
			t.Errorf("Margins.Bottom not overridden: got %v, want %v", cfg.Margins.Bottom, 30.0)
		}
		if cfg.Margins.Left != 25.0 {
			t.Errorf("Margins.Left not overridden: got %v, want %v", cfg.Margins.Left, 25.0)
		}
	})

	t.Run("partial template overrides only set fields", func(t *testing.T) {
		cfg := &config.Config{
			Logo:         "original.png",
			LogoPosition: "top-left",
			LogoWidth:    100.0,
			Font:         "Arial",
			FontSize:     12.0,
			AccentColor:  "#FF0000",
			Header:       true,
			Footer:       "Page {page}",
			Margins: config.Margins{
				Top:    20.0,
				Right:  15.0,
				Bottom: 20.0,
				Left:   15.0,
			},
		}
		tmpl := &Template{
			Font:        "Helvetica",
			AccentColor: "#00FF00",
			Margins: Margins{
				Top:  30.0,
				Left: 25.0,
			},
		}
		Apply(cfg, tmpl)

		// Should be overridden
		if cfg.Font != "Helvetica" {
			t.Errorf("Font not overridden: got %q", cfg.Font)
		}
		if cfg.AccentColor != "#00FF00" {
			t.Errorf("AccentColor not overridden: got %q", cfg.AccentColor)
		}
		if cfg.Margins.Top != 30.0 {
			t.Errorf("Margins.Top not overridden: got %v", cfg.Margins.Top)
		}
		if cfg.Margins.Left != 25.0 {
			t.Errorf("Margins.Left not overridden: got %v", cfg.Margins.Left)
		}

		// Should remain unchanged
		if cfg.Logo != "original.png" {
			t.Errorf("Logo changed unexpectedly: got %q", cfg.Logo)
		}
		if cfg.LogoPosition != "top-left" {
			t.Errorf("LogoPosition changed unexpectedly: got %q", cfg.LogoPosition)
		}
		if cfg.LogoWidth != 100.0 {
			t.Errorf("LogoWidth changed unexpectedly: got %v", cfg.LogoWidth)
		}
		if cfg.FontSize != 12.0 {
			t.Errorf("FontSize changed unexpectedly: got %v", cfg.FontSize)
		}
		if cfg.Header != true {
			t.Errorf("Header changed unexpectedly: got %v", cfg.Header)
		}
		if cfg.Footer != "Page {page}" {
			t.Errorf("Footer changed unexpectedly: got %q", cfg.Footer)
		}
		if cfg.Margins.Right != 15.0 {
			t.Errorf("Margins.Right changed unexpectedly: got %v", cfg.Margins.Right)
		}
		if cfg.Margins.Bottom != 20.0 {
			t.Errorf("Margins.Bottom changed unexpectedly: got %v", cfg.Margins.Bottom)
		}
	})

	t.Run("header can be explicitly set to false", func(t *testing.T) {
		cfg := &config.Config{
			Header: true,
		}
		tmpl := &Template{
			Header: boolPtr(false),
		}
		Apply(cfg, tmpl)

		if cfg.Header != false {
			t.Errorf("Header not overridden to false: got %v", cfg.Header)
		}
	})

	t.Run("header can be explicitly set to true", func(t *testing.T) {
		cfg := &config.Config{
			Header: false,
		}
		tmpl := &Template{
			Header: boolPtr(true),
		}
		Apply(cfg, tmpl)

		if cfg.Header != true {
			t.Errorf("Header not overridden to true: got %v", cfg.Header)
		}
	})

	t.Run("nil header pointer does not change config", func(t *testing.T) {
		cfg := &config.Config{
			Header: true,
		}
		tmpl := &Template{
			Header: nil, // not set
		}
		Apply(cfg, tmpl)

		if cfg.Header != true {
			t.Errorf("Header changed unexpectedly: got %v", cfg.Header)
		}
	})
}

func TestApplyStringFields(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		setValue  func(*Template, string)
		getValue  func(*config.Config) string
	}{
		{
			name:      "Logo",
			fieldName: "Logo",
			setValue:  func(t *Template, v string) { t.Logo = v },
			getValue:  func(c *config.Config) string { return c.Logo },
		},
		{
			name:      "LogoPosition",
			fieldName: "LogoPosition",
			setValue:  func(t *Template, v string) { t.LogoPosition = v },
			getValue:  func(c *config.Config) string { return c.LogoPosition },
		},
		{
			name:      "Font",
			fieldName: "Font",
			setValue:  func(t *Template, v string) { t.Font = v },
			getValue:  func(c *config.Config) string { return c.Font },
		},
		{
			name:      "AccentColor",
			fieldName: "AccentColor",
			setValue:  func(t *Template, v string) { t.AccentColor = v },
			getValue:  func(c *config.Config) string { return c.AccentColor },
		},
		{
			name:      "Footer",
			fieldName: "Footer",
			setValue:  func(t *Template, v string) { t.Footer = v },
			getValue:  func(c *config.Config) string { return c.Footer },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" overrides when set", func(t *testing.T) {
			cfg := &config.Config{}
			tmpl := &Template{}
			tt.setValue(tmpl, "new-value")
			Apply(cfg, tmpl)

			if got := tt.getValue(cfg); got != "new-value" {
				t.Errorf("%s not overridden: got %q, want %q", tt.fieldName, got, "new-value")
			}
		})

		t.Run(tt.name+" does not override when empty", func(t *testing.T) {
			cfg := &config.Config{}
			switch tt.fieldName {
			case "Logo":
				cfg.Logo = "original"
			case "LogoPosition":
				cfg.LogoPosition = "original"
			case "Font":
				cfg.Font = "original"
			case "AccentColor":
				cfg.AccentColor = "original"
			case "Footer":
				cfg.Footer = "original"
			}

			tmpl := &Template{} // empty template
			Apply(cfg, tmpl)

			if got := tt.getValue(cfg); got != "original" {
				t.Errorf("%s changed unexpectedly: got %q, want %q", tt.fieldName, got, "original")
			}
		})
	}
}

func TestApplyFloat64Fields(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		setValue  func(*Template, float64)
		getValue  func(*config.Config) float64
	}{
		{
			name:      "LogoWidth",
			fieldName: "LogoWidth",
			setValue:  func(t *Template, v float64) { t.LogoWidth = v },
			getValue:  func(c *config.Config) float64 { return c.LogoWidth },
		},
		{
			name:      "FontSize",
			fieldName: "FontSize",
			setValue:  func(t *Template, v float64) { t.FontSize = v },
			getValue:  func(c *config.Config) float64 { return c.FontSize },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" overrides when set", func(t *testing.T) {
			cfg := &config.Config{}
			tmpl := &Template{}
			tt.setValue(tmpl, 42.5)
			Apply(cfg, tmpl)

			if got := tt.getValue(cfg); got != 42.5 {
				t.Errorf("%s not overridden: got %v, want %v", tt.fieldName, got, 42.5)
			}
		})

		t.Run(tt.name+" does not override when zero", func(t *testing.T) {
			cfg := &config.Config{}
			switch tt.fieldName {
			case "LogoWidth":
				cfg.LogoWidth = 100.0
			case "FontSize":
				cfg.FontSize = 12.0
			}

			tmpl := &Template{} // zero values
			Apply(cfg, tmpl)

			var want float64
			switch tt.fieldName {
			case "LogoWidth":
				want = 100.0
			case "FontSize":
				want = 12.0
			}

			if got := tt.getValue(cfg); got != want {
				t.Errorf("%s changed unexpectedly: got %v, want %v", tt.fieldName, got, want)
			}
		})
	}
}

func TestApplyMargins(t *testing.T) {
	tests := []struct {
		name     string
		setMargin func(*Template, float64)
		getMargin func(*config.Config) float64
	}{
		{
			name:      "Top",
			setMargin: func(t *Template, v float64) { t.Margins.Top = v },
			getMargin: func(c *config.Config) float64 { return c.Margins.Top },
		},
		{
			name:      "Right",
			setMargin: func(t *Template, v float64) { t.Margins.Right = v },
			getMargin: func(c *config.Config) float64 { return c.Margins.Right },
		},
		{
			name:      "Bottom",
			setMargin: func(t *Template, v float64) { t.Margins.Bottom = v },
			getMargin: func(c *config.Config) float64 { return c.Margins.Bottom },
		},
		{
			name:      "Left",
			setMargin: func(t *Template, v float64) { t.Margins.Left = v },
			getMargin: func(c *config.Config) float64 { return c.Margins.Left },
		},
	}

	for _, tt := range tests {
		t.Run("Margins."+tt.name+" overrides when set", func(t *testing.T) {
			cfg := &config.Config{
				Margins: config.Margins{
					Top:    10.0,
					Right:  10.0,
					Bottom: 10.0,
					Left:   10.0,
				},
			}
			tmpl := &Template{}
			tt.setMargin(tmpl, 50.0)
			Apply(cfg, tmpl)

			if got := tt.getMargin(cfg); got != 50.0 {
				t.Errorf("Margins.%s not overridden: got %v, want %v", tt.name, got, 50.0)
			}
		})

		t.Run("Margins."+tt.name+" does not override when zero", func(t *testing.T) {
			cfg := &config.Config{
				Margins: config.Margins{
					Top:    10.0,
					Right:  10.0,
					Bottom: 10.0,
					Left:   10.0,
				},
			}
			tmpl := &Template{} // zero margins
			Apply(cfg, tmpl)

			if got := tt.getMargin(cfg); got != 10.0 {
				t.Errorf("Margins.%s changed unexpectedly: got %v, want %v", tt.name, got, 10.0)
			}
		})
	}
}
