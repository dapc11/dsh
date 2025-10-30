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
	currentPage  int              // Current page number (0-based)
	itemsPerPage int              // Items per page (default 10)
}

// NewCompletionRenderer creates a new completion renderer
func NewCompletionRenderer(term terminal.TerminalInterface) *CompletionRenderer {
	return &CompletionRenderer{
		videoBuf:     NewVideoBuffer(term),
		terminal:     term,
		itemsPerPage: 10,
		currentPage:  0,
	}
}

// ShowCompletion displays completion menu with minimal rendering
func (cr *CompletionRenderer) ShowCompletion(items []CompletionItem, selected int) {
	if len(items) == 0 {
		return
	}

	// Store all items for pagination
	cr.lastItems = items
	cr.currentPage = selected / cr.itemsPerPage // Calculate which page the selected item is on

	// Save cursor position and move to next line
	_, _ = cr.terminal.WriteString("\033[s") // Save cursor for UpdateSelectionHighlight
	_, _ = cr.terminal.WriteString("\r\n")   // New line

	cr.renderCurrentPage(selected)
	cr.active = true
}

// renderCurrentPage renders the current page of items
func (cr *CompletionRenderer) renderCurrentPage(selected int) {
	if len(cr.lastItems) == 0 {
		return
	}

	// Calculate page bounds
	startIdx := cr.currentPage * cr.itemsPerPage
	endIdx := startIdx + cr.itemsPerPage
	if endIdx > len(cr.lastItems) {
		endIdx = len(cr.lastItems)
	}

	// Get items for current page
	pageItems := cr.lastItems[startIdx:endIdx]

	cols := 2
	for i, item := range pageItems {
		if i > 0 && i%cols == 0 {
			_, _ = cr.terminal.WriteString("\r\n")
		}

		// Show selection indicator for selected item
		globalIdx := startIdx + i
		if globalIdx == selected {
			_, _ = cr.terminal.WriteString(fmt.Sprintf("> %-33s", item.Text))
		} else {
			_, _ = cr.terminal.WriteString(fmt.Sprintf("  %-33s", item.Text))
		}

		if i%cols != cols-1 {
			_, _ = cr.terminal.WriteString("  ")
		}
	}

	_, _ = cr.terminal.WriteString("\r\n")
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
	if !cr.active || cr.lastItems == nil || oldSelected == newSelected {
		return
	}

	// Check if we need to change pages
	oldPage := oldSelected / cr.itemsPerPage
	newPage := newSelected / cr.itemsPerPage

	if oldPage != newPage {
		// Page change required - re-render entire menu
		cr.currentPage = newPage
		cr.clearMenu()
		cr.renderCurrentPage(newSelected)
		return
	}

	// Same page - just update selection indicators
	cr.updateSelectionOnSamePage(oldSelected, newSelected)
}

// clearMenu clears the current menu display
func (cr *CompletionRenderer) clearMenu() {
	// Move back to saved cursor position and clear menu area
	_, _ = cr.terminal.WriteString("\033[u") // Restore cursor
	_, _ = cr.terminal.WriteString("\033[J") // Clear from cursor to end of screen
	_, _ = cr.terminal.WriteString("\r\n")   // Move to next line for new menu
}

// updateSelectionOnSamePage updates selection when staying on the same page
func (cr *CompletionRenderer) updateSelectionOnSamePage(oldSelected, newSelected int) {
	// Calculate positions in the 2-column layout
	cols := 2
	startIdx := cr.currentPage * cr.itemsPerPage

	// Update old selection (remove "> ")
	if oldSelected >= 0 && oldSelected < len(cr.lastItems) {
		pageIdx := oldSelected - startIdx
		if pageIdx >= 0 && pageIdx < cr.itemsPerPage {
			row := (pageIdx / cols) + 1  // +1 for first menu line
			col := (pageIdx%cols)*37 + 1 // 37 chars per column

			_, _ = cr.terminal.WriteString("\033[u")                                    // Restore cursor
			_, _ = cr.terminal.WriteString(fmt.Sprintf("\033[%dB\033[%dG  ", row, col)) // Clear selection
		}
	}

	// Update new selection (add "> ")
	if newSelected >= 0 && newSelected < len(cr.lastItems) {
		pageIdx := newSelected - startIdx
		if pageIdx >= 0 && pageIdx < cr.itemsPerPage {
			row := (pageIdx / cols) + 1  // +1 for first menu line
			col := (pageIdx%cols)*37 + 1 // 37 chars per column

			_, _ = cr.terminal.WriteString("\033[u")                                    // Restore cursor
			_, _ = cr.terminal.WriteString(fmt.Sprintf("\033[%dB\033[%dG> ", row, col)) // Add selection
		}
	}
}
