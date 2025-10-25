package readline

import (
	"testing"
)

func TestTerminal_New(t *testing.T) {
	// This test may fail in non-terminal environments, so we'll be lenient
	terminal, err := NewTerminal()
	if err != nil {
		// In test environments, terminal initialization may fail
		t.Skipf("Terminal initialization failed (expected in test env): %v", err)
	}
	
	if terminal == nil {
		t.Error("NewTerminal returned nil")
	}
}

func TestTerminal_GetTerminalSize(t *testing.T) {
	terminal, err := NewTerminal()
	if err != nil {
		t.Skipf("Terminal initialization failed: %v", err)
	}
	
	width, height := terminal.GetTerminalSize()
	
	// Should return reasonable defaults even if detection fails
	if width <= 0 {
		t.Errorf("Width should be positive, got %d", width)
	}
	if height <= 0 {
		t.Errorf("Height should be positive, got %d", height)
	}
}

func TestColor_New(t *testing.T) {
	color := NewColor()
	
	if color == nil {
		t.Error("NewColor returned nil")
	}
}

func TestColor_Colorize(t *testing.T) {
	color := NewColor()
	
	text := "test text"
	
	// Test with valid color
	result := color.Colorize(text, "red")
	if result == "" {
		t.Error("Colorize returned empty string")
	}
	
	// Should contain original text
	if !contains(result, text) {
		t.Errorf("Colorized text should contain original text '%s'", text)
	}
	
	// Test with invalid color
	result = color.Colorize(text, "invalidcolor")
	if result == "" {
		t.Error("Colorize with invalid color returned empty string")
	}
}

func TestColor_DisabledColors(t *testing.T) {
	color := &Color{enabled: false}
	
	text := "test text"
	result := color.Colorize(text, "red")
	
	// When disabled, should return original text
	if result != text {
		t.Errorf("Disabled color should return original text, got '%s'", result)
	}
}

func TestSupportsColor(t *testing.T) {
	// This function checks environment variables and terminal capabilities
	// We'll just ensure it doesn't panic and returns a boolean
	result := supportsColor()
	
	// Should return a boolean (true or false, both are valid)
	_ = result
}

// Helper function since strings.Contains might not be available
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findSubstring(s, substr) >= 0)
}

func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
