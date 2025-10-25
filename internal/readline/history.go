package readline

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// History manages command history with persistent storage.
type History struct {
	items    []string
	pos      int
	file     string
	maxSize  int
	modified bool
}

// NewHistory creates a new history manager with persistent storage.
func NewHistory() *History {
	homeDir, _ := os.UserHomeDir()
	histFile := filepath.Join(homeDir, ".dsh_history")

	h := &History{
		items:    make([]string, 0, 1000),
		pos:      0,
		file:     histFile,
		maxSize:  1000,
		modified: false,
	}

	h.load()

	return h
}

// Add adds a command to history and saves to disk.
func (h *History) Add(line string) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return
	}

	// Remove any existing occurrence of this command
	for i := len(h.items) - 1; i >= 0; i-- {
		if h.items[i] == line {
			h.items = append(h.items[:i], h.items[i+1:]...)
			break
		}
	}

	// Add to end (most recent)
	h.items = append(h.items, line)
	h.pos = len(h.items)
	h.modified = true

	// Trim history if too large
	if len(h.items) > h.maxSize {
		h.items = h.items[len(h.items)-h.maxSize:]
		h.pos = len(h.items)
	}

	// Save immediately for shared sessions
	h.save()
}

// Previous moves to previous history item.
func (h *History) Previous() string {
	if h.pos > 0 {
		h.pos--
		return h.items[h.pos]
	}

	// Return first item if we're at the beginning
	if len(h.items) > 0 {
		return h.items[0]
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

// Search performs substring search in history (zsh-like).
func (h *History) Search(query string) []string {
	var matches []string
	query = strings.ToLower(query)

	// Search backwards through history
	for i := len(h.items) - 1; i >= 0; i-- {
		if strings.Contains(strings.ToLower(h.items[i]), query) {
			matches = append(matches, h.items[i])
			if len(matches) >= 10 { // Limit results
				break
			}
		}
	}

	return matches
}

// GetSuggestion returns auto-suggestion based on current input.
func (h *History) GetSuggestion(input string) string {
	if input == "" {
		return ""
	}

	input = strings.ToLower(input)

	// Find most recent command starting with input
	for i := len(h.items) - 1; i >= 0; i-- {
		cmd := strings.ToLower(h.items[i])
		if strings.HasPrefix(cmd, input) && len(h.items[i]) > len(input) {
			return h.items[i][len(input):] // Return the suggestion part
		}
	}

	return ""
}

// Reload refreshes history from disk (for shared sessions).
func (h *History) Reload() {
	currentPos := h.pos
	h.items = h.items[:0] // Clear current items
	h.load()

	// Restore position if possible
	if currentPos < len(h.items) {
		h.pos = currentPos
	} else {
		h.pos = len(h.items)
	}
}

// ResetPosition resets history position to end.
func (h *History) ResetPosition() {
	h.pos = len(h.items)
}

// PreviousWithPrefix finds previous history item starting with prefix.
func (h *History) PreviousWithPrefix(prefix string) string {
	if prefix == "" {
		return h.Previous()
	}

	for i := h.pos - 1; i >= 0; i-- {
		if strings.HasPrefix(h.items[i], prefix) {
			h.pos = i

			return h.items[i]
		}
	}

	return ""
}

// NextWithPrefix finds next history item starting with prefix.
func (h *History) NextWithPrefix(prefix string) string {
	if prefix == "" {
		return h.Next()
	}

	for i := h.pos + 1; i < len(h.items); i++ {
		if strings.HasPrefix(h.items[i], prefix) {
			h.pos = i

			return h.items[i]
		}
	}

	// If no match found, go to end
	h.pos = len(h.items)

	return ""
}

// load reads history from disk.
func (h *History) load() {
	file, err := os.Open(h.file)
	if err != nil {
		return // File doesn't exist yet
	}
	defer func() {
		_ = file.Close() // Ignore close error on read-only file
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			h.items = append(h.items, line)
		}
	}

	h.pos = len(h.items)
}

// save writes history to disk.
func (h *History) save() {
	if !h.modified {
		return
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(h.file)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return
	}

	// Write with timestamp for shared sessions
	file, err := os.OpenFile(h.file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close() // Ignore close error after successful write
	}()

	// Only write the last item (newly added)
	if len(h.items) > 0 {
		lastItem := h.items[len(h.items)-1]
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		_, _ = fmt.Fprintf(file, "# %s\n%s\n", timestamp, lastItem)
	}

	h.modified = false
}
