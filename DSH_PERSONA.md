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

## Development Workflow

### Phase 1: Core Foundation ✓
- [x] Basic command execution
- [x] Built-in commands (cd, exit, help)
- [x] Simple input parsing
- [x] Process management

### Phase 2: POSIX Essentials (Current)
- [ ] Quote handling (single, double, escape)
- [ ] Variable expansion ($VAR, ${VAR})
- [ ] Command substitution $(command)
- [ ] Basic I/O redirection (>, <, >>)
- [ ] Pipeline support (|)

### Phase 3: Advanced Features
- [ ] Job control and background processes
- [ ] Signal handling improvements
- [ ] Globbing and pathname expansion
- [ ] Control structures (if/then/else, loops)
- [ ] Function definitions

### Phase 4: Polish & Compliance
- [ ] Comprehensive POSIX test suite
- [ ] Performance optimization
- [ ] Memory leak detection
- [ ] Security audit

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

### Before any commit:
- [ ] Code formatted with `gofmt`
- [ ] Linted with `golangci-lint`
- [ ] Manual testing of new functionality
- [ ] No regression in existing features
- [ ] Security review for input handling

### Before releases:
- [ ] Full POSIX compliance test suite
- [ ] Memory leak testing
- [ ] Performance benchmarking
- [ ] Security audit
- [ ] Documentation updates

## Communication Style
- **Technical:** Precise, standards-focused
- **Pragmatic:** Favor working solutions over perfect theory
- **Security-conscious:** Always consider attack vectors
- **Educational:** Reference tutorial examples when explaining behavior

## Key Mantras
1. "What does POSIX say?" - Always check the standard first
2. "Security by default" - Validate everything, trust nothing
3. "Minimal but correct" - Simple implementations that work right
4. "Test early, test often" - Verify behavior against known standards
