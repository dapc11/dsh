package readline

import (
	"fmt"
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
	lastItems    []CompletionItem // Store items for redraw
}

// NewCompletionRenderer creates a new completion renderer
func NewCompletionRenderer(term terminal.TerminalInterface) *CompletionRenderer {
	return &CompletionRenderer{
		videoBuf: NewVideoBuffer(term),
		terminal: term,
	}
}

// ShowCompletion displays completion menu with minimal rendering
func (cr *CompletionRenderer) ShowCompletion(items []CompletionItem, selected int) {
	if len(items) == 0 {
		return
	}

	// Store items for redraw
	cr.lastItems = items
	
	// Save cursor and move to next line
	cr.terminal.WriteString("\033[s") // Save cursor
	cr.terminal.WriteString("\r\n")   // New line
	
	// Simple menu rendering - just show items in columns
	maxItems := 10 // Limit items to prevent excessive output
	if len(items) > maxItems {
		items = items[:maxItems]
	}
	
	cols := 2
	for i, item := range items {
		if i > 0 && i%cols == 0 {
			cr.terminal.WriteString("\r\n")
		}
		
		// Highlight selected item
		if i == selected {
			cr.terminal.WriteString(fmt.Sprintf("\033[7m%-35s\033[0m", item.Text))
		} else {
			cr.terminal.WriteString(fmt.Sprintf("%-35s", item.Text))
		}
		
		if i%cols != cols-1 {
			cr.terminal.WriteString("  ")
		}
	}
	
	cr.terminal.WriteString("\r\n")
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
	cr.active = false
}

// IsActive returns whether completion menu is currently active
func (cr *CompletionRenderer) IsActive() bool {
	return cr.active
}

// UpdateSelectionHighlight updates only the selection highlighting efficiently
func (cr *CompletionRenderer) UpdateSelectionHighlight(oldSelected, newSelected int) {
	if !cr.active || cr.lastItems == nil {
		return
	}
	
	// Clear from cursor and redraw menu with new selection
	cr.terminal.WriteString("\033[u") // Restore cursor
	cr.terminal.WriteString("\033[J") // Clear from cursor to end
	cr.terminal.WriteString("\033[s") // Save cursor again for next time
	cr.terminal.WriteString("\r\n")   // New line
	
	// Calculate which page of items to show based on selection
	maxItems := 10
	pageStart := (newSelected / maxItems) * maxItems
	pageEnd := pageStart + maxItems
	if pageEnd > len(cr.lastItems) {
		pageEnd = len(cr.lastItems)
	}
	
	// Show the correct page of items
	pageItems := cr.lastItems[pageStart:pageEnd]
	cols := 2
	
	for i, item := range pageItems {
		if i > 0 && i%cols == 0 {
			cr.terminal.WriteString("\r\n")
		}
		
		// Highlight selected item (adjust index for page)
		globalIndex := pageStart + i
		if globalIndex == newSelected {
			cr.terminal.WriteString(fmt.Sprintf("\033[7m%-35s\033[0m", item.Text))
		} else {
			cr.terminal.WriteString(fmt.Sprintf("%-35s", item.Text))
		}
		
		if i%cols != cols-1 {
			cr.terminal.WriteString("  ")
		}
	}
	
	cr.terminal.WriteString("\r\n")
}
