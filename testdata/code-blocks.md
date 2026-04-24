# Code Blocks

## Inline Code

Use the `fmt.Println()` function to print output. The `main` package is the entry point.

## Go

```go
package main

import "fmt"

func main() {
    items := []string{"markdown", "pdf", "cli"}
    for i, item := range items {
        fmt.Printf("%d. %s\n", i+1, item)
    }
}
```

## JavaScript

```javascript
async function fetchData(url) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }
  return response.json();
}
```

## Python

```python
def fibonacci(n):
    if n <= 1:
        return n
    a, b = 0, 1
    for _ in range(2, n + 1):
        a, b = b, a + b
    return b
```

## YAML

```yaml
version: "1"
font: "Helvetica"
margins:
  top: 20
  right: 20
  bottom: 20
  left: 20
```

## Plain Code Block

```
No language specified.
This is a plain code block.
It should render in a monospace font.
```
