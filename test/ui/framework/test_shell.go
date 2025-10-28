package framework

import (
	"dsh/internal/readline"
	"dsh/internal/terminal"
	"dsh/test/rendering"
)

// TestShell wraps readline functionality for UI testing
type TestShell struct {
	readline *readline.Readline
	mockTerm *rendering.MockTerminalInterface
	buffer   []rune
	cursor   int
}

// NewTestShell creates a shell instance for testing
func NewTestShell(mockTerm *rendering.MockTerminalInterface) *TestShell {
	return &TestShell{
		mockTerm: mockTerm,
		buffer:   make([]rune, 0),
		cursor:   0,
	}
}

// ProcessKey processes a key event and updates the shell state
func (s *TestShell) ProcessKey(keyEvent terminal.KeyEvent) {
	switch keyEvent.Key {
	case terminal.KeyTab:
		s.handleTabCompletion()
	case terminal.KeyArrowDown:
		s.handleArrowDown()
	case terminal.KeyArrowUp:
		s.handleArrowUp()
	case terminal.KeyEnter:
		s.handleEnter()
	case terminal.KeyCtrlA:
		s.cursor = 0
	case terminal.KeyCtrlU:
		s.buffer = s.buffer[:0]
		s.cursor = 0
	default:
		if keyEvent.Rune != 0 {
			s.insertRune(keyEvent.Rune)
		}
	}
}

// insertRune adds a character to the buffer at cursor position
func (s *TestShell) insertRune(r rune) {
	if s.cursor >= len(s.buffer) {
		// Insert at end
		s.buffer = append(s.buffer, r)
	} else {
		// Insert at cursor position
		newBuffer := make([]rune, len(s.buffer)+1)
		copy(newBuffer[:s.cursor], s.buffer[:s.cursor])
		newBuffer[s.cursor] = r
		copy(newBuffer[s.cursor+1:], s.buffer[s.cursor:])
		s.buffer = newBuffer
	}
	s.cursor++
}

// handleTabCompletion simulates tab completion
func (s *TestShell) handleTabCompletion() {
	input := string(s.buffer)

	// Simple completion logic for testing
	if input == "e" {
		// Simulate showing completion menu
		s.mockTerm.SaveCursor()
		s.mockTerm.WriteString("\r\n")
		s.mockTerm.WriteString("\033[7mecho\033[0m  \033[32mexit\033[0m  \033[34mhelp\033[0m")
		s.mockTerm.WriteString("\r\n")
	}
}

// handleArrowDown simulates arrow down navigation
func (s *TestShell) handleArrowDown() {
	// Simulate menu navigation - update selection
	output := s.mockTerm.GetOutput()
	if len(output) > 0 {
		// Menu is visible, simulate navigation
		s.mockTerm.WriteString("\033[36mexit\033[0m") // Highlight next item
	}
}

// handleArrowUp simulates arrow up navigation
func (s *TestShell) handleArrowUp() {
	// Similar to arrow down but in reverse
}

// handleEnter simulates enter key
func (s *TestShell) handleEnter() {
	// Complete the current selection or execute command
}

// GetBuffer returns current input buffer
func (s *TestShell) GetBuffer() string {
	return string(s.buffer)
}

// SetBuffer sets the input buffer
func (s *TestShell) SetBuffer(text string) {
	s.buffer = []rune(text)
	s.cursor = len(s.buffer)
}
