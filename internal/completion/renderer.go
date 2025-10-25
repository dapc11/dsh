package completion

import (
	"fmt"
	"strings"
)

// Renderer handles the visual display of completion menus.
type Renderer struct {
	colorizer ColorProvider
	terminal  TerminalProvider
}

// ColorProvider provides color formatting capabilities.
type ColorProvider interface {
	Colorize(text, color string) string
}

// TerminalProvider provides terminal size information.
type TerminalProvider interface {
	GetTerminalSize() (width, height int)
}

// NewRenderer creates a new menu renderer.
func NewRenderer(colorizer ColorProvider, terminal TerminalProvider) *Renderer {
	return &Renderer{
		colorizer: colorizer,
		terminal:  terminal,
	}
}

// Render displays the completion menu.
func (r *Renderer) Render(menu *Menu) {
	if !menu.IsDisplayed() || !menu.HasItems() {
		return
	}

	width, height := r.terminal.GetTerminalSize()
	maxItemWidth := 0

	// Find max item width
	for _, item := range menu.items {
		if len(item.Text) > maxItemWidth {
			maxItemWidth = len(item.Text)
		}
	}

	itemWidth := maxItemWidth + 2
	cols := width / itemWidth
	if cols < 1 {
		cols = 1
	}

	// Calculate available rows
	availableRows := height - 5
	if availableRows < 3 {
		availableRows = 3
	}

	maxRows := availableRows
	if maxRows > menu.maxRows {
		maxRows = menu.maxRows
	}

	itemsPerPage := maxRows * cols
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

	// Display items in grid
	_, _ = fmt.Print("\r\n") //nolint:forbidigo
	for i := 0; i < maxRows; i++ {
		for j := 0; j < cols; j++ {
			idx := i*cols + j
			if idx >= len(pageItems) {
				break
			}

			item := pageItems[idx]
			globalIdx := startIdx + idx
			text := item.Text

			var displayText string
			if globalIdx == menu.selected {
				displayText = r.colorizer.Colorize(text, "reverse")
			} else {
				switch item.Type {
				case "builtin":
					displayText = r.colorizer.Colorize(text, "cyan")
				case "command":
					displayText = r.colorizer.Colorize(text, "green")
				case "directory":
					displayText = r.colorizer.Colorize(text, "blue")
				default:
					displayText = text
				}
			}

			padding := itemWidth - len(text)
			_, _ = fmt.Print(displayText + strings.Repeat(" ", padding)) //nolint:forbidigo
		}
		_, _ = fmt.Print("\r\n") //nolint:forbidigo
	}

	// Show pagination info
	if totalPages > 1 {
		pageInfo := fmt.Sprintf("Page %d/%d (%d items)", menu.page+1, totalPages, len(menu.items))
		_, _ = fmt.Print(r.colorizer.Colorize(pageInfo, "gray") + "\r\n") //nolint:forbidigo
	}

	// Restore cursor to original position
	_, _ = fmt.Print("\0338") //nolint:forbidigo // Restore cursor position

	// Mark menu as displayed
	menu.displayed = true
	menu.linesDrawn = maxRows + 1 // menu rows + pagination line
}

// Clear removes the completion menu from display.
func (r *Renderer) Clear(menu *Menu) {
	if !menu.displayed {
		return
	}

	// Clear from current position to end of screen
	_, _ = fmt.Print("\r\n\033[0J\033[A") //nolint:forbidigo // Move to next line, clear to end, move back up

	menu.displayed = false
	menu.linesDrawn = 0
}
