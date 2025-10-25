package readline

import (
	"testing"
)

func createTestReadline() *Readline {
	r, _ := New("test> ")
	return r
}

func TestCursor_Movement(t *testing.T) {
	r := createTestReadline()
	r.buffer = []rune("hello world")
	r.cursor = 5 // Between "hello" and " world"

	// Test cursor left
	originalCursor := r.cursor
	r.moveCursorLeft()
	if r.cursor != originalCursor-1 {
		t.Errorf("moveCursorLeft: cursor = %d, want %d", r.cursor, originalCursor-1)
	}

	// Test cursor right
	r.moveCursorRight()
	if r.cursor != originalCursor {
		t.Errorf("moveCursorRight: cursor = %d, want %d", r.cursor, originalCursor)
	}

	// Test move to start
	r.moveCursorToStart()
	if r.cursor != 0 {
		t.Errorf("moveCursorToStart: cursor = %d, want 0", r.cursor)
	}

	// Test move to end
	r.moveCursorToEnd()
	if r.cursor != len(r.buffer) {
		t.Errorf("moveCursorToEnd: cursor = %d, want %d", r.cursor, len(r.buffer))
	}
}

func TestCursor_WordMovement(t *testing.T) {
	r := createTestReadline()
	r.buffer = []rune("hello world test")
	r.cursor = 0

	// Test word forward
	r.moveWordForward()
	if r.cursor != 6 { // Should be after "hello "
		t.Errorf("moveWordForward: cursor = %d, want 6", r.cursor)
	}

	r.moveWordForward()
	if r.cursor != 12 { // Should be after "world "
		t.Errorf("moveWordForward: cursor = %d, want 12", r.cursor)
	}

	// Test word backward
	r.moveWordBackward()
	if r.cursor != 6 { // Should be at start of "world"
		t.Errorf("moveWordBackward: cursor = %d, want 6", r.cursor)
	}
}

func TestCursor_Boundaries(t *testing.T) {
	r := createTestReadline()
	r.buffer = []rune("test")
	r.cursor = 0

	// Test left at beginning
	r.moveCursorLeft()
	if r.cursor != 0 {
		t.Errorf("moveCursorLeft at start: cursor = %d, want 0", r.cursor)
	}

	// Test right at end
	r.cursor = len(r.buffer)
	r.moveCursorRight()
	if r.cursor != len(r.buffer) {
		t.Errorf("moveCursorRight at end: cursor = %d, want %d", r.cursor, len(r.buffer))
	}
}
