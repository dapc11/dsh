# DSH Project Work Persona

## Project Identity
**Name:** DSH (Daniel's Shell)  
**Mission:** Build a minimal, secure, POSIX-compliant shell from scratch in Go  
**Philosophy:** Simplicity, correctness, and standards compliance over feature bloat

## Knowledge Base References
- **POSIX Tutorial:** `./Sh.html` - Grymoire's comprehensive POSIX shell tutorial
- **Best Practices:** `./BEST_PRACTICES.md` - Our curated shell implementation guidelines
- **Current Implementation:** `./main.go` - Minimal working shell foundation

## Core Principles

### 1. POSIX First
- Every feature must align with POSIX shell specification
- Reference `./Sh.html` for canonical behavior examples
- Test against standard POSIX shell test suites
- When in doubt, match bash/dash behavior in POSIX mode

### 2. Security by Design
- Input validation at every boundary
- Proper signal handling (SIGINT, SIGTERM, SIGCHLD)
- Safe process execution without injection vulnerabilities
- Resource cleanup and file descriptor management

### 3. Go Best Practices
- Leverage Go's concurrency for process management
- Use standard library (`os/exec`, `syscall`) appropriately
- Implement proper error handling with meaningful messages
- Follow Go formatting and linting standards
- **ALWAYS run `golangci-lint run` before commits** using `.golangci.yml`
- **ALWAYS format with `gofmt`** for consistency

## Development Workflow

### MANDATORY: Every Code Change Must Follow This Process

#### 1. Before Writing Any Code:
- [ ] **Run existing tests** to ensure baseline: `make test`
- [ ] **Check current lint status**: `golangci-lint run`
- [ ] **Verify formatting**: `gofmt -d .` (should show no output)

#### 2. During Implementation:
- [ ] **Write tests FIRST** for new functionality (TDD approach)
- [ ] **Run tests frequently** during development: `go test ./...`
- [ ] **Format code after every change**: `gofmt -w .`
- [ ] **Check linting after every change**: `golangci-lint run`
- [ ] **Fix ALL lint warnings immediately** - never commit with warnings

#### 3. Before Every Commit:
- [ ] **Run full test suite**: `make test`
- [ ] **Ensure zero lint issues**: `golangci-lint run` (must show no output)
- [ ] **Verify formatting**: `gofmt -d .` (must show no output)
- [ ] **Fix any race conditions** in tests
- [ ] **Manual testing** of new functionality
- [ ] **No regression** in existing features

#### 4. Quality Gates (Non-Negotiable):
- **ZERO** linting issues allowed
- **ALL** tests must pass
- **NO** race conditions in tests
- **PROPER** formatting (gofmt)
- **COMPREHENSIVE** test coverage for new code

### Current Implementation Status

#### Phase 1-3: Complete ✅
- Core command execution, I/O redirection, background processes
- Quote handling, command chaining, comment support
- Emacs-like line editing with full key bindings
- Tab completion with interactive menus
- History management with fuzzy search
- Kill ring operations and auto-suggestions

#### Phase 4: Next Priority
- Pipeline support (|)
- Variable expansion ($VAR, ${VAR})
- Command substitution $(command)
- Job control and signal handling
- Globbing and pathname expansion

## Decision Framework

### When implementing new features:
1. **Check POSIX:** Does `./Sh.html` cover this behavior?
2. **Security Review:** Does this introduce vulnerabilities?
3. **Best Practices:** Does this align with `./BEST_PRACTICES.md`?
4. **Minimal Implementation:** What's the simplest correct approach?

### When debugging:
1. **Reference Behavior:** How does the tutorial say this should work?
2. **Test Against Standards:** Does our behavior match POSIX shells?
3. **Security Implications:** Could this be exploited?
4. **Resource Management:** Are we cleaning up properly?

## Quality Gates

### CRITICAL: Never Skip These Steps

#### Every Code Change:
1. **Format immediately**: `gofmt -w .`
2. **Lint immediately**: `golangci-lint run`
3. **Test immediately**: `go test ./...`
4. **Fix ALL issues before continuing**

#### Every Commit:
1. **Full test suite**: `make test` (must pass 100%)
2. **Zero lint issues**: `golangci-lint run` (no output)
3. **Perfect formatting**: `gofmt -d .` (no output)
4. **No race conditions**: Tests must pass with `-race` flag
5. **Manual verification**: Test new functionality works

#### Every Feature:
1. **Write tests first** (TDD)
2. **Achieve good test coverage**
3. **Document behavior changes**
4. **Verify POSIX compliance** where applicable

### Automation Commands
```bash
# Quick quality check (run after every change)
make fmt lint test

# Full verification (run before commit)
make test-all

# Individual steps
gofmt -w .
golangci-lint run
go test -race ./...
```

## Communication Style
- **Technical:** Precise, standards-focused
- **Pragmatic:** Favor working solutions over perfect theory
- **Security-conscious:** Always consider attack vectors
- **Educational:** Reference tutorial examples when explaining behavior

## Key Mantras
1. **"Format, Lint, Test"** - After every single change
2. **"Red, Green, Refactor"** - Write failing tests first
3. **"Zero tolerance for warnings"** - Fix lint issues immediately
4. **"What does POSIX say?"** - Always check the standard first
5. **"Security by default"** - Validate everything, trust nothing
6. **"Test early, test often"** - Verify behavior continuously
