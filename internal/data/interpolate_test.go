package data

import (
	"strings"
	"testing"
)

func TestInterpolate_SimpleVariable(t *testing.T) {
	template := "Hello, {{ name }}!"
	data := map[string]any{
		"name": "World",
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_MultipleVariables(t *testing.T) {
	template := "{{ greeting }}, {{ name }}! Today is {{ day }}."
	data := map[string]any{
		"greeting": "Hello",
		"name":     "Alice",
		"day":      "Monday",
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "Hello, Alice! Today is Monday."
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_NestedPath(t *testing.T) {
	template := "Client: {{ client.name }}"
	data := map[string]any{
		"client": map[string]any{
			"name": "Acme Corp",
		},
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "Client: Acme Corp"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_DeeplyNestedPath(t *testing.T) {
	template := "City: {{ client.address.city }}"
	data := map[string]any{
		"client": map[string]any{
			"address": map[string]any{
				"city": "New York",
			},
		},
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "City: New York"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_Loop(t *testing.T) {
	template := "Items:{{#items}}\n- {{ . }}{{/items}}"
	data := map[string]any{
		"items": []any{"Apple", "Banana", "Cherry"},
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "Items:\n- Apple\n- Banana\n- Cherry"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_LoopWithObjects(t *testing.T) {
	template := "Products:{{#products}}\n- {{ name }}: ${{ price }}{{/products}}"
	data := map[string]any{
		"products": []any{
			map[string]any{"name": "Widget", "price": "19.99"},
			map[string]any{"name": "Gadget", "price": "29.99"},
		},
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "Products:\n- Widget: $19.99\n- Gadget: $29.99"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_ConditionalTruthy(t *testing.T) {
	template := "{{#show}}This is shown{{/show}}"
	data := map[string]any{
		"show": true,
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "This is shown"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_ConditionalFalsy(t *testing.T) {
	template := "{{#show}}This is shown{{/show}}"
	data := map[string]any{
		"show": false,
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := ""
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_InvertedSection(t *testing.T) {
	template := "{{^show}}This is shown when false{{/show}}"
	data := map[string]any{
		"show": false,
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "This is shown when false"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_InvertedSectionWithTruthyValue(t *testing.T) {
	template := "{{^show}}Hidden{{/show}}"
	data := map[string]any{
		"show": true,
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := ""
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_UndefinedVariable(t *testing.T) {
	template := "Hello, {{ name }}!"
	data := map[string]any{}

	_, err := Interpolate(template, data)
	if err == nil {
		t.Fatal("Interpolate() should return error for undefined variable")
	}

	undefinedErr, ok := err.(*UndefinedVariableError)
	if !ok {
		t.Fatalf("error should be *UndefinedVariableError, got %T: %v", err, err)
	}

	if undefinedErr.Variable != "name" {
		t.Errorf("Variable = %q, want %q", undefinedErr.Variable, "name")
	}
}

func TestInterpolate_UndefinedVariableWithLineNumber(t *testing.T) {
	template := `Line 1
Line 2
Hello, {{ missing }}!
Line 4`
	data := map[string]any{}

	_, err := Interpolate(template, data)
	if err == nil {
		t.Fatal("Interpolate() should return error for undefined variable")
	}

	undefinedErr, ok := err.(*UndefinedVariableError)
	if !ok {
		t.Fatalf("error should be *UndefinedVariableError, got %T: %v", err, err)
	}

	if undefinedErr.Variable != "missing" {
		t.Errorf("Variable = %q, want %q", undefinedErr.Variable, "missing")
	}

	if undefinedErr.Line != 3 {
		t.Errorf("Line = %d, want %d", undefinedErr.Line, 3)
	}

	// Error message should include line number
	if !strings.Contains(err.Error(), "line 3") {
		t.Errorf("error should mention 'line 3': %v", err)
	}
}

func TestInterpolate_UndefinedNestedVariable(t *testing.T) {
	template := "Client: {{ client.address.city }}"
	data := map[string]any{
		"client": map[string]any{
			"name": "Acme",
		},
	}

	_, err := Interpolate(template, data)
	if err == nil {
		t.Fatal("Interpolate() should return error for undefined nested variable")
	}

	undefinedErr, ok := err.(*UndefinedVariableError)
	if !ok {
		t.Fatalf("error should be *UndefinedVariableError, got %T: %v", err, err)
	}

	// Should report the undefined part of the path
	if !strings.Contains(undefinedErr.Variable, "address") {
		t.Errorf("Variable = %q, should contain 'address'", undefinedErr.Variable)
	}
}

func TestInterpolate_MalformedTemplate(t *testing.T) {
	template := "Hello, {{ name"
	data := map[string]any{
		"name": "World",
	}

	_, err := Interpolate(template, data)
	if err == nil {
		t.Fatal("Interpolate() should return error for malformed template")
	}

	templateErr, ok := err.(*TemplateError)
	if !ok {
		t.Fatalf("error should be *TemplateError, got %T: %v", err, err)
	}

	if !strings.Contains(templateErr.Message, "tag") || !strings.Contains(strings.ToLower(templateErr.Message), "close") {
		// Accept various phrasings of "unclosed tag" errors
		if !strings.Contains(strings.ToLower(err.Error()), "tag") {
			t.Errorf("error message should mention unclosed tag: %v", err)
		}
	}
}

func TestInterpolate_UnclosedSection(t *testing.T) {
	template := "{{#items}}Item{{/wrong}}"
	data := map[string]any{
		"items": []any{"a", "b"},
	}

	_, err := Interpolate(template, data)
	if err == nil {
		t.Fatal("Interpolate() should return error for mismatched section tags")
	}

	_, ok := err.(*TemplateError)
	if !ok {
		t.Fatalf("error should be *TemplateError, got %T: %v", err, err)
	}
}

func TestInterpolate_EmptyTemplate(t *testing.T) {
	template := ""
	data := map[string]any{"name": "test"}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	if result != "" {
		t.Errorf("Interpolate() = %q, want empty string", result)
	}
}

func TestInterpolate_NoVariables(t *testing.T) {
	template := "Just plain text"
	data := map[string]any{}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	if result != "Just plain text" {
		t.Errorf("Interpolate() = %q, want %q", result, "Just plain text")
	}
}

func TestInterpolate_NilData(t *testing.T) {
	template := "Hello!"
	var data map[string]any

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	if result != "Hello!" {
		t.Errorf("Interpolate() = %q, want %q", result, "Hello!")
	}
}

func TestInterpolate_EmptyList(t *testing.T) {
	template := "Items:{{#items}}\n- {{ . }}{{/items}}"
	data := map[string]any{
		"items": []any{},
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	expected := "Items:"
	if result != expected {
		t.Errorf("Interpolate() = %q, want %q", result, expected)
	}
}

func TestInterpolate_NumericValues(t *testing.T) {
	template := "Count: {{ count }}, Price: {{ price }}"
	data := map[string]any{
		"count": 42,
		"price": 19.99,
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	// Numbers should be converted to strings
	if !strings.Contains(result, "42") || !strings.Contains(result, "19.99") {
		t.Errorf("Interpolate() = %q, should contain numeric values", result)
	}
}

func TestInterpolate_BooleanInConditional(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		expected string
	}{
		{"true shows content", map[string]any{"active": true}, "Active"},
		{"false hides content", map[string]any{"active": false}, ""},
		{"truthy string shows content", map[string]any{"active": "yes"}, "Active"},
		{"empty string hides content", map[string]any{"active": ""}, ""},
	}

	template := "{{#active}}Active{{/active}}"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Interpolate(template, tt.data)
			if err != nil {
				t.Fatalf("Interpolate() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Interpolate() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestInterpolate_HTMLEscaping(t *testing.T) {
	template := "Content: {{ content }}"
	data := map[string]any{
		"content": "<script>alert('xss')</script>",
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	// Mustache escapes HTML by default
	if strings.Contains(result, "<script>") {
		t.Errorf("Interpolate() should escape HTML: %q", result)
	}
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Errorf("Interpolate() should contain escaped HTML: %q", result)
	}
}

func TestInterpolate_UnescapedContent(t *testing.T) {
	template := "Content: {{{ content }}}"
	data := map[string]any{
		"content": "<b>Bold</b>",
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	// Triple braces should not escape
	if !strings.Contains(result, "<b>Bold</b>") {
		t.Errorf("Interpolate() with triple braces should not escape HTML: %q", result)
	}
}

func TestInterpolate_MarkdownDocument(t *testing.T) {
	template := `# Invoice for {{ client.name }}

**Date:** {{ date }}
**Invoice Number:** {{ invoice_number }}

## Items

{{#items}}
- {{ description }}: ${{ amount }}
{{/items}}

---

**Total:** ${{ total }}

{{#notes}}
### Notes
{{ . }}
{{/notes}}
`

	data := map[string]any{
		"client": map[string]any{
			"name": "Acme Corp",
		},
		"date":           "2024-03-15",
		"invoice_number": "INV-001",
		"items": []any{
			map[string]any{"description": "Consulting", "amount": "500.00"},
			map[string]any{"description": "Development", "amount": "1500.00"},
		},
		"total": "2000.00",
		"notes": "Payment due within 30 days",
	}

	result, err := Interpolate(template, data)
	if err != nil {
		t.Fatalf("Interpolate() error = %v, want nil", err)
	}

	// Verify key parts of the output
	if !strings.Contains(result, "# Invoice for Acme Corp") {
		t.Errorf("result should contain client name in heading")
	}
	if !strings.Contains(result, "**Date:** 2024-03-15") {
		t.Errorf("result should contain date")
	}
	if !strings.Contains(result, "- Consulting: $500.00") {
		t.Errorf("result should contain first item")
	}
	if !strings.Contains(result, "- Development: $1500.00") {
		t.Errorf("result should contain second item")
	}
	if !strings.Contains(result, "**Total:** $2000.00") {
		t.Errorf("result should contain total")
	}
	if !strings.Contains(result, "Payment due within 30 days") {
		t.Errorf("result should contain notes")
	}
}

func TestTemplateError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *TemplateError
		contains string
	}{
		{
			"with line number",
			&TemplateError{Message: "syntax error", Line: 5},
			"line 5",
		},
		{
			"without line number",
			&TemplateError{Message: "syntax error", Line: 0},
			"syntax error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.err.Error(), tt.contains) {
				t.Errorf("Error() = %q, should contain %q", tt.err.Error(), tt.contains)
			}
		})
	}
}

func TestUndefinedVariableError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *UndefinedVariableError
		contains []string
	}{
		{
			"with line number",
			&UndefinedVariableError{Variable: "foo", Line: 12},
			[]string{"undefined variable", "foo", "line 12"},
		},
		{
			"without line number",
			&UndefinedVariableError{Variable: "bar", Line: 0},
			[]string{"undefined variable", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, s := range tt.contains {
				if !strings.Contains(errMsg, s) {
					t.Errorf("Error() = %q, should contain %q", errMsg, s)
				}
			}
		})
	}
}
