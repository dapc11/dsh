package readline

import (
	"dsh/internal/terminal"
)

// VideoElement represents a single character on screen with attributes
type VideoElement struct {
	Char rune
	Attr int // Color/style attributes
}

// VideoBuffer represents the screen state (like zsh's nbuf/obuf)
type VideoBuffer struct {
	terminal terminal.TerminalInterface
	width    int
	height   int
	newBuf   [][]VideoElement // New screen state (like zsh's nbuf)
	oldBuf   [][]VideoElement // Old screen state (like zsh's obuf)
	cursorX  int
	cursorY  int
}

// NewVideoBuffer creates a new video buffer system
func NewVideoBuffer(term terminal.TerminalInterface) *VideoBuffer {
	width, height := term.Size()

	vb := &VideoBuffer{
		terminal: term,
		width:    width,
		height:   height,
		newBuf:   make([][]VideoElement, height),
		oldBuf:   make([][]VideoElement, height),
	}

	// Initialize buffers
	for i := 0; i < height; i++ {
		vb.newBuf[i] = make([]VideoElement, width)
		vb.oldBuf[i] = make([]VideoElement, width)
	}

	return vb
}

// Clear clears the new buffer
func (vb *VideoBuffer) Clear() {
	for y := 0; y < vb.height; y++ {
		for x := 0; x < vb.width; x++ {
			vb.newBuf[y][x] = VideoElement{Char: ' ', Attr: 0}
		}
	}
}

// WriteString writes text to the video buffer at current cursor position
func (vb *VideoBuffer) WriteString(text string, attr int) {
	for _, char := range text {
		if char == '\n' {
			vb.cursorY++
			vb.cursorX = 0
			continue
		}
		if char == '\r' {
			vb.cursorX = 0
			continue
		}

		if vb.cursorY < vb.height && vb.cursorX < vb.width {
			vb.newBuf[vb.cursorY][vb.cursorX] = VideoElement{Char: char, Attr: attr}
			vb.cursorX++
		}
	}
}

// MoveCursor moves the cursor position in the buffer
func (vb *VideoBuffer) MoveCursor(x, y int) {
	vb.cursorX = x
	vb.cursorY = y
}

// Refresh updates the terminal to match the new buffer (like zsh's refreshline)
func (vb *VideoBuffer) Refresh() {
	for y := 0; y < vb.height; y++ {
		vb.refreshLine(y)
	}
	vb.swapBuffers()
}

// refreshLine refreshes a single line (like zsh's refreshline function)
func (vb *VideoBuffer) refreshLine(line int) {
	if line >= vb.height {
		return
	}

	newLine := vb.newBuf[line]
	oldLine := vb.oldBuf[line]

	// Find first difference
	start := 0
	for start < vb.width && newLine[start] == oldLine[start] {
		start++
	}

	// Find last difference
	end := vb.width - 1
	for end >= start && newLine[end] == oldLine[end] {
		end--
	}

	if start > end {
		return // No changes on this line
	}

	// Move to start of changes
	vb.terminal.MoveCursor(start, line)

	// Write changed characters
	currentAttr := -1
	for x := start; x <= end; x++ {
		elem := newLine[x]

		// Handle attribute changes
		if elem.Attr != currentAttr {
			vb.applyAttributes(elem.Attr)
			currentAttr = elem.Attr
		}

		vb.terminal.WriteString(string(elem.Char))
	}

	// Reset attributes
	if currentAttr != 0 {
		vb.terminal.WriteString("\033[0m")
	}
}

// applyAttributes applies color/style attributes
func (vb *VideoBuffer) applyAttributes(attr int) {
	switch attr {
	case 1: // Selected/highlighted
		vb.terminal.WriteString("\033[7m") // Reverse video
	case 2: // Command
		vb.terminal.WriteString("\033[32m") // Green
	case 3: // Directory
		vb.terminal.WriteString("\033[34m") // Blue
	case 4: // Builtin
		vb.terminal.WriteString("\033[36m") // Cyan
	default:
		vb.terminal.WriteString("\033[0m") // Reset
	}
}

// swapBuffers swaps new and old buffers (like zsh's bufswap)
func (vb *VideoBuffer) swapBuffers() {
	vb.newBuf, vb.oldBuf = vb.oldBuf, vb.newBuf
}

// ClearFromLine clears from specified line to end of screen
func (vb *VideoBuffer) ClearFromLine(line int) {
	for y := line; y < vb.height; y++ {
		for x := 0; x < vb.width; x++ {
			vb.newBuf[y][x] = VideoElement{Char: ' ', Attr: 0}
		}
	}
}

// GetCursorPosition returns current cursor position
func (vb *VideoBuffer) GetCursorPosition() (int, int) {
	return vb.cursorX, vb.cursorY
}
