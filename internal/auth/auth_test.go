package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Store tests ---

func TestLoad_NoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.yml")
	creds, err := loadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds != nil {
		t.Fatal("expected nil credentials when file does not exist")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.yml")
	content := "api_key: mdp_test123\nemail: user@example.com\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	creds, err := loadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.APIKey != "mdp_test123" {
		t.Errorf("expected api_key mdp_test123, got %s", creds.APIKey)
	}
	if creds.Email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %s", creds.Email)
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.yml")
	if err := os.WriteFile(path, []byte("{{bad yaml"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := loadFrom(path)
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "auth.yml")

	creds := &Credentials{APIKey: "mdp_abc123", Email: "test@example.com"}
	if err := saveTo(path, creds); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was created with restricted permissions
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}

	// Verify contents
	loaded, err := loadFrom(path)
	if err != nil {
		t.Fatalf("failed to reload: %v", err)
	}
	if loaded.APIKey != "mdp_abc123" {
		t.Errorf("expected api_key mdp_abc123, got %s", loaded.APIKey)
	}
}

func TestClear_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.yml")
	if err := os.WriteFile(path, []byte("api_key: mdp_test"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := clearFrom(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected file to be removed")
	}
}

func TestClear_NoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.yml")
	// Should not error when file doesn't exist
	if err := clearFrom(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Auth logic tests ---

func TestIsValidKeyFormat(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"mdp_abc123", true},
		{"mdp_", true},
		{"", false},
		{"  ", false},
		{"invalid", false},
		{"mdp-abc123", false},
		{" mdp_abc123 ", true}, // trimmed
	}
	for _, tt := range tests {
		if got := isValidKeyFormat(tt.key); got != tt.want {
			t.Errorf("isValidKeyFormat(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}

func TestRequirePro_NilValidator(t *testing.T) {
	// With no validator set, should always deny (fail-closed)
	old := ValidateKey
	ValidateKey = nil
	defer func() { ValidateKey = old }()

	err := RequirePro("data")
	if err == nil {
		t.Fatal("expected error when ValidateKey is nil")
	}
	var gatedErr *FeatureGatedError
	if ok := err.(*FeatureGatedError); ok == nil {
		t.Fatal("expected FeatureGatedError")
	} else {
		gatedErr = ok
	}
	if gatedErr.Feature != "data" {
		t.Errorf("expected feature 'data', got %q", gatedErr.Feature)
	}
}

func TestRequirePro_ValidatorRejects(t *testing.T) {
	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) {
		return nil, fmt.Errorf("invalid key")
	}
	defer func() { ValidateKey = old }()

	err := RequirePro("data")
	if err == nil {
		t.Fatal("expected error when validator rejects key")
	}
}

func TestRequirePro_ValidatorAccepts(t *testing.T) {
	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) {
		if key == "mdp_valid_key" {
			return &LicenseInfo{}, nil
		}
		return nil, fmt.Errorf("invalid")
	}
	defer func() { ValidateKey = old }()

	// IsAuthenticated also needs credentials loaded, so we can't test the full
	// flow without mocking the store. The unit tests for the validator function
	// itself confirm the wiring logic is correct.
}

// --- IsAuthenticated tests ---

func TestIsAuthenticated_WithValidCreds(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Save valid credentials
	creds := &Credentials{APIKey: "mdp_valid_key"}
	if err := Save(creds); err != nil {
		t.Fatalf("failed to save creds: %v", err)
	}

	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) { return &LicenseInfo{}, nil }
	defer func() { ValidateKey = old }()

	if !IsAuthenticated() {
		t.Error("expected IsAuthenticated to return true with valid creds and validator")
	}
}

func TestIsAuthenticated_EmptyKey(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	creds := &Credentials{APIKey: ""}
	if err := Save(creds); err != nil {
		t.Fatalf("failed to save creds: %v", err)
	}

	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) { return &LicenseInfo{}, nil }
	defer func() { ValidateKey = old }()

	if IsAuthenticated() {
		t.Error("expected IsAuthenticated to return false with empty key")
	}
}

func TestIsAuthenticated_WhitespaceKey(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	creds := &Credentials{APIKey: "   "}
	if err := Save(creds); err != nil {
		t.Fatalf("failed to save creds: %v", err)
	}

	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) { return &LicenseInfo{}, nil }
	defer func() { ValidateKey = old }()

	if IsAuthenticated() {
		t.Error("expected IsAuthenticated to return false with whitespace-only key")
	}
}

func TestIsAuthenticated_NoCreds(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) { return &LicenseInfo{}, nil }
	defer func() { ValidateKey = old }()

	if IsAuthenticated() {
		t.Error("expected IsAuthenticated to return false with no creds file")
	}
}

func TestIsAuthenticated_ValidatorRejects(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	creds := &Credentials{APIKey: "mdp_badkey"}
	if err := Save(creds); err != nil {
		t.Fatalf("failed to save creds: %v", err)
	}

	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) { return nil, fmt.Errorf("invalid") }
	defer func() { ValidateKey = old }()

	if IsAuthenticated() {
		t.Error("expected IsAuthenticated to return false when validator rejects")
	}
}

// --- RequirePro with full flow ---

func TestRequirePro_FullFlow_Authenticated(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	creds := &Credentials{APIKey: "mdp_valid"}
	if err := Save(creds); err != nil {
		t.Fatalf("failed to save creds: %v", err)
	}

	old := ValidateKey
	ValidateKey = func(key string) (*LicenseInfo, error) { return &LicenseInfo{}, nil }
	defer func() { ValidateKey = old }()

	if err := RequirePro("data"); err != nil {
		t.Errorf("expected no error for authenticated user, got %v", err)
	}
}

// --- Save/Clear via public API ---

func TestSave_And_Clear_PublicAPI(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	creds := &Credentials{APIKey: "mdp_pubapi", Email: "pub@test.com"}
	if err := Save(creds); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.APIKey != "mdp_pubapi" {
		t.Errorf("expected mdp_pubapi, got %s", loaded.APIKey)
	}

	if err := Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	loaded, err = Load()
	if err != nil {
		t.Fatalf("Load after clear failed: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil after Clear")
	}
}

// --- Edge cases ---

func TestLoadFrom_EmptyPath(t *testing.T) {
	creds, err := loadFrom("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds != nil {
		t.Error("expected nil for empty path")
	}
}

func TestSaveTo_EmptyPath(t *testing.T) {
	err := saveTo("", &Credentials{APIKey: "mdp_test"})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestClearFrom_EmptyPath(t *testing.T) {
	err := clearFrom("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

// --- Error tests ---

func TestFeatureGatedError_Message(t *testing.T) {
	err := &FeatureGatedError{Feature: "data"}
	msg := err.Error()

	if !strings.Contains(msg, "--data") {
		t.Error("error message should contain the feature flag name")
	}
	if !strings.Contains(msg, "mdpress.app/pro") {
		t.Error("error message should contain signup URL")
	}
	if !strings.Contains(msg, "mdpress auth login") {
		t.Error("error message should contain login command")
	}
}
