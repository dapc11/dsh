package readline

import (
	"fmt"
	"os"
	"strings"

	"github.com/sahilm/fuzzy"
	"golang.org/x/term"
)

// CustomFzf implements a custom fzf-style interface
type CustomFzf struct {
	items         []string
	matches       []fuzzy.Match
	query         string
	selected      int
	offset        int
	width         int
	height        int
	oldState      *term.State
	lastDrawnLines int
}

// NewCustomFzf creates a new custom fzf interface
func NewCustomFzf(items []string) *CustomFzf {
	return &CustomFzf{
		items:    items,
		matches:  make([]fuzzy.Match, len(items)),
		selected: 0,
		offset:   0,
	}
}

// Run starts the custom fzf interface
func (f *CustomFzf) Run() (string, error) {
	// Initialize matches with all items
	for i, item := range f.items {
		f.matches[i] = fuzzy.Match{Str: item, Index: i}
	}

	// Enter raw mode
	if err := f.enterRawMode(); err != nil {
		return "", err
	}
	defer f.exitRawMode()

	for {
		f.drawInline()

		// Read key
		key, err := f.readKey()
		if err != nil {
			return "", err
		}

		switch key {
		case 13: // Enter
			f.clearInline()
			if len(f.matches) > 0 && f.selected < len(f.matches) {
				return f.matches[f.selected].Str, nil
			}
			return "", fmt.Errorf("no selection")
		case 3, 27: // Ctrl-C or Escape
			f.clearInline()
			return "", fmt.Errorf("cancelled")
		case 16, 1000: // Ctrl-P or Up arrow - navigate up
			if f.selected > 0 {
				f.selected--
				f.adjustOffset()
			}
		case 14, 1001: // Ctrl-N or Down arrow - navigate down  
			if f.selected < len(f.matches)-1 {
				f.selected++
				f.adjustOffset()
			}
		case 18: // Ctrl-R - cycle to next match
			if len(f.matches) > 0 {
				f.selected = (f.selected + 1) % len(f.matches)
				f.adjustOffset()
			}
		case 127: // Backspace
			if len(f.query) > 0 {
				f.query = f.query[:len(f.query)-1]
				f.updateMatches()
			}
		default:
			if key >= 32 && key < 127 { // Printable characters
				f.query += string(rune(key))
				f.updateMatches()
			}
		}
	}
}

// enterRawMode puts terminal in raw mode
func (f *CustomFzf) enterRawMode() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	f.oldState = oldState
	return nil
}

// exitRawMode restores terminal mode
func (f *CustomFzf) exitRawMode() {
	if f.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), f.oldState)
	}
}

// updateSize gets current terminal size
func (f *CustomFzf) updateSize() {
	f.width, f.height, _ = term.GetSize(int(os.Stdin.Fd()))
}

// readKey reads a single key from stdin, handling escape sequences
func (f *CustomFzf) readKey() (int, error) {
	var buf [3]byte
	n, err := os.Stdin.Read(buf[:1])
	if err != nil || n == 0 {
		return 0, err
	}

	// Handle escape sequences (arrow keys)
	if buf[0] == 27 { // ESC
		n, err := os.Stdin.Read(buf[1:3])
		if err != nil || n < 2 {
			return 27, nil // Just escape
		}
		
		if buf[1] == '[' {
			switch buf[2] {
			case 'A': // Up arrow
				return 1000, nil
			case 'B': // Down arrow  
				return 1001, nil
			}
		}
		return 27, nil // Unhandled escape sequence
	}

	return int(buf[0]), nil
}

// updateMatches performs fuzzy search and updates matches
func (f *CustomFzf) updateMatches() {
	if f.query == "" {
		// Show all items when no query
		f.matches = make([]fuzzy.Match, len(f.items))
		for i, item := range f.items {
			f.matches[i] = fuzzy.Match{Str: item, Index: i}
		}
	} else {
		f.matches = fuzzy.Find(f.query, f.items)
	}
	
	// Reset selection
	f.selected = 0
	f.offset = 0
}

// adjustOffset adjusts scroll offset to keep selected item visible
func (f *CustomFzf) adjustOffset() {
	maxVisible := 8
	
	if f.selected < f.offset {
		f.offset = f.selected
	} else if f.selected >= f.offset+maxVisible {
		f.offset = f.selected - maxVisible + 1
	}
}

// drawInline renders the interface inline below current position
func (f *CustomFzf) drawInline() {
	// Clear current line and show current match
	fmt.Print("\r\033[K")
	
	if len(f.matches) == 0 {
		fmt.Printf("🔍 %s (no matches)", f.query)
		return
	}
	
	// Show current selection
	current := f.matches[f.selected].Str
	fmt.Printf("🔍 %s (%d/%d) > %s", f.query, f.selected+1, len(f.matches), current)
}

// clearInline clears the inline display
func (f *CustomFzf) clearInline() {
	fmt.Print("\r\033[K") // Clear current line
}

// FuzzyHistorySearchCustom uses the custom fzf implementation
func (r *Readline) FuzzyHistorySearchCustom() string {
	if len(r.history.items) == 0 {
		return ""
	}

	// Get unique history items (most recent first)
	items := make([]string, 0, len(r.history.items))
	seen := make(map[string]bool)
	
	for i := len(r.history.items) - 1; i >= 0; i-- {
		item := strings.TrimSpace(r.history.items[i])
		if item != "" && !seen[item] {
			items = append(items, item)
			seen[item] = true
		}
	}

	if len(items) == 0 {
		return ""
	}

	fzf := NewCustomFzf(items)
	result, err := fzf.Run()
	if err != nil {
		return ""
	}
	
	return result
}
