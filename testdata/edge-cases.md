# Edge Cases

## Empty Sections

### This Section is Empty

### This One Has Content

Some content here.

## Very Long Paragraph

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.

## Special Characters

Ampersand: &
Less than: <
Greater than: >
Quotes: "double" and 'single'
Copyright: (c) 2026
Em dash: ---
Ellipsis: ...

## Long Table Cell Content

| Short | This column has a much longer header title that might wrap |
|-------|-----------------------------------------------------------|
| A     | This cell contains a longer piece of text to test how the table handles wrapping and overflow in PDF output |
| B     | Short cell |

## Consecutive Code Blocks

```go
func first() {}
```

```go
func second() {}
```

## Deep Nesting

1. Level 1
   1. Level 2
      1. Level 3
         - Level 4 unordered
         - Another level 4

## Single Item List

- Just one item

## Unicode

Emoji: Hello World
Accented: café, résumé, naïve
CJK: 你好世界
Arabic: مرحبا

## Very Long Code Line

```
This is a single line of code that is intentionally very long to test how the PDF renderer handles horizontal overflow in code blocks without wrapping or truncation issues in the output
```
