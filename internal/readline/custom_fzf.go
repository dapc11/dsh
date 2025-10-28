package readline

import (
	"errors"
	"strings"

	"dsh/internal/terminal"
	"github.com/sahilm/fuzzy"
)

var (
	// ErrNoSelection indicates no item was selected during fuzzy search.
	ErrNoSelection = errors.New("no selection")
	// ErrCancelled indicates the fuzzy search was cancelled by user.
	ErrCancelled = errors.New("cancelled")
)

// CustomFzf implements a custom fzf-style interface.
type CustomFzf struct {
	items          []string
	matches        []fuzzy.Match
	query          string
	selected       int
	offset         int
	terminal       *terminal.Interface
	lastDrawnLines int
}

// NewCustomFzf creates a new custom fzf interface.
func NewCustomFzf(items []string) *CustomFzf {
	return &CustomFzf{
		items:    items,
		matches:  make([]fuzzy.Match, len(items)),
		selected: 0,
		offset:   0,
		terminal: terminal.NewInterface(),
	}
}

// Run starts the custom fzf interface.
func (f *CustomFzf) Run() (string, error) {
	// Initialize matches with all items
	for i, item := range f.items {
		f.matches[i] = fuzzy.Match{Str: item, Index: i}
	}

	// Enter raw mode
	err := f.terminal.EnableRawMode()
	if err != nil {
		return "", err
	}
	defer f.terminal.DisableRawMode()

	for {
		f.drawInline()

		// Read key
		keyEvent, err := f.terminal.ReadKey()
		if err != nil {
			return "", err
		}

		switch keyEvent.Key {
		case terminal.KeyEnter:
			f.clearInline()
			if len(f.matches) > 0 && f.selected < len(f.matches) {
				return f.matches[f.selected].Str, nil
			}
			return "", ErrNoSelection
		case terminal.KeyCtrlC, terminal.KeyEscape:
			f.clearInline()
			return "", ErrCancelled
		case terminal.KeyCtrlP, terminal.KeyArrowUp:
			if f.selected > 0 {
				f.selected--
				f.adjustOffset()
			}
		case terminal.KeyCtrlN, terminal.KeyArrowDown:
			if f.selected < len(f.matches)-1 {
				f.selected++
				f.adjustOffset()
			}
		case terminal.KeyCtrlR:
			if len(f.matches) > 0 {
				f.selected = (f.selected + 1) % len(f.matches)
				f.adjustOffset()
			}
		case terminal.KeyBackspace:
			if len(f.query) > 0 {
				f.query = f.query[:len(f.query)-1]
				f.updateMatches()
			}
		default:
			if keyEvent.Rune != 0 && keyEvent.Rune >= 32 && keyEvent.Rune < 127 {
				f.query += string(keyEvent.Rune)
				f.updateMatches()
			}
		}
	}
}

// updateMatches performs fuzzy search and updates matches.
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

// adjustOffset adjusts scroll offset to keep selected item visible.
func (f *CustomFzf) adjustOffset() {
	maxVisible := 5

	if f.selected < f.offset {
		f.offset = f.selected
	} else if f.selected >= f.offset+maxVisible {
		f.offset = f.selected - maxVisible + 1
	}
}

// drawInline renders the interface inline below current position.
func (f *CustomFzf) drawInline() {
	// Always clear exactly 6 lines (header + 5 matches) if we've drawn before
	if f.lastDrawnLines > 0 {
		for range 6 {
			f.terminal.WriteString("\033[1A\033[2K")
		}
	}

	lines := 0

	// Header line
	if len(f.matches) == 0 {
		f.terminal.Printf("🔍 %s (no matches)\r\n", f.query)
		lines++
	} else {
		f.terminal.Printf("🔍 %s (%d/%d)\r\n", f.query, len(f.matches), len(f.items))
		lines++

		// Show max 5 matches
		maxVisible := 5
		endIdx := f.offset + maxVisible
		if endIdx > len(f.matches) {
			endIdx = len(f.matches)
		}

		for i := f.offset; i < endIdx; i++ {
			// Truncate long commands to prevent wrapping
			cmd := f.matches[i].Str
			if len(cmd) > 70 {
				cmd = cmd[:67] + "..."
			}

			if i == f.selected {
				f.terminal.WriteString(f.terminal.StyleText("> "+cmd, terminal.Style{Reverse: true}) + "\r\n")
			} else {
				f.terminal.Printf("  %s\r\n", cmd)
			}
			lines++
		}

		// Fill remaining lines with empty lines to maintain consistent clearing
		for i := lines; i < 6; i++ {
			f.terminal.WriteString("\r\n")
			lines++
		}
	}

	f.lastDrawnLines = lines
}

// clearInline clears the inline display.
func (f *CustomFzf) clearInline() {
	if f.lastDrawnLines > 0 {
		for range 6 {
			f.terminal.WriteString("\033[1A\033[2K")
		}
	}
}

// FuzzyHistorySearchCustom uses the custom fzf implementation.
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
