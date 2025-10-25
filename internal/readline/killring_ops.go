package readline

// Kill and yank operations.
func (r *Readline) killToEnd() {
	if r.cursor < len(r.buffer) {
		killed := string(r.buffer[r.cursor:])
		r.killRing.Add(killed)
		r.buffer = r.buffer[:r.cursor]
		r.redraw()
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

	start := r.cursor
	// Move back to find word boundary
	for r.cursor > 0 && r.buffer[r.cursor-1] == ' ' {
		r.cursor--
	}
	for r.cursor > 0 && r.buffer[r.cursor-1] != ' ' {
		r.cursor--
	}

	if r.cursor < start {
		killed := string(r.buffer[r.cursor:start])
		r.killRing.Add(killed)
		r.buffer = append(r.buffer[:r.cursor], r.buffer[start:]...)
		r.redraw()
	}
}

func (r *Readline) killWordForward() {
	if r.cursor >= len(r.buffer) {
		return
	}

	start := r.cursor
	// Move forward to find word boundary
	for r.cursor < len(r.buffer) && r.buffer[r.cursor] != ' ' {
		r.cursor++
	}
	for r.cursor < len(r.buffer) && r.buffer[r.cursor] == ' ' {
		r.cursor++
	}

	if r.cursor > start {
		killed := string(r.buffer[start:r.cursor])
		r.killRing.Add(killed)
		r.buffer = append(r.buffer[:start], r.buffer[r.cursor:]...)
		r.cursor = start
		r.redraw()
	}
}

func (r *Readline) yank() {
	yankText := r.killRing.Yank()
	if yankText == "" {
		return
	}

	// Insert yanked text at cursor
	yankRunes := []rune(yankText)
	if r.cursor == len(r.buffer) {
		r.buffer = append(r.buffer, yankRunes...)
	} else {
		// Insert in middle
		newBuffer := make([]rune, len(r.buffer)+len(yankRunes))
		copy(newBuffer, r.buffer[:r.cursor])
		copy(newBuffer[r.cursor:], yankRunes)
		copy(newBuffer[r.cursor+len(yankRunes):], r.buffer[r.cursor:])
		r.buffer = newBuffer
	}

	r.cursor += len(yankRunes)
	r.killRing.SetLastYank(len(yankRunes))
	r.redraw()
}

func (r *Readline) yankPop() {
	r.yankCycle(1) // Forward direction
}

func (r *Readline) yankCycle(direction int) {
	lastYank := r.killRing.GetLastYank()
	if lastYank == 0 {
		return // No previous yank
	}

	// Remove the previously yanked text
	start := r.cursor - lastYank
	if start < 0 {
		start = 0
	}
	r.buffer = append(r.buffer[:start], r.buffer[r.cursor:]...)
	r.cursor = start

	// Get next item from kill ring using Cycle
	yankText := r.killRing.Cycle(direction)

	if yankText == "" {
		return
	}

	// Insert new yanked text
	yankRunes := []rune(yankText)
	if r.cursor == len(r.buffer) {
		r.buffer = append(r.buffer, yankRunes...)
	} else {
		newBuffer := make([]rune, len(r.buffer)+len(yankRunes))
		copy(newBuffer, r.buffer[:r.cursor])
		copy(newBuffer[r.cursor:], yankRunes)
		copy(newBuffer[r.cursor+len(yankRunes):], r.buffer[r.cursor:])
		r.buffer = newBuffer
	}

	r.cursor += len(yankRunes)
	r.killRing.SetLastYank(len(yankRunes))
	r.redraw()
}
