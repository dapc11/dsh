package readline

import "fmt"

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
	r.redraw()
}

func (r *Readline) clearScreen() {
	_, _ = fmt.Print("\033[2J\033[H") //nolint:forbidigo // Clear screen and move to top
	r.displayPrompt()
	r.redraw()
}

func (r *Readline) redraw() {
	// Update suggestion before drawing
	r.updateSuggestion()

	// Clear line and redraw
	_, _ = fmt.Print("\r\033[K") //nolint:forbidigo
	r.displayPrompt()

	// Print buffer with suggestion
	_, _ = fmt.Print(string(r.buffer)) //nolint:forbidigo
	if r.suggestion != "" {
		_, _ = fmt.Print(r.color.Colorize(r.suggestion, Gray)) //nolint:forbidigo
	}

	// Position cursor correctly
	r.setCursorPosition()
}
