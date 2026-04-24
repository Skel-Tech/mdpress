package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mdpress "github.com/skel-tech/mdpress-core"
	"github.com/skel-tech/mdpress/internal/config"
	"github.com/spf13/cobra"
)

// --- stripExt tests ---

func TestStripExt_Markdown(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"document.md", "document"},
		{"path/to/file.md", "path/to/file"},
		{"file.markdown", "file"},
		{"no-extension", "no-extension"},
		{"multiple.dots.md", "multiple.dots"},
		{"./relative.md", "./relative"},
	}

	for _, tt := range tests {
		got := stripExt(tt.input)
		if got != tt.want {
			t.Errorf("stripExt(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// --- applyConfigToOptions tests ---

func TestApplyConfigToOptions_Font(t *testing.T) {
	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()

	loadedConfig = &config.Config{Font: "Courier"}
	opts := mdpress.DefaultOptions()
	applyConfigToOptions(&opts)

	if opts.DefaultFamily != "Courier" {
		t.Errorf("expected font Courier, got %s", opts.DefaultFamily)
	}
}

func TestApplyConfigToOptions_Margins(t *testing.T) {
	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()

	loadedConfig = &config.Config{
		Margins: config.Margins{
			Top:    30,
			Right:  25,
			Bottom: 35,
			Left:   20,
		},
	}
	opts := mdpress.DefaultOptions()
	applyConfigToOptions(&opts)

	if opts.MarginTop != 30 {
		t.Errorf("expected margin top 30, got %f", opts.MarginTop)
	}
	if opts.MarginRight != 25 {
		t.Errorf("expected margin right 25, got %f", opts.MarginRight)
	}
	if opts.MarginBottom != 35 {
		t.Errorf("expected margin bottom 35, got %f", opts.MarginBottom)
	}
	if opts.MarginLeft != 20 {
		t.Errorf("expected margin left 20, got %f", opts.MarginLeft)
	}
}

func TestApplyConfigToOptions_ZeroMarginSkipped(t *testing.T) {
	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()

	loadedConfig = &config.Config{
		Margins: config.Margins{Top: 0, Right: 0, Bottom: 0, Left: 0},
	}
	opts := mdpress.DefaultOptions()
	defaultTop := opts.MarginTop
	applyConfigToOptions(&opts)

	if opts.MarginTop != defaultTop {
		t.Error("zero margins should not override defaults")
	}
}

func TestApplyConfigToOptions_EmptyFontSkipped(t *testing.T) {
	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()

	loadedConfig = &config.Config{Font: ""}
	opts := mdpress.DefaultOptions()
	defaultFont := opts.DefaultFamily
	applyConfigToOptions(&opts)

	if opts.DefaultFamily != defaultFont {
		t.Error("empty font should not override default")
	}
}

func TestApplyConfigToOptions_LogoFallback(t *testing.T) {
	oldConfig := loadedConfig
	oldLogo := logo
	oldLogoPos := logoPos
	oldLogoWidth := logoWidth
	oldAccentColor := accentColor
	defer func() {
		loadedConfig = oldConfig
		logo = oldLogo
		logoPos = oldLogoPos
		logoWidth = oldLogoWidth
		accentColor = oldAccentColor
	}()

	// Reset CLI flags
	logo = ""
	logoPos = ""
	logoWidth = 0
	accentColor = ""

	loadedConfig = &config.Config{
		Logo:         "logos/test.png",
		LogoPosition: "top-right",
		LogoWidth:    100,
		AccentColor:  "#336699",
	}
	opts := mdpress.DefaultOptions()
	applyConfigToOptions(&opts)

	if logo != "logos/test.png" {
		t.Errorf("expected logo from config, got %q", logo)
	}
	if logoPos != "top-right" {
		t.Errorf("expected logoPos from config, got %q", logoPos)
	}
	if logoWidth != 100 {
		t.Errorf("expected logoWidth from config, got %f", logoWidth)
	}
	if accentColor != "#336699" {
		t.Errorf("expected accentColor from config, got %q", accentColor)
	}
}

func TestApplyConfigToOptions_CLIFlagsTakePrecedence(t *testing.T) {
	oldConfig := loadedConfig
	oldLogo := logo
	oldLogoPos := logoPos
	defer func() {
		loadedConfig = oldConfig
		logo = oldLogo
		logoPos = oldLogoPos
	}()

	// CLI flags already set
	logo = "cli-logo.png"
	logoPos = "bottom-left"

	loadedConfig = &config.Config{
		Logo:         "config-logo.png",
		LogoPosition: "top-right",
	}
	opts := mdpress.DefaultOptions()
	applyConfigToOptions(&opts)

	// CLI flags should NOT be overridden
	if logo != "cli-logo.png" {
		t.Errorf("CLI logo flag should take precedence, got %q", logo)
	}
	if logoPos != "bottom-left" {
		t.Errorf("CLI logoPos flag should take precedence, got %q", logoPos)
	}
}

// --- applyFlagOverrides tests ---

func TestApplyFlagOverrides_FontChanged(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("font", "", "")
	cmd.Flags().Set("font", "Times")

	oldFont := font
	defer func() { font = oldFont }()
	font = "Times"

	opts := mdpress.DefaultOptions()
	applyFlagOverrides(cmd, &opts)

	if opts.DefaultFamily != "Times" {
		t.Errorf("expected font Times, got %s", opts.DefaultFamily)
	}
}

func TestApplyFlagOverrides_MarginsChanged(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Float64("margin-top", 0, "")
	cmd.Flags().Float64("margin-right", 0, "")
	cmd.Flags().Float64("margin-bottom", 0, "")
	cmd.Flags().Float64("margin-left", 0, "")
	cmd.Flags().Set("margin-top", "50")
	cmd.Flags().Set("margin-right", "15")
	cmd.Flags().Set("margin-bottom", "40")
	cmd.Flags().Set("margin-left", "10")

	oldTop, oldRight, oldBot, oldLeft := marginTop, marginRight, marginBot, marginLeft
	defer func() { marginTop, marginRight, marginBot, marginLeft = oldTop, oldRight, oldBot, oldLeft }()
	marginTop, marginRight, marginBot, marginLeft = 50, 15, 40, 10

	opts := mdpress.DefaultOptions()
	applyFlagOverrides(cmd, &opts)

	if opts.MarginTop != 50 {
		t.Errorf("expected margin top 50, got %f", opts.MarginTop)
	}
	if opts.MarginRight != 15 {
		t.Errorf("expected margin right 15, got %f", opts.MarginRight)
	}
	if opts.MarginBottom != 40 {
		t.Errorf("expected margin bottom 40, got %f", opts.MarginBottom)
	}
	if opts.MarginLeft != 10 {
		t.Errorf("expected margin left 10, got %f", opts.MarginLeft)
	}
}

func TestApplyFlagOverrides_UnchangedFlagsIgnored(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("font", "", "")
	cmd.Flags().Float64("margin-top", 0, "")
	cmd.Flags().Float64("margin-right", 0, "")
	cmd.Flags().Float64("margin-bottom", 0, "")
	cmd.Flags().Float64("margin-left", 0, "")
	// Don't set any flags

	opts := mdpress.DefaultOptions()
	defaultFont := opts.DefaultFamily
	defaultTop := opts.MarginTop
	applyFlagOverrides(cmd, &opts)

	if opts.DefaultFamily != defaultFont {
		t.Error("unchanged font flag should not affect options")
	}
	if opts.MarginTop != defaultTop {
		t.Error("unchanged margin flag should not affect options")
	}
}

// --- render command error cases ---

func TestRender_MissingArgument(t *testing.T) {
	// The render command requires exactly 1 arg
	cmd := &cobra.Command{Use: "render [file]", Args: cobra.ExactArgs(1)}
	cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no file argument provided")
	}
}

func TestRender_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Set up minimal config to avoid config loading errors
	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()
	loadedConfig = &config.Config{Version: "1"}

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{
		Use:  "render [file]",
		Args: cobra.ExactArgs(1),
		RunE: renderCmd.RunE,
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "")
	cmd.Flags().StringVarP(&templateFlag, "template", "t", "", "")
	cmd.Flags().StringVarP(&dataFile, "data", "d", "", "")
	cmd.Flags().StringVar(&font, "font", "", "")
	cmd.Flags().StringVar(&accentColor, "accent-color", "", "")
	cmd.Flags().StringVar(&logoPos, "logo-position", "", "")
	cmd.Flags().StringVar(&logo, "logo", "", "")
	cmd.Flags().Float64Var(&logoWidth, "logo-width", 0, "")
	cmd.Flags().Float64Var(&marginTop, "margin-top", 0, "")
	cmd.Flags().Float64Var(&marginRight, "margin-right", 0, "")
	cmd.Flags().Float64Var(&marginBot, "margin-bottom", 0, "")
	cmd.Flags().Float64Var(&marginLeft, "margin-left", 0, "")

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"/nonexistent/file.md"})

	// Reset dataFile to avoid pro auth check
	dataFile = ""

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestRender_BasicMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create a simple markdown file
	mdPath := filepath.Join(tmpDir, "test.md")
	os.WriteFile(mdPath, []byte("# Hello\n\nThis is a test.\n"), 0644)

	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()
	loadedConfig = &config.Config{Version: "1"}

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{
		Use:  "render [file]",
		Args: cobra.ExactArgs(1),
		RunE: renderCmd.RunE,
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "")
	cmd.Flags().StringVarP(&templateFlag, "template", "t", "", "")
	cmd.Flags().StringVarP(&dataFile, "data", "d", "", "")
	cmd.Flags().StringVar(&font, "font", "", "")
	cmd.Flags().StringVar(&accentColor, "accent-color", "", "")
	cmd.Flags().StringVar(&logoPos, "logo-position", "", "")
	cmd.Flags().StringVar(&logo, "logo", "", "")
	cmd.Flags().Float64Var(&logoWidth, "logo-width", 0, "")
	cmd.Flags().Float64Var(&marginTop, "margin-top", 0, "")
	cmd.Flags().Float64Var(&marginRight, "margin-right", 0, "")
	cmd.Flags().Float64Var(&marginBot, "margin-bottom", 0, "")
	cmd.Flags().Float64Var(&marginLeft, "margin-left", 0, "")

	outPath := filepath.Join(tmpDir, "output.pdf")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{mdPath})

	// Reset global state
	dataFile = ""
	output = outPath
	templateFlag = ""

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify PDF was created
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output PDF should not be empty")
	}

	// Verify output message
	if !strings.Contains(buf.String(), "Written to") {
		t.Errorf("expected 'Written to' message, got %q", buf.String())
	}
}

func TestRender_DefaultOutputName(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	mdPath := filepath.Join(tmpDir, "document.md")
	os.WriteFile(mdPath, []byte("# Test\n"), 0644)

	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()
	loadedConfig = &config.Config{Version: "1"}

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{
		Use:  "render [file]",
		Args: cobra.ExactArgs(1),
		RunE: renderCmd.RunE,
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "")
	cmd.Flags().StringVarP(&templateFlag, "template", "t", "", "")
	cmd.Flags().StringVarP(&dataFile, "data", "d", "", "")
	cmd.Flags().StringVar(&font, "font", "", "")
	cmd.Flags().StringVar(&accentColor, "accent-color", "", "")
	cmd.Flags().StringVar(&logoPos, "logo-position", "", "")
	cmd.Flags().StringVar(&logo, "logo", "", "")
	cmd.Flags().Float64Var(&logoWidth, "logo-width", 0, "")
	cmd.Flags().Float64Var(&marginTop, "margin-top", 0, "")
	cmd.Flags().Float64Var(&marginRight, "margin-right", 0, "")
	cmd.Flags().Float64Var(&marginBot, "margin-bottom", 0, "")
	cmd.Flags().Float64Var(&marginLeft, "margin-left", 0, "")

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{mdPath})

	// Reset global state - empty output means auto-derive from input
	dataFile = ""
	output = ""
	templateFlag = ""

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should create document.pdf
	expectedPath := filepath.Join(tmpDir, "document.pdf")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected auto-generated output at %s", expectedPath)
	}
}

func TestRender_CustomFont(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	mdPath := filepath.Join(tmpDir, "test.md")
	os.WriteFile(mdPath, []byte("# Test\n"), 0644)

	oldConfig := loadedConfig
	defer func() { loadedConfig = oldConfig }()
	loadedConfig = &config.Config{Version: "1"}

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{
		Use:  "render [file]",
		Args: cobra.ExactArgs(1),
		RunE: renderCmd.RunE,
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "")
	cmd.Flags().StringVarP(&templateFlag, "template", "t", "", "")
	cmd.Flags().StringVarP(&dataFile, "data", "d", "", "")
	cmd.Flags().StringVar(&font, "font", "", "")
	cmd.Flags().StringVar(&accentColor, "accent-color", "", "")
	cmd.Flags().StringVar(&logoPos, "logo-position", "", "")
	cmd.Flags().StringVar(&logo, "logo", "", "")
	cmd.Flags().Float64Var(&logoWidth, "logo-width", 0, "")
	cmd.Flags().Float64Var(&marginTop, "margin-top", 0, "")
	cmd.Flags().Float64Var(&marginRight, "margin-right", 0, "")
	cmd.Flags().Float64Var(&marginBot, "margin-bottom", 0, "")
	cmd.Flags().Float64Var(&marginLeft, "margin-left", 0, "")

	outPath := filepath.Join(tmpDir, "out.pdf")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{mdPath, "--font", "Courier"})

	dataFile = ""
	output = outPath
	templateFlag = ""

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("render with custom font failed: %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Error("output file should exist")
	}
}
