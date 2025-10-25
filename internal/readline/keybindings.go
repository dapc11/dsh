package readline

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (r *Readline) handleKey(ch byte) bool { //nolint:cyclop,funlen // Key handling naturally requires many branches and statements
	switch ch {
	case '\r', '\n':
		if r.menuMode {
			r.acceptMenuSelection()
			return true
		}
		return false // Signal completion
	case KeyCtrlA:
		r.killRing.ResetYank()
		r.moveCursorToStart()
	case KeyCtrlE:
		r.killRing.ResetYank()
		r.moveCursorToEnd()
	case KeyCtrlB:
		r.killRing.ResetYank()
		r.moveCursorLeft()
	case KeyCtrlC:
		r.killRing.ResetYank()
		if r.menuMode {
			r.clearCompletionMenu()
			r.resetCompletion()
		}
		r.clearLine()
		_, _ = fmt.Print("^C\r\n") //nolint:forbidigo
		r.displayPrompt()
	case KeyCtrlF:
		r.killRing.ResetYank()
		if r.suggestion != "" {
			r.acceptSuggestion()
		} else {
			r.moveCursorRight()
		}
	case KeyCtrlD:
		r.killRing.ResetYank()
		if len(r.buffer) == 0 {
			// EOF case - let ReadLine handle it
			return false
		}
		r.deleteChar()
	case KeyCtrlK:
		r.killToEnd()
	case KeyCtrlL:
		r.killRing.ResetYank()
		r.clearScreen()
	case KeyCtrlN:
		r.killRing.ResetYank()
		r.historyNext()
	case KeyCtrlP:
		r.killRing.ResetYank()
		r.historyPrevious()
	case KeyCtrlR:
		r.killRing.ResetYank()
		if selected := r.FuzzyHistorySearchCustom(); selected != "" {
			r.buffer = []rune(selected)
			r.cursor = len(r.buffer)
			r.redraw()
		}
	case KeyCtrlU:
		r.killLine()
	case KeyCtrlW:
		r.killWordBackward()
	case KeyCtrlY:
		r.yank()
	case KeyCtrlZ:
		r.killRing.ResetYank()
		_, _ = fmt.Print("^Z\r\n") //nolint:forbidigo
		r.displayPrompt()
	case KeyTab:
		r.killRing.ResetYank()
		if r.menuMode {
			// Cycle to next option and apply it
			r.menuSelected = (r.menuSelected + 1) % len(r.completionList)
			r.applyCycleCompletion()
		} else {
			// Start new completion
			r.performCompletion()
		}
	case KeyEscape:
		// Always try to handle escape sequences first (including shift-tab)
		err := r.handleEscapeSequence()
		if err != nil {
			// If escape sequence handling fails and we're in menu mode, exit menu
			if r.menuMode {
				r.clearCompletionMenu()
				r.resetCompletion()
				r.redraw()
			}
		}
	case KeyBackspace:
		r.killRing.ResetYank()
		r.backspace()
		r.searchPrefix = ""
		r.browseMode = false // Exit browse mode when editing
		r.updateSuggestion()
	default:
		if ch >= 32 && ch < 127 {
			r.killRing.ResetYank()
			if r.menuMode {
				r.clearCompletionMenu()
			}
			r.insertChar(rune(ch))
			r.searchPrefix = ""
			r.browseMode = false // Exit browse mode when editing
			r.updateSuggestion()
		}
	}

	return true
}

func (r *Readline) handleEscapeSequence() error { //nolint:gocognit,cyclop,funlen // Terminal input handling requires complexity
	ch1, err := r.readChar()
	if err != nil {
		return fmt.Errorf("failed to read escape sequence: %w", err)
	}

	switch ch1 {
	case 127: // Alt+Backspace (ESC + DEL)
		r.killRing.ResetYank()
		r.killWordBackward()
	case '[':
		ch2, err := r.readChar()
		if err != nil {
			return fmt.Errorf("failed to read escape sequence: %w", err)
		}

		switch ch2 {
		case 'A': // Up arrow
			if r.menuMode {
				r.navigateMenu(-1)
			} else {
				r.historyPrevious()
			}
		case 'B': // Down arrow
			if r.menuMode {
				r.navigateMenu(1)
			} else {
				r.historyNext()
			}
		case 'C': // Right arrow
			if r.menuMode {
				r.navigateMenuHorizontal(1)
			} else {
				r.moveCursorRight()
			}
		case 'D': // Left arrow
			if r.menuMode {
				r.navigateMenuHorizontal(-1)
			} else {
				r.moveCursorLeft()
			}
		case 'Z': // Shift-Tab (ESC[Z)
			if r.menuMode {
				// Navigate backward through completion menu
				r.navigateMenu(-1)
			}
		case '1':
			ch3, err := r.readChar()
			if err == nil && ch3 == ';' {
				ch4, err := r.readChar()
				if err == nil && ch4 == '5' {
					ch5, err := r.readChar()
					if err == nil {
						switch ch5 {
						case 'C': // Ctrl+Right
							r.moveWordForward()
						case 'D': // Ctrl+Left
							r.moveWordBackward()
						}
					}
				}
			}
		}
	case 'y': // Alt+Y (yank pop forward)
		r.yankPop()
	case 'Y': // Alt+Shift+Y (yank pop backward)
		r.yankCycle(-1)
	case 'd': // Alt+D (delete word forward)
		r.killRing.ResetYank()
		r.killWordForward()
	}

	return nil
}

// performCompletion performs tab completion with navigable menu.
func (r *Readline) performCompletion() {
	input := string(r.buffer)

	// Save cursor position for menu display
	_, _ = fmt.Print("\0337") //nolint:forbidigo // Save cursor position

	// Start new completion
	matches, completion := r.completion.Complete(input, r.cursor)

	if len(matches) == 1 && completion != "" {
		// Single match - apply directly
		r.resetCompletion()
		completionRunes := []rune(completion)
		r.buffer = append(r.buffer, completionRunes...)
		r.cursor = len(r.buffer)
		r.redraw()
	} else if len(matches) > 1 {
		// Multiple matches - show menu
		r.completionList = matches
		r.completionBase = input
		r.menuMode = true
		r.menuSelected = 0
		r.showCompletionMenu()
	}
}

// showCompletionMenu displays the completion menu with pagination.
func (r *Readline) showCompletionMenu() { //nolint:cyclop,funlen // Complex UI rendering is acceptable
	width, height := r.terminal.GetTerminalSize()
	maxItemWidth := 0

	// Find max item width (text only, no color codes)
	for _, item := range r.completionList {
		if len(item.Text) > maxItemWidth {
			maxItemWidth = len(item.Text)
		}
	}

	itemWidth := maxItemWidth + 2
	cols := width / itemWidth
	if cols < 1 {
		cols = 1
	}

	// Calculate available rows (leave space for prompt and pagination info)
	availableRows := height - 5
	if availableRows < 3 {
		availableRows = 3
	}

	maxRows := availableRows
	if maxRows > r.menuMaxRows {
		maxRows = r.menuMaxRows
	}

	itemsPerPage := maxRows * cols
	totalPages := (len(r.completionList) + itemsPerPage - 1) / itemsPerPage

	// Adjust page if selection moved
	selectedPage := r.menuSelected / itemsPerPage
	if selectedPage != r.menuPage {
		r.menuPage = selectedPage
	}

	// Calculate items to show on current page
	startIdx := r.menuPage * itemsPerPage
	endIdx := startIdx + itemsPerPage
	if endIdx > len(r.completionList) {
		endIdx = len(r.completionList)
	}

	pageItems := r.completionList[startIdx:endIdx]

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
			if globalIdx == r.menuSelected {
				displayText = r.color.Colorize(text, "reverse")
			} else {
				switch item.Type {
				case "builtin":
					displayText = r.color.Colorize(text, "cyan")
				case "command":
					displayText = r.color.Colorize(text, "green")
				case "directory":
					displayText = r.color.Colorize(text, "blue")
				default:
					displayText = text
				}
			}

			// Calculate padding based on original text length (not colored text)
			padding := itemWidth - len(text)
			_, _ = fmt.Print(displayText + strings.Repeat(" ", padding)) //nolint:forbidigo
		}
		_, _ = fmt.Print("\r\n") //nolint:forbidigo
	}

	// Show pagination info
	if totalPages > 1 {
		pageInfo := fmt.Sprintf("Page %d/%d (%d items)", r.menuPage+1, totalPages, len(r.completionList))
		_, _ = fmt.Print(r.color.Colorize(pageInfo, "gray") + "\r\n") //nolint:forbidigo
	}

	// Restore cursor to original position
	_, _ = fmt.Print("\0338") //nolint:forbidigo // Restore cursor position

	// Mark menu as displayed
	r.menuDisplayed = true
	r.menuLinesDrawn = maxRows + 1 // menu rows + pagination line
}

// acceptMenuSelection accepts the selected completion.
func (r *Readline) acceptMenuSelection() {
	if !r.menuMode || r.menuSelected >= len(r.completionList) {
		return
	}

	selected := r.completionList[r.menuSelected].Text

	// Apply the selected completion
	words := strings.Fields(r.completionBase)
	r.buffer = r.applySelectedCompletion(selected, words)
	r.cursor = len(r.buffer)

	// Clear the menu before resetting completion
	r.clearCompletionMenu()

	// Redraw command line and position cursor at end
	_, _ = fmt.Print("\r\033[K") //nolint:forbidigo
	r.displayPrompt()
	_, _ = fmt.Print(string(r.buffer)) //nolint:forbidigo

	// Reset completion AFTER using completionBase
	r.resetCompletion()
}

// applySelectedCompletion applies the selected completion to create new buffer.
func (r *Readline) applySelectedCompletion(selected string, words []string) []rune {
	if len(words) == 1 && !strings.HasSuffix(r.completionBase, " ") {
		// Completing command
		return []rune(selected)
	}

	// Completing file - check for trailing space
	if strings.HasSuffix(r.completionBase, " ") {
		// Trailing space - append file to command
		return []rune(r.completionBase + selected)
	}

	// No trailing space - replace last word
	lastWord := ""
	if len(words) > 0 {
		lastWord = words[len(words)-1]
	}
	if strings.Contains(lastWord, "/") {
		return []rune(r.completionBase[:len(r.completionBase)-len(filepath.Base(lastWord))] + selected)
	}
	return []rune(r.completionBase[:len(r.completionBase)-len(lastWord)] + selected)
}

// applyCycleCompletion applies completion while keeping menu active for cycling.
func (r *Readline) applyCycleCompletion() {
	if !r.menuMode || r.menuSelected >= len(r.completionList) {
		return
	}

	selected := r.completionList[r.menuSelected].Text

	// Apply the selected completion but keep menu active
	words := strings.Fields(r.completionBase)
	r.buffer = r.applySelectedCompletion(selected, words)
	r.cursor = len(r.buffer)

	// Just redraw command line and menu - let it overwrite
	_, _ = fmt.Print("\r\033[K") //nolint:forbidigo
	r.displayPrompt()
	_, _ = fmt.Print(string(r.buffer)) //nolint:forbidigo

	// Save cursor position before showing menu
	_, _ = fmt.Print("\0337") //nolint:forbidigo // Save cursor position

	// Show menu (will overwrite previous)
	r.showCompletionMenu()
}

// navigateMenu moves selection up/down in completion menu.
func (r *Readline) navigateMenu(direction int) {
	if !r.menuMode || len(r.completionList) == 0 {
		return
	}

	if direction > 0 {
		r.menuSelected = (r.menuSelected + 1) % len(r.completionList)
	} else {
		r.menuSelected = (r.menuSelected - 1 + len(r.completionList)) % len(r.completionList)
	}

	r.showCompletionMenu()
}

// navigateMenuHorizontal moves selection left/right in completion menu.
func (r *Readline) navigateMenuHorizontal(direction int) {
	if !r.menuMode || len(r.completionList) == 0 {
		return
	}

	width, _ := r.terminal.GetTerminalSize()
	maxItemWidth := 0
	for _, item := range r.completionList {
		if len(item.Text) > maxItemWidth {
			maxItemWidth = len(item.Text)
		}
	}

	itemWidth := maxItemWidth + 2
	cols := width / itemWidth
	if cols < 1 {
		cols = 1
	}

	currentRow := r.menuSelected / cols
	currentCol := r.menuSelected % cols

	if direction > 0 {
		// Move right
		newCol := (currentCol + 1) % cols
		newSelected := currentRow*cols + newCol
		if newSelected < len(r.completionList) {
			r.menuSelected = newSelected
		}
	} else {
		// Move left
		newCol := (currentCol - 1 + cols) % cols
		newSelected := currentRow*cols + newCol
		if newSelected < len(r.completionList) {
			r.menuSelected = newSelected
		}
	}

	r.showCompletionMenu()
}

// clearCompletionMenu clears the displayed completion menu.
func (r *Readline) clearCompletionMenu() {
	if !r.menuDisplayed {
		return
	}

	// Clear from current position to end of screen
	_, _ = fmt.Print("\r\n\033[0J\033[A") //nolint:forbidigo // Move to next line, clear to end, move back up

	r.menuDisplayed = false
	r.menuLinesDrawn = 0
}

// resetCompletion clears completion state.
func (r *Readline) resetCompletion() {
	r.completionList = nil
	r.completionIdx = -1
	r.completionBase = ""
	r.menuMode = false
	r.menuSelected = 0
	r.menuDisplayed = false
	r.menuPage = 0
}

// History operations.
func (r *Readline) historyPrevious() {
	// Start browse mode if buffer is empty
	if len(r.buffer) == 0 && r.searchPrefix == "" {
		r.browseMode = true
		r.searchPrefix = ""
	}

	var line string
	if r.browseMode {
		// Use regular history navigation
		line = r.history.Previous()
	} else {
		// Use substring search - fallback to regular navigation
		if r.searchPrefix == "" {
			r.searchPrefix = string(r.buffer)
		}
		line = r.history.Previous()
	}

	if line != "" {
		r.setBufferFromHistory(line)
	}
}

func (r *Readline) historyNext() {
	// Use regular history navigation in browse mode
	if r.browseMode {
		line := r.history.Next()
		if line != "" {
			r.setBufferFromHistory(line)
		} else {
			// End of history - clear buffer and exit browse mode
			r.buffer = r.buffer[:0]
			r.cursor = 0
			r.browseMode = false
			r.searchPrefix = ""
			r.redraw()
		}
	} else if r.searchPrefix != "" {
		// Use regular navigation for now
		line := r.history.Next()
		if line != "" {
			r.setBufferFromHistory(line)
		}
	}
}

func (r *Readline) setBufferFromHistory(line string) {
	r.buffer = []rune(line)
	r.cursor = len(r.buffer)
	r.redraw()
}

// updateSuggestion updates the autosuggestion based on current input.
func (r *Readline) updateSuggestion() {
	// Disable autosuggestions - use tab completion instead
	r.suggestion = ""
}

// acceptSuggestion accepts the current autosuggestion.
func (r *Readline) acceptSuggestion() {
	if r.suggestion != "" {
		r.buffer = append(r.buffer, []rune(r.suggestion)...)
		r.cursor = len(r.buffer)
		r.suggestion = ""
		r.redraw()
	}
}
