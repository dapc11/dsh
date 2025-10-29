package readline

import "dsh/internal/terminal"

func (r *Readline) insertRune(ch rune) {
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
	r.redraw()
}

func (r *Readline) clearScreen() {
	if r.terminal != nil {
		_, _ = r.terminal.WriteString("\033[2J\033[H") // Clear screen and move to top
		r.displayPrompt()
	}
	r.redraw()
}

func (r *Readline) redraw() {
	if r.terminal == nil {
		return
	}

	// Update suggestion before drawing
	r.updateSuggestion()

	// Clear line and redraw
	_, _ = r.terminal.WriteString("\r\033[K")
	r.displayPrompt()

	// Print buffer with suggestion
	_, _ = r.terminal.WriteString(string(r.buffer))
	if r.suggestion != "" {
		_, _ = r.terminal.WriteString(r.terminal.Colorize(r.suggestion, terminal.ColorBrightBlack))
	}

	// Position cursor correctly
	r.setCursorPosition()
}
