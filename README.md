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

### Phase 3 🚧 (Planned)
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
- **Lexer** (`lexer.go`) - Tokenizes input with quote and escape handling
- **Parser** (`parser.go`) - Parses tokens into command structures
- **Executor** (`main.go`) - Executes commands with I/O redirection
- **Built-ins** - Integrated built-in command implementations

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

**Current Phase:** 2 Complete ✅  
**Code Quality:** 0 linting issues ✅  
**POSIX Features:** Basic compliance ✅  
**Production Ready:** Phase 2 features ✅
