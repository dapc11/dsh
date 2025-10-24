package readline

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Key codes for special keys.
const (
	KeyCtrlA     = 1
	KeyCtrlB     = 2
	KeyCtrlC     = 3
	KeyCtrlD     = 4
	KeyCtrlE     = 5
	KeyCtrlF     = 6
	KeyCtrlK     = 11
	KeyCtrlL     = 12
	KeyCtrlN     = 14
	KeyCtrlP     = 16
	KeyCtrlU     = 21
	KeyCtrlW     = 23
	KeyCtrlY     = 25
	KeyCtrlZ     = 26
	KeyEscape    = 27
	KeyTab       = 9
	KeyBackspace = 127
	KeyDelete    = 127
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
		// Position cursor at end of prompt
		_, _ = fmt.Print("\r") //nolint:forbidigo
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
			return false // Will be handled as EOF in caller
		}
		r.deleteChar()
	case KeyCtrlK:
		r.killRing.ResetYank()
		r.killToEnd()
	case KeyCtrlU:
		r.killRing.ResetYank()
		r.killLine()
	case KeyCtrlW:
		r.killRing.ResetYank()
		r.killWordBackward()
	case KeyCtrlY:
		r.yank()
	case KeyCtrlL:
		r.killRing.ResetYank()
		r.clearScreen()
	case KeyCtrlP:
		r.killRing.ResetYank()
		r.historyPrevious()
	case KeyCtrlN:
		r.killRing.ResetYank()
		r.historyNext()
	case KeyCtrlZ:
		r.killRing.ResetYank()
		r.clearLine()
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
	case KeyBackspace:
		r.killRing.ResetYank()
		r.backspace()
		r.searchPrefix = ""
		r.browseMode = false // Exit browse mode when editing
		r.updateSuggestion()
	case KeyEscape:
		if r.menuMode {
			r.clearCompletionMenu()
			r.resetCompletion()
			r.redraw()
		} else {
			err := r.handleEscapeSequence()
			if err != nil {
				return true // Continue on error
			}
		}
	default:
		if ch >= 32 && ch < 127 {
			r.killRing.ResetYank()
			if r.menuMode {
				r.clearCompletionMenu()
			}
			r.insertChar(rune(ch))
			r.searchPrefix = ""  // Reset search prefix on new input
			r.browseMode = false // Exit browse mode when typing
			r.resetCompletion()  // Reset completion cycling
			r.updateSuggestion()
		}
	}

	return true // Continue reading
}

func (r *Readline) handleEscapeSequence() error { //nolint:cyclop // Escape sequence handling requires many branches
	ch1, err := r.readChar()
	if err != nil {
		return err
	}

	switch ch1 {
	case 127: // Alt+Backspace (ESC + DEL)
		r.killRing.ResetYank()
		r.killWordBackward()
	case '[':
		ch2, err := r.readChar()
		if err != nil {
			return err
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
	case 'y': // Alt+Y (forward)
		r.yankPop()
	case 'Y': // Alt+Shift+Y (backward)
		r.yankCycle(-1)
	case 'd': // Alt+D (delete word forward)
		r.killRing.ResetYank()
		r.killWordForward()
	}

	return nil
}

// Movement and editing functions.
func (r *Readline) moveCursorLeft() {
	if r.cursor > 0 {
		r.cursor--
		_, _ = fmt.Print("\b") //nolint:forbidigo
	}
}

func (r *Readline) moveCursorRight() {
	if r.cursor < len(r.buffer) {
		_, _ = fmt.Printf("%c", r.buffer[r.cursor]) //nolint:forbidigo
		r.cursor++
	}
}

func (r *Readline) moveCursorToStart() {
	for r.cursor > 0 {
		r.moveCursorLeft()
	}
}

func (r *Readline) moveCursorToEnd() {
	for r.cursor < len(r.buffer) {
		r.moveCursorRight()
	}
}

func (r *Readline) insertChar(ch rune) {
	if r.cursor == len(r.buffer) {
		r.buffer = append(r.buffer, ch)
	} else {
		r.buffer = append(r.buffer[:r.cursor+1], r.buffer[r.cursor:]...)
		r.buffer[r.cursor] = ch
	}
	r.cursor++
	r.redraw()
}

func (r *Readline) deleteChar() {
	if r.cursor < len(r.buffer) {
		r.buffer = append(r.buffer[:r.cursor], r.buffer[r.cursor+1:]...)
		r.redraw()
	}
}

func (r *Readline) backspace() {
	if r.cursor > 0 {
		r.buffer = append(r.buffer[:r.cursor-1], r.buffer[r.cursor:]...)
		r.cursor--
		r.redraw()
	}
}

func (r *Readline) clearLine() {
	r.buffer = r.buffer[:0]
	r.cursor = 0
}

func (r *Readline) redraw() {
	// Update suggestion before drawing
	r.updateSuggestion()

	// Clear entire line with sufficient space
	_, _ = fmt.Print("\r\033[K") //nolint:forbidigo // Clear to end of line

	// Display prompt and buffer
	_, _ = fmt.Print(r.prompt + string(r.buffer)) //nolint:forbidigo

	// Display autosuggestion with colors if available
	if r.suggestion != "" {
		suggestion := r.color.Colorize(r.suggestion, Gray)
		_, _ = fmt.Print(suggestion) //nolint:forbidigo
	}

	r.setCursorPosition()

	// Don't show menu here - it's handled by specific navigation functions
}

func (r *Readline) setCursorPosition() {
	// Move cursor to correct position using ANSI escape sequences
	promptLen := len(r.prompt)
	cursorPos := promptLen + r.cursor
	_, _ = fmt.Printf("\r\033[%dC", cursorPos) //nolint:forbidigo // Move cursor to position
}

func (r *Readline) clearScreen() {
	_, _ = fmt.Print("\033[2J\033[H") //nolint:forbidigo // Clear screen and move to top
	r.displayPrompt()
	_, _ = fmt.Print(string(r.buffer)) //nolint:forbidigo
	r.setCursorPosition()
}

// Kill and yank operations.
func (r *Readline) killToEnd() {
	if r.cursor < len(r.buffer) {
		killed := string(r.buffer[r.cursor:])
		r.killRing.Add(killed)
		r.buffer = r.buffer[:r.cursor]
		_, _ = fmt.Print("\033[K") //nolint:forbidigo // Clear to end of line
	}
}

func (r *Readline) killLine() {
	killed := string(r.buffer)
	r.killRing.Add(killed)
	r.buffer = r.buffer[:0]
	r.cursor = 0
	r.redraw()
}

func (r *Readline) killWordBackward() {
	if r.cursor == 0 {
		return
	}

	start := r.cursor - 1
	for start > 0 && r.buffer[start] == ' ' {
		start--
	}
	for start > 0 && r.buffer[start] != ' ' {
		start--
	}
	if r.buffer[start] == ' ' {
		start++
	}

	killed := string(r.buffer[start:r.cursor])
	r.killRing.Add(killed)
	r.buffer = append(r.buffer[:start], r.buffer[r.cursor:]...)
	r.cursor = start
	r.redraw()
}

func (r *Readline) killWordForward() {
	if r.cursor >= len(r.buffer) {
		return
	}

	end := r.cursor
	// Skip whitespace
	for end < len(r.buffer) && r.buffer[end] == ' ' {
		end++
	}
	// Skip word characters
	for end < len(r.buffer) && r.buffer[end] != ' ' {
		end++
	}

	if end > r.cursor {
		killed := string(r.buffer[r.cursor:end])
		r.killRing.Add(killed)
		r.buffer = append(r.buffer[:r.cursor], r.buffer[end:]...)
		r.redraw()
	}
}

func (r *Readline) yank() {
	yankText := r.killRing.Yank()
	if yankText == "" {
		return
	}

	yankRunes := []rune(yankText)
	r.killRing.SetLastYank(len(yankRunes))

	if r.cursor >= len(r.buffer) {
		// At end of buffer
		r.buffer = append(r.buffer, yankRunes...)
		_, _ = fmt.Print(yankText) //nolint:forbidigo
	} else {
		// In middle of buffer - make space and insert
		newBuffer := make([]rune, len(r.buffer)+len(yankRunes))
		copy(newBuffer, r.buffer[:r.cursor])
		copy(newBuffer[r.cursor:], yankRunes)
		copy(newBuffer[r.cursor+len(yankRunes):], r.buffer[r.cursor:])
		r.buffer = newBuffer
		r.redraw()
	}
	r.cursor += len(yankRunes)
}

func (r *Readline) yankPop() {
	r.yankCycle(1) // Forward direction
}

func (r *Readline) yankCycle(direction int) {
	lastYank := r.killRing.GetLastYank()
	if lastYank == 0 {
		return
	}

	yankText := r.killRing.Cycle(direction)
	if yankText == "" {
		return
	}

	// Calculate start position of previously yanked text
	startPos := r.cursor - lastYank
	if startPos < 0 {
		startPos = 0
	}

	// Remove the previously yanked text
	r.buffer = append(r.buffer[:startPos], r.buffer[r.cursor:]...)
	r.cursor = startPos

	yankRunes := []rune(yankText)
	r.killRing.SetLastYank(len(yankRunes))

	// Insert new text at cursor position
	if r.cursor >= len(r.buffer) {
		// At end of buffer
		r.buffer = append(r.buffer, yankRunes...)
	} else {
		// In middle of buffer - make space and insert
		newBuffer := make([]rune, len(r.buffer)+len(yankRunes))
		copy(newBuffer, r.buffer[:r.cursor])
		copy(newBuffer[r.cursor:], yankRunes)
		copy(newBuffer[r.cursor+len(yankRunes):], r.buffer[r.cursor:])
		r.buffer = newBuffer
	}

	r.cursor += len(yankRunes)
	r.redraw()
}

// Word movement.
func (r *Readline) moveWordForward() {
	for r.cursor < len(r.buffer) && r.buffer[r.cursor] != ' ' {
		r.moveCursorRight()
	}
	for r.cursor < len(r.buffer) && r.buffer[r.cursor] == ' ' {
		r.moveCursorRight()
	}
}

func (r *Readline) moveWordBackward() {
	if r.cursor == 0 {
		return
	}

	r.moveCursorLeft()
	for r.cursor > 0 && r.buffer[r.cursor] == ' ' {
		r.moveCursorLeft()
	}
	for r.cursor > 0 && r.buffer[r.cursor] != ' ' {
		r.moveCursorLeft()
	}
	if r.cursor > 0 && r.buffer[r.cursor] == ' ' {
		r.moveCursorRight()
	}
}

// History operations.
func (r *Readline) historyPrevious() {
	// Start browse mode if buffer is empty
	if len(r.buffer) == 0 && r.searchPrefix == "" {
		r.browseMode = true
	}

	// Use regular history navigation in browse mode
	if r.browseMode {
		line := r.history.Previous()
		if line != "" {
			r.setBufferFromHistory(line)
		}
		return
	}

	// Use prefix search if there's input
	prefix := r.searchPrefix
	if prefix == "" {
		prefix = string(r.buffer)
	}

	line := r.history.PreviousWithPrefix(prefix)
	if line != "" {
		r.searchPrefix = prefix
		r.setBufferFromHistory(line)
	}
}

func (r *Readline) historyNext() {
	// Use regular history navigation in browse mode
	if r.browseMode {
		line := r.history.Next()
		r.setBufferFromHistory(line)
		return
	}

	// Use prefix search if there's input
	prefix := r.searchPrefix
	if prefix == "" {
		prefix = string(r.buffer)
	}

	line := r.history.NextWithPrefix(prefix)
	r.searchPrefix = prefix
	r.setBufferFromHistory(line)
}

func (r *Readline) setBufferFromHistory(line string) {
	r.buffer = []rune(line)
	r.cursor = len(r.buffer)
	r.suggestion = ""
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
		r.updateSuggestion()
		r.redraw()
	}
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
func (r *Readline) showCompletionMenu() {
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
	rows := (len(pageItems) + cols - 1) / cols

	// Clear previous menu
	if r.menuDisplayed {
		clearLines := maxRows + 3 // menu + pagination info
		for i := 0; i < clearLines; i++ {
			_, _ = fmt.Print("\033[A\033[K") //nolint:forbidigo
		}
	}

	r.menuDisplayed = true

	// Show menu
	_, _ = fmt.Print("\r\n") //nolint:forbidigo

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			idx := row + col*rows
			if idx >= len(pageItems) {
				break
			}

			globalIdx := startIdx + idx
			item := pageItems[idx]
			text := item.Text
			displayText := text

			// Apply type-based coloring
			switch item.Type {
			case "builtin":
				displayText = r.color.Colorize(text, Yellow)
			case "command":
				displayText = r.color.Colorize(text, Green)
			case "directory":
				displayText = r.color.Colorize(text, Blue)
			case "file":
				displayText = r.color.Colorize(text, Cyan)
			}

			// Highlight selected item
			if globalIdx == r.menuSelected {
				displayText = "\033[7m" + displayText + "\033[0m"
			}

			padding := itemWidth - len(text)
			_, _ = fmt.Print(displayText + strings.Repeat(" ", padding)) //nolint:forbidigo
		}
		_, _ = fmt.Print("\r\n") //nolint:forbidigo
	}

	// Show pagination info
	if totalPages > 1 {
		pageInfo := fmt.Sprintf("Page %d/%d (%d items)", r.menuPage+1, totalPages, len(r.completionList))
		_, _ = fmt.Print(r.color.Colorize(pageInfo, Gray) + "\r\n") //nolint:forbidigo
	}

	// Restore cursor to original position
	_, _ = fmt.Print("\0338") //nolint:forbidigo // Restore cursor position

	// Mark menu as displayed
	r.menuDisplayed = true
	r.menuLinesDrawn = maxRows + 1 // menu rows + pagination line
}

// applyCompletion applies a completion option.
func (r *Readline) applyCompletion(match string) {
	// Restore base and apply new completion
	r.buffer = []rune(r.completionBase)

	// Find the part to complete
	words := strings.Fields(r.completionBase)
	if len(words) == 0 {
		return
	}

	if len(words) == 1 && !strings.HasSuffix(r.completionBase, " ") {
		// Completing command
		r.buffer = []rune(match)
	} else {
		// Completing file
		lastWord := ""
		if !strings.HasSuffix(r.completionBase, " ") && len(words) > 0 {
			lastWord = words[len(words)-1]
		}

		if strings.Contains(lastWord, "/") {
			r.buffer = []rune(r.completionBase[:len(r.completionBase)-len(filepath.Base(lastWord))] + match)
		} else {
			r.buffer = []rune(r.completionBase[:len(r.completionBase)-len(lastWord)] + match)
		}
	}

	r.cursor = len(r.buffer)
	r.redraw()
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
	rows := (len(r.completionList) + cols - 1) / cols

	currentRow := r.menuSelected % rows
	currentCol := r.menuSelected / rows

	if direction > 0 {
		currentCol = (currentCol + 1) % cols
	} else {
		currentCol = (currentCol - 1 + cols) % cols
	}

	newIdx := currentRow + currentCol*rows
	if newIdx < len(r.completionList) {
		r.menuSelected = newIdx
		r.showCompletionMenu()
	}
}

// acceptMenuSelection accepts the selected completion.
func (r *Readline) acceptMenuSelection() {
	if !r.menuMode || r.menuSelected >= len(r.completionList) {
		return
	}

	selected := r.completionList[r.menuSelected].Text
	r.resetCompletion()

	// Apply the selected completion
	words := strings.Fields(r.completionBase)
	if len(words) == 1 && !strings.HasSuffix(r.completionBase, " ") {
		// Completing command
		r.buffer = []rune(selected)
	} else {
		// Completing file
		lastWord := ""
		if !strings.HasSuffix(r.completionBase, " ") && len(words) > 0 {
			lastWord = words[len(words)-1]
		}

		if strings.Contains(lastWord, "/") {
			r.buffer = []rune(r.completionBase[:len(r.completionBase)-len(filepath.Base(lastWord))] + selected)
		} else {
			r.buffer = []rune(r.completionBase[:len(r.completionBase)-len(lastWord)] + selected)
		}
	}

	r.cursor = len(r.buffer)
	r.redraw()
}

// applyCycleCompletion applies completion while keeping menu active for cycling.
func (r *Readline) applyCycleCompletion() {
	if !r.menuMode || r.menuSelected >= len(r.completionList) {
		return
	}

	selected := r.completionList[r.menuSelected].Text

	// Apply the selected completion but keep menu active
	words := strings.Fields(r.completionBase)
	if len(words) == 1 && !strings.HasSuffix(r.completionBase, " ") {
		// Completing command
		r.buffer = []rune(selected)
	} else {
		// Completing file
		lastWord := ""
		if !strings.HasSuffix(r.completionBase, " ") && len(words) > 0 {
			lastWord = words[len(words)-1]
		}

		if strings.Contains(lastWord, "/") {
			r.buffer = []rune(r.completionBase[:len(r.completionBase)-len(filepath.Base(lastWord))] + selected)
		} else {
			r.buffer = []rune(r.completionBase[:len(r.completionBase)-len(lastWord)] + selected)
		}
	}

	r.cursor = len(r.buffer)

	// Just redraw command line and menu - let it overwrite
	_, _ = fmt.Print("\r\033[K") //nolint:forbidigo
	r.displayPrompt()
	_, _ = fmt.Print(string(r.buffer)) //nolint:forbidigo

	// Show menu (will overwrite previous)
	r.showCompletionMenu()
}

// clearCompletionMenu clears the displayed completion menu.
func (r *Readline) clearCompletionMenu() {
	if !r.menuDisplayed {
		return
	}

	// Move to start of next line (where menu begins) and clear to end
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
