# mdpress Syntax Guide

This document covers the core syntax elements for authoring mdpress-compatible Markdown documents.

---

## Markdown Basics

mdpress uses standard Markdown (CommonMark) for all formatting. All standard Markdown features work as expected:

### Headings

```markdown
# Heading 1
## Heading 2
### Heading 3
#### Heading 4
```

### Text Formatting

```markdown
**bold text**
*italic text*
***bold and italic***
`inline code`
```

### Lists

```markdown
- Unordered item
- Another item
  - Nested item

1. Ordered item
2. Another item
   1. Nested ordered
```

### Links and Images

```markdown
[Link text](https://example.com)
![Alt text](image.png)
```

### Tables

```markdown
| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |
| Cell 3   | Cell 4   |
```

### Code Blocks

````markdown
```javascript
function hello() {
  console.log("Hello, world!");
}
```
````

### Blockquotes

```markdown
> This is a blockquote.
> It can span multiple lines.
```

---

## Variable Syntax

Variables use Mustache templating syntax with double curly braces. Variables are a **Pro feature** requiring the `--data` flag and Pro license.

### Basic Variable

```
{{ variable_name }}
```

Inserts the value from your data file.

**GOOD:**
```markdown
Hello, {{ name }}!
Invoice #{{ invoice_number }}
```

**BAD:**
```markdown
{{name}}              <- works but less readable, prefer spaces
{{ Name }}            <- case must match data exactly
{{ name | upper }}    <- filters not supported
```

### Nested Object Access

```
{{ object.property }}
{{ object.nested.path }}
```

Access nested data using dot notation.

**GOOD:**
```markdown
{{ client.name }}
{{ client.address.city }}
```

**BAD:**
```markdown
{{ client['name'] }}     <- bracket notation not supported
{{ client?.name }}       <- optional chaining not supported
```

### Loops and Sections

```
{{#array}}
Content repeated for each item
{{/array}}
```

Inside loops, access properties directly or use `{{ . }}` for simple values.

**GOOD:**
```markdown
{{#items}}
- {{ name }}: ${{ price }}
{{/items}}
```

**BAD:**
```markdown
{% for item in items %}   <- Jinja syntax not supported
{{#items}}{{ items.name }}{{/items}}   <- inside loop, access properties directly
```

### Conditionals

```
{{#condition}}
Content shown if truthy
{{/condition}}

{{^condition}}
Content shown if falsy
{{/condition}}
```

**GOOD:**
```markdown
{{#show_notes}}
## Notes
{{ notes }}
{{/show_notes}}

{{^paid}}
**Payment pending**
{{/paid}}
```

**BAD:**
```markdown
{{#if show_notes}}        <- "if" is literal, not a keyword
{{#paid == true}}         <- comparison operators not supported
```

---

## Directive Syntax

Directives control layout and semantic structure. They use the `{{@name}}` syntax.

### Key Rules

1. **Directives must be on their own line** - never inline with text
2. **No spaces inside the tag** - use `{{@pagebreak}}` not `{{@ pagebreak }}`
3. **Use only documented directives** - unknown directives cause errors

### Self-Closing Directives

Stand-alone directives that don't wrap content:

```markdown
{{@pagebreak}}
{{@toc}}
```

### Block Directives

Directives that wrap content with opening and closing tags:

```markdown
{{@note}}
Content inside the note.
{{/note}}

{{@warning}}
Warning content here.
{{/warning}}
```

### Directive vs Variable

- **Directive:** starts with `{{@` (e.g., `{{@pagebreak}}`)
- **Variable:** starts with `{{` without `@` (e.g., `{{ name }}`)

**GOOD:**
```markdown
## Invoice for {{ client.name }}

{{@note}}
Payment due within 30 days.
{{/note}}

{{@pagebreak}}

## Line Items
{{#items}}
- {{ description }}: ${{ amount }}
{{/items}}
```

**BAD:**
```markdown
Total: {{@highlight}}$500{{/highlight}}   <- inline directive (wrong)
{{@custom}}                                <- unknown directive (error)
<style>body{color:red}</style>             <- raw CSS not supported
```

---

## Summary Table

| Syntax | Purpose | Example |
|--------|---------|---------|
| `{{ var }}` | Variable interpolation | `{{ name }}` |
| `{{ a.b }}` | Nested access | `{{ client.name }}` |
| `{{#arr}}...{{/arr}}` | Loop | `{{#items}}...{{/items}}` |
| `{{#bool}}...{{/bool}}` | Conditional (truthy) | `{{#paid}}...{{/paid}}` |
| `{{^bool}}...{{/bool}}` | Conditional (falsy) | `{{^paid}}...{{/paid}}` |
| `{{{ var }}}` | Unescaped output | `{{{ html }}}` |
| `{{@directive}}` | Self-closing directive | `{{@pagebreak}}` |
| `{{@dir}}...{{/dir}}` | Block directive | `{{@note}}...{{/note}}` |

---

## Related Documentation

- [Directives Reference](directives.md) - Complete directive documentation
- [Variables Guide](variables.md) - Mustache templating details
- [MCP Integration](mcp.md) - Tool integration for AI assistants
