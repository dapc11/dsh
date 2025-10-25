package readline

import (
	"testing"
)

func TestKeyHandling_BasicKeys(t *testing.T) {
	r := createTestReadline()

	// Test Ctrl+A (move to start)
	r.buffer = []rune("hello")
	r.cursor = 3
	r.handleKey(KeyCtrlA)

	if r.cursor != 0 {
		t.Errorf("Ctrl+A: cursor = %d, want 0", r.cursor)
	}

	// Test Ctrl+E (move to end)
	r.handleKey(KeyCtrlE)
	if r.cursor != len(r.buffer) {
		t.Errorf("Ctrl+E: cursor = %d, want %d", r.cursor, len(r.buffer))
	}

	// Test Ctrl+U (clear line)
	r.handleKey(KeyCtrlU)
	if len(r.buffer) != 0 {
		t.Errorf("Ctrl+U: buffer length = %d, want 0", len(r.buffer))
	}
}

func TestKeyHandling_EOF(t *testing.T) {
	r := createTestReadline()

	// Ctrl+D on empty buffer should return false (EOF)
	r.buffer = []rune{}
	r.cursor = 0

	result := r.handleKey(KeyCtrlD)
	if result != false {
		t.Errorf("Ctrl+D on empty buffer: got %v, want false (EOF)", result)
	}

	// Ctrl+D on non-empty buffer should delete char and return true
	r.buffer = []rune("hello")
	r.cursor = 1

	result = r.handleKey(KeyCtrlD)
	if result != true {
		t.Errorf("Ctrl+D on non-empty buffer: got %v, want true", result)
	}

	expected := "hllo"
	if string(r.buffer) != expected {
		t.Errorf("Ctrl+D delete: buffer = %s, want %s", string(r.buffer), expected)
	}
}

func TestKeyHandling_EnterKey(t *testing.T) {
	r := createTestReadline()

	// Enter without menu should return false (complete line)
	result := r.handleKey('\r')
	if result != false {
		t.Errorf("Enter without menu: got %v, want false", result)
	}

	// Enter with menu should return true (accept selection)
	r.menuMode = true
	r.completionList = []CompletionItem{{Text: "test", Type: "command"}}
	r.menuSelected = 0
	r.completionBase = "t"

	result = r.handleKey('\r')
	if result != true {
		t.Errorf("Enter with menu: got %v, want true", result)
	}
}

func TestKeyHandling_PrintableChars(t *testing.T) {
	r := createTestReadline()

	// Test inserting printable characters
	chars := []rune{'h', 'e', 'l', 'l', 'o'}
	for _, ch := range chars {
		r.handleKey(byte(ch))
	}

	expected := "hello"
	if string(r.buffer) != expected {
		t.Errorf("Printable chars: buffer = %s, want %s", string(r.buffer), expected)
	}

	if r.cursor != len(expected) {
		t.Errorf("Printable chars: cursor = %d, want %d", r.cursor, len(expected))
	}
}

func TestKeyHandling_NonPrintableChars(t *testing.T) {
	r := createTestReadline()
	originalBuffer := string(r.buffer)

	// Test non-printable characters (should be ignored)
	r.handleKey(1)   // Ctrl+A is handled separately
	r.handleKey(255) // Non-printable

	// Buffer should be unchanged for non-printable chars outside the switch
	if string(r.buffer) != originalBuffer {
		t.Errorf("Non-printable chars changed buffer: %s", string(r.buffer))
	}
}
