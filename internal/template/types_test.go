package template

import (
	"testing"
)

func TestTemplateStruct(t *testing.T) {
	// Test that Template struct can be instantiated with all fields
	headerVal := true
	tmpl := Template{
		Name:         "Test Template",
		Version:      "1",
		Description:  "A test template",
		Logo:         "logo.png",
		LogoPosition: "top-left",
		LogoWidth:    100,
		Font:         "Arial",
		FontSize:     12,
		AccentColor:  "#FF0000",
		Footer:       "Page {page}",
		Header:       &headerVal,
		Margins: Margins{
			Top:    20,
			Right:  15,
			Bottom: 20,
			Left:   15,
		},
		SourcePath: "/path/to/template.yml",
	}

	if tmpl.Name != "Test Template" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Test Template")
	}
	if tmpl.Version != "1" {
		t.Errorf("Version = %q, want %q", tmpl.Version, "1")
	}
	if tmpl.Header == nil || *tmpl.Header != true {
		t.Errorf("Header = %v, want true", tmpl.Header)
	}
}

func TestTemplateHeaderPointerDistinguishesNotSetFromFalse(t *testing.T) {
	// Test that nil Header (not set) is different from false Header
	tmplNotSet := Template{
		Name:    "No Header Set",
		Version: "1",
		Header:  nil, // not set
	}

	headerFalse := false
	tmplExplicitlyFalse := Template{
		Name:    "Header False",
		Version: "1",
		Header:  &headerFalse, // explicitly false
	}

	// nil means "not set"
	if tmplNotSet.Header != nil {
		t.Error("Header should be nil when not set")
	}

	// non-nil with false value means "explicitly set to false"
	if tmplExplicitlyFalse.Header == nil {
		t.Error("Header should not be nil when explicitly set to false")
	}
	if *tmplExplicitlyFalse.Header != false {
		t.Error("Header should be false when explicitly set to false")
	}
}

func TestMarginsStruct(t *testing.T) {
	margins := Margins{
		Top:    10,
		Right:  20,
		Bottom: 30,
		Left:   40,
	}

	if margins.Top != 10 {
		t.Errorf("Top = %v, want %v", margins.Top, 10)
	}
	if margins.Right != 20 {
		t.Errorf("Right = %v, want %v", margins.Right, 20)
	}
	if margins.Bottom != 30 {
		t.Errorf("Bottom = %v, want %v", margins.Bottom, 30)
	}
	if margins.Left != 40 {
		t.Errorf("Left = %v, want %v", margins.Left, 40)
	}
}
