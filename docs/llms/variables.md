# mdpress Variables Guide

Variables enable data-driven document generation using Mustache templating syntax.

---

## Pro Feature Requirement

Variables are a **Pro feature**. To use variables you need:

1. **Pro license** - Sign in with `mdpress auth login`
2. **Data file** - Pass data via the `--data` flag

```bash
# Requires Pro authentication
mdpress render template.md --data values.json
```

Without Pro authentication, the `--data` flag will return an error.

---

## Basic Interpolation

**Syntax:** `{{ variable }}`

Inserts the value of a variable. HTML entities are escaped by default.

**Template:**
```markdown
Hello, {{ name }}!
```

**Data (JSON):**
```json
{
  "name": "World"
}
```

**Output:**
```
Hello, World!
```

### Examples

**GOOD:**
```markdown
{{ client_name }}
{{ invoice_number }}
{{ date }}
```

**BAD:**
```markdown
{{client_name}}        <- works but less readable, prefer spaces
{{ CLIENT_NAME }}      <- case must match data exactly
{{ name | upper }}     <- filters not supported
{{ name | default: "N/A" }}  <- defaults not supported
```

---

## Nested Object Access

**Syntax:** `{{ object.property }}` or `{{ object.nested.path }}`

Access properties within nested objects using dot notation.

**Template:**
```markdown
Client: {{ client.name }}
City: {{ client.address.city }}
```

**Data (JSON):**
```json
{
  "client": {
    "name": "Acme Corporation",
    "address": {
      "city": "San Francisco",
      "state": "CA"
    }
  }
}
```

**Output:**
```
Client: Acme Corporation
City: San Francisco
```

### Examples

**GOOD:**
```markdown
{{ client.name }}
{{ client.address.city }}
{{ company.contact.email }}
```

**BAD:**
```markdown
{{ client['name'] }}         <- bracket notation not supported
{{ client?.name }}           <- optional chaining not supported
{{ client.address[0] }}      <- array index in path not supported
```

---

## Loops and Sections

**Syntax:** `{{#array}}...{{/array}}`

Iterate over arrays. Inside the loop:
- Use `{{ . }}` for simple values (strings, numbers)
- Access object properties directly (no prefix needed)

### Simple Array

**Template:**
```markdown
{{#colors}}
- {{ . }}
{{/colors}}
```

**Data:**
```json
{
  "colors": ["Red", "Green", "Blue"]
}
```

**Output:**
```
- Red
- Green
- Blue
```

### Object Array

**Template:**
```markdown
| Item | Price |
|------|-------|
{{#items}}
| {{ name }} | ${{ price }} |
{{/items}}
```

**Data:**
```json
{
  "items": [
    { "name": "Widget", "price": "19.99" },
    { "name": "Gadget", "price": "29.99" }
  ]
}
```

**Output:**
```
| Item | Price |
|------|-------|
| Widget | $19.99 |
| Gadget | $29.99 |
```

### Examples

**GOOD:**
```markdown
{{#items}}
- {{ description }}: ${{ amount }}
{{/items}}

{{#tags}}
- {{ . }}
{{/tags}}
```

**BAD:**
```markdown
{% for item in items %}      <- Jinja syntax not supported
{{#items}}{{ items.name }}{{/items}}   <- inside loop, access directly: {{ name }}
{{ items[0].name }}          <- array indexing not supported
```

---

## Conditionals

**Syntax:** `{{#condition}}...{{/condition}}`

Renders content only if the value is truthy:
- `true`
- Non-empty string
- Non-empty array
- Object with properties

**Template:**
```markdown
{{#premium}}
## Premium Member Benefits
- Priority support
- Extended features
{{/premium}}
```

**Data (truthy):**
```json
{ "premium": true }
```

**Output:**
```
## Premium Member Benefits
- Priority support
- Extended features
```

**Data (falsy):**
```json
{ "premium": false }
```

**Output:** *(empty - section not rendered)*

### Examples

**GOOD:**
```markdown
{{#show_notes}}
## Notes
{{ notes }}
{{/show_notes}}

{{#items}}
## Items
(loop content)
{{/items}}
```

**BAD:**
```markdown
{{#if show_notes}}          <- "if" is literal, not a keyword
{{#show_notes == true}}     <- comparison operators not supported
{{#count > 0}}              <- comparisons not supported
```

---

## Inverted Conditionals

**Syntax:** `{{^condition}}...{{/condition}}`

Renders content only if the value is falsy:
- `false`
- Empty string `""`
- Empty array `[]`
- `null` or undefined

**Template:**
```markdown
{{^items}}
No items found.
{{/items}}
```

**Data:**
```json
{ "items": [] }
```

**Output:**
```
No items found.
```

### Examples

**GOOD:**
```markdown
{{^paid}}
**Payment pending**
{{/paid}}

{{^items}}
No items to display.
{{/items}}
```

**BAD:**
```markdown
{{!paid}}              <- wrong syntax (! is for comments)
{{^paid == false}}     <- comparison not supported
{{#!paid}}             <- wrong syntax
```

---

## Unescaped Content

**Syntax:** `{{{ variable }}}` (triple braces)

Outputs content without HTML escaping. Use when your data contains intentional HTML or Markdown formatting.

**Template:**
```markdown
{{{ formatted_content }}}
```

**Data:**
```json
{
  "formatted_content": "**Bold** and *italic* text"
}
```

**Output:**
```
**Bold** and *italic* text
```

### Examples

**GOOD:**
```markdown
{{{ html_content }}}
{{{ markdown_snippet }}}
{{{ preformatted }}}
```

**BAD:**
```markdown
{{{ user_input }}}     <- security risk if content is untrusted
```

---

## Data File Formats

mdpress accepts data in JSON or YAML format.

### JSON Example

**data.json:**
```json
{
  "title": "Monthly Report",
  "date": "2024-03-15",
  "author": {
    "name": "Jane Smith",
    "email": "jane@example.com"
  },
  "items": [
    { "name": "Task 1", "status": "complete" },
    { "name": "Task 2", "status": "pending" }
  ]
}
```

### YAML Example

**data.yaml:**
```yaml
title: Monthly Report
date: 2024-03-15
author:
  name: Jane Smith
  email: jane@example.com
items:
  - name: Task 1
    status: complete
  - name: Task 2
    status: pending
```

### Usage

```bash
mdpress render report.md --data data.json
mdpress render report.md --data data.yaml
```

---

## Error Handling

### Undefined Variables

Referencing a variable that doesn't exist in your data will output an empty string (no error).

**Template:**
```markdown
Name: {{ name }}
Email: {{ email }}
```

**Data:**
```json
{ "name": "Alice" }
```

**Output:**
```
Name: Alice
Email:
```

### Best Practice

Ensure all variables in your template have corresponding keys in your data file. Missing variables produce empty output, which may not be the intended behavior.

**Recommended:** Use inverted conditionals to handle missing data:

```markdown
{{#email}}
Email: {{ email }}
{{/email}}
{{^email}}
Email: Not provided
{{/email}}
```

---

## Quick Reference

| Syntax | Purpose | Example |
|--------|---------|---------|
| `{{ var }}` | Basic interpolation | `{{ name }}` |
| `{{ a.b }}` | Nested access | `{{ user.email }}` |
| `{{#arr}}...{{/arr}}` | Loop over array | `{{#items}}...{{/items}}` |
| `{{#bool}}...{{/bool}}` | Conditional (truthy) | `{{#active}}...{{/active}}` |
| `{{^bool}}...{{/bool}}` | Inverted conditional | `{{^active}}...{{/active}}` |
| `{{{ var }}}` | Unescaped output | `{{{ html }}}` |

---

## Related Documentation

- [Syntax Guide](syntax.md) - Core syntax overview
- [Directives Reference](directives.md) - Layout and semantic directives
- [MCP Integration](mcp.md) - Tool integration for AI assistants
