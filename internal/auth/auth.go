package auth

import "strings"

// ValidateKey is the function used to verify an API key.
// It must be set by the application (wired to the core engine's validator).
// If nil, all pro feature checks will be denied (fail-closed).
// Returns license information on success, or an error if the key is invalid.
var ValidateKey func(key string) (*LicenseInfo, error)

// IsAuthenticated returns true if valid credentials exist and the key passes validation.
func IsAuthenticated() bool {
	creds, err := Load()
	if err != nil || creds == nil {
		return false
	}
	key := strings.TrimSpace(creds.APIKey)
	if key == "" {
		return false
	}
	if ValidateKey == nil {
		return false
	}
	_, err = ValidateKey(key)
	return err == nil
}

// RequirePro checks that the user has a valid Pro license.
// Returns nil if authenticated and the key is valid, or a FeatureGatedError if not.
func RequirePro(feature string) error {
	if IsAuthenticated() {
		return nil
	}
	return &FeatureGatedError{Feature: feature}
}

// isValidKeyFormat checks that an API key has the expected prefix format.
// This is a superficial check — real validation happens in ValidateKey.
func isValidKeyFormat(key string) bool {
	key = strings.TrimSpace(key)
	return key != "" && strings.HasPrefix(key, "mdp_")
}
