package terminal

import (
	"bytes"
	"testing"
)

func TestInputReader_DELCharacterHandling(t *testing.T) {
	// Test that DEL character (127) is properly mapped to KeyBackspace
	buf := bytes.NewBuffer([]byte{127})
	reader := NewInputReader(buf)

	event, err := reader.ReadKey()
	if err != nil {
		t.Fatalf("ReadKey failed: %v", err)
	}

	if event.Key != KeyBackspace {
		t.Errorf("DEL character (127) should map to KeyBackspace, got Key=%v, Rune=%v", event.Key, event.Rune)
	}

	if event.Rune != 0 {
		t.Errorf("DEL character should not set Rune field, got %v", event.Rune)
	}
}

func TestInputReader_BackspaceCharacterHandling(t *testing.T) {
	// Test that BS character (8) is properly mapped to KeyBackspace
	buf := bytes.NewBuffer([]byte{8})
	reader := NewInputReader(buf)

	event, err := reader.ReadKey()
	if err != nil {
		t.Fatalf("ReadKey failed: %v", err)
	}

	if event.Key != KeyBackspace {
		t.Errorf("BS character (8) should map to KeyBackspace, got Key=%v, Rune=%v", event.Key, event.Rune)
	}

	if event.Rune != 0 {
		t.Errorf("BS character should not set Rune field, got %v", event.Rune)
	}
}

func TestInputReader_PrintableCharacters(t *testing.T) {
	// Test that printable characters are handled as runes
	testCases := []struct {
		input    byte
		expected rune
	}{
		{32, ' '},  // space
		{65, 'A'},  // A
		{97, 'a'},  // a
		{126, '~'}, // tilde (last printable ASCII)
	}

	for _, tc := range testCases {
		buf := bytes.NewBuffer([]byte{tc.input})
		reader := NewInputReader(buf)

		event, err := reader.ReadKey()
		if err != nil {
			t.Fatalf("ReadKey failed for input %d: %v", tc.input, err)
		}

		if event.Key != KeyNone {
			t.Errorf("Printable character %d should have Key=KeyNone, got %v", tc.input, event.Key)
		}

		if event.Rune != tc.expected {
			t.Errorf("Printable character %d should have Rune=%v, got %v", tc.input, tc.expected, event.Rune)
		}
	}
}

func TestInputReader_ControlCharacterBoundaries(t *testing.T) {
	// Test boundary conditions for control character detection
	testCases := []struct {
		input       byte
		shouldBeKey bool
		expectedKey Key
	}{
		{1, true, KeyCtrlA},       // Ctrl+A
		{8, true, KeyBackspace},   // BS
		{31, true, KeyNone},       // Unmapped control char
		{32, false, KeyNone},      // Space (first printable)
		{126, false, KeyNone},     // Tilde (last printable)
		{127, true, KeyBackspace}, // DEL (special case)
	}

	for _, tc := range testCases {
		buf := bytes.NewBuffer([]byte{tc.input})
		reader := NewInputReader(buf)

		event, err := reader.ReadKey()
		if err != nil {
			t.Fatalf("ReadKey failed for input %d: %v", tc.input, err)
		}

		if tc.shouldBeKey {
			if tc.expectedKey != KeyNone && event.Key != tc.expectedKey {
				t.Errorf("Control character %d should map to %v, got %v", tc.input, tc.expectedKey, event.Key)
			}
			if event.Rune != 0 {
				t.Errorf("Control character %d should not set Rune, got %v", tc.input, event.Rune)
			}
		} else {
			if event.Key != KeyNone {
				t.Errorf("Printable character %d should have Key=KeyNone, got %v", tc.input, event.Key)
			}
			if event.Rune == 0 {
				t.Errorf("Printable character %d should set Rune field", tc.input)
			}
		}
	}
}
