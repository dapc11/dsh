package terminal

import (
	"strings"
	"testing"
)

func TestColorManager(t *testing.T) {
	cm := NewColorManager()

	// Test basic colorization
	red := cm.Colorize("error", ColorRed)
	if !strings.Contains(red, "error") {
		t.Error("Colorize should contain original text")
	}

	// Test styling
	style := Style{
		Foreground: ColorGreen,
		Bold:       true,
	}
	styled := cm.StyleText("success", style)
	if !strings.Contains(styled, "success") {
		t.Error("StyleText should contain original text")
	}
}

func TestTerminal(t *testing.T) {
	term := New()

	width, height := term.Size()
	if width <= 0 || height <= 0 {
		t.Error("Terminal size should be positive")
	}
}

func TestInterface(t *testing.T) {
	iface := NewInterface()

	if iface.Terminal == nil {
		t.Error("Interface should have Terminal")
	}
	if iface.ColorManager == nil {
		t.Error("Interface should have ColorManager")
	}
	if iface.InputReader == nil {
		t.Error("Interface should have InputReader")
	}
}
