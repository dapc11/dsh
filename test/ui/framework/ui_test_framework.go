package framework

import (
	"dsh/test/rendering"
)

// UITestFramework provides automated UI testing for DSH shell
type UITestFramework struct {
	shell    *TestShell
	mockTerm *rendering.MockTerminalInterface
	recorder *InteractionRecorder
	history  []string
}

// NewUITestFramework creates a new UI test framework instance
func NewUITestFramework() *UITestFramework {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	shell := NewTestShell(mockTerm)

	return &UITestFramework{
		shell:    shell,
		mockTerm: mockTerm,
		recorder: NewInteractionRecorder(),
		history:  make([]string, 0),
	}
}

// SetPrompt sets the shell prompt for testing
func (f *UITestFramework) SetPrompt(prompt string) *UITestFramework {
	// For now, just store the prompt - could be used for display
	return f
}

// AddHistory adds commands to shell history
func (f *UITestFramework) AddHistory(commands ...string) *UITestFramework {
	for _, cmd := range commands {
		f.history = append(f.history, cmd)
	}
	return f
}

// GetOutput returns the current terminal output
func (f *UITestFramework) GetOutput() string {
	return f.mockTerm.GetOutput()
}

// ClearOutput clears the terminal output buffer
func (f *UITestFramework) ClearOutput() *UITestFramework {
	f.mockTerm.ClearOutput()
	return f
}

// Reset resets the framework to initial state
func (f *UITestFramework) Reset() *UITestFramework {
	f.mockTerm = rendering.NewMockTerminalInterface(80, 24)
	f.shell = NewTestShell(f.mockTerm)
	f.recorder = NewInteractionRecorder()
	f.history = make([]string, 0)
	return f
}

// GetShell returns the test shell for direct access
func (f *UITestFramework) GetShell() *TestShell {
	return f.shell
}
