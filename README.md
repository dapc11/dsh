# DSH - Daniel's Shell

A minimal, secure, POSIX-compatible shell implementation written in Go.

## Features

### Phase 1 ✅ (Complete)
- Basic command execution
- Built-in commands (`cd`, `exit`, `help`, `pwd`)
- Simple input parsing
- Process management

### Phase 2 ✅ (Complete)
- **Quote handling** - Single and double quotes with escape sequences
- **I/O redirection** - Output (`>`), input (`<`), append (`>>`)
- **Background processes** - Command execution with `&`
- **Command chaining** - Multiple commands with `;`
- **Comment support** - Lines starting with `#`
- **Enhanced parsing** - Proper lexer/parser architecture
- **Clean code standards** - 0 linting issues, comprehensive error handling

### Phase 3 ✅ (Complete)
- **Enhanced Line Editing** - Emacs-like readline functionality
  - Cursor movement (Ctrl+A/E, Ctrl+B/F, arrows) ✅
  - Word movement (Ctrl+←/→, Alt+D) ✅
  - Command history (↑/↓, Ctrl+P/N) ✅
  - Line editing (Ctrl+D/K/U/W, backspace) ✅
  - Screen control (Ctrl+L) ✅
  - Kill ring operations (Ctrl+Y, Alt+Y) ✅
- **Tab Completion** - Interactive command and file completion
  - Command completion from PATH ✅
  - File and directory completion ✅
  - Interactive menu navigation ✅
  - Fuzzy matching support ✅
- **History Features** - Advanced history management
  - Persistent command history ✅
  - History navigation (↑/↓) ✅
  - Fuzzy history search (Ctrl+R) ✅
  - Auto-suggestions based on history ✅

### Phase 4 🚧 (In Progress)
- Pipeline support (`|`)
- Job control and signal handling
- Variable expansion (`$VAR`, `${VAR}`)
- Command substitution (`$(command)`)
- Globbing and pathname expansion
- Control structures (if/then/else, loops)

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd dsh

# Build the shell
make build

# Or use Go directly
go build -o dsh .
```

## Usage

```bash
# Start the shell
./dsh

# Example commands
dsh> echo "Hello, World!"
dsh> ls -la > output.txt
dsh> cat < input.txt
dsh> sleep 5 &
dsh> echo "first"; echo "second"
dsh> # This is a comment
dsh> pwd
dsh> cd /tmp
dsh> exit
```

## Development

### Prerequisites
- Go 1.21+
- golangci-lint

### Code Quality
```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Build
make build
```

### Architecture
- **Lexer** (`internal/lexer/`) - Tokenizes input with quote and escape handling
- **Parser** (`internal/parser/`) - Parses tokens into command structures  
- **Executor** (`internal/executor/`) - Executes commands with I/O redirection
- **Built-ins** (`internal/builtins/`) - Built-in command implementations
- **Readline** (`internal/readline/`) - Emacs-like line editing with history

## Documentation

- [`CODING_STANDARDS.md`](CODING_STANDARDS.md) - Code quality guidelines
- [`BEST_PRACTICES.md`](BEST_PRACTICES.md) - Shell implementation best practices
- [`DSH_PERSONA.md`](DSH_PERSONA.md) - Project development persona
- [`Sh.html`](Sh.html) - POSIX shell reference tutorial

## POSIX Compliance

DSH aims for POSIX shell compatibility with focus on:
- Standard command execution
- Proper quote handling (strong vs weak quoting)
- I/O redirection operators
- Background process execution
- Environment variable handling
- Exit status propagation

## Security

- Input validation at all boundaries
- Safe process execution without injection vulnerabilities
- Proper file descriptor management
- Resource cleanup and error handling
- Security-focused linting with gosec

## Contributing

1. Follow the coding standards in `CODING_STANDARDS.md`
2. Ensure all tests pass: `make test`
3. Run linter: `make lint` (must have 0 issues)
4. Format code: `make fmt`
5. Update documentation as needed

## License

[Add your license here]

## Status

**Current Phase:** 3 Complete ✅  
**Code Quality:** Minor race condition in tests ⚠️  
**POSIX Features:** Basic compliance ✅  
**Production Ready:** Phase 3 features ✅
