package readline

// KillRing manages killed text for yank operations.
type KillRing struct {
	items    []string
	index    int
	lastYank int
}

// NewKillRing creates a new kill ring.
func NewKillRing() *KillRing {
	return &KillRing{
		items:    make([]string, 0, 10),
		index:    0,
		lastYank: 0,
	}
}

// Add adds text to the kill ring.
func (k *KillRing) Add(text string) {
	if text == "" {
		return
	}

	// Add to front of kill ring
	k.items = append([]string{text}, k.items...)

	// Limit kill ring size
	if len(k.items) > 10 {
		k.items = k.items[:10]
	}

	k.index = 0
}

// Yank returns the current kill ring item.
func (k *KillRing) Yank() string {
	if len(k.items) == 0 {
		return ""
	}

	k.index = 0

	return k.items[k.index]
}

// Cycle moves to the next item in the kill ring.
func (k *KillRing) Cycle(direction int) string {
	if len(k.items) <= 1 {
		return ""
	}

	// Cycle through kill ring (with proper wrap-around)
	newIndex := k.index + direction
	if newIndex < 0 {
		newIndex = len(k.items) - 1
	} else if newIndex >= len(k.items) {
		newIndex = 0
	}
	k.index = newIndex

	return k.items[k.index]
}

// SetLastYank sets the length of the last yanked text.
func (k *KillRing) SetLastYank(length int) {
	k.lastYank = length
}

// GetLastYank returns the length of the last yanked text.
func (k *KillRing) GetLastYank() int {
	return k.lastYank
}

// ResetYank resets the yank state.
func (k *KillRing) ResetYank() {
	k.lastYank = 0
}
