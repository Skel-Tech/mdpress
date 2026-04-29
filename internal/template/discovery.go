package template

import (
	"os"
	"path/filepath"
	"strings"
)

// Source constants for template origins.
const (
	SourceProject = "project"
	SourceGlobal  = "global"
	SourceCloud   = "cloud"
)

// TemplateInfo contains metadata about an available template.
type TemplateInfo struct {
	Name        string // Template name (from YAML name field)
	Description string // Template description (from YAML description field)
	Source      string // "project", "global", or "cloud"
	Path        string // Full path to the template file
	Free        bool   // Whether the template is free (false = requires Pro subscription)
}

// GlobalTemplateDir returns the path to the global templates directory.
func GlobalTemplateDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "mdpress", "templates")
}

// ProjectTemplateDir returns the path to the project templates directory.
func ProjectTemplateDir() string {
	return "templates"
}

// ListTemplates returns a list of all available templates from both
// global and project directories.
// Project templates are listed first, followed by global templates.
// If a template appears in both directories, both versions are returned
// (caller should prefer project over global when resolving by name).
func ListTemplates() ([]TemplateInfo, error) {
	var templates []TemplateInfo

	// List project templates first (they take precedence)
	projectTemplates, err := listTemplatesFromDir(ProjectTemplateDir(), SourceProject)
	if err != nil {
		return nil, err
	}
	templates = append(templates, projectTemplates...)

	// List global templates
	globalTemplates, err := listTemplatesFromDir(GlobalTemplateDir(), SourceGlobal)
	if err != nil {
		return nil, err
	}
	templates = append(templates, globalTemplates...)

	return templates, nil
}

// listTemplatesFromDir scans a directory for YAML template files and returns their info.
func listTemplatesFromDir(dir, source string) ([]TemplateInfo, error) {
	var templates []TemplateInfo

	// If directory doesn't exist, return empty list (not an error)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return templates, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Only process .yml and .yaml files
		if !strings.HasSuffix(name, ".yml") && !strings.HasSuffix(name, ".yaml") {
			continue
		}

		path := filepath.Join(dir, name)
		tmpl, err := LoadFromFile(path)
		if err != nil {
			// Skip invalid templates but continue scanning
			continue
		}

		templates = append(templates, TemplateInfo{
			Name:        tmpl.Name,
			Description: tmpl.Description,
			Source:      source,
			Path:        path,
		})
	}

	return templates, nil
}

// TemplateExistsLocal checks if a template with the given name exists locally
// in either the project or global templates directory.
func TemplateExistsLocal(name string) bool {
	// Check project templates first
	projectTemplates, err := listTemplatesFromDir(ProjectTemplateDir(), SourceProject)
	if err == nil {
		for _, t := range projectTemplates {
			if t.Name == name {
				return true
			}
		}
	}

	// Check global templates
	globalTemplates, err := listTemplatesFromDir(GlobalTemplateDir(), SourceGlobal)
	if err == nil {
		for _, t := range globalTemplates {
			if t.Name == name {
				return true
			}
		}
	}

	return false
}
