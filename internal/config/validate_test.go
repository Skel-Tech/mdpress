package config

import (
	"strings"
	"testing"
)

func TestValidate_ValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := Validate(&cfg); err != nil {
		t.Errorf("Validate(DefaultConfig()) = %v, want nil", err)
	}
}

func TestValidate_ValidConfigWithAccentColor(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AccentColor = "#FF5500"
	if err := Validate(&cfg); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"valid version 1", "1", false},
		{"invalid version 2", "2", true},
		{"invalid empty", "", true},
		{"invalid version 0", "0", true},
		{"invalid version string", "one", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "version") {
				t.Errorf("error should contain field name 'version': %v", err)
			}
		})
	}
}

func TestValidateLogoPosition(t *testing.T) {
	tests := []struct {
		name     string
		position string
		wantErr  bool
	}{
		{"valid top-left", "top-left", false},
		{"valid top-right", "top-right", false},
		{"valid bottom-left", "bottom-left", false},
		{"valid bottom-right", "bottom-right", false},
		{"invalid center", "center", true},
		{"invalid empty", "", true},
		{"invalid top", "top", true},
		{"invalid left", "left", true},
		{"invalid topleft no hyphen", "topleft", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogoPosition(tt.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogoPosition(%q) error = %v, wantErr %v", tt.position, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "logo_position") {
				t.Errorf("error should contain field name 'logo_position': %v", err)
			}
		})
	}
}

func TestValidateLogoWidth(t *testing.T) {
	tests := []struct {
		name    string
		width   float64
		wantErr bool
	}{
		{"valid 100", 100, false},
		{"valid 1", 1, false},
		{"valid 120", 120, false},
		{"invalid zero", 0, true},
		{"invalid negative", -10, true},
		{"invalid decimal", 10.5, true},
		{"invalid small decimal", 0.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogoWidth(tt.width)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogoWidth(%v) error = %v, wantErr %v", tt.width, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "logo_width") {
				t.Errorf("error should contain field name 'logo_width': %v", err)
			}
		})
	}
}

func TestValidateFontSize(t *testing.T) {
	tests := []struct {
		name    string
		size    float64
		wantErr bool
	}{
		{"valid 12", 12, false},
		{"valid 10.5", 10.5, false},
		{"valid 1", 1, false},
		{"valid large", 72, false},
		{"invalid zero", 0, true},
		{"invalid negative", -10, true},
		{"invalid small negative", -0.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFontSize(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFontSize(%v) error = %v, wantErr %v", tt.size, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "font_size") {
				t.Errorf("error should contain field name 'font_size': %v", err)
			}
		})
	}
}

func TestValidateAccentColor(t *testing.T) {
	tests := []struct {
		name    string
		color   string
		wantErr bool
	}{
		{"valid empty", "", false},
		{"valid 6-digit hex", "#FF5500", false},
		{"valid 3-digit hex", "#F50", false},
		{"valid lowercase", "#ff5500", false},
		{"valid mixed case", "#Ff5500", false},
		{"invalid no hash", "FF5500", true},
		{"invalid 4-digit", "#FF55", true},
		{"invalid 5-digit", "#FF550", true},
		{"invalid 7-digit", "#FF55001", true},
		{"invalid non-hex chars", "#GGGGGG", true},
		{"invalid rgb format", "rgb(255,0,0)", true},
		{"invalid named color", "red", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccentColor(tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccentColor(%q) error = %v, wantErr %v", tt.color, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "accent_color") {
				t.Errorf("error should contain field name 'accent_color': %v", err)
			}
		})
	}
}

func TestValidateMargins(t *testing.T) {
	tests := []struct {
		name       string
		margins    Margins
		wantErr    bool
		fieldInErr string
	}{
		{
			name:    "valid all positive",
			margins: Margins{Top: 20, Right: 20, Bottom: 20, Left: 20},
			wantErr: false,
		},
		{
			name:    "valid all zero",
			margins: Margins{Top: 0, Right: 0, Bottom: 0, Left: 0},
			wantErr: false,
		},
		{
			name:       "invalid negative top",
			margins:    Margins{Top: -1, Right: 20, Bottom: 20, Left: 20},
			wantErr:    true,
			fieldInErr: "margins.top",
		},
		{
			name:       "invalid negative right",
			margins:    Margins{Top: 20, Right: -1, Bottom: 20, Left: 20},
			wantErr:    true,
			fieldInErr: "margins.right",
		},
		{
			name:       "invalid negative bottom",
			margins:    Margins{Top: 20, Right: 20, Bottom: -1, Left: 20},
			wantErr:    true,
			fieldInErr: "margins.bottom",
		},
		{
			name:       "invalid negative left",
			margins:    Margins{Top: 20, Right: 20, Bottom: 20, Left: -1},
			wantErr:    true,
			fieldInErr: "margins.left",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMargins(tt.margins)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMargins() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.fieldInErr != "" && !strings.Contains(err.Error(), tt.fieldInErr) {
				t.Errorf("error should contain field name %q: %v", tt.fieldInErr, err)
			}
		})
	}
}

func TestValidate_InvalidVersion(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Version = "2"
	err := Validate(&cfg)
	if err == nil {
		t.Error("Validate() should return error for invalid version")
	}
	if !strings.Contains(err.Error(), "version") {
		t.Errorf("error should contain field name 'version': %v", err)
	}
}

func TestValidate_InvalidLogoPosition(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LogoPosition = "center"
	err := Validate(&cfg)
	if err == nil {
		t.Error("Validate() should return error for invalid logo position")
	}
	if !strings.Contains(err.Error(), "logo_position") {
		t.Errorf("error should contain field name 'logo_position': %v", err)
	}
}

func TestValidate_InvalidLogoWidth(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LogoWidth = -10
	err := Validate(&cfg)
	if err == nil {
		t.Error("Validate() should return error for invalid logo width")
	}
	if !strings.Contains(err.Error(), "logo_width") {
		t.Errorf("error should contain field name 'logo_width': %v", err)
	}
}

func TestValidate_InvalidFontSize(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FontSize = 0
	err := Validate(&cfg)
	if err == nil {
		t.Error("Validate() should return error for invalid font size")
	}
	if !strings.Contains(err.Error(), "font_size") {
		t.Errorf("error should contain field name 'font_size': %v", err)
	}
}

func TestValidate_InvalidAccentColor(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AccentColor = "invalid"
	err := Validate(&cfg)
	if err == nil {
		t.Error("Validate() should return error for invalid accent color")
	}
	if !strings.Contains(err.Error(), "accent_color") {
		t.Errorf("error should contain field name 'accent_color': %v", err)
	}
}

func TestValidate_InvalidMargins(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Margins.Top = -5
	err := Validate(&cfg)
	if err == nil {
		t.Error("Validate() should return error for invalid margins")
	}
	if !strings.Contains(err.Error(), "margins") {
		t.Errorf("error should contain field name 'margins': %v", err)
	}
}
