package auth

import "fmt"

// FeatureGatedError is returned when a Pro feature is used without authentication.
type FeatureGatedError struct {
	Feature string
}

func (e *FeatureGatedError) Error() string {
	return fmt.Sprintf(
		"The --%s flag requires an mdpress Pro account.\n\n"+
			"  Sign up at: https://mdpress.app/pro\n\n"+
			"  Already have an account? Run: mdpress auth login",
		e.Feature,
	)
}
