package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain builds the binary once before all E2E tests.
var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary to a temp location
	tmp, err := os.MkdirTemp("", "mdpress-e2e-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	binaryPath = filepath.Join(tmp, "mdpress")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	os.Exit(m.Run())
}

// run executes the mdpress binary with the given args in an isolated HOME.
func run(t *testing.T, home string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(),
		"HOME="+home,
	)

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = -1
	}

	return outBuf.String(), errBuf.String(), exitCode
}

// --- Version ---

func TestE2E_Version(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "version")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.HasPrefix(stdout, "mdpress ") {
		t.Errorf("expected output starting with 'mdpress ', got %q", stdout)
	}
	if !strings.Contains(stdout, "commit:") {
		t.Error("version output should contain 'commit:'")
	}
	if !strings.Contains(stdout, "built:") {
		t.Error("version output should contain 'built:'")
	}
}

// --- Help ---

func TestE2E_Help(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "--help")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "Markdown") || !strings.Contains(stdout, "PDF") {
		t.Error("help should contain description mentioning Markdown and PDF")
	}
	if !strings.Contains(stdout, "render") {
		t.Error("help should list render command")
	}
	if !strings.Contains(stdout, "init") {
		t.Error("help should list init command")
	}
	if !strings.Contains(stdout, "auth") {
		t.Error("help should list auth command")
	}
	if !strings.Contains(stdout, "templates") {
		t.Error("help should list templates command")
	}
}

func TestE2E_RenderHelp(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "render", "--help")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "--output") {
		t.Error("render help should show --output flag")
	}
	if !strings.Contains(stdout, "--template") {
		t.Error("render help should show --template flag")
	}
	if !strings.Contains(stdout, "--data") {
		t.Error("render help should show --data flag")
	}
	if !strings.Contains(stdout, "--font") {
		t.Error("render help should show --font flag")
	}
	if !strings.Contains(stdout, "--margin-top") {
		t.Error("render help should show --margin-top flag")
	}
}

// --- Init ---

func TestE2E_InitGlobal(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "init")

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "Created") {
		t.Error("init should output 'Created' messages")
	}

	// Verify files exist
	configPath := filepath.Join(home, ".config", "mdpress", "mdpress.yml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("global config should exist at %s", configPath)
	}

	logosDir := filepath.Join(home, ".config", "mdpress", "logos")
	if info, err := os.Stat(logosDir); err != nil || !info.IsDir() {
		t.Error("logos directory should exist")
	}

	templatePath := filepath.Join(home, ".config", "mdpress", "templates", "default.yml")
	if _, err := os.Stat(templatePath); err != nil {
		t.Error("default template should exist")
	}
}

func TestE2E_InitIdempotent(t *testing.T) {
	home := t.TempDir()

	// First run
	_, _, code := run(t, home, "init")
	if code != 0 {
		t.Fatal("first init failed")
	}

	// Second run
	stdout, _, code := run(t, home, "init")
	if code != 0 {
		t.Fatal("second init failed")
	}
	if !strings.Contains(stdout, "already initialized") {
		t.Errorf("second init should say 'already initialized', got %q", stdout)
	}
}

func TestE2E_InitProject(t *testing.T) {
	home := t.TempDir()
	projectDir := t.TempDir()

	cmd := exec.Command(binaryPath, "init", "--project")
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "HOME="+home)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("init --project failed: %v\n%s", err, out)
	}

	configPath := filepath.Join(projectDir, "mdpress.yml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal("project config not created")
	}
	if !strings.Contains(string(content), "mdpress project configuration") {
		t.Error("project config should contain project header")
	}
}

// --- Render ---

func TestE2E_RenderBasic(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create markdown file
	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Hello World\n\nThis is a test document.\n"), 0644)

	stdout, _, code := run(t, home, "render", mdPath)
	if code != 0 {
		t.Fatalf("render failed with exit code %d", code)
	}

	expectedPDF := filepath.Join(workDir, "test.pdf")
	if !strings.Contains(stdout, "Written to") {
		t.Errorf("expected 'Written to' message, got %q", stdout)
	}
	if _, err := os.Stat(expectedPDF); err != nil {
		t.Errorf("PDF should be created at %s", expectedPDF)
	}
}

func TestE2E_RenderCustomOutput(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "input.md")
	os.WriteFile(mdPath, []byte("# Test\n"), 0644)

	outPath := filepath.Join(workDir, "custom-output.pdf")
	stdout, _, code := run(t, home, "render", mdPath, "-o", outPath)
	if code != 0 {
		t.Fatalf("render failed with exit code %d", code)
	}

	if !strings.Contains(stdout, "custom-output.pdf") {
		t.Errorf("output should mention custom filename, got %q", stdout)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Error("custom output PDF should exist")
	}
}

func TestE2E_RenderWithFont(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Font Test\n"), 0644)

	outPath := filepath.Join(workDir, "out.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath, "--font", "Courier")
	if code != 0 {
		t.Fatal("render with --font should succeed")
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatal("output PDF should exist")
	}
	if info.Size() == 0 {
		t.Error("PDF should not be empty")
	}
}

func TestE2E_RenderWithMargins(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Margin Test\n"), 0644)

	outPath := filepath.Join(workDir, "out.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath,
		"--margin-top", "30", "--margin-right", "25",
		"--margin-bottom", "30", "--margin-left", "25")
	if code != 0 {
		t.Fatal("render with margins should succeed")
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output PDF should exist")
	}
}

func TestE2E_RenderNoArgs(t *testing.T) {
	home := t.TempDir()
	_, stderr, code := run(t, home, "render")

	if code == 0 {
		t.Fatal("render with no args should fail")
	}
	if !strings.Contains(stderr, "accepts 1 arg") {
		t.Errorf("error should mention arg count, got stderr: %q", stderr)
	}
}

func TestE2E_RenderNonexistentFile(t *testing.T) {
	home := t.TempDir()
	_, _, code := run(t, home, "render", "/nonexistent/file.md")

	if code == 0 {
		t.Fatal("render of nonexistent file should fail")
	}
}

func TestE2E_RenderWithTemplate(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create a template
	templateDir := filepath.Join(home, ".config", "mdpress", "templates")
	os.MkdirAll(templateDir, 0755)
	tmplContent := `name: "test-tmpl"
version: "1"
font: "Times"
`
	os.WriteFile(filepath.Join(templateDir, "test-tmpl.yml"), []byte(tmplContent), 0644)

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Template Test\n"), 0644)

	outPath := filepath.Join(workDir, "out.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath, "-t", "test-tmpl")
	if code != 0 {
		t.Fatal("render with template should succeed")
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output PDF should exist")
	}
}

func TestE2E_RenderWithConfig(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create global config with custom font
	configDir := filepath.Join(home, ".config", "mdpress")
	os.MkdirAll(configDir, 0755)
	configContent := `version: "1"
font: "Times"
margins:
  top: 30
  right: 25
  bottom: 30
  left: 25
`
	os.WriteFile(filepath.Join(configDir, "mdpress.yml"), []byte(configContent), 0644)

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Config Test\n"), 0644)

	outPath := filepath.Join(workDir, "out.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath)
	if code != 0 {
		t.Fatal("render with config should succeed")
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output PDF should exist")
	}
}

// --- Templates ---

func TestE2E_TemplatesEmpty(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	cmd := exec.Command(binaryPath, "templates")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "HOME="+home)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("templates command failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "No templates found") {
		t.Errorf("expected 'No templates found', got %q", string(out))
	}
}

func TestE2E_TemplatesListing(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create a global template
	templateDir := filepath.Join(home, ".config", "mdpress", "templates")
	os.MkdirAll(templateDir, 0755)
	tmplContent := `name: "invoice"
version: "1"
description: "Invoice document preset"
font: "Times"
`
	os.WriteFile(filepath.Join(templateDir, "invoice.yml"), []byte(tmplContent), 0644)

	cmd := exec.Command(binaryPath, "templates")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "HOME="+home)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("templates command failed: %v\n%s", err, out)
	}

	outStr := string(out)
	if !strings.Contains(outStr, "Available templates") {
		t.Error("should show 'Available templates' header")
	}
	if !strings.Contains(outStr, "invoice") {
		t.Error("should list 'invoice' template")
	}
	if !strings.Contains(outStr, "Invoice document preset") {
		t.Error("should show template description")
	}
}

// --- Auth ---

func TestE2E_AuthWhoami_NotLoggedIn(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "auth", "whoami")

	if code != 0 {
		t.Fatalf("auth whoami should succeed, got exit %d", code)
	}
	if !strings.Contains(stdout, "Not logged in") {
		t.Errorf("expected 'Not logged in', got %q", stdout)
	}
}

func TestE2E_AuthLogout_Clean(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "auth", "logout")

	if code != 0 {
		t.Fatalf("auth logout should succeed even with no auth file, got exit %d", code)
	}
	if !strings.Contains(stdout, "Logged out") {
		t.Errorf("expected logout message, got %q", stdout)
	}
}

func TestE2E_AuthHelp(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "auth", "--help")

	if code != 0 {
		t.Fatal("auth --help should succeed")
	}
	if !strings.Contains(stdout, "login") {
		t.Error("auth help should mention login")
	}
	if !strings.Contains(stdout, "logout") {
		t.Error("auth help should mention logout")
	}
	if !strings.Contains(stdout, "whoami") {
		t.Error("auth help should mention whoami")
	}
}

// --- Config Debug ---

func TestE2E_ConfigDebug(t *testing.T) {
	home := t.TempDir()
	_, stderr, code := run(t, home, "--config-debug", "version")

	if code != 0 {
		t.Fatalf("config-debug should succeed, got exit %d", code)
	}
	if !strings.Contains(stderr, "Config Debug") {
		t.Error("should output config debug header")
	}
	if !strings.Contains(stderr, "Config files:") {
		t.Error("should show config files section")
	}
	if !strings.Contains(stderr, "Resolved config values:") {
		t.Error("should show resolved values section")
	}
}

// --- Unknown command ---

func TestE2E_UnknownCommand(t *testing.T) {
	home := t.TempDir()
	_, _, code := run(t, home, "nonexistent-command")

	if code == 0 {
		t.Fatal("unknown command should fail")
	}
}

// --- Templates List ---

func TestE2E_TemplatesListLocalOnly(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create a global template
	templateDir := filepath.Join(home, ".config", "mdpress", "templates")
	os.MkdirAll(templateDir, 0755)
	os.WriteFile(filepath.Join(templateDir, "report.yml"), []byte("name: \"report\"\nversion: \"1\"\ndescription: \"Report preset\"\n"), 0644)

	cmd := exec.Command(binaryPath, "templates", "list", "--local-only")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "HOME="+home)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("templates list --local-only failed: %v\n%s", err, out)
	}

	outStr := string(out)
	if !strings.Contains(outStr, "report") {
		t.Error("should list 'report' template")
	}
	if !strings.Contains(outStr, "global") {
		t.Error("should show 'global' source")
	}
}

func TestE2E_TemplatesListProjectAndGlobal(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create global template
	globalDir := filepath.Join(home, ".config", "mdpress", "templates")
	os.MkdirAll(globalDir, 0755)
	os.WriteFile(filepath.Join(globalDir, "global-tmpl.yml"), []byte("name: \"global-tmpl\"\nversion: \"1\"\ndescription: \"Global\"\n"), 0644)

	// Create project template
	projectDir := filepath.Join(workDir, "templates")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "project-tmpl.yml"), []byte("name: \"project-tmpl\"\nversion: \"1\"\ndescription: \"Project\"\n"), 0644)

	cmd := exec.Command(binaryPath, "templates", "list", "--local-only")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "HOME="+home)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("templates list failed: %v\n%s", err, out)
	}

	outStr := string(out)
	if !strings.Contains(outStr, "global-tmpl") {
		t.Error("should list global template")
	}
	if !strings.Contains(outStr, "project-tmpl") {
		t.Error("should list project template")
	}
}

func TestE2E_TemplatesListSubcommand(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// 'templates' with no subcommand should default to 'list'
	// Both should produce the same output
	cmd1 := exec.Command(binaryPath, "templates")
	cmd1.Dir = workDir
	cmd1.Env = append(os.Environ(), "HOME="+home)
	out1, _ := cmd1.CombinedOutput()

	cmd2 := exec.Command(binaryPath, "templates", "list", "--local-only")
	cmd2.Dir = workDir
	cmd2.Env = append(os.Environ(), "HOME="+home)
	out2, _ := cmd2.CombinedOutput()

	// Both should mention "No templates found" since there are none
	if !strings.Contains(string(out1), "No templates found") {
		t.Error("'templates' should show no templates found")
	}
	if !strings.Contains(string(out2), "No templates found") {
		t.Error("'templates list --local-only' should show no templates found")
	}
}

// --- Templates Pull ---

func TestE2E_TemplatesPullNoArgs(t *testing.T) {
	home := t.TempDir()
	_, stderr, code := run(t, home, "templates", "pull")

	if code == 0 {
		t.Fatal("templates pull with no args should fail")
	}
	if !strings.Contains(stderr, "accepts 1 arg") {
		t.Errorf("error should mention arg count, got stderr: %q", stderr)
	}
}

func TestE2E_TemplatesPullHelp(t *testing.T) {
	home := t.TempDir()
	stdout, _, code := run(t, home, "templates", "pull", "--help")

	if code != 0 {
		t.Fatal("templates pull --help should succeed")
	}
	if !strings.Contains(stdout, "--force") {
		t.Error("pull help should show --force flag")
	}
}

// --- Init edge cases ---

func TestE2E_InitThenRender(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Init first
	_, _, code := run(t, home, "init")
	if code != 0 {
		t.Fatal("init should succeed")
	}

	// Then render - config should be loaded from init
	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# After Init\n\nRendered with default config.\n"), 0644)

	outPath := filepath.Join(workDir, "test.pdf")
	stdout, _, code := run(t, home, "render", mdPath, "-o", outPath)
	if code != 0 {
		t.Fatal("render after init should succeed")
	}
	if !strings.Contains(stdout, "Written to") {
		t.Error("should show success message")
	}

	// Verify PDF is valid
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal("output should exist")
	}
	if string(data[:5]) != "%PDF-" {
		t.Error("output should be valid PDF")
	}
}

func TestE2E_InitProjectThenRender(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Init project config with custom settings
	cmd := exec.Command(binaryPath, "init", "--project")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "HOME="+home)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("init --project failed: %v\n%s", err, out)
	}

	// Write a project config with custom font
	configPath := filepath.Join(workDir, "mdpress.yml")
	os.WriteFile(configPath, []byte("version: \"1\"\nfont: \"Courier\"\n"), 0644)

	// Render from that directory
	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Project Test\n"), 0644)

	outPath := filepath.Join(workDir, "test.pdf")
	renderCmd := exec.Command(binaryPath, "render", mdPath, "-o", outPath)
	renderCmd.Dir = workDir
	renderCmd.Env = append(os.Environ(), "HOME="+home)

	if out, err := renderCmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output PDF should exist")
	}
}

// --- Render edge cases ---

func TestE2E_RenderEmptyFile(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "empty.md")
	os.WriteFile(mdPath, []byte(""), 0644)

	outPath := filepath.Join(workDir, "empty.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath)
	if code != 0 {
		t.Fatal("render of empty file should succeed")
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal("output should exist")
	}
	if string(data[:5]) != "%PDF-" {
		t.Error("output should be valid PDF even for empty input")
	}
}

func TestE2E_RenderAllFonts(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Font Test\n\nBody text.\n"), 0644)

	for _, fontName := range []string{"Helvetica", "Times", "Courier"} {
		t.Run(fontName, func(t *testing.T) {
			outPath := filepath.Join(workDir, fontName+".pdf")
			_, _, code := run(t, home, "render", mdPath, "-o", outPath, "--font", fontName)
			if code != 0 {
				t.Fatalf("render with --%s should succeed", fontName)
			}
			info, err := os.Stat(outPath)
			if err != nil {
				t.Fatal("output should exist")
			}
			if info.Size() == 0 {
				t.Error("PDF should not be empty")
			}
		})
	}
}

func TestE2E_RenderWithAccentColor(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Accent Color Test\n\nBody text.\n"), 0644)

	outPath := filepath.Join(workDir, "accent.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath, "--accent-color", "#FF5733")
	if code != 0 {
		t.Fatal("render with --accent-color should succeed")
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output should exist")
	}
}

func TestE2E_RenderWithAllFlags(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# All Flags Test\n\nBody.\n"), 0644)

	outPath := filepath.Join(workDir, "all-flags.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath,
		"--font", "Times",
		"--accent-color", "#336699",
		"--margin-top", "30",
		"--margin-right", "25",
		"--margin-bottom", "30",
		"--margin-left", "25",
	)
	if code != 0 {
		t.Fatal("render with all flags should succeed")
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatal("output should exist")
	}
	if info.Size() < 500 {
		t.Error("PDF seems too small")
	}
}

func TestE2E_RenderFlagOverridesConfig(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Create config with Helvetica
	configDir := filepath.Join(home, ".config", "mdpress")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "mdpress.yml"), []byte("version: \"1\"\nfont: \"Helvetica\"\n"), 0644)

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Override Test\n"), 0644)

	// Render with --font Courier (should override config's Helvetica)
	outPath := filepath.Join(workDir, "override.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath, "--font", "Courier")
	if code != 0 {
		t.Fatal("render with flag override should succeed")
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output should exist")
	}
}

func TestE2E_RenderMultipleFiles(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Render two separate files in sequence
	for _, name := range []string{"doc1", "doc2"} {
		mdPath := filepath.Join(workDir, name+".md")
		os.WriteFile(mdPath, []byte("# "+name+"\n\nContent for "+name+".\n"), 0644)

		outPath := filepath.Join(workDir, name+".pdf")
		_, _, code := run(t, home, "render", mdPath, "-o", outPath)
		if code != 0 {
			t.Fatalf("render %s should succeed", name)
		}

		if _, err := os.Stat(outPath); err != nil {
			t.Fatalf("%s.pdf should exist", name)
		}
	}
}

func TestE2E_RenderAutoOutputName(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "my-document.md")
	os.WriteFile(mdPath, []byte("# Auto Name\n"), 0644)

	stdout, _, code := run(t, home, "render", mdPath)
	if code != 0 {
		t.Fatal("render should succeed")
	}

	// Should auto-derive output as my-document.pdf
	expectedPDF := filepath.Join(workDir, "my-document.pdf")
	if _, err := os.Stat(expectedPDF); err != nil {
		t.Error("should create my-document.pdf automatically")
	}
	if !strings.Contains(stdout, "my-document.pdf") {
		t.Error("output message should mention the auto-derived filename")
	}
}

func TestE2E_RenderNonMarkdownExtension(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	// Should work with any text file, not just .md
	txtPath := filepath.Join(workDir, "notes.txt")
	os.WriteFile(txtPath, []byte("# Notes\n\nSome notes here.\n"), 0644)

	outPath := filepath.Join(workDir, "notes.pdf")
	_, _, code := run(t, home, "render", txtPath, "-o", outPath)
	if code != 0 {
		t.Fatal("render should work with non-.md files")
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("output should exist")
	}
}

// --- Auth edge cases ---

func TestE2E_AuthWhoamiWithCredentials(t *testing.T) {
	home := t.TempDir()

	// Create auth file
	authDir := filepath.Join(home, ".config", "mdpress")
	os.MkdirAll(authDir, 0755)
	os.WriteFile(filepath.Join(authDir, "auth.yml"), []byte("api_key: mdp_testkey1234567890\nemail: test@example.com\n"), 0600)

	stdout, _, code := run(t, home, "auth", "whoami")
	if code != 0 {
		t.Fatal("whoami should succeed")
	}

	if !strings.Contains(stdout, "test@example.com") {
		t.Error("should show email")
	}
	if !strings.Contains(stdout, "API Key:") {
		t.Error("should show masked API key")
	}
	if strings.Contains(stdout, "mdp_testkey1234567890") {
		t.Error("should NOT show full API key")
	}
}

func TestE2E_AuthLogoutThenWhoami(t *testing.T) {
	home := t.TempDir()

	// Create auth file
	authDir := filepath.Join(home, ".config", "mdpress")
	os.MkdirAll(authDir, 0755)
	os.WriteFile(filepath.Join(authDir, "auth.yml"), []byte("api_key: mdp_test\n"), 0600)

	// Logout
	stdout, _, code := run(t, home, "auth", "logout")
	if code != 0 {
		t.Fatal("logout should succeed")
	}
	if !strings.Contains(stdout, "Logged out") {
		t.Error("should confirm logout")
	}

	// Whoami should show not logged in
	stdout, _, code = run(t, home, "auth", "whoami")
	if code != 0 {
		t.Fatal("whoami should succeed")
	}
	if !strings.Contains(stdout, "Not logged in") {
		t.Error("should show not logged in after logout")
	}
}

// --- Render data flag without auth ---

func TestE2E_RenderDataRequiresAuth(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# Hello {{name}}\n"), 0644)

	dataPath := filepath.Join(workDir, "data.json")
	os.WriteFile(dataPath, []byte(`{"name": "World"}`), 0644)

	_, stderr, code := run(t, home, "render", mdPath, "-d", dataPath)
	if code == 0 {
		t.Fatal("render with --data should fail without auth")
	}
	if !strings.Contains(stderr, "Pro") || !strings.Contains(stderr, "data") {
		t.Errorf("error should mention Pro and data feature, got %q", stderr)
	}
}

// --- PDF output validation ---

func TestE2E_RenderOutputIsPDF(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	mdPath := filepath.Join(workDir, "test.md")
	os.WriteFile(mdPath, []byte("# PDF Check\n\nBody text here.\n"), 0644)

	outPath := filepath.Join(workDir, "test.pdf")
	_, _, code := run(t, home, "render", mdPath, "-o", outPath)
	if code != 0 {
		t.Fatal("render should succeed")
	}

	// Read first bytes to verify PDF magic number
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Error("output file should start with PDF magic number (%PDF-)")
	}
}

// --- Render complex markdown ---

func TestE2E_RenderComplexMarkdown(t *testing.T) {
	home := t.TempDir()
	workDir := t.TempDir()

	content := `# Main Heading

## Section 1

This is a paragraph with **bold** and *italic* text.

- Item 1
- Item 2
- Item 3

## Section 2

> A blockquote goes here.

### Code Example

` + "```\nfmt.Println(\"hello\")\n```" + `

## Section 3

| Column A | Column B |
|----------|----------|
| Cell 1   | Cell 2   |
| Cell 3   | Cell 4   |

---

Final paragraph.
`

	mdPath := filepath.Join(workDir, "complex.md")
	os.WriteFile(mdPath, []byte(content), 0644)

	outPath := filepath.Join(workDir, "complex.pdf")
	stdout, _, code := run(t, home, "render", mdPath, "-o", outPath)
	if code != 0 {
		t.Fatal("render of complex markdown should succeed")
	}

	if !strings.Contains(stdout, "Written to") {
		t.Error("should output success message")
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatal("output PDF should exist")
	}
	// Complex document should be reasonably sized
	if info.Size() < 1000 {
		t.Errorf("complex PDF seems too small: %d bytes", info.Size())
	}
}
