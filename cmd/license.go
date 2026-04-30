package cmd

import (
	mdpress "github.com/skel-tech/mdpress-core"
	"github.com/skel-tech/mdpress/internal/auth"
)

func init() {
	// Wire the core engine's license validator into the auth package.
	// The actual validation logic lives in mdpress-core (private) so it
	// cannot be inspected or bypassed from the public CLI source.
	auth.ValidateKey = func(key string) (*auth.LicenseInfo, error) {
		info, err := mdpress.ValidateLicense(key)
		if err != nil {
			return nil, err
		}
		return &auth.LicenseInfo{
			UserID:    info.UserID,
			Email:     info.Email,
			Plan:      info.Plan,
			Features:  info.Features,
			ExpiresAt: info.ExpiresAt,
		}, nil
	}
}
