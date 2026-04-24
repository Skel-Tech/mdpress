package config

// DefaultConfig returns a Config with sensible default values.
func DefaultConfig() Config {
	return Config{
		Version:      "1",
		LogoPosition: "top-right",
		LogoWidth:    120,
		Font:         "Helvetica",
		FontSize:     12,
		Header:       false,
		Footer:       "",
		Margins: Margins{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		},
	}
}
