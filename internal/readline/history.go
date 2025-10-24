package readline

// History manages command history.
type History struct {
	items []string
	pos   int
}

// NewHistory creates a new history manager.
func NewHistory() *History {
	return &History{
		items: make([]string, 0, 100),
		pos:   0,
	}
}

// Add adds a command to history.
func (h *History) Add(line string) {
	if len(h.items) == 0 || h.items[len(h.items)-1] != line {
		h.items = append(h.items, line)
		if len(h.items) > 100 {
			h.items = h.items[1:]
		}
	}
}

// ResetPosition resets history position to end.
func (h *History) ResetPosition() {
	h.pos = len(h.items)
}

// Previous moves to previous history item.
func (h *History) Previous() string {
	if h.pos > 0 {
		h.pos--

		return h.items[h.pos]
	}

	return ""
}

// Next moves to next history item.
func (h *History) Next() string {
	if h.pos < len(h.items) {
		h.pos++
		if h.pos < len(h.items) {
			return h.items[h.pos]
		}
	}

	return ""
}
