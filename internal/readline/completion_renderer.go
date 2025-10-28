package readline

import (
	"strings"

	"dsh/internal/terminal"
)

// CompletionRenderer handles completion display using zsh-style video buffers
type CompletionRenderer struct {
	videoBuf     *VideoBuffer
	terminal     terminal.TerminalInterface
	savedCursorX int
	savedCursorY int
	menuStartY   int
	menuLines    int
	active       bool
}

// NewCompletionRenderer creates a new completion renderer
func NewCompletionRenderer(term terminal.TerminalInterface) *CompletionRenderer {
	return &CompletionRenderer{
		videoBuf: NewVideoBuffer(term),
		terminal: term,
	}
}

// ShowCompletion displays completion menu (like zsh's compprintlist)
func (cr *CompletionRenderer) ShowCompletion(items []CompletionItem, selected int) {
	if len(items) == 0 {
		return
	}

	// Save current cursor position (like zsh does)
	cr.savedCursorX, cr.savedCursorY = cr.videoBuf.GetCursorPosition()

	// Calculate menu layout
	maxWidth := 0
	for _, item := range items {
		if len(item.Text) > maxWidth {
			maxWidth = len(item.Text)
		}
	}

	itemWidth := maxWidth + 2
	cols := cr.videoBuf.width / itemWidth
	if cols < 1 {
		cols = 1
	}

	rows := (len(items) + cols - 1) / cols

	// Position menu after current line
	cr.menuStartY = cr.savedCursorY + 1
	cr.menuLines = rows

	// Clear menu area in video buffer
	for y := cr.menuStartY; y < cr.menuStartY+cr.menuLines && y < cr.videoBuf.height; y++ {
		for col := 0; col < cr.videoBuf.width; col++ {
			cr.videoBuf.newBuf[y][col] = VideoElement{Char: ' ', Attr: 0}
		}
	}

	// Render completion items to video buffer
	cr.videoBuf.MoveCursor(0, cr.menuStartY)

	for i, item := range items {
		if i > 0 && i%cols == 0 {
			// Move to next line
			_, y := cr.videoBuf.GetCursorPosition()
			cr.videoBuf.MoveCursor(0, y+1)
		}

		// Determine attributes
		attr := 0
		if i == selected {
			attr = 1 // Highlighted
		} else {
			switch item.Type {
			case "command":
				attr = 2
			case "directory":
				attr = 3
			case "builtin":
				attr = 4
			}
		}

		// Write item text with padding
		text := item.Text + strings.Repeat(" ", itemWidth-len(item.Text))
		cr.videoBuf.WriteString(text, attr)
	}

	// Refresh only the menu area
	cr.refreshMenuArea()

	// Restore cursor to original position
	cr.terminal.MoveCursor(cr.savedCursorX, cr.savedCursorY)
	cr.active = true
}

// UpdateSelection updates the selected item (like zsh's singledraw)
func (cr *CompletionRenderer) UpdateSelection(items []CompletionItem, oldSelected, newSelected int) {
	if !cr.active || oldSelected == newSelected {
		return
	}

	// Calculate layout (same as ShowCompletion)
	maxWidth := 0
	for _, item := range items {
		if len(item.Text) > maxWidth {
			maxWidth = len(item.Text)
		}
	}

	itemWidth := maxWidth + 2
	cols := cr.videoBuf.width / itemWidth
	if cols < 1 {
		cols = 1
	}

	// Update old selection (remove highlight)
	if oldSelected >= 0 && oldSelected < len(items) {
		cr.updateItem(items[oldSelected], oldSelected, cols, itemWidth, false)
	}

	// Update new selection (add highlight)
	if newSelected >= 0 && newSelected < len(items) {
		cr.updateItem(items[newSelected], newSelected, cols, itemWidth, true)
	}

	// Refresh only changed areas
	cr.refreshMenuArea()

	// Restore cursor
	cr.terminal.MoveCursor(cr.savedCursorX, cr.savedCursorY)
}

// updateItem updates a single completion item in the video buffer
func (cr *CompletionRenderer) updateItem(item CompletionItem, index, cols, itemWidth int, selected bool) {
	row := index / cols
	col := index % cols

	y := cr.menuStartY + row
	x := col * itemWidth

	if y >= cr.videoBuf.height {
		return
	}

	// Determine attributes
	attr := 0
	if selected {
		attr = 1 // Highlighted
	} else {
		switch item.Type {
		case "command":
			attr = 2
		case "directory":
			attr = 3
		case "builtin":
			attr = 4
		}
	}

	// Update video buffer
	text := item.Text + strings.Repeat(" ", itemWidth-len(item.Text))
	for i, char := range text {
		if x+i < cr.videoBuf.width {
			cr.videoBuf.newBuf[y][x+i] = VideoElement{Char: char, Attr: attr}
		}
	}
}

// refreshMenuArea refreshes only the menu area (optimized like zsh)
func (cr *CompletionRenderer) refreshMenuArea() {
	for y := cr.menuStartY; y < cr.menuStartY+cr.menuLines && y < cr.videoBuf.height; y++ {
		cr.videoBuf.refreshLine(y)
	}
	cr.videoBuf.swapBuffers()
}

// HideCompletion clears the completion menu
func (cr *CompletionRenderer) HideCompletion() {
	if !cr.active {
		return
	}

	// Clear menu area
	cr.videoBuf.ClearFromLine(cr.menuStartY)
	cr.refreshMenuArea()

	// Restore cursor
	cr.terminal.MoveCursor(cr.savedCursorX, cr.savedCursorY)
	cr.active = false
}

// IsActive returns whether completion menu is currently active
func (cr *CompletionRenderer) IsActive() bool {
	return cr.active
}
