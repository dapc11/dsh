# DSH Integration Tests

This directory contains integration tests for DSH that verify end-to-end functionality and prevent regressions during refactoring.

## Test Strategy

### Timeout Protection
- All tests use `context.WithTimeout` to prevent hanging
- Default timeout: 3 seconds per command
- Tests fail fast if commands hang

### Test Categories

1. **Core Commands** - Basic shell functionality (echo, pwd, help)
2. **Command Chaining** - Semicolon-separated commands
3. **File Redirection** - I/O redirection operators
4. **Quote Handling** - Single and double quote processing
5. **Error Handling** - Invalid commands and file operations
6. **Workflow Integration** - Realistic usage scenarios

### Non-Interactive Testing
- Uses `dsh -c "command"` for deterministic testing
- Avoids interactive features that are hard to test
- Focuses on core shell functionality

## Running Tests

```bash
# Run all integration tests
make test-integration

# Run with verbose output
go test -v ./test/integration/

# Run with custom timeout
go test -timeout 30s ./test/integration/

# Run specific test
go test -run TestCoreCommands ./test/integration/
```

## Test Design Principles

1. **Fast and Reliable** - Tests complete quickly and don't flake
2. **Isolated** - Each test uses temporary directories
3. **Focused** - Tests specific functionality without complex setup
4. **Timeout Protected** - Prevents hanging in CI/CD environments
5. **Realistic** - Tests actual user workflows

## Adding New Tests

When adding new DSH features:

1. Add integration tests for the new functionality
2. Use timeout protection for all command execution
3. Test both success and error cases
4. Include the feature in workflow integration tests
5. Ensure tests are deterministic and don't depend on external state

## Limitations

- Interactive features (tab completion, history) require separate testing approaches
- Some commands may timeout in certain environments
- Tests focus on core functionality rather than edge cases
