# DSH Rendering Tests

Automated tests for terminal rendering functionality that validate ANSI escape sequences and visual behavior without requiring interactive sessions.

## Approach

Instead of testing interactive terminal sessions (which hang), these tests validate:

1. **ANSI Sequence Generation** - Ensures correct escape sequences are produced
2. **Rendering Logic** - Tests the visual output patterns 
3. **Sequence Validation** - Verifies sequences match expected patterns
4. **Diagnostic Information** - Provides detailed output for debugging

## Test Categories

### ANSI Sequence Tests
- Cursor movement (left, right, up, down, home)
- Screen operations (clear line, clear screen)
- Cursor save/restore for tab completion
- Color codes (red, green, blue, bold, reset)

### Rendering Pattern Tests
- Prompt display validation
- Tab completion sequence patterns
- Line editing operations
- Color output validation

### Diagnostic Tests
- Sequence extraction and analysis
- Pattern matching validation
- Environment information logging

## Benefits

✅ **Fast Execution** - Tests complete in 3ms  
✅ **No Hanging** - Pure unit tests, no interactive sessions  
✅ **Detailed Diagnostics** - Logs exact sequences for debugging  
✅ **Pattern Validation** - Regex matching ensures correct format  
✅ **Comprehensive Coverage** - Tests all common terminal operations  

## Running Tests

```bash
# Run rendering tests
make test-rendering

# Run with verbose output for diagnostics
go test -v ./test/rendering/

# Run specific test category
go test -run TestANSISequenceGeneration ./test/rendering/
```

## Sample Output

```
✓ Move cursor left: "\x1b[D" matches "\\x1b\\[D"
✓ Tab completion: "\x1b[s\necho  exit  help\x1b[u"
✓ Color output: "\x1b[32mgreen\x1b[0m"
```

## Debugging Rendering Issues

When rendering issues occur:

1. **Check Sequence Generation** - Verify correct ANSI codes are produced
2. **Validate Patterns** - Ensure sequences match expected regex patterns  
3. **Review Diagnostics** - Examine detailed test output logs
4. **Test Individual Operations** - Run specific test categories

This approach provides confidence that DSH generates correct terminal sequences without the complexity of interactive testing.
