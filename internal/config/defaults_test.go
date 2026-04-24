package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Version != "1" {
		t.Errorf("Version = %q, want %q", cfg.Version, "1")
	}

	if cfg.LogoPosition != "top-right" {
		t.Errorf("LogoPosition = %q, want %q", cfg.LogoPosition, "top-right")
	}

	if cfg.LogoWidth != 120 {
		t.Errorf("LogoWidth = %v, want %v", cfg.LogoWidth, 120)
	}

	if cfg.Font != "Helvetica" {
		t.Errorf("Font = %q, want %q", cfg.Font, "Helvetica")
	}

	if cfg.FontSize != 12 {
		t.Errorf("FontSize = %v, want %v", cfg.FontSize, 12)
	}

	if cfg.Header != false {
		t.Errorf("Header = %v, want %v", cfg.Header, false)
	}

	if cfg.Footer != "" {
		t.Errorf("Footer = %q, want empty string", cfg.Footer)
	}

	if cfg.Margins.Top != 20 {
		t.Errorf("Margins.Top = %v, want %v", cfg.Margins.Top, 20)
	}

	if cfg.Margins.Right != 20 {
		t.Errorf("Margins.Right = %v, want %v", cfg.Margins.Right, 20)
	}

	if cfg.Margins.Bottom != 20 {
		t.Errorf("Margins.Bottom = %v, want %v", cfg.Margins.Bottom, 20)
	}

	if cfg.Margins.Left != 20 {
		t.Errorf("Margins.Left = %v, want %v", cfg.Margins.Left, 20)
	}
}

func TestDefaultConfigHasEmptyOptionalFields(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Logo != "" {
		t.Errorf("Logo = %q, want empty string", cfg.Logo)
	}

	if cfg.DefaultTemplate != "" {
		t.Errorf("DefaultTemplate = %q, want empty string", cfg.DefaultTemplate)
	}

	if cfg.AccentColor != "" {
		t.Errorf("AccentColor = %q, want empty string", cfg.AccentColor)
	}

	if cfg.SourcePath != "" {
		t.Errorf("SourcePath = %q, want empty string", cfg.SourcePath)
	}
}
