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
	items    []string
	matches  []fuzzy.Match
	query    string
	selected int
	offset   int
	width    int
	height   int
	oldState *term.State
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

	// Enter alternate screen
	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l")

	// Get terminal size
	f.updateSize()

	for {
		f.draw()

		// Read key
		key, err := f.readKey()
		if err != nil {
			return "", err
		}

		switch key {
		case 13: // Enter
			if len(f.matches) > 0 && f.selected < len(f.matches) {
				return f.matches[f.selected].Str, nil
			}
			return "", fmt.Errorf("no selection")
		case 3, 27: // Ctrl-C or Escape
			return "", fmt.Errorf("cancelled")
		case 16: // Ctrl-P - up
			if f.selected > 0 {
				f.selected--
				f.adjustOffset()
			}
		case 14: // Ctrl-N - down  
			if f.selected < len(f.matches)-1 {
				f.selected++
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

// readKey reads a single key from stdin
func (f *CustomFzf) readKey() (byte, error) {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	return buf[0], err
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
	maxVisible := f.height - 4 // Leave space for prompt, header, counter
	
	if f.selected < f.offset {
		f.offset = f.selected
	} else if f.selected >= f.offset+maxVisible {
		f.offset = f.selected - maxVisible + 1
	}
}

// draw renders the interface
func (f *CustomFzf) draw() {
	// Clear screen and reset cursor
	fmt.Print("\033[2J\033[H")
	
	// Calculate visible area
	maxVisible := f.height - 3 // Header + counter + prompt
	if maxVisible < 1 {
		maxVisible = 1
	}
	
	// Header
	fmt.Print("History Search (Ctrl-R/S: navigate, Enter: select, Esc: cancel)\r\n")
	
	// Counter
	fmt.Printf("%d/%d\r\n", len(f.matches), len(f.items))
	
	// Items
	endIdx := f.offset + maxVisible
	if endIdx > len(f.matches) {
		endIdx = len(f.matches)
	}
	
	for i := f.offset; i < endIdx; i++ {
		if i == f.selected {
			fmt.Printf("\033[7m> %s\033[0m\r\n", f.matches[i].Str)
		} else {
			fmt.Printf("  %s\r\n", f.matches[i].Str)
		}
	}
	
	// Fill remaining lines
	for i := endIdx - f.offset; i < maxVisible; i++ {
		fmt.Print("\r\n")
	}
	
	// Prompt at bottom
	fmt.Printf("> %s", f.query)
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
