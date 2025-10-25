package readline

import "fmt"

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

func (r *Readline) setCursorPosition() {
	// Move cursor to correct position using ANSI escape sequences
	promptLen := len(r.prompt)
	totalPos := promptLen + r.cursor + 1
	_, _ = fmt.Printf("\033[%dG", totalPos) //nolint:forbidigo
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

	// Move back one position first
	r.moveCursorLeft()

	// Skip spaces
	for r.cursor > 0 && r.buffer[r.cursor] == ' ' {
		r.moveCursorLeft()
	}

	// Move to start of word
	for r.cursor > 0 && r.buffer[r.cursor-1] != ' ' {
		r.moveCursorLeft()
	}
}
