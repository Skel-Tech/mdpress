package template

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Resolve resolves a template name to a loaded Template.
// Project templates take precedence over global templates with the same name.
// Returns a helpful error listing available templates if the name is not found.
func Resolve(name string) (*Template, error) {
	// If it looks like a path, use LoadFromFile directly
	if IsPath(name) {
		return LoadFromFile(name)
	}

	// List all available templates
	templates, err := ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	// Search for matching template by name
	// Since ListTemplates returns project templates first, we'll naturally
	// get project precedence by taking the first match
	for _, info := range templates {
		if info.Name == name {
			return LoadFromFile(info.Path)
		}
	}

	// Template not found - build helpful error message
	return nil, buildNotFoundError(name, templates)
}

// IsPath determines whether a value should be treated as a file path
// rather than a template name.
// Values containing "/" or starting with "." are treated as paths.
func IsPath(value string) bool {
	if value == "" {
		return false
	}
	// Contains path separator (works for both Unix and Windows)
	if strings.Contains(value, "/") || strings.Contains(value, string(filepath.Separator)) {
		return true
	}
	// Starts with . (e.g., ./template.yml or ../template.yml)
	if strings.HasPrefix(value, ".") {
		return true
	}
	return false
}

// buildNotFoundError creates a helpful error message listing available templates.
func buildNotFoundError(name string, templates []TemplateInfo) error {
	if len(templates) == 0 {
		return fmt.Errorf("template %q not found: no templates available", name)
	}

	var available []string
	for _, t := range templates {
		available = append(available, fmt.Sprintf("  - %s (%s)", t.Name, t.Source))
	}

	return fmt.Errorf("template %q not found\n\nAvailable templates:\n%s",
		name, strings.Join(available, "\n"))
}
