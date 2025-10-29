package readline

import (
	"fmt"
	"strings"

	"dsh/internal/terminal"
)

const (
	itemTypeBuiltin   = "builtin"
	itemTypeCommand   = "command"
	itemTypeDirectory = "directory"
)

func (r *Readline) handleKeyEvent(keyEvent terminal.KeyEvent) bool { //nolint:cyclop,funlen // Key handling naturally requires many branches and statements
	// Handle printable characters first
	if keyEvent.Rune != 0 && keyEvent.Key == terminal.KeyNone {
		r.killRing.ResetYank()
		if r.completionMenu.IsActive() {
			r.clearTabCompletion()
		}
		r.insertRune(keyEvent.Rune)
		r.updateSuggestion()
		return true
	}

	switch keyEvent.Key {
	case terminal.KeyEnter:
		if r.completionMenu.IsActive() {
			r.acceptTabCompletion()
			return true
		}
		return false // Signal completion
	case terminal.KeyCtrlA:
		r.killRing.ResetYank()
		r.moveCursorToStart()
	case terminal.KeyCtrlE:
		r.killRing.ResetYank()
		r.moveCursorToEnd()
	case terminal.KeyCtrlB:
		r.killRing.ResetYank()
		r.moveCursorLeft()
	case terminal.KeyCtrlC:
		r.killRing.ResetYank()
		if r.completionMenu.IsActive() {
			r.clearTabCompletion()
		}
		r.clearLine()
		r.terminal.WriteString("^C\r\n")
		r.displayPrompt()
	case terminal.KeyCtrlF:
		r.killRing.ResetYank()
		if r.suggestion != "" {
			r.acceptSuggestion()
		} else {
			r.moveCursorRight()
		}
	case terminal.KeyCtrlD:
		r.killRing.ResetYank()
		if len(r.buffer) == 0 {
			// EOF case - let ReadLine handle it
			return false
		}
		r.deleteChar()
	case terminal.KeyCtrlK:
		r.killToEnd()
	case terminal.KeyCtrlL:
		r.killRing.ResetYank()
		r.clearScreen()
	case terminal.KeyCtrlN:
		r.killRing.ResetYank()
		r.historyNext()
	case terminal.KeyCtrlP:
		r.killRing.ResetYank()
		r.historyPrevious()
	case terminal.KeyCtrlR:
		r.killRing.ResetYank()
		// Clear current line before fuzzy search
		r.moveCursorToStart()
		r.terminal.WriteString("\033[K") // Clear line
		if selected := r.FuzzyHistorySearchCustom(); selected != "" {
			r.buffer = []rune(selected)
			r.cursor = len(r.buffer)
			r.redraw()
		} else {
			// Restore prompt if cancelled
			r.redraw()
		}
	case terminal.KeyCtrlU:
		r.killLine()
	case terminal.KeyCtrlW:
		r.killWordBackward()
	case terminal.KeyCtrlY:
		r.yank()
	case terminal.KeyCtrlZ:
		r.killRing.ResetYank()
		r.terminal.WriteString("^Z\r\n")
		r.displayPrompt()
	case terminal.KeyTab:
		r.killRing.ResetYank()
		if r.completionMenu.IsActive() {
			r.navigateTabCompletion(1)
		} else {
			r.handleTabCompletion()
		}
	case terminal.KeyEscape:
		// Always try to handle escape sequences first (including shift-tab)
		err := r.handleEscapeSequence()
		if err != nil {
			// If escape sequence handling fails and we're in menu mode, exit menu
			if r.completionMenu.IsActive() {
				r.clearTabCompletion()
				r.redraw()
			}
		}
	case terminal.KeyBackspace:
		r.killRing.ResetYank()
		r.backspace()
		r.searchPrefix = ""
		r.browseMode = false // Exit browse mode when editing
		r.updateSuggestion()
	case terminal.KeyArrowUp:
		r.killRing.ResetYank()
		r.historyPrevious()
	case terminal.KeyArrowDown:
		r.killRing.ResetYank()
		r.historyNext()
	case terminal.KeyArrowLeft:
		r.killRing.ResetYank()
		r.moveCursorLeft()
	case terminal.KeyArrowRight:
		r.killRing.ResetYank()
		if r.suggestion != "" {
			r.acceptSuggestion()
		} else {
			r.moveCursorRight()
		}
	default:
		// Handle unimplemented keys
	}

	return true
}

func (r *Readline) handleEscapeSequence() error { //nolint:gocognit,cyclop,funlen // Terminal input handling requires complexity
	// Read additional key events for complex sequences
	keyEvent, err := r.terminal.ReadKey()
	if err != nil {
		return fmt.Errorf("failed to read escape sequence: %w", err)
	}

	// Handle Alt+key combinations
	if keyEvent.Alt {
		switch keyEvent.Rune {
		case 127: // Alt+Backspace
			r.killRing.ResetYank()
			r.killWordBackward()
		case 'd': // Alt+D
			r.killWordForward()
		case 'y': // Alt+Y
			r.yankPop()
		}
		return nil
	}

	// Handle special sequences
	switch keyEvent.Key {
	case terminal.KeyArrowUp:
		if r.completionMenu.IsActive() {
			r.navigateTabCompletion(-1)
		} else {
			r.historyPrevious()
		}
	case terminal.KeyArrowDown:
		if r.completionMenu.IsActive() {
			r.navigateTabCompletion(1)
		} else {
			r.historyNext()
		}
	case terminal.KeyArrowLeft:
		if r.completionMenu.IsActive() {
			r.navigateTabCompletion(-1)
		} else {
			r.moveCursorLeft()
		}
	case terminal.KeyArrowRight:
		if r.completionMenu.IsActive() {
			r.navigateTabCompletion(1)
		} else {
			r.moveCursorRight()
		}
	default:
		// Handle unimplemented keys
	}

	return nil
}

// handleTabCompletion handles tab key press for completion.
func (r *Readline) handleTabCompletion() {
	input := string(r.buffer)
	matches, completion := r.completion.Complete(input, r.cursor)

	if len(matches) == 1 && completion != "" {
		// Single match - apply directly
		completionRunes := []rune(completion)
		r.buffer = append(r.buffer, completionRunes...)
		r.cursor = len(r.buffer)
	} else if len(matches) > 1 {
		if r.completionMenu.IsActive() {
			// Menu already active - just navigate, don't re-render everything
			r.navigateTabCompletion(1)
		} else {
			// Show menu for first time only
			r.completionMenu.Show(matches)
			r.completionMenu.Render(r.bufferManager, r.terminal)
		}
	}
}

// clearTabCompletion clears the tab completion menu.
func (r *Readline) clearTabCompletion() {
	if r.completionMenu.IsActive() {
		r.completionMenu.Hide()
		if r.bufferManager != nil {
			r.bufferManager.CleanupAll()
		}

		// Just restore cursor and clear menu area - don't redraw prompt
		r.terminal.WriteString("\033[u")         // Restore cursor to saved position
		r.terminal.WriteString("\033[J")         // Clear from cursor to end of screen
		r.terminal.WriteString(string(r.buffer)) // Redraw just the buffer
		r.setCursorPosition()
	}
}

// navigateTabCompletion navigates the completion menu.
func (r *Readline) navigateTabCompletion(direction int) {
	if !r.completionMenu.IsActive() {
		return
	}

	if direction > 0 {
		r.completionMenu.Next()
	} else {
		r.completionMenu.Prev()
	}

	// Use incremental update instead of full re-render
	r.completionMenu.UpdateSelectionOnly()
}

// acceptTabCompletion accepts the current completion selection.
func (r *Readline) acceptTabCompletion() {
	if !r.completionMenu.IsActive() {
		return
	}

	selected, ok := r.completionMenu.GetSelected()
	if !ok {
		return
	}

	// Apply completion
	input := string(r.buffer)
	words := strings.Fields(input)

	if len(words) == 1 && !strings.HasSuffix(input, " ") {
		// Completing command
		r.buffer = []rune(selected.Text)
	} else {
		// Completing file/argument
		if strings.HasSuffix(input, " ") {
			r.buffer = []rune(input + selected.Text)
		} else {
			lastWord := ""
			if len(words) > 0 {
				lastWord = words[len(words)-1]
			}
			r.buffer = []rune(input[:len(input)-len(lastWord)] + selected.Text)
		}
	}

	r.cursor = len(r.buffer)

	// Hide completion menu and cleanup buffers
	r.completionMenu.Hide()
	if r.bufferManager != nil {
		r.bufferManager.CleanupAll()
	}

	// Only clear the completion menu area, not the entire screen
	// Move cursor back to the input line and clear from there down
	r.terminal.WriteString("\033[u") // Restore to saved cursor position (input line)
	r.terminal.WriteString("\033[J") // Clear from cursor to end of screen (removes menu)
	// Redraw the current input line cleanly
	r.terminal.WriteString("\033[2K") // Clear the current line
	r.terminal.WriteString("\r")      // Move to beginning of line
	fullLine := r.prompt + string(r.buffer)
	r.terminal.WriteString(fullLine)
}

// History operations.
func (r *Readline) historyPrevious() {
	if len(r.buffer) == 0 && r.searchPrefix == "" {
		r.browseMode = true
		r.searchPrefix = ""
	}

	var line string
	if r.browseMode {
		line = r.history.Previous()
	} else {
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
	if r.browseMode {
		line := r.history.Next()
		if line != "" {
			r.setBufferFromHistory(line)
		} else {
			r.buffer = r.buffer[:0]
			r.cursor = 0
			r.browseMode = false
			r.searchPrefix = ""
			r.redraw()
		}
	} else if r.searchPrefix != "" {
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

func (r *Readline) updateSuggestion() {
	r.suggestion = ""
}

func (r *Readline) acceptSuggestion() {
	if r.suggestion != "" {
		r.buffer = append(r.buffer, []rune(r.suggestion)...)
		r.cursor = len(r.buffer)
		r.suggestion = ""
		r.redraw()
	}
}
