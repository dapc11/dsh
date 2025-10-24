package readline

import (
	"fmt"
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
		r.acceptSuggestion()
	case KeyBackspace:
		r.killRing.ResetYank()
		r.backspace()
		r.searchPrefix = ""
		r.browseMode = false // Exit browse mode when editing
		r.updateSuggestion()
	case KeyEscape:
		err := r.handleEscapeSequence()
		if err != nil {
			return true // Continue on error
		}
	default:
		if ch >= 32 && ch < 127 {
			r.killRing.ResetYank()
			r.insertChar(rune(ch))
			r.searchPrefix = ""  // Reset search prefix on new input
			r.browseMode = false // Exit browse mode when typing
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
			r.historyPrevious()
		case 'B': // Down arrow
			r.historyNext()
		case 'C': // Right arrow
			r.moveCursorRight()
		case 'D': // Left arrow
			r.moveCursorLeft()
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
		_, _ = fmt.Printf("%c", ch) //nolint:forbidigo
	} else {
		r.buffer = append(r.buffer[:r.cursor+1], r.buffer[r.cursor:]...)
		r.buffer[r.cursor] = ch
		r.redraw()
	}
	r.cursor++
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

	// Display autosuggestion with simple text indicator
	if r.suggestion != "" {
		_, _ = fmt.Print(" -> " + r.suggestion + " [TAB]") //nolint:forbidigo // Simple text suggestion
	}

	r.setCursorPosition()
}

func (r *Readline) setCursorPosition() {
	if r.cursor < len(r.buffer) {
		_, _ = fmt.Printf("\r%s%s\r%s%s", //nolint:forbidigo
			r.prompt,
			string(r.buffer),
			r.prompt,
			string(r.buffer[:r.cursor]))
	}
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
	r.suggestion = ""
	
	if len(r.buffer) == 0 {
		return
	}

	currentInput := string(r.buffer)
	r.suggestion = r.history.GetSuggestion(currentInput)
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
