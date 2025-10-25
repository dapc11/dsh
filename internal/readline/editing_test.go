package readline

import (
	"testing"
)

func TestEditing_InsertChar(t *testing.T) {
	t.Parallel()
	r := createTestReadline()
	r.buffer = []rune("hllo")
	r.cursor = 1

	// Insert 'e' between 'h' and 'l'
	r.insertChar('e')

	expected := "hello"
	if string(r.buffer) != expected {
		t.Errorf("insertChar: buffer = %s, want %s", string(r.buffer), expected)
	}

	if r.cursor != 2 {
		t.Errorf("insertChar: cursor = %d, want 2", r.cursor)
	}
}

func TestEditing_DeleteChar(t *testing.T) {
	t.Parallel()
	r := createTestReadline()
	r.buffer = []rune("hello")
	r.cursor = 1 // At 'e'

	r.deleteChar()

	expected := "hllo"
	if string(r.buffer) != expected {
		t.Errorf("deleteChar: buffer = %s, want %s", string(r.buffer), expected)
	}

	// Cursor should stay at same position
	if r.cursor != 1 {
		t.Errorf("deleteChar: cursor = %d, want 1", r.cursor)
	}
}

func TestEditing_Backspace(t *testing.T) {
	t.Parallel()
	r := createTestReadline()
	r.buffer = []rune("hello")
	r.cursor = 2 // After 'e'

	r.backspace()

	expected := "hllo"
	if string(r.buffer) != expected {
		t.Errorf("backspace: buffer = %s, want %s", string(r.buffer), expected)
	}

	// Cursor should move back
	if r.cursor != 1 {
		t.Errorf("backspace: cursor = %d, want 1", r.cursor)
	}
}

func TestEditing_ClearLine(t *testing.T) {
	t.Parallel()
	r := createTestReadline()
	r.buffer = []rune("hello world")
	r.cursor = 5

	r.clearLine()

	if len(r.buffer) != 0 {
		t.Errorf("clearLine: buffer length = %d, want 0", len(r.buffer))
	}

	if r.cursor != 0 {
		t.Errorf("clearLine: cursor = %d, want 0", r.cursor)
	}
}

func TestEditing_BoundaryConditions(t *testing.T) {
	t.Parallel()
	r := createTestReadline()

	// Test backspace at beginning
	r.buffer = []rune("hello")
	r.cursor = 0
	originalBuffer := string(r.buffer)
	r.backspace()

	if string(r.buffer) != originalBuffer {
		t.Errorf("backspace at start: buffer changed from %s to %s", originalBuffer, string(r.buffer))
	}

	// Test delete at end
	r.cursor = len(r.buffer)
	originalBuffer = string(r.buffer)
	r.deleteChar()

	if string(r.buffer) != originalBuffer {
		t.Errorf("deleteChar at end: buffer changed from %s to %s", originalBuffer, string(r.buffer))
	}
}
