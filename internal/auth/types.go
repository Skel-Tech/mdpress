// Package auth provides authentication and feature-gating for mdpress Pro features.
package auth

import "time"

// Credentials holds the authentication state stored in ~/.config/mdpress/auth.yml.
type Credentials struct {
	APIKey    string `yaml:"api_key"`
	Email     string `yaml:"email,omitempty"`
	ExpiresAt string `yaml:"expires_at,omitempty"`
}

// LicenseInfo contains the decoded and verified payload from a license key.
type LicenseInfo struct {
	UserID    string
	Email     string
	Plan      string
	Features  []string
	ExpiresAt time.Time
}
