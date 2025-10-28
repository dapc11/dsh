package framework

import (
	"dsh/internal/readline"
	"dsh/test/rendering"
)

// UITestFramework provides automated UI testing for DSH shell
type UITestFramework struct {
	readline *readline.Readline
	mockTerm *rendering.MockTerminalInterface
	recorder *InteractionRecorder
	history  []string
}

// NewUITestFramework creates a new UI test framework instance
func NewUITestFramework() *UITestFramework {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)

	// Create readline with mock terminal using the test helper
	rl := readline.NewTestReadline(mockTerm)

	return &UITestFramework{
		readline: rl,
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
		// Also add to the readline's history
		if f.readline != nil && f.readline.GetHistory() != nil {
			f.readline.GetHistory().Add(cmd)
		}
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
	f.readline = readline.NewTestReadline(f.mockTerm)
	f.recorder = NewInteractionRecorder()
	f.history = make([]string, 0)
	return f
}

// GetShell returns the readline for direct access
func (f *UITestFramework) GetShell() *readline.Readline {
	return f.readline
}

// GetMockTerminal returns the mock terminal for direct access
func (f *UITestFramework) GetMockTerminal() *rendering.MockTerminalInterface {
	return f.mockTerm
}
