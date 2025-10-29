package readline

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/rendering"
)

func TestKeyBindings_CtrlKeys(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)
	rl.buffer = []rune("hello world")
	rl.cursor = 6

	tests := []struct {
		name     string
		key      terminal.Key
		expected string
		cursor   int
	}{
		{"Ctrl+A", terminal.KeyCtrlA, "hello world", 0},
		{"Ctrl+E", terminal.KeyCtrlE, "hello world", 11},
		{"Ctrl+B", terminal.KeyCtrlB, "hello world", 5},
		{"Ctrl+F", terminal.KeyCtrlF, "hello world", 7},
		{"Ctrl+K", terminal.KeyCtrlK, "hello ", 6},
		{"Ctrl+U", terminal.KeyCtrlU, "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl.buffer = []rune("hello world")
			rl.cursor = 6
			rl.handleKeyEvent(terminal.KeyEvent{Key: tt.key})
			
			if string(rl.buffer) != tt.expected {
				t.Errorf("Expected buffer %q, got %q", tt.expected, string(rl.buffer))
			}
			if rl.cursor != tt.cursor {
				t.Errorf("Expected cursor %d, got %d", tt.cursor, rl.cursor)
			}
		})
	}
}

func TestKeyBindings_ArrowKeys(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)
	rl.buffer = []rune("test")
	rl.cursor = 2

	tests := []struct {
		name   string
		key    terminal.Key
		cursor int
	}{
		{"Left", terminal.KeyArrowLeft, 1},
		{"Right", terminal.KeyArrowRight, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl.cursor = 2
			rl.handleKeyEvent(terminal.KeyEvent{Key: tt.key})
			
			if rl.cursor != tt.cursor {
				t.Errorf("Expected cursor %d, got %d", tt.cursor, rl.cursor)
			}
		})
	}
}

func TestKeyBindings_BackspaceDelete(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)

	t.Run("Backspace", func(t *testing.T) {
		rl.buffer = []rune("hello")
		rl.cursor = 3
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyBackspace})
		
		if string(rl.buffer) != "helo" {
			t.Errorf("Expected buffer %q, got %q", "helo", string(rl.buffer))
		}
		if rl.cursor != 2 {
			t.Errorf("Expected cursor %d, got %d", 2, rl.cursor)
		}
	})

	t.Run("Ctrl+D", func(t *testing.T) {
		rl.buffer = []rune("hello")
		rl.cursor = 2
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyCtrlD})
		
		if string(rl.buffer) != "helo" {
			t.Errorf("Expected buffer %q, got %q", "helo", string(rl.buffer))
		}
		if rl.cursor != 2 {
			t.Errorf("Expected cursor %d, got %d", 2, rl.cursor)
		}
	})
}

func TestKeyBindings_WordMovement(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)
	rl.buffer = []rune("hello world test")
	rl.cursor = 8

	t.Run("Ctrl+W", func(t *testing.T) {
		rl.buffer = []rune("hello world")
		rl.cursor = 11
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyCtrlW})
		
		if string(rl.buffer) != "hello " {
			t.Errorf("Expected buffer %q, got %q", "hello ", string(rl.buffer))
		}
	})
}

func TestKeyBindings_KillYank(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)

	t.Run("Ctrl+Y", func(t *testing.T) {
		rl.buffer = []rune("hello")
		rl.cursor = 5
		rl.killRing.Add("world")
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyCtrlY})
		
		if string(rl.buffer) != "helloworld" {
			t.Errorf("Expected buffer %q, got %q", "helloworld", string(rl.buffer))
		}
	})
}

func TestKeyBindings_SpecialKeys(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)

	t.Run("Ctrl+L", func(t *testing.T) {
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyCtrlL})
		// Should not crash
	})

	t.Run("Ctrl+Z", func(t *testing.T) {
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyCtrlZ})
		// Should not crash
	})
}

func TestKeyBindings_PrintableChars(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)
	rl.buffer = []rune("hello")
	rl.cursor = 5

	rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyNone, Rune: ' '})
	rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyNone, Rune: 'w'})
	
	if string(rl.buffer) != "hello w" {
		t.Errorf("Expected buffer %q, got %q", "hello w", string(rl.buffer))
	}
	if rl.cursor != 7 {
		t.Errorf("Expected cursor %d, got %d", 7, rl.cursor)
	}
}

func TestKeyBindings_EnterKey(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)

	result := rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyEnter})
	if result != false {
		t.Errorf("Expected Enter to return false, got %v", result)
	}
}

func TestKeyBindings_TabCompletion(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)
	rl.buffer = []rune("ec")
	rl.cursor = 2

	rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyTab})
	
	if string(rl.buffer) != "echo" {
		t.Errorf("Expected buffer %q, got %q", "echo", string(rl.buffer))
	}
}

func TestKeyBindings_EscapeSequences(t *testing.T) {
	mockTerm := rendering.NewMockTerminalInterface(80, 24)
	rl := NewTestReadline(mockTerm)

	t.Run("Escape", func(t *testing.T) {
		rl.handleKeyEvent(terminal.KeyEvent{Key: terminal.KeyEscape})
		// Should not crash
	})
}
