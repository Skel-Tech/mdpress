package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

// LoadFile loads data from a JSON or YAML file and returns it as a map.
// The file format is detected by the file extension:
//   - .json for JSON files
//   - .yaml or .yml for YAML files
//
// Returns an error if the file cannot be read, contains invalid content,
// or has an unsupported file extension.
func LoadFile(path string) (map[string]any, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("data file not found: %s", path)
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file %s: %w", path, err)
	}

	// Determine format by extension
	ext := strings.ToLower(filepath.Ext(path))

	var data map[string]any

	switch ext {
	case ".json":
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("invalid YAML in %s: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported data file extension %q in %s: expected .json, .yaml, or .yml", ext, path)
	}

	return data, nil
}
