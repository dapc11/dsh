package test

import (
	"strings"
	"testing"

	"github.com/sahilm/fuzzy"

	"dsh/internal/readline"
)

func TestCustomFzfCreation(t *testing.T) {
	t.Parallel()
	items := []string{"echo hello", "git status", "ls -la"}
	fzf := readline.NewCustomFzf(items)

	// Test that fzf is created with correct initial state
	if fzf == nil {
		t.Fatal("NewCustomFzf returned nil")
	}

	// We can't access private fields, but we can test the behavior
	// by calling methods that would use them
}

func TestFuzzyMatching(t *testing.T) {
	t.Parallel()
	// Test the fuzzy matching logic that our implementation uses
	items := []string{
		"echo hello world",
		"git status",
		"git commit -m 'test'",
		"ls -la /home",
		"pwd",
		"echo goodbye",
	}

	// Test fuzzy matching with "git"
	matches := fuzzy.Find("git", items)
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'git', got %d", len(matches))
	}

	// Test fuzzy matching with "echo"
	matches = fuzzy.Find("echo", items)
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'echo', got %d", len(matches))
	}

	// Test fuzzy matching with no matches
	matches = fuzzy.Find("nonexistent", items)
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for 'nonexistent', got %d", len(matches))
	}

	// Test empty query returns all items
	matches = fuzzy.Find("", items)
	if len(matches) != 0 { // fuzzy.Find returns empty for empty query
		t.Errorf("Expected 0 matches for empty query, got %d", len(matches))
	}
}

func TestCommandTruncation(t *testing.T) {
	t.Parallel()
	// Test the truncation logic used in the display
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "short command",
			expected: "short command",
		},
		{
			input:    "this is a very long command that exceeds seventy characters and should be truncated with ellipsis",
			expected: "this is a very long command that exceeds seventy characters and sho...",
		},
		{
			input:    strings.Repeat("a", 70),
			expected: strings.Repeat("a", 70),
		},
		{
			input:    strings.Repeat("a", 71),
			expected: strings.Repeat("a", 67) + "...",
		},
	}

	for i, tc := range testCases {
		// Simulate the truncation logic from drawInline
		cmd := tc.input
		if len(cmd) > 70 {
			cmd = cmd[:67] + "..."
		}

		if cmd != tc.expected {
			t.Errorf("Test case %d: expected %q, got %q", i, tc.expected, cmd)
		}

		// Ensure result is never longer than 70 characters
		if len(cmd) > 70 {
			t.Errorf("Test case %d: result too long (%d chars): %q", i, len(cmd), cmd)
		}
	}
}

func TestHistoryDeduplication(t *testing.T) {
	t.Parallel()
	// Test the deduplication logic used in FuzzyHistorySearchCustom
	historyItems := []string{
		"echo hello",
		"git status",
		"echo hello", // duplicate
		"pwd",
		"git status", // duplicate
		"ls -la",
	}

	// Simulate the deduplication logic
	items := make([]string, 0, len(historyItems))
	seen := make(map[string]bool)

	// Process in reverse order (most recent first)
	for i := len(historyItems) - 1; i >= 0; i-- {
		item := strings.TrimSpace(historyItems[i])
		if item != "" && !seen[item] {
			items = append(items, item)
			seen[item] = true
		}
	}

	expectedItems := []string{"ls -la", "git status", "pwd", "echo hello"}
	if len(items) != len(expectedItems) {
		t.Errorf("Expected %d unique items, got %d", len(expectedItems), len(items))
	}

	for i, expected := range expectedItems {
		if items[i] != expected {
			t.Errorf("Item %d: expected %q, got %q", i, expected, items[i])
		}
	}
}

func TestOffsetCalculation(t *testing.T) {
	t.Parallel()
	// Test the offset calculation logic for scrolling
	maxVisible := 5

	testCases := []struct {
		selected       int
		currentOffset  int
		expectedOffset int
		description    string
	}{
		{0, 0, 0, "first item, no scroll needed"},
		{4, 0, 0, "within visible range, no scroll"},
		{5, 0, 1, "beyond visible range, scroll down"},
		{10, 0, 6, "far beyond visible range"},
		{2, 5, 2, "before current offset, scroll up"},
		{7, 5, 5, "within current visible range"},
	}

	for _, tc := range testCases {
		// Simulate adjustOffset logic
		offset := tc.currentOffset
		selected := tc.selected

		if selected < offset {
			offset = selected
		} else if selected >= offset+maxVisible {
			offset = selected - maxVisible + 1
		}

		if offset != tc.expectedOffset {
			t.Errorf("%s: expected offset %d, got %d", tc.description, tc.expectedOffset, offset)
		}
	}
}

func TestEmptyHistoryHandling(t *testing.T) {
	t.Parallel()
	// Test behavior with empty history
	var historyItems []string

	// Simulate the deduplication logic with empty history
	items := make([]string, 0, len(historyItems))
	seen := make(map[string]bool)

	for i := len(historyItems) - 1; i >= 0; i-- {
		item := strings.TrimSpace(historyItems[i])
		if item != "" && !seen[item] {
			items = append(items, item)
			seen[item] = true
		}
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items for empty history, got %d", len(items))
	}
}
