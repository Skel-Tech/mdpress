# mdpress Directives Reference

Directives control document layout and semantic structure in mdpress. They use the `{{@name}}` syntax and must appear on their own line.

---

## General Rules

1. **Own line only** - Directives cannot appear inline within text
2. **No spaces in tags** - Use `{{@pagebreak}}` not `{{@ pagebreak }}`
3. **Exact names only** - Unknown directives cause parse errors
4. **Block directives require closing** - `{{@note}}...{{/note}}`

---

## Layout Directives

### {{@pagebreak}}

Forces a page break at this location.

**Type:** Self-closing

**Use case:** Separate chapters, start new sections on fresh pages.

**GOOD:**
```markdown
## Chapter 1

Content for chapter one...

{{@pagebreak}}

## Chapter 2

Content for chapter two...
```

**BAD:**
```markdown
Content {{@pagebreak}} more content   <- must be on own line
{{@ pagebreak }}                       <- no spaces inside tag
{{@page-break}}                        <- wrong name (hyphen)
{{@break}}                             <- wrong name
```

---

### {{@columns}}...{{/columns}}

Renders enclosed content in a multi-column layout. Use horizontal rules (`---`) to separate columns.

**Type:** Block

**Use case:** Side-by-side content, comparison layouts.

**GOOD:**
```markdown
{{@columns}}
**Left Column**

Content for the left side of the page.

---

**Right Column**

Content for the right side of the page.
{{/columns}}
```

**BAD:**
```markdown
{{@columns}}Left{{/columns}}{{@columns}}Right{{/columns}}  <- use one block with ---
{{@columns cols=3}}                                          <- attributes not supported
{{@columns}}inline content{{/columns}}                       <- should span multiple lines
```

---

### {{@toc}}

Generates a table of contents based on document headings.

**Type:** Self-closing

**Use case:** Documents with multiple sections, reports, long-form content.

**GOOD:**
```markdown
# Document Title

{{@toc}}

## Introduction

Introduction content...

## Methods

Methods content...

## Results

Results content...
```

**BAD:**
```markdown
{{@toc depth=2}}         <- attributes not supported
{{@tableofcontents}}     <- wrong name
{{@table-of-contents}}   <- wrong name
```

---

## Semantic Directives

### {{@note}}...{{/note}}

Renders content as a note callout box. Use for helpful tips, additional information, or non-critical context.

**Type:** Block

**Use case:** Tips, supplementary information, helpful reminders.

**GOOD:**
```markdown
{{@note}}
Remember to save your work frequently. Auto-save is enabled by default.
{{/note}}
```

```markdown
{{@note}}
**Tip:** You can use keyboard shortcuts for faster navigation.

- `Ctrl+S` - Save
- `Ctrl+Z` - Undo
- `Ctrl+Y` - Redo
{{/note}}
```

**BAD:**
```markdown
{{@note}}Inline note{{/note}}          <- should be on own lines
> **Note:** Remember to save...         <- use directive for consistent styling
**Note:** Some text                     <- use directive instead
```

---

### {{@warning}}...{{/warning}}

Renders content as a warning callout box. Use for cautions, important alerts, or information that requires attention.

**Type:** Block

**Use case:** Cautions, critical information, alerts, destructive action warnings.

**GOOD:**
```markdown
{{@warning}}
This action cannot be undone. Please backup your data before proceeding.
{{/warning}}
```

```markdown
{{@warning}}
**Breaking Change**

Version 2.0 removes support for legacy configuration files.
Migrate your settings before upgrading.
{{/warning}}
```

**BAD:**
```markdown
{{@warning}}Danger!{{/warning}}        <- should be on own lines
**WARNING:** This action...             <- use directive for consistent styling
{{@danger}}...{{/danger}}              <- wrong name, use warning
{{@alert}}...{{/alert}}                <- wrong name, use warning
```

---

### {{@highlight}}...{{/highlight}}

Renders content with visual emphasis/highlighting. Use for key takeaways, important conclusions, or content that should stand out.

**Type:** Block

**Use case:** Key findings, important conclusions, executive summaries.

**GOOD:**
```markdown
{{@highlight}}
Key takeaway: Always validate user input before processing.
{{/highlight}}
```

```markdown
{{@highlight}}
**Revenue increased 23% year-over-year**, marking the strongest Q1 in company history.
{{/highlight}}
```

**BAD:**
```markdown
The {{@highlight}}important{{/highlight}} part   <- must be block level, not inline
{{@highlight}}One liner{{/highlight}}            <- works but prefer multiline for readability
{{@emphasis}}...{{/emphasis}}                     <- wrong name
{{@callout}}...{{/callout}}                       <- wrong name
```

---

## Quick Reference

| Directive | Type | Purpose |
|-----------|------|---------|
| `{{@pagebreak}}` | Self-closing | Force page break |
| `{{@toc}}` | Self-closing | Generate table of contents |
| `{{@columns}}...{{/columns}}` | Block | Multi-column layout |
| `{{@note}}...{{/note}}` | Block | Note/tip callout |
| `{{@warning}}...{{/warning}}` | Block | Warning/alert callout |
| `{{@highlight}}...{{/highlight}}` | Block | Emphasized content block |

---

## Common Mistakes

### Inline Directives

Directives cannot appear inline within text.

**WRONG:**
```markdown
The total is {{@highlight}}$500{{/highlight}} due today.
```

**CORRECT:**
```markdown
{{@highlight}}
The total is $500 due today.
{{/highlight}}
```

### Unknown Directives

Only use documented directives. These will cause errors:

```markdown
{{@custom}}          <- not a real directive
{{@sidebar}}         <- not a real directive
{{@footer}}          <- not a real directive
{{@header}}          <- not a real directive
{{@section}}         <- not a real directive
```

### Attribute Syntax

Directives do not support attributes or parameters:

```markdown
{{@toc depth=2}}           <- attributes not supported
{{@columns count=3}}       <- attributes not supported
{{@pagebreak type=soft}}   <- attributes not supported
```

---

## Related Documentation

- [Syntax Guide](syntax.md) - Core syntax overview
- [Variables Guide](variables.md) - Mustache templating
- [MCP Integration](mcp.md) - Tool integration for AI assistants
