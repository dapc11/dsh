package readline

import (
	"os"
	"strconv"
	"strings"
)

// Color represents terminal colors.
type Color struct {
	enabled bool
}

// NewColor creates a new color instance with terminal detection.
func NewColor() *Color {
	return &Color{
		enabled: supportsColor(),
	}
}

// supportsColor detects if terminal supports ANSI colors.
func supportsColor() bool {
	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	// Check for common color-supporting terminals
	colorTerms := []string{"xterm", "screen", "tmux", "color", "ansi"}
	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	// Check COLORTERM
	if os.Getenv("COLORTERM") != "" {
		return true
	}

	return true // Default to true for modern terminals
}

// ANSI color codes.
const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[90m"
	BrightRed = "\033[91m"
	BrightBlue = "\033[94m"
)

// Colorize applies color if supported.
func (c *Color) Colorize(text, color string) string {
	if !c.enabled {
		return text
	}
	return color + text + Reset
}

// Hex converts hex color to ANSI 256-color.
func (c *Color) Hex(text, hexColor string) string {
	if !c.enabled {
		return text
	}
	
	// Convert hex to RGB
	if len(hexColor) != 7 || hexColor[0] != '#' {
		return text // Invalid hex
	}
	
	r, _ := strconv.ParseInt(hexColor[1:3], 16, 64)
	g, _ := strconv.ParseInt(hexColor[3:5], 16, 64)
	b, _ := strconv.ParseInt(hexColor[5:7], 16, 64)
	
	// Convert RGB to 256-color approximation
	colorCode := 16 + (36 * (r / 51)) + (6 * (g / 51)) + (b / 51)
	
	return "\033[38;5;" + strconv.FormatInt(colorCode, 10) + "m" + text + Reset
}
