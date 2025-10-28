package terminal

import (
	"os"
)

// TerminalInterface defines the methods needed for terminal operations.
type TerminalInterface interface {
	// Terminal methods
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
	Printf(format string, args ...interface{}) (int, error)
	Size() (width, height int)
	MoveCursor(x, y int)
	ClearLine()
	ClearToEnd()
	ClearFromCursor()
	SaveCursor()
	RestoreCursor()
	HideCursor()
	ShowCursor()

	// ColorManager methods
	Colorize(text string, color Color) string
	StyleText(text string, style Style) string

	// InputReader methods
	ReadKey() (KeyEvent, error)

	// Interface methods
	EnableRawMode() error
	DisableRawMode() error
	IsRawMode() bool
	Cleanup()
}

// Interface provides a unified terminal abstraction.
type Interface struct {
	*Terminal
	*ColorManager
	*InputReader

	rawMode bool
}

// NewInterface creates a complete terminal interface.
func NewInterface() *Interface {
	return &Interface{
		Terminal:     New(),
		ColorManager: NewColorManager(),
		InputReader:  NewInputReader(os.Stdin),
		rawMode:      false,
	}
}

// EnableRawMode enables raw terminal mode for character-by-character input.
func (i *Interface) EnableRawMode() error {
	// TODO: Implement proper raw mode using syscalls
	i.rawMode = true
	return nil
}

// DisableRawMode restores normal terminal mode.
func (i *Interface) DisableRawMode() error {
	// TODO: Implement proper mode restoration
	i.rawMode = false
	return nil
}

// IsRawMode returns whether raw mode is enabled.
func (i *Interface) IsRawMode() bool {
	return i.rawMode
}

// Cleanup performs terminal cleanup.
func (i *Interface) Cleanup() {
	i.ShowCursor()
	_ = i.DisableRawMode() // Ignore error during cleanup
}
