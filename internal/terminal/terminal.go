// Package terminal provides a robust abstraction for terminal I/O operations.
package terminal

import (
	"fmt"
	"io"
	"os"
)

// Terminal provides safe terminal operations.
type Terminal struct {
	input  io.Reader
	output io.Writer
	width  int
	height int
}

// New creates a new terminal instance.
func New() *Terminal {
	t := &Terminal{
		input:  os.Stdin,
		output: os.Stdout,
	}
	t.updateSize()
	return t
}

// Write writes data to terminal.
func (t *Terminal) Write(data []byte) (int, error) {
	return t.output.Write(data)
}

// WriteString writes a string to terminal.
func (t *Terminal) WriteString(s string) (int, error) {
	return io.WriteString(t.output, s)
}

// Printf writes formatted output.
func (t *Terminal) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(t.output, format, args...)
}

// Size returns terminal dimensions.
func (t *Terminal) Size() (width, height int) {
	t.updateSize()
	return t.width, t.height
}

// MoveCursor moves cursor to position.
func (t *Terminal) MoveCursor(x, y int) {
	_, _ = t.WriteString(fmt.Sprintf("\033[%d;%dH", y+1, x+1))
}

// ClearLine clears current line.
func (t *Terminal) ClearLine() {
	_, _ = t.WriteString("\033[2K")
}

// ClearToEnd clears from cursor to end of line.
func (t *Terminal) ClearToEnd() {
	_, _ = t.WriteString("\033[K")
}

// SaveCursor saves cursor position.
func (t *Terminal) SaveCursor() {
	_, _ = t.WriteString("\033[7")
}

// RestoreCursor restores cursor position.
func (t *Terminal) RestoreCursor() {
	_, _ = t.WriteString("\033[8")
}

// HideCursor hides the cursor.
func (t *Terminal) HideCursor() {
	_, _ = t.WriteString("\033[?25l")
}

// ShowCursor shows the cursor.
func (t *Terminal) ShowCursor() {
	_, _ = t.WriteString("\033[?25h")
}

// ClearFromCursor clears from cursor to end of screen.
func (t *Terminal) ClearFromCursor() {
	_, _ = t.WriteString("\033[0J")
}

// updateSize gets current terminal size.
func (t *Terminal) updateSize() {
	// Simple fallback - could be enhanced with proper syscalls
	t.width = 80
	t.height = 24
}
