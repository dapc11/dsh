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

	// Get terminal size
	f.updateSize()
	
	// Limit display to max 10 lines
	maxLines := 10
	if f.height < maxLines+3 {
		maxLines = f.height - 3
	}

	for {
		f.drawCompact(maxLines)

		// Read key
		key, err := f.readKey()
		if err != nil {
			return "", err
		}

		switch key {
		case 13: // Enter
			f.clearCompact(maxLines)
			if len(f.matches) > 0 && f.selected < len(f.matches) {
				return f.matches[f.selected].Str, nil
			}
			return "", fmt.Errorf("no selection")
		case 3, 27: // Ctrl-C or Escape
			f.clearCompact(maxLines)
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
		case 1002: // Ctrl-Shift-R - cycle to previous match
			if len(f.matches) > 0 {
				f.selected = (f.selected - 1 + len(f.matches)) % len(f.matches)
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
	var buf [10]byte
	n, err := os.Stdin.Read(buf[:1])
	if err != nil || n == 0 {
		return 0, err
	}

	// Handle escape sequences
	if buf[0] == 27 { // ESC
		n, err := os.Stdin.Read(buf[1:10])
		if err != nil || n < 1 {
			return 27, nil // Just escape
		}
		
		if buf[1] == '[' {
			switch buf[2] {
			case 'A': // Up arrow
				return 1000, nil
			case 'B': // Down arrow  
				return 1001, nil
			case '1':
				// Check for various Ctrl-Shift-R patterns
				if n >= 4 && buf[3] == ';' {
					if (buf[4] == '6' || buf[4] == '5') && buf[5] == 'R' {
						return 1002, nil // Ctrl-Shift-R
					}
				}
			}
		}
		// Check for other Ctrl-Shift-R patterns like ESC[18;2~
		if buf[1] == '[' && buf[2] == '1' && buf[3] == '8' && buf[4] == ';' && buf[5] == '2' && buf[6] == '~' {
			return 1002, nil // Another Ctrl-Shift-R pattern
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
	maxVisible := 9 // Fixed for compact display
	
	if f.selected < f.offset {
		f.offset = f.selected
	} else if f.selected >= f.offset+maxVisible {
		f.offset = f.selected - maxVisible + 1
	}
}

// drawCompact renders the interface using limited screen space
func (f *CustomFzf) drawCompact(maxLines int) {
	// Save cursor position on first draw
	if f.lastDrawnLines == 0 {
		fmt.Print("\033[s") // Save cursor position
	} else {
		// Restore cursor and clear from there
		fmt.Print("\033[u") // Restore cursor position
		fmt.Print("\033[J") // Clear from cursor to end of screen
	}
	
	lines := 0
	
	// Header
	fmt.Print("🔍 ")
	if f.query != "" {
		fmt.Printf("'%s' ", f.query)
	}
	fmt.Printf("(%d/%d)\r\n", len(f.matches), len(f.items))
	lines++
	
	// Items (limited)
	displayCount := maxLines - 1
	if displayCount > len(f.matches) {
		displayCount = len(f.matches)
	}
	
	endIdx := f.offset + displayCount
	if endIdx > len(f.matches) {
		endIdx = len(f.matches)
	}
	
	for i := f.offset; i < endIdx; i++ {
		if i == f.selected {
			fmt.Printf("\033[7m> %s\033[0m\r\n", f.matches[i].Str)
		} else {
			fmt.Printf("  %s\r\n", f.matches[i].Str)
		}
		lines++
	}
	
	f.lastDrawnLines = lines
}

// clearCompact clears the compact display
func (f *CustomFzf) clearCompact(maxLines int) {
	if f.lastDrawnLines > 0 {
		fmt.Print("\033[u") // Restore cursor position
		fmt.Print("\033[J") // Clear from cursor to end of screen
	}
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
