// Package completion provides terminal-based completion menu rendering functionality.
package completion

import (
	"strings"

	"dsh/internal/terminal"
)

// Renderer handles the visual display of completion menus with proper buffer management.
type Renderer struct {
	terminal terminal.TerminalInterface
}

// NewRenderer creates a new menu renderer.
func NewRenderer(term terminal.TerminalInterface) *Renderer {
	return &Renderer{
		terminal: term,
	}
}

// Render displays the completion menu with proper cursor management.
func (r *Renderer) Render(menu *Menu) {
	if !menu.IsDisplayed() || !menu.HasItems() {
		return
	}

	// Save cursor position before rendering
	r.terminal.SaveCursor()

	itemWidth, cols, maxRows, itemsPerPage := r.calculateLayout(menu)
	totalPages := (len(menu.items) + itemsPerPage - 1) / itemsPerPage

	// Adjust page if selection moved
	selectedPage := menu.selected / itemsPerPage
	if selectedPage != menu.page {
		menu.page = selectedPage
	}

	// Calculate items to show on current page
	startIdx := menu.page * itemsPerPage
	endIdx := startIdx + itemsPerPage
	if endIdx > len(menu.items) {
		endIdx = len(menu.items)
	}

	pageItems := menu.items[startIdx:endIdx]

	// Move to next line to start rendering
	r.terminal.WriteString("\r\n")

	// Display items in grid
	linesUsed := 0
	for i := range maxRows {
		for j := range cols {
			idx := i*cols + j
			if idx >= len(pageItems) {
				break
			}

			item := pageItems[idx]
			globalIdx := startIdx + idx
			text := item.Text

			var displayText string
			if globalIdx == menu.selected {
				displayText = r.terminal.StyleText(text, terminal.Style{Reverse: true})
			} else {
				switch item.Type {
				case "builtin":
					displayText = r.terminal.Colorize(text, terminal.ColorCyan)
				case "command":
					displayText = r.terminal.Colorize(text, terminal.ColorGreen)
				case "directory":
					displayText = r.terminal.Colorize(text, terminal.ColorBlue)
				default:
					displayText = text
				}
			}

			padding := itemWidth - len(text)
			r.terminal.WriteString(displayText + strings.Repeat(" ", padding))
		}
		r.terminal.WriteString("\r\n")
		linesUsed++
	}

	// Show pagination info
	if totalPages > 1 {
		pageInfo := "Page " + string(rune(menu.page+1+'0')) + "/" + string(rune(totalPages+'0'))
		r.terminal.WriteString(r.terminal.Colorize(pageInfo, terminal.ColorBrightBlack) + "\r\n")
		linesUsed++
	}

	// Store lines used for cleanup
	menu.linesDrawn = linesUsed

	// Mark menu as displayed
	menu.displayed = true
}

// Clear removes the completion menu from display with proper cleanup.
func (r *Renderer) Clear(menu *Menu) {
	if !menu.displayed {
		return
	}

	// Clear from cursor to end of screen
	r.terminal.ClearFromCursor()

	// Restore cursor to original position
	r.terminal.RestoreCursor()

	// Mark menu as not displayed
	menu.displayed = false
	menu.linesDrawn = 0
}

// calculateLayout calculates the layout parameters for the menu.
func (r *Renderer) calculateLayout(menu *Menu) (itemWidth, cols, maxRows, itemsPerPage int) {
	width, height := r.terminal.Size()
	maxItemWidth := 0

	// Find max item width
	for _, item := range menu.items {
		if len(item.Text) > maxItemWidth {
			maxItemWidth = len(item.Text)
		}
	}

	itemWidth = maxItemWidth + 2
	cols = width / itemWidth
	if cols < 1 {
		cols = 1
	}

	// Calculate available rows
	availableRows := height - 5
	if availableRows < 3 {
		availableRows = 3
	}

	maxRows = availableRows
	if maxRows > menu.maxRows {
		maxRows = menu.maxRows
	}

	itemsPerPage = maxRows * cols
	return itemWidth, cols, maxRows, itemsPerPage
}
