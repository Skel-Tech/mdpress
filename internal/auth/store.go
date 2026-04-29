package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

// AuthFilePath returns the path to the auth credentials file.
func AuthFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "mdpress", "auth.yml")
}

// Load reads credentials from ~/.config/mdpress/auth.yml.
// Returns nil if the file does not exist.
func Load() (*Credentials, error) {
	return loadFrom(AuthFilePath())
}

// loadFrom reads credentials from the given path. Exported for testing via Load.
func loadFrom(path string) (*Credentials, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading auth file: %w", err)
	}

	var creds Credentials
	if err := yaml.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("parsing auth file: %w", err)
	}

	return &creds, nil
}

// Save writes credentials to ~/.config/mdpress/auth.yml.
func Save(creds *Credentials) error {
	return saveTo(AuthFilePath(), creds)
}

// saveTo writes credentials to the given path.
func saveTo(path string, creds *Credentials) error {
	if path == "" {
		return fmt.Errorf("could not determine auth file path")
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(creds)
	if err != nil {
		return fmt.Errorf("encoding credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing auth file: %w", err)
	}

	return nil
}

// Clear removes the auth credentials file.
func Clear() error {
	return clearFrom(AuthFilePath())
}

// clearFrom removes the auth file at the given path.
func clearFrom(path string) error {
	if path == "" {
		return fmt.Errorf("could not determine auth file path")
	}

	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil // Already gone, not an error
	}
	return err
}
