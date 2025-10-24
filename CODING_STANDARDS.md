# DSH Coding Standards

## Core Principles

### Clean Code
- **Meaningful Names**: Functions and variables should clearly express intent
- **Small Functions**: Max 20 lines, single responsibility
- **No Magic Numbers**: Use named constants
- **Clear Error Messages**: Include context and actionable information

### KISS (Keep It Simple, Stupid)
- **Minimal Implementation**: Solve the problem, nothing more
- **Avoid Premature Optimization**: Correct first, fast second
- **Simple Data Structures**: Prefer basic types over complex abstractions
- **Linear Flow**: Avoid deep nesting (max 3 levels)

### DRY (Don't Repeat Yourself)
- **Extract Common Logic**: Shared functionality in separate functions
- **Configuration Constants**: Define once, use everywhere
- **Error Handling Patterns**: Consistent error handling approach
- **No Copy-Paste**: If you copy code twice, extract it

### Clean Architecture
- **Separation of Concerns**: Lexer → Parser → Executor
- **Dependency Direction**: Core logic doesn't depend on I/O
- **Interface Boundaries**: Clear contracts between components
- **Testable Units**: Each component can be tested in isolation

## Go-Specific Standards

### Formatting and Linting
- **gofmt**: All code must be formatted with `gofmt`
- **golangci-lint**: Use existing `.golangci.yml` configuration for linting
- **Pre-commit**: Run `make fmt lint` before every commit
- **CI Integration**: Linting failures block merges

### golangci-lint Configuration
The project uses `.golangci.yml` with comprehensive linting rules:

**Enabled Features:**
- **All linters by default** with strategic exclusions
- **Dependency guard** prevents problematic imports
- **Exhaustive checking** for switch statements
- **Security scanning** with gosec (disabled in tests)
- **Code complexity** and duplication detection

**Key Exclusions:**
- `lll` (line length) - intentionally excluded
- `mnd` (magic numbers) - evaluated, not worth it
- `nonamedreturns` - intentionally excluded
- `tagliatelle` - evaluated, not worth it
- `testpackage` - intentionally excluded

**Test-Specific Rules:**
- Relaxed `funlen`, `goconst`, `gosec` in test files
- Allows longer functions and constants in tests

### Naming Conventions
- **Packages**: lowercase, single word
- **Functions**: camelCase, verb-based for actions
- **Types**: PascalCase
- **Constants**: PascalCase or UPPER_CASE for exported
- **Variables**: camelCase, descriptive

## Quality Gates

### Before Every Commit
```bash
make fmt     # Format code with gofmt
make lint    # Run golangci-lint with .golangci.yml
make test    # Run all tests
```

### Linting Rules Enforced
- **All linters enabled** by default with strategic exclusions
- **Dependency guard**: Prevents problematic package imports
- **Security scanning**: gosec enabled (relaxed in tests)
- **Code quality**: Complexity, duplication, unused code detection
- **Import standards**: Enforced import organization
- **Exhaustive checking**: Complete switch statement coverage

## Architecture Rules

### Layer Separation
1. **Input Layer** (`main.go`): Handle user input/output
2. **Parsing Layer** (`lexer.go`, `parser.go`): Convert text to structures  
3. **Execution Layer** (`executor.go`, `builtins.go`): Execute commands
4. **System Layer**: OS interactions (contained within execution)

### Dependencies
- Core logic (parser, lexer) has no OS dependencies
- Only executor layer calls OS functions
- No circular dependencies between layers

## Security Standards

### Input Validation
- **Validate all user input** before processing
- **Sanitize file paths** to prevent directory traversal
- **Limit input size** to prevent DoS
- **Escape shell metacharacters** in external commands

### gosec Integration
The `.golangci.yml` enables comprehensive security scanning:
- **Enabled in production code** for security vulnerability detection
- **Disabled in test files** to allow test-specific patterns
- **Dependency restrictions** prevent problematic imports
- **Standard library preference** enforced via depguard rules

## Development Workflow

### Setup
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verify configuration
golangci-lint config path  # Should show .golangci.yml
```

### Daily Workflow
```bash
# Before coding
git pull origin main

# During development
make fmt lint    # Fix formatting and linting issues

# Before commit
make test        # Ensure all tests pass
git add .
git commit -m "Add feature: implement X"
```

### CI/CD Integration
The `.golangci.yml` configuration ensures:
- Consistent linting across all environments
- Same rules for all contributors
- Automated quality checks in CI pipeline
- Blocking of non-compliant code

## Error Handling Standards

### Pattern
```go
// Good: Specific error with context
if err := os.Chdir(path); err != nil {
    return fmt.Errorf("cd: cannot change to %s: %w", path, err)
}
```

### errcheck Compliance
- **No ignored errors**: All function returns must be checked
- **Explicit error handling**: Use `_ = err` if intentionally ignored
- **Error wrapping**: Use `fmt.Errorf` with `%w` verb for context
