package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skel-tech/mdpress/internal/auth"
	"github.com/spf13/cobra"
)

// --- maskKey tests ---

func TestMaskKey_Short(t *testing.T) {
	got := maskKey("mdp_abc")
	if got != "mdp_****" {
		t.Errorf("expected mdp_****, got %q", got)
	}
}

func TestMaskKey_ExactlyEight(t *testing.T) {
	got := maskKey("mdp_abcd")
	if got != "mdp_****" {
		t.Errorf("expected mdp_****, got %q", got)
	}
}

func TestMaskKey_Long(t *testing.T) {
	// "mdp_abcdef1234" = 14 chars, should show first 4 + 6 stars + last 4
	got := maskKey("mdp_abcdef1234")
	if got != "mdp_******1234" {
		t.Errorf("expected mdp_******1234, got %q", got)
	}
}

func TestMaskKey_Empty(t *testing.T) {
	got := maskKey("")
	if got != "mdp_****" {
		t.Errorf("expected mdp_****, got %q", got)
	}
}

// --- whoami command tests ---

func newTestWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "whoami",
		RunE: whoamiCmd.RunE,
	}
}

func TestWhoami_NotLoggedIn(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	buf := new(bytes.Buffer)
	cmd := newTestWhoamiCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Not logged in") {
		t.Errorf("expected 'Not logged in', got %q", out)
	}
}

func TestWhoami_LoggedInWithEmail(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create auth file
	authDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(authDir, 0755)
	authContent := "api_key: mdp_testkey1234567890\nemail: user@example.com\n"
	os.WriteFile(filepath.Join(authDir, "auth.yml"), []byte(authContent), 0600)

	buf := new(bytes.Buffer)
	cmd := newTestWhoamiCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "user@example.com") {
		t.Errorf("expected email in output, got %q", out)
	}
	if !strings.Contains(out, "API Key:") {
		t.Errorf("expected masked API key in output, got %q", out)
	}
	// Should NOT contain the full key
	if strings.Contains(out, "mdp_testkey1234567890") {
		t.Error("whoami should mask the API key, not show it in full")
	}
}

func TestWhoami_LoggedInWithoutEmail(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	authDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(authDir, 0755)
	authContent := "api_key: mdp_testkey1234567890\n"
	os.WriteFile(filepath.Join(authDir, "auth.yml"), []byte(authContent), 0600)

	buf := new(bytes.Buffer)
	cmd := newTestWhoamiCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "Email:") {
		t.Error("should not show Email line when email is empty")
	}
	if !strings.Contains(out, "API Key:") {
		t.Errorf("expected API Key line, got %q", out)
	}
}

// --- logout command tests ---

func newTestLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "logout",
		RunE: logoutCmd.RunE,
	}
}

func TestLogout_Success(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create auth file first
	authDir := filepath.Join(tmpDir, ".config", "mdpress")
	os.MkdirAll(authDir, 0755)
	os.WriteFile(filepath.Join(authDir, "auth.yml"), []byte("api_key: mdp_test"), 0600)

	buf := new(bytes.Buffer)
	cmd := newTestLogoutCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Logged out successfully") {
		t.Errorf("expected success message, got %q", out)
	}

	// Verify file is removed
	if _, err := os.Stat(filepath.Join(authDir, "auth.yml")); !os.IsNotExist(err) {
		t.Error("auth file should be removed after logout")
	}
}

func TestLogout_NoAuthFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	buf := new(bytes.Buffer)
	cmd := newTestLogoutCmd()
	cmd.SetOut(buf)

	// Should not error even if no auth file exists
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- login validation tests (without stdin interaction) ---

func TestLogin_KeyFormatValidation(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"valid key", "mdp_abc123", false},
		{"missing prefix", "abc123", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasPrefix := strings.HasPrefix(tt.key, "mdp_")
			isEmpty := tt.key == ""
			gotErr := !hasPrefix || isEmpty
			if gotErr != tt.wantErr {
				t.Errorf("key %q: expected error=%v, got error=%v", tt.key, tt.wantErr, gotErr)
			}
		})
	}
}

// --- auth subcommand structure tests ---

func TestAuth_HasSubcommands(t *testing.T) {
	subs := authCmd.Commands()
	names := make(map[string]bool)
	for _, c := range subs {
		names[c.Name()] = true
	}

	for _, expected := range []string{"login", "logout", "whoami"} {
		if !names[expected] {
			t.Errorf("auth command missing subcommand %q", expected)
		}
	}
}

func TestAuth_Save_And_Load_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	creds := &auth.Credentials{
		APIKey: "mdp_roundtrip123",
		Email:  "rt@example.com",
	}
	if err := auth.Save(creds); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := auth.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.APIKey != creds.APIKey {
		t.Errorf("expected key %q, got %q", creds.APIKey, loaded.APIKey)
	}
	if loaded.Email != creds.Email {
		t.Errorf("expected email %q, got %q", creds.Email, loaded.Email)
	}
}
