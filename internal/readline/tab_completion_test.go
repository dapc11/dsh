package readline

import (
	"strings"
	"testing"
)

func TestTabCompletion_SingleMatch(t *testing.T) {
	term := NewMockTerminal()
	rl, err := newTestReadline(term)
	if err != nil {
		t.Fatal(err)
	}

	// Set up buffer with partial command
	rl.buffer = []rune("ech")
	rl.cursor = 3

	// Trigger tab completion
	rl.handleTabCompletion()

	// Should complete to "echo" without showing menu
	expected := "echo"
	actual := string(rl.buffer)
	if actual != expected {
		t.Errorf("Expected buffer '%s', got '%s'", expected, actual)
	}

	// Should not create temporary buffer for single match
	if rl.bufferManager.HasActiveBuffers() {
		t.Error("Should not have active buffers for single match")
	}
}

func TestTabCompletion_MultipleMatches(t *testing.T) {
	term := NewMockTerminal()
	rl, err := newTestReadline(term)
	if err != nil {
		t.Fatal(err)
	}

	// Set up buffer with ambiguous prefix
	rl.buffer = []rune("e")
	rl.cursor = 1

	// Trigger tab completion
	rl.handleTabCompletion()

	// Should show completion menu (basic check that completion works)
	output := term.GetOutput()
	if !strings.Contains(output, "exit") || !strings.Contains(output, "e2") {
		t.Error("Should display completion options")
	}
}

func TestTabCompletion_ClearOnEscape(t *testing.T) {
	term := NewMockTerminal()
	rl, err := newTestReadline(term)
	if err != nil {
		t.Fatal(err)
	}

	// Set up completion menu
	rl.buffer = []rune("e")
	rl.cursor = 1
	rl.handleTabCompletion()

	// Verify menu is active
	if !rl.bufferManager.HasActiveBuffers() {
		t.Error("Should have active buffer")
	}

	// Press escape to clear (don't reset output before this)
	rl.clearTabCompletion()

	// Should clear all buffers
	if rl.bufferManager.HasActiveBuffers() {
		t.Error("Should clear all buffers on escape")
	}

	// Should restore cursor
	output := term.GetOutput()
	if !strings.Contains(output, "\033[u") {
		t.Error("Should restore cursor position")
	}
}

func TestTabCompletion_NavigateMenu(t *testing.T) {
	term := NewMockTerminal()
	rl, err := newTestReadline(term)
	if err != nil {
		t.Fatal(err)
	}

	// Set up completion menu
	rl.buffer = []rune("e")
	rl.cursor = 1
	rl.handleTabCompletion()

	// Navigate to next item
	rl.navigateTabCompletion(1)

	// Should update menu display
	output := term.GetOutput()
	if !strings.Contains(output, "> ") { // Selection indicator
		t.Error("Should show selection highlight")
	}
}

func TestTabCompletion_AcceptSelection(t *testing.T) {
	term := NewMockTerminal()
	rl, err := newTestReadline(term)
	if err != nil {
		t.Fatal(err)
	}

	// Set up completion menu
	rl.buffer = []rune("e")
	rl.cursor = 1
	rl.handleTabCompletion()

	// Accept current selection
	rl.acceptTabCompletion()

	// Should update buffer with selection
	result := string(rl.buffer)
	if !strings.HasPrefix(result, "e") {
		t.Error("Should preserve original prefix")
	}

	// Should clear temporary buffer
	if rl.bufferManager.HasActiveBuffers() {
		t.Error("Should clear buffer after acceptance")
	}
}

func TestTabCompletion_RepeatedTab(t *testing.T) {
	term := NewMockTerminal()
	rl, err := newTestReadline(term)
	if err != nil {
		t.Fatal(err)
	}

	// Set up buffer with ambiguous prefix
	rl.buffer = []rune("e")
	rl.cursor = 1

	// First TAB - show menu
	rl.handleTabCompletion()
	if !rl.completionMenu.IsActive() {
		t.Error("Should show menu on first TAB")
	}

	// Get first selection
	firstSelected, _ := rl.completionMenu.GetSelected()

	term.Reset()

	// Second TAB - should navigate to next item
	rl.handleTabCompletion()
	if !rl.completionMenu.IsActive() {
		t.Error("Menu should still be active after second TAB")
	}

	// Should have moved to next item
	secondSelected, _ := rl.completionMenu.GetSelected()
	if firstSelected.Text == secondSelected.Text {
		t.Error("Should navigate to different item on repeated TAB")
	}

	// Should have some output from tab completion
	output := term.GetOutput()
	if len(output) == 0 {
		t.Error("Should have some output from tab completion")
	}
}

// Helper function to create test readline instance.
func newTestReadline(term *MockTerminal) (*Readline, error) {
	// Create a minimal readline for testing
	rl := &Readline{
		buffer:        make([]rune, 0, 256),
		cursor:        0,
		completion:    NewCompletion(),
		bufferManager: NewBufferManager(term),
		terminal:      term, // Set terminal field
	}
	rl.completionMenu = NewCompletionMenu(term)
	return rl, nil
}
