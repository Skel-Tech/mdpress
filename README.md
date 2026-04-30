# mdpress

**Turn your Markdown and AI output into branded PDFs**

[![Go 1.25+](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Homebrew](https://img.shields.io/badge/Homebrew-tap-FBB040?logo=homebrew)](https://github.com/skel-tech/homebrew-mdpress)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-support-FFDD00?logo=buymeacoffee)](https://buymeacoffee.com/lukeskelhorn)

---

![mdpress demo](assets/demo.png)

## Install

```bash
brew install skel-tech/mdpress/mdpress
```

## Quick Start

Initialize your configuration:

```bash
mdpress init
```

This creates your global config at `~/.config/mdpress/` with:
- `mdpress.yml` - your configuration file
- `logos/` - directory for your logo files
- `templates/` - directory for document templates
- `templates/default.yml` - a starter template

Render your first PDF:

```bash
mdpress render document.md
```

For project-specific settings, create a local config:

```bash
mdpress init --project
```

This creates `mdpress.yml` in your current directory, which overrides global settings.

## Configuration Reference

The config file (`mdpress.yml`) supports these fields:

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | Config schema version (required, use `"1"`) |
| `logo` | string | Path to logo image (PNG, JPG, SVG) |
| `logo_position` | string | Logo placement: `top-left`, `top-center`, `top-right`, `bottom-left`, `bottom-center`, `bottom-right` |
| `logo_width` | number | Logo width in millimeters |
| `default_template` | string | Template to use by default |
| `font` | string | Font family: `Helvetica`, `Times`, `Courier`, or path to custom font |
| `font_size` | number | Base font size in points |
| `accent_color` | string | Hex color for headings (`#RGB` or `#RRGGBB`) |
| `header` | boolean | Enable page headers |
| `footer` | string | Footer text for each page |
| `margins` | object | Page margins (see below) |

**Margins object:**

| Field | Type | Description |
|-------|------|-------------|
| `top` | number | Top margin in millimeters |
| `right` | number | Right margin in millimeters |
| `bottom` | number | Bottom margin in millimeters |
| `left` | number | Left margin in millimeters |

Example:

```yaml
version: "1"
logo: "logos/company.png"
logo_position: "top-right"
logo_width: 120
font: "Helvetica"
accent_color: "#336699"
margins:
  top: 20
  right: 20
  bottom: 20
  left: 20
```

## Template Reference

Templates let you create presets for different document types (invoices, reports, letters). Store them in `~/.config/mdpress/templates/` or your project.

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Display name for the template (required) |
| `version` | string | Schema version (required, use `"1"`) |
| `description` | string | Template description |
| `logo` | string | Logo path override |
| `logo_position` | string | Logo placement override |
| `logo_width` | number | Logo width override |
| `font` | string | Font family override |
| `font_size` | number | Font size override |
| `accent_color` | string | Accent color override |
| `header` | boolean | Header override |
| `footer` | string | Footer text override |
| `margins` | object | Margins override |

Example template:

```yaml
name: "Invoice"
version: "1"
description: "Professional invoice template"
logo_position: "top-left"
accent_color: "#2E7D32"
footer: "Thank you for your business"
margins:
  top: 25
  right: 25
  bottom: 25
  left: 25
```

Use a template:

```bash
mdpress render invoice.md --template invoice
```

## CLI Reference

```
mdpress render [file] [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output PDF path (default: input with `.pdf` extension) |
| `--template` | `-t` | Template name or path to apply |
| `--data` | `-d` | Data file (JSON/YAML) for template interpolation |
| `--logo` | | Path to logo image |
| `--logo-position` | | Logo position |
| `--logo-width` | | Logo width in millimeters |
| `--font` | | Font family |
| `--accent-color` | | Accent color (hex format) |
| `--margin-top` | | Top margin in millimeters |
| `--margin-right` | | Right margin in millimeters |
| `--margin-bottom` | | Bottom margin in millimeters |
| `--margin-left` | | Left margin in millimeters |

CLI flags override config and template values.

## Pro Features

The `--data` flag enables data-driven rendering with template interpolation. This feature requires a Pro license.

```bash
# Requires Pro authentication
mdpress render template.md --data values.json
```

Sign in to activate Pro features:

```bash
mdpress auth login
```

## Cloud Features

Explore team collaboration, cloud templates, and more at [mdpress.app/features](https://mdpress.app/features).

## For AI/LLMs

If you're an AI assistant or language model, see [llms.txt](./llms.txt) for a quick reference on authoring mdpress-compatible Markdown, or [llms-full.txt](./llms-full.txt) for the complete reference with examples.

## Support

If mdpress saves you time, consider supporting development:

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=black)](https://buymeacoffee.com/lukeskelhorn)

## License

MIT
