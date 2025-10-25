package completion

// Item represents a completion with its type.
type Item struct {
	Text string
	Type string // "builtin", "command", "file", "directory"
}

// Menu handles the display and navigation of completion options.
type Menu struct {
	items      []Item
	selected   int
	displayed  bool
	linesDrawn int
	maxRows    int
	page       int
	base       string
}

// NewMenu creates a new completion menu.
func NewMenu() *Menu {
	return &Menu{
		maxRows: 10,
	}
}

// Show displays the completion menu with the given items.
func (m *Menu) Show(items []Item, base string) {
	m.items = items
	m.base = base
	m.selected = 0
	m.displayed = true
	m.page = 0
}

// Hide clears the completion menu.
func (m *Menu) Hide() {
	m.displayed = false
	m.items = nil
	m.selected = 0
	m.linesDrawn = 0
}

// IsDisplayed returns whether the menu is currently shown.
func (m *Menu) IsDisplayed() bool {
	return m.displayed
}

// GetSelected returns the currently selected item.
func (m *Menu) GetSelected() (Item, bool) {
	if !m.displayed || m.selected >= len(m.items) {
		return Item{}, false
	}
	return m.items[m.selected], true
}

// GetBase returns the completion base string.
func (m *Menu) GetBase() string {
	return m.base
}

// NextItem moves selection to the next item.
func (m *Menu) NextItem() {
	if len(m.items) > 0 {
		m.selected = (m.selected + 1) % len(m.items)
	}
}

// PrevItem moves selection to the previous item.
func (m *Menu) PrevItem() {
	if len(m.items) > 0 {
		m.selected = (m.selected - 1 + len(m.items)) % len(m.items)
	}
}

// HasItems returns whether the menu has any items.
func (m *Menu) HasItems() bool {
	return len(m.items) > 0
}
