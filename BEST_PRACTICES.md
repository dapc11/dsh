# Shell Implementation Best Practices

## Development Workflow (CRITICAL)

### Every Code Change Must Follow:
1. **Format immediately**: `gofmt -w .`
2. **Lint immediately**: `golangci-lint run` 
3. **Test immediately**: `go test ./...`
4. **Fix ALL issues before proceeding**

### Test-Driven Development:
- **Write tests FIRST** for new functionality
- **Run tests frequently** during development
- **Achieve comprehensive coverage** for new code
- **Fix race conditions immediately**

### Quality Automation:
```bash
# After every change
make fmt lint test

# Before every commit  
make test-all
golangci-lint run  # Must show zero issues
gofmt -d .         # Must show no output
```

## Security & Safety
- [x] Validate input and sanitize commands
- [ ] Implement proper signal handling (SIGINT, SIGTERM, SIGCHLD)
- [x] Use `execve()` family functions correctly to prevent injection
- [x] Handle file descriptors properly to avoid leaks
- [ ] Implement proper process group management
- [x] Prevent command injection attacks
- [ ] Handle setuid/setgid programs safely

## Performance & Resource Management
- [x] Use efficient parsing algorithms
- [x] Implement proper memory management
- [ ] Handle large command outputs gracefully
- [ ] Avoid blocking on I/O operations
- [x] Implement job control for background processes
- [x] Optimize for common use cases
- [x] Implement proper cleanup on exit

## Error Handling
- [x] Provide meaningful error messages
- [x] Handle edge cases (empty input, invalid commands)
- [x] Implement proper cleanup on errors
- [x] Handle system call failures gracefully
- [x] Return appropriate exit codes
- [ ] Handle out-of-memory conditions

## Architecture
- [x] Separate lexer, parser, and executor components
- [x] Use state machines for complex parsing
- [x] Implement modular built-in command system
- [x] Design for extensibility
- [x] Maintain clean separation of concerns
- [x] Use proper abstraction layers

## User Experience (Phase 3 Complete)
- [x] Emacs-like line editing with cursor movement
- [x] Command history with navigation and search
- [x] Tab completion for commands and files
- [x] Interactive completion menus
- [x] Kill ring operations (cut/copy/paste)
- [x] Auto-suggestions based on history
- [x] Fuzzy search in history (Ctrl+R)
- [x] Word movement and editing operations

## Testing
- [x] Unit tests for individual components
- [x] Integration tests for complete workflows
- [ ] Test with edge cases and malformed input
- [ ] Verify POSIX compliance with standard test suites
- [ ] Test signal handling and process management
- [ ] Performance testing with large inputs
- [ ] Fix race condition in executor tests

## Code Quality
- [x] Zero linting issues (golangci-lint)
- [x] Proper Go formatting (gofmt)
- [x] Comprehensive error handling
- [x] Clean separation of concerns
- [x] Modular architecture
- [ ] Fix test race conditions
