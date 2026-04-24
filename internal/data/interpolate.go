package data

import (
	"regexp"
	"strings"
	"sync"

	"github.com/cbroglie/mustache"
)

// mu protects access to the global AllowMissingVariables setting
var mu sync.Mutex

// Interpolate renders a Mustache template with the provided data.
// It supports:
//   - Variable interpolation: {{ variable }} and {{ nested.path }}
//   - Sections/loops: {{#items}}...{{/items}}
//   - Conditionals: {{#condition}}...{{/condition}} (truthy) and {{^condition}}...{{/condition}} (falsy)
//
// Returns an error if the template is malformed or references undefined variables.
// Error messages include line numbers where possible.
func Interpolate(template string, data map[string]any) (string, error) {
	// Parse the template first (before modifying global state)
	tmpl, err := mustache.ParseString(template)
	if err != nil {
		return "", parseErrorToTemplateError(err, template)
	}

	// Lock and set strict mode for rendering
	mu.Lock()
	oldSetting := mustache.AllowMissingVariables
	mustache.AllowMissingVariables = false
	defer func() {
		mustache.AllowMissingVariables = oldSetting
		mu.Unlock()
	}()

	// Render with the data
	result, err := tmpl.Render(data)
	if err != nil {
		return "", convertRenderError(err, template)
	}

	return result, nil
}

// missingVarRegex matches "missing variable "name"" errors
var missingVarRegex = regexp.MustCompile(`missing variable "([^"]+)"`)

// convertRenderError converts a mustache render error to the appropriate error type.
func convertRenderError(err error, template string) error {
	errMsg := err.Error()

	// Check if this is a missing variable error
	if matches := missingVarRegex.FindStringSubmatch(errMsg); len(matches) > 1 {
		varName := matches[1]
		line := findVariableLine(template, varName)
		return &UndefinedVariableError{
			Variable: varName,
			Line:     line,
		}
	}

	// Generic template error
	return &TemplateError{
		Message: errMsg,
		Line:    extractLineNumber(errMsg),
	}
}

// parseErrorToTemplateError converts a mustache parse error to a TemplateError.
func parseErrorToTemplateError(err error, template string) error {
	msg := err.Error()
	return &TemplateError{
		Message: msg,
		Line:    extractLineNumber(msg),
	}
}

// lineRegex matches "line X" in error messages
var lineRegex = regexp.MustCompile(`line (\d+)`)

// extractLineNumber tries to extract a line number from an error message.
func extractLineNumber(msg string) int {
	if matches := lineRegex.FindStringSubmatch(msg); len(matches) > 1 {
		var n int
		for _, c := range matches[1] {
			n = n*10 + int(c-'0')
		}
		return n
	}
	return 0
}

// findVariableLine finds the line number where a variable is first referenced.
func findVariableLine(template, variable string) int {
	lines := strings.Split(template, "\n")
	for i, line := range lines {
		// Check for {{ variable }} or {{variable}} or {{{variable}}} patterns
		// Also handle nested paths - look for the full variable or just the first part
		if strings.Contains(line, "{{") && containsVariable(line, variable) {
			return i + 1
		}
	}
	return 0
}

// containsVariable checks if a line contains a reference to the given variable.
func containsVariable(line, variable string) bool {
	// Direct match
	if strings.Contains(line, variable) {
		return true
	}
	// For nested paths like "client.address", also check for the path parts
	parts := strings.Split(variable, ".")
	if len(parts) > 1 {
		// Check if the first part of the path is in the line
		return strings.Contains(line, parts[0])
	}
	return false
}
