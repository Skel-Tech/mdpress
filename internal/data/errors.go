package data

import "fmt"

// TemplateError represents an error that occurred during template interpolation.
// It includes the line number where the error occurred when available.
type TemplateError struct {
	Message string
	Line    int
}

// Error implements the error interface.
func (e *TemplateError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("template error: %s at line %d", e.Message, e.Line)
	}
	return fmt.Sprintf("template error: %s", e.Message)
}

// UndefinedVariableError represents an error when a variable is not defined in the data.
type UndefinedVariableError struct {
	Variable string
	Line     int
}

// Error implements the error interface.
func (e *UndefinedVariableError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("template error: undefined variable %q at line %d", e.Variable, e.Line)
	}
	return fmt.Sprintf("template error: undefined variable %q", e.Variable)
}
