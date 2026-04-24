package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFile_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.json")

	content := `{
  "name": "Acme Corp",
  "year": 2024,
  "active": true
}`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	data, err := LoadFile(dataPath)
	if err != nil {
		t.Fatalf("LoadFile() error = %v, want nil", err)
	}

	if data["name"] != "Acme Corp" {
		t.Errorf("name = %q, want %q", data["name"], "Acme Corp")
	}
	// JSON numbers are float64
	if data["year"] != float64(2024) {
		t.Errorf("year = %v, want %v", data["year"], 2024)
	}
	if data["active"] != true {
		t.Errorf("active = %v, want %v", data["active"], true)
	}
}

func TestLoadFile_YAML(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.yaml")

	content := `
name: "Acme Corp"
year: 2024
active: true
`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	data, err := LoadFile(dataPath)
	if err != nil {
		t.Fatalf("LoadFile() error = %v, want nil", err)
	}

	if data["name"] != "Acme Corp" {
		t.Errorf("name = %q, want %q", data["name"], "Acme Corp")
	}
	if data["year"] != 2024 {
		t.Errorf("year = %v, want %v", data["year"], 2024)
	}
	if data["active"] != true {
		t.Errorf("active = %v, want %v", data["active"], true)
	}
}

func TestLoadFile_YML(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.yml")

	content := `
company: "Example Inc"
employees: 50
`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	data, err := LoadFile(dataPath)
	if err != nil {
		t.Fatalf("LoadFile() error = %v, want nil", err)
	}

	if data["company"] != "Example Inc" {
		t.Errorf("company = %q, want %q", data["company"], "Example Inc")
	}
	if data["employees"] != 50 {
		t.Errorf("employees = %v, want %v", data["employees"], 50)
	}
}

func TestLoadFile_NestedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "nested.json")

	content := `{
  "client": {
    "name": "Acme Corp",
    "address": {
      "city": "New York",
      "zip": "10001"
    }
  },
  "items": ["item1", "item2", "item3"]
}`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	data, err := LoadFile(dataPath)
	if err != nil {
		t.Fatalf("LoadFile() error = %v, want nil", err)
	}

	// Check nested client object
	client, ok := data["client"].(map[string]any)
	if !ok {
		t.Fatalf("client should be a map, got %T", data["client"])
	}
	if client["name"] != "Acme Corp" {
		t.Errorf("client.name = %q, want %q", client["name"], "Acme Corp")
	}

	// Check deeply nested address
	address, ok := client["address"].(map[string]any)
	if !ok {
		t.Fatalf("client.address should be a map, got %T", client["address"])
	}
	if address["city"] != "New York" {
		t.Errorf("client.address.city = %q, want %q", address["city"], "New York")
	}
	if address["zip"] != "10001" {
		t.Errorf("client.address.zip = %q, want %q", address["zip"], "10001")
	}

	// Check items array
	items, ok := data["items"].([]any)
	if !ok {
		t.Fatalf("items should be an array, got %T", data["items"])
	}
	if len(items) != 3 {
		t.Errorf("items length = %d, want %d", len(items), 3)
	}
	if items[0] != "item1" {
		t.Errorf("items[0] = %q, want %q", items[0], "item1")
	}
}

func TestLoadFile_NestedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "nested.yaml")

	content := `
client:
  name: "Acme Corp"
  address:
    city: "New York"
    zip: "10001"
items:
  - item1
  - item2
  - item3
`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	data, err := LoadFile(dataPath)
	if err != nil {
		t.Fatalf("LoadFile() error = %v, want nil", err)
	}

	// Check nested client object
	client, ok := data["client"].(map[string]any)
	if !ok {
		t.Fatalf("client should be a map, got %T", data["client"])
	}
	if client["name"] != "Acme Corp" {
		t.Errorf("client.name = %q, want %q", client["name"], "Acme Corp")
	}

	// Check deeply nested address
	address, ok := client["address"].(map[string]any)
	if !ok {
		t.Fatalf("client.address should be a map, got %T", client["address"])
	}
	if address["city"] != "New York" {
		t.Errorf("client.address.city = %q, want %q", address["city"], "New York")
	}
	if address["zip"] != "10001" {
		t.Errorf("client.address.zip = %q, want %q", address["zip"], "10001")
	}

	// Check items array
	items, ok := data["items"].([]any)
	if !ok {
		t.Fatalf("items should be an array, got %T", data["items"])
	}
	if len(items) != 3 {
		t.Errorf("items length = %d, want %d", len(items), 3)
	}
	if items[0] != "item1" {
		t.Errorf("items[0] = %q, want %q", items[0], "item1")
	}
}

func TestLoadFile_FileNotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/data.json")
	if err == nil {
		t.Fatal("LoadFile() should return error for non-existent file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention 'not found': %v", err)
	}
	if !strings.Contains(err.Error(), "/nonexistent/path/data.json") {
		t.Errorf("error should include file path: %v", err)
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "invalid.json")

	content := `{
  "name": "Acme Corp",
  invalid json here
}`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFile(dataPath)
	if err == nil {
		t.Fatal("LoadFile() should return error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("error should mention 'invalid JSON': %v", err)
	}
	if !strings.Contains(err.Error(), dataPath) {
		t.Errorf("error should include file path: %v", err)
	}
}

func TestLoadFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "invalid.yaml")

	content := `
name: "Bad YAML"
items: [invalid
  - yaml
`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFile(dataPath)
	if err == nil {
		t.Fatal("LoadFile() should return error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "invalid YAML") {
		t.Errorf("error should mention 'invalid YAML': %v", err)
	}
	if !strings.Contains(err.Error(), dataPath) {
		t.Errorf("error should include file path: %v", err)
	}
}

func TestLoadFile_UnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.txt")

	content := `some text content`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFile(dataPath)
	if err == nil {
		t.Fatal("LoadFile() should return error for unsupported extension")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("error should mention 'unsupported': %v", err)
	}
	if !strings.Contains(err.Error(), ".txt") {
		t.Errorf("error should mention the extension '.txt': %v", err)
	}
	if !strings.Contains(err.Error(), dataPath) {
		t.Errorf("error should include file path: %v", err)
	}
}

func TestLoadFile_NoExtension(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data")

	content := `{"name": "test"}`
	if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFile(dataPath)
	if err == nil {
		t.Fatal("LoadFile() should return error for file without extension")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("error should mention 'unsupported': %v", err)
	}
}

func TestLoadFile_EmptyJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "empty.json")

	if err := os.WriteFile(dataPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Empty JSON is invalid (requires at least {} or null)
	_, err := LoadFile(dataPath)
	if err == nil {
		t.Fatal("LoadFile() should return error for empty JSON file")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("error should mention 'invalid JSON': %v", err)
	}
}

func TestLoadFile_EmptyYAMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "empty.yaml")

	if err := os.WriteFile(dataPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Empty YAML results in nil map, not an error
	data, err := LoadFile(dataPath)
	if err != nil {
		t.Fatalf("LoadFile() error = %v, want nil", err)
	}
	if data != nil {
		t.Errorf("data = %v, want nil for empty YAML file", data)
	}
}

func TestLoadFile_CaseInsensitiveExtension(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name string
		ext  string
	}{
		{"uppercase JSON", ".JSON"},
		{"uppercase YAML", ".YAML"},
		{"uppercase YML", ".YML"},
		{"mixed case Json", ".Json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataPath := filepath.Join(tmpDir, "data"+tt.ext)
			var content string
			if strings.Contains(strings.ToLower(tt.ext), "json") {
				content = `{"name": "test"}`
			} else {
				content = `name: "test"`
			}
			if err := os.WriteFile(dataPath, []byte(content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			data, err := LoadFile(dataPath)
			if err != nil {
				t.Fatalf("LoadFile() error = %v, want nil", err)
			}
			if data["name"] != "test" {
				t.Errorf("name = %q, want %q", data["name"], "test")
			}
		})
	}
}
