package template

import (
	"fmt"
	"os"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// LoadFromFile loads a template from a YAML file.
// It returns an error if the file cannot be read, contains invalid YAML,
// has unknown fields, or fails validation.
func LoadFromFile(path string) (*Template, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file not found: %s", path)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", path, err)
	}

	var tmpl Template
	decoderConfig := mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &tmpl,
		TagName:     "mapstructure",
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(v.AllSettings()); err != nil {
		return nil, fmt.Errorf("invalid template %s: %w", path, err)
	}

	tmpl.SourcePath = path

	// Validate the template
	if err := validate(&tmpl); err != nil {
		return nil, fmt.Errorf("invalid template %s: %w", path, err)
	}

	return &tmpl, nil
}

// validate checks that the template has all required fields and valid values.
func validate(tmpl *Template) error {
	// Validate required fields
	if tmpl.Name == "" {
		return fmt.Errorf("name is required")
	}

	if tmpl.Version == "" {
		return fmt.Errorf("version is required")
	}

	if tmpl.Version != "1" {
		return fmt.Errorf("version must be \"1\", got %q", tmpl.Version)
	}

	return nil
}
