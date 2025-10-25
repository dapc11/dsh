package readline

import (
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

// FuzzyHistorySearch opens fzf-style window for history search.
func (r *Readline) FuzzyHistorySearch() string {
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

	// Use fuzzyfinder with enhanced visuals
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPromptString("🔍 "),
		fuzzyfinder.WithHeader("History Search - Press Ctrl-C to cancel"),
	)

	if err != nil {
		// User cancelled or error occurred
		return ""
	}

	return items[idx]
}

// FuzzyFileSearch opens fzf-style window for file search.
func (r *Readline) FuzzyFileSearch() string {
	// TODO: Implement Ctrl-T file search
	return ""
}
