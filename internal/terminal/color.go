package terminal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Color represents terminal colors.
type Color int

// Terminal color constants.
const (
	ColorReset Color = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// Style represents text styling.
type Style struct {
	Foreground Color
	Background Color
	Bold       bool
	Italic     bool
	Underline  bool
	Reverse    bool
}

// ColorManager handles color operations.
type ColorManager struct {
	enabled bool
}

// NewColorManager creates a color manager.
func NewColorManager() *ColorManager {
	return &ColorManager{
		enabled: supportsColor(),
	}
}

// supportsColor detects color support.
func supportsColor() bool {
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	colorTerms := []string{"xterm", "screen", "tmux", "color", "ansi"}
	for _, ct := range colorTerms {
		if strings.Contains(term, ct) {
			return true
		}
	}

	return os.Getenv("COLORTERM") != ""
}

// Colorize applies color to text.
func (c *ColorManager) Colorize(text string, fg Color) string {
	if !c.enabled {
		return text
	}
	return c.colorCode(fg) + text + c.resetCode()
}

// StyleText applies full styling to text.
func (c *ColorManager) StyleText(text string, style Style) string {
	if !c.enabled {
		return text
	}

	codes := []string{}

	if style.Bold {
		codes = append(codes, "1")
	}
	if style.Italic {
		codes = append(codes, "3")
	}
	if style.Underline {
		codes = append(codes, "4")
	}
	if style.Reverse {
		codes = append(codes, "7")
	}

	if style.Foreground != ColorReset {
		codes = append(codes, c.fgCode(style.Foreground))
	}
	if style.Background != ColorReset {
		codes = append(codes, c.bgCode(style.Background))
	}

	if len(codes) == 0 {
		return text
	}

	return fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(codes, ";"), text)
}

// colorCode returns ANSI color code.
func (c *ColorManager) colorCode(color Color) string {
	return "\033[" + c.fgCode(color) + "m"
}

// fgCode returns foreground color code.
func (c *ColorManager) fgCode(color Color) string {
	switch color {
	case ColorBlack:
		return "30"
	case ColorRed:
		return "31"
	case ColorGreen:
		return "32"
	case ColorYellow:
		return "33"
	case ColorBlue:
		return "34"
	case ColorMagenta:
		return "35"
	case ColorCyan:
		return "36"
	case ColorWhite:
		return "37"
	case ColorBrightBlack:
		return "90"
	case ColorBrightRed:
		return "91"
	case ColorBrightGreen:
		return "92"
	case ColorBrightYellow:
		return "93"
	case ColorBrightBlue:
		return "94"
	case ColorBrightMagenta:
		return "95"
	case ColorBrightCyan:
		return "96"
	case ColorBrightWhite:
		return "97"
	default:
		return "39" // default
	}
}

// bgCode returns background color code.
func (c *ColorManager) bgCode(color Color) string {
	fg := c.fgCode(color)
	if code, err := strconv.Atoi(fg); err == nil {
		return strconv.Itoa(code + 10)
	}
	return "49" // default
}

// resetCode returns reset sequence.
func (c *ColorManager) resetCode() string {
	return "\033[0m"
}
