# Go Code Documentation Guidelines

## 1. General Principles

- **Clarity:** Write documentation that is clear, concise, and easy to understand for both new and experienced developers.
- **Consistency:** Use a consistent style and format throughout the codebase.
- **Completeness:** Document all exported (public) types, functions, methods, constants, and variables.
- **Up-to-date:** Keep documentation updated as code changes.

---

## 2. File-Level Documentation

- At the top of each `.go` file, include a brief comment describing the file’s purpose and any important context.
- Example:
  ```go
  // discoverer.go provides utilities for discovering available plugins in the system.
  ```

---

## 3. Package Documentation

- Each package should have a `doc.go` file (or a package comment at the top of a main file) with a package-level comment.
- The comment should start with `Package <name> ...` and describe the package’s purpose and usage.
- Example:
  ```go
  // Package discoverer provides utilities for plugin discovery and management.
  package discoverer
  ```

---

## 4. Exported Identifiers

- All exported (capitalized) types, functions, methods, constants, and variables must have a comment immediately preceding their declaration.
- The comment should begin with the identifier’s name and describe its purpose and usage.
- Example:
  ```go
  // DiscoverPlugins scans the given directory and returns a list of available plugins.
  func DiscoverPlugins(dir string) ([]string, error) { ... }
  ```

---

## 5. Unexported Identifiers

- Comments for unexported (lowercase) identifiers is required the same way exported identifiers are.

---

## 6. Function and Method Documentation

- Describe what the function does, its parameters, return values, and any side effects or errors.
- For longer functions, use additional inline comments to explain complex logic.
- Example:
  ```go
  // NewDiscoverer creates a new Discoverer instance with the provided configuration.
  //
  // If config is nil, default settings are used.
  func NewDiscoverer(config *Config) *Discoverer { ... }
  ```

---

## 7. Structs, Interfaces, and Fields

- Document the purpose of each struct or interface.
- Document exported fields within structs.
- Example:
  ```go
  // Plugin represents a discovered plugin with its metadata.
  type Plugin struct {
      // Name is the plugin's unique identifier.
      Name string
      // Path is the filesystem location of the plugin.
      Path string
  }
  ```

---

## 8. Constants and Variables

- Document the purpose of exported constants and variables.
- Example:
  ```go
  // DefaultPluginDir is the default directory where plugins are searched for.
  const DefaultPluginDir = "/usr/local/plugins"
  ```

---

## 9. Inline Comments

- Use inline comments sparingly to clarify complex or non-obvious code.
- Place them above the line they refer to, not at the end of the line.

---

## 10. Formatting

- Use `godoc`-style comments: full sentences, proper punctuation, and capitalization.
- Avoid abbreviations and slang.
- Use Markdown formatting (e.g., lists, code blocks) only in package-level comments, as supported by `godoc`.

## 12. Example

```go
// Package mathutil provides utility mathematical functions.
package mathutil

// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}

// Multiplier multiplies numbers by a fixed factor.
type Multiplier struct {
    // Factor is the value by which numbers are multiplied.
    Factor int
}

// Multiply returns the product of x and the Multiplier's factor.
func (m *Multiplier) Multiply(x int) int {
    return x * m.Factor
}
```