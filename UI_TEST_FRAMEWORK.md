# DSH UI Test Framework Specification

## Overview
Automated testing framework for DSH shell's user interface, leveraging existing MockTerminal infrastructure for deterministic, fast, and reliable UI/UX testing.

## Architecture

### Core Components
```
UITestFramework
├── InteractionRecorder    # Records user input sequences
├── OutputValidator       # Validates ANSI sequences and screen state
├── ScenarioRunner       # Executes test scenarios
└── AssertionEngine      # UI-specific assertions
```

### Test Structure
```go
type UITest struct {
    Name        string
    Setup       func(*UITestFramework)
    Scenario    []UIAction
    Assertions  []UIAssertion
    Cleanup     func(*UITestFramework)
}

type UIAction struct {
    Type     ActionType  // KeyPress, Wait, Validate
    Data     interface{} // KeyEvent, Duration, ValidationFunc
    Expected string      // Expected output pattern
}
```

## Requirements

### Functional Requirements
1. **Input Simulation**: Simulate keyboard input (keys, sequences, timing)
2. **Output Capture**: Capture all ANSI sequences and terminal output
3. **State Validation**: Validate cursor position, screen content, colors
4. **Interaction Recording**: Record and replay user interaction sequences
5. **Assertion Library**: Rich assertions for UI elements and behaviors

### Non-Functional Requirements
1. **Deterministic**: No timing dependencies, reproducible results
2. **Fast**: Tests should run in milliseconds, not seconds
3. **Isolated**: Each test runs in clean environment
4. **Readable**: Clear, expressive test syntax
5. **Maintainable**: Easy to add new test scenarios

## Implementation Plan

### Phase 1: Core Framework (Priority 1)
- [ ] Create `UITestFramework` struct with MockTerminal integration
- [ ] Implement `InteractionRecorder` for keystroke sequences
- [ ] Build `OutputValidator` for ANSI sequence validation
- [ ] Create basic assertion methods (`AssertCursorAt`, `AssertOutput`, `AssertColors`)

### Phase 2: Action System (Priority 1)
- [ ] Define `UIAction` types (KeyPress, KeySequence, Wait, Validate)
- [ ] Implement `ScenarioRunner` to execute action sequences
- [ ] Add timing simulation (optional delays between actions)
- [ ] Create action builders (`Type()`, `Press()`, `Expect()`)

### Phase 3: Rich Assertions (Priority 2)
- [ ] Screen content assertions (`AssertLineContains`, `AssertScreenState`)
- [ ] Menu and completion assertions (`AssertMenuVisible`, `AssertSelectedItem`)
- [ ] Color and styling assertions (`AssertTextColor`, `AssertHighlight`)
- [ ] Cursor and positioning assertions (`AssertCursorPosition`, `AssertPromptState`)

### Phase 4: Test Scenarios (Priority 2)
- [ ] Tab completion workflows
- [ ] History navigation (Ctrl+R, arrows)
- [ ] Line editing operations (Ctrl+A/E/K/U/W)
- [ ] Menu navigation and selection
- [ ] Multi-line input handling

### Phase 5: Advanced Features (Priority 3)
- [ ] Test recording from real sessions
- [ ] Visual diff for screen state changes
- [ ] Performance benchmarking for UI operations
- [ ] Regression testing against baseline outputs

## Test Categories

### 1. Tab Completion Tests
```go
func TestTabCompletionWorkflow() {
    test := UITest{
        Name: "Tab completion shows menu and navigates",
        Scenario: []UIAction{
            Type("e"),
            Press(KeyTab),
            AssertMenuVisible(),
            AssertMenuContains("echo", "exit"),
            Press(KeyArrowDown),
            AssertSelectedItem("exit"),
            Press(KeyEnter),
            AssertBuffer("exit"),
        },
    }
}
```

### 2. History Navigation Tests
```go
func TestHistorySearch() {
    test := UITest{
        Name: "Ctrl+R fuzzy history search",
        Setup: func(f *UITestFramework) {
            f.AddHistory("git status", "ls -la", "echo hello")
        },
        Scenario: []UIAction{
            Press(KeyCtrlR),
            AssertPromptContains("search:"),
            Type("git"),
            AssertSuggestion("git status"),
            Press(KeyEnter),
            AssertBuffer("git status"),
        },
    }
}
```

### 3. Line Editing Tests
```go
func TestLineEditingOperations() {
    test := UITest{
        Name: "Emacs-like line editing",
        Scenario: []UIAction{
            Type("hello world"),
            Press(KeyCtrlA),           // Move to beginning
            AssertCursorAt(0),
            Press(KeyCtrlK),           // Kill to end
            AssertBuffer(""),
            Press(KeyCtrlY),           // Yank back
            AssertBuffer("hello world"),
        },
    }
}
```

## File Structure
```
test/ui/
├── framework/
│   ├── ui_test_framework.go    # Main framework
│   ├── interaction_recorder.go # Input recording
│   ├── output_validator.go     # Output validation
│   ├── scenario_runner.go      # Test execution
│   └── assertions.go           # UI assertions
├── scenarios/
│   ├── tab_completion_test.go  # Tab completion tests
│   ├── history_test.go         # History navigation tests
│   ├── line_editing_test.go    # Line editing tests
│   └── menu_navigation_test.go # Menu interaction tests
└── fixtures/
    ├── test_history.txt        # Sample history data
    └── expected_outputs/       # Baseline outputs
```

## API Design

### Framework Initialization
```go
framework := NewUITestFramework()
framework.SetShellPrompt("dsh> ")
framework.LoadHistory("fixtures/test_history.txt")
```

### Test Definition
```go
test := framework.NewTest("Tab completion workflow").
    Type("e").
    Press(KeyTab).
    AssertMenuVisible().
    AssertMenuContains("echo", "exit").
    Press(KeyArrowDown).
    AssertSelectedItem("exit").
    Press(KeyEnter).
    AssertBuffer("exit")

framework.Run(test)
```

### Fluent Assertions
```go
framework.
    AssertOutput().Contains("echo").
    AssertCursor().IsAt(5, 0).
    AssertMenu().IsVisible().HasItems(3).
    AssertColors().HasGreen("echo").HasReverse("selected")
```

## Success Criteria

### Immediate Goals
1. Framework can simulate basic keyboard input
2. Framework can capture and validate ANSI output
3. At least 5 core UI scenarios are automated
4. Tests run in under 100ms each
5. Zero flaky tests (100% deterministic)

### Long-term Goals
1. 90%+ coverage of interactive features
2. Regression testing prevents UI breakage
3. New features include UI tests by default
4. Framework is reusable for other terminal applications

## Implementation Prompt

When implementing this framework, follow these principles:

1. **Start Simple**: Begin with basic keystroke simulation and output capture
2. **Build Incrementally**: Add one assertion type at a time
3. **Test the Framework**: Write tests for the test framework itself
4. **Keep It Fast**: Optimize for speed - no sleeps or waits
5. **Make It Readable**: Tests should read like user stories
6. **Leverage Existing Code**: Build on MockTerminal and BufferManager
7. **Follow DSH Patterns**: Use same coding standards and architecture

## Next Steps

1. Create `test/ui/framework/` directory structure
2. Implement basic `UITestFramework` with MockTerminal integration
3. Add simple keystroke simulation (`Type()`, `Press()`)
4. Create basic output assertions (`AssertOutput()`, `AssertCursor()`)
5. Write first test scenario (tab completion)
6. Iterate and expand based on results

This framework will provide comprehensive, fast, and reliable UI testing for DSH while maintaining the project's high code quality standards.
