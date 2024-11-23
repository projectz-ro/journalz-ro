
# **Style Guide for JournalZ-ro**

## **Table of Contents**
- [File Organization](#file-organization)
- [Naming Conventions](#naming-conventions)
- [Code Formatting](#code-formatting)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Documentation](#documentation)
- [Security](#security)

---

## **File Organization**

```
    cmd
   └  journalz-ro
     │  find.go              
     │  main.go              
     └  new.go              
   󰉖 config
    internal
   │  core
   │ │  entries
   │ │ │  entries.go        
   │ │ └  entries_test.go    
   │ │  volumes
   │ │ │  volumes.go         
   │ │ └  volumes_test.go    
   │ └  errors.go            
   │  database
   │ │  sqlite
   │ │ │ 󰉖 migrations
   │ │ │  migrations.go      
   │ │ │  schema.sql        
   │ │ └  sqlite.go          
   │ │  database.go          
   │ └  errors.go            
   └  ui
     │  components           
     │  input                
     │  views                
     │  errors.go            
     │  keybindings.go       
     │  renderer.go          
     └  ui.go                
    pkg
   └ 󰉖 utils
```

---

## **Naming Conventions**

### **Files and Directories**
- **Directories**: Use lowercase and hyphenated words, e.g., `cmd`, `config`, `internal`, `ui`, `pkg`.
- **Files**: Use lowercase and hyphenated naming for consistency. Avoid using underscores in file names (e.g., `find.go` instead of `find_entry.go`).
- **camelCase**: Used for private variables and functions
- **PascalCase**: Used for types, structs and public variables and functions.

### **Package Naming**
- **Packages** should have descriptive names that indicate their purpose (e.g., `core`, `ui`, `database`, `utils`).
- Avoid using names that are too general, like `main` or `handler`, to ensure clarity.

---

## **Code Formatting**

### **General Formatting**
- **Indentation**: Use **tabs** for indentation, not spaces.
- **Line Length**: Keep lines to **80 characters or fewer** for better readability.
- **Spacing**:
  - Place a blank line between top-level functions or methods.
  - Add a blank line between function parameters and the body of the function.

### **Imports**
- Import groups should be organized into the following order:
  - Standard library imports (grouped together).
  - Third-party imports (grouped together).
  - Internal imports (grouped together).
  
  Example:
  ```go
  import (
      "fmt"
      "log"

      "github.com/spf13/viper"

      "yourapp/internal/core"
      "yourapp/internal/database"
  )
  ```

### **Comments**
- Use **full sentences** and keep comments clear, concise, and grammatically correct.
- Add comments before **complex logic** and **function definitions**.
- Avoid obvious comments, such as `// Increment counter`, when the code is self-explanatory.
- For exported functions and types, use **doc comments**.

---

## **Error Handling**

### Error Organization
1. **Error Naming Convention**:
   - Use the pattern: `Err[Subject][SubjectPart][Violation]`.
   - Ensure that errors are **actionable** and contain useful information for debugging.
   - Examples:
     ```go
     ErrEntryContentTooLong   // Subject: Entry, Violation: TooLong
     ErrVolumeTitleEmpty      // Subject: Volume, Violation: Empty
     ```

2. **Error Wrapping**:
   - Use `fmt.Errorf()` with `%w` to wrap errors and preserve context.

### Example:
```go
if err := validateTags(tags); err != nil {
    return fmt.Errorf("failed to validate tags: %w", err)
}
```

---

## **Testing**

### **General Guidelines**
- Use the **testing** package for unit and integration tests.
- Tests should be written alongside the code they test (i.e., in the same package).
- Write tests for **all exported functions**, especially when they have non-trivial logic or affect the program state.
- **Test for edge cases**, such as empty inputs, invalid data, and unexpected behavior.

### Test Structure
1. **Test File Location**: 
   - Like errors.go every package should have a tests.go 
2. **Test Function Naming Convention**:
   - Use descriptive names: `TestSubject_Scenario` or `TestSubject_Scenario_Condition`.
   - Examples:
     ```go
     func TestCreateEntry_ValidTags(t *testing.T)
     func TestCreateVolume_EmptyTitle(t *testing.T)
     ```
### Table-Driven Tests
- Use table-driven tests for cases with multiple input scenarios.

### Example:
```go
func TestValidateTags(t *testing.T) {
    tests := []struct {
        name    string
        tags    []string
        wantErr bool
    }{
        {"valid tags", []string{"tag1", "tag-2"}, false},
        {"invalid tag", []string{"tag 1", "tag!@#"}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := ValidateTags(tt.tags); (err != nil) != tt.wantErr {
                t.Errorf("ValidateTags() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### **Unit Testing**
- **Mock dependencies** for unit tests (e.g., mocking the database calls).
- Test **functions** in isolation to ensure each part of the application is functioning independently.
  
### **Integration Testing**
- For tests that require database interaction, use an **in-memory SQLite** database to test database operations.
- Perform **end-to-end testing** for commands like `find` and `new` to ensure they work together as expected.

### **Bubble Tea UI Testing**
- For testing Bubble Tea-based UI components, simulate user interactions using Bubble Tea’s test helpers (e.g., sending simulated key events to the model).

---

## **Documentation**

### **General Documentation**
- Maintain up-to-date documentation in `README.md` for setup, usage instructions, and how to contribute.
- Provide detailed **doc comments** for exported functions and types to describe their purpose and behavior.

### **In-code Documentation**
- Use **godoc**-style comments above exported functions, variables, and structs.
  Example:
  ```go
  // CreateNewEntry creates a new journal entry and stores it in the database.
  func CreateNewEntry(entry Entry) error {
      // Implementation
  }
  ```

---

## **Security**

### **Data Validation**
- **Sanitize inputs** and **validate user inputs** (e.g., tags cannot have spaces or special characters).
- Make sure that the user cannot inject malicious data through inputs like tags or entry titles.

### **Sensitive Data Handling**
- Avoid storing sensitive information in plaintext (e.g., passwords, private keys).
- Use secure methods for storing sensitive data and ensure it's properly encrypted if necessary.

### **Error Disclosure**
- Do not disclose sensitive information in error messages.
- Use generalized error messages in production and log full error details on the server for debugging purposes.

### **Dependencies**
- Regularly update and audit dependencies to prevent known vulnerabilities.
- Use tools like `go list -m all` to check for outdated dependencies.

### Input Validation and Sanitization
- **Validation**: Use a centralized validation function to check inputs such as tags and content length.
  
  Example:
  ```go
  // ValidateTags checks that each tag is valid (no special characters or spaces).
  func ValidateTags(tags []string) error {
      for _, tag := range tags {
          if len(tag) == 0 || strings.ContainsAny(tag, " !@#$%^&*()") {
              return ErrInvalidTag
          }
      }
      return nil
  }
  ```

- **Sanitization**: Ensure that paths, filenames, and any user-generated content are sanitized before being used.

  Example:
  ```go
  // SanitizeFilePath ensures that the file path is safe for use.
  func SanitizeFilePath(path string) string {
      return filepath.Clean(path)
  }
  ```

### File Operations
- **Ensure Safe File Operations**: Validate the existence of files and directories, and sanitize filenames to avoid directory traversal vulnerabilities.

  Example:
  ```go
  func ensureSafeFile(filename string) error {
      if _, err := os.Stat(filename); err == nil {
          return ErrFileExists
      }

      if !hasValidExtension(filename) {
          return ErrInvalidFileExtension
      }

      dir := filepath.Dir(filename)
      if err := os.MkdirAll(dir, 0755); err != nil {
          return fmt.Errorf("creating directory: %w", err)
      }

      return nil
  }
 
---
 ```
