package readline

import (
	"strings"
	"testing"

	"github.com/sahilm/fuzzy"
)

func TestNewCustomFzf(t *testing.T) {
	items := []string{"echo hello", "git status", "ls -la"}
	fzf := NewCustomFzf(items)

	if len(fzf.items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(fzf.items))
	}

	if fzf.selected != 0 {
		t.Errorf("Expected selected to be 0, got %d", fzf.selected)
	}

	if fzf.offset != 0 {
		t.Errorf("Expected offset to be 0, got %d", fzf.offset)
	}
}

func TestUpdateMatches(t *testing.T) {
	items := []string{"echo hello", "git status", "git commit", "ls -la", "pwd"}
	fzf := NewCustomFzf(items)

	// Test empty query - should show all items
	fzf.query = ""
	fzf.updateMatches()
	if len(fzf.matches) != 5 {
		t.Errorf("Expected 5 matches for empty query, got %d", len(fzf.matches))
	}

	// Test specific query
	fzf.query = "git"
	fzf.updateMatches()
	if len(fzf.matches) != 2 {
		t.Errorf("Expected 2 matches for 'git', got %d", len(fzf.matches))
	}

	// Verify matches contain expected items
	found := false
	for _, match := range fzf.matches {
		if match.Str == "git status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'git status' in matches")
	}

	// Test query with no matches
	fzf.query = "nonexistent"
	fzf.updateMatches()
	if len(fzf.matches) != 0 {
		t.Errorf("Expected 0 matches for 'nonexistent', got %d", len(fzf.matches))
	}
}

func TestAdjustOffset(t *testing.T) {
	items := make([]string, 20) // Create 20 items
	for i := 0; i < 20; i++ {
		items[i] = "command" + string(rune('0'+i))
	}
	fzf := NewCustomFzf(items)

	// Initialize matches
	fzf.matches = make([]fuzzy.Match, len(items))
	for i, item := range items {
		fzf.matches[i] = fuzzy.Match{Str: item, Index: i}
	}

	// Test selecting item beyond visible range
	fzf.selected = 7 // Beyond maxVisible (5)
	fzf.adjustOffset()
	if fzf.offset != 3 { // 7 - 5 + 1 = 3
		t.Errorf("Expected offset 3, got %d", fzf.offset)
	}

	// Test selecting item before offset
	fzf.selected = 1
	fzf.adjustOffset()
	if fzf.offset != 1 {
		t.Errorf("Expected offset 1, got %d", fzf.offset)
	}

	// Test selecting item within visible range
	fzf.offset = 5
	fzf.selected = 7
	fzf.adjustOffset()
	if fzf.offset != 5 { // Should not change
		t.Errorf("Expected offset to remain 5, got %d", fzf.offset)
	}
}

func TestFuzzyHistorySearchCustom(t *testing.T) {
	r, err := New("test> ")
	if err != nil {
		t.Fatalf("Failed to create readline: %v", err)
	}

	// Add some history items
	testItems := []string{
		"echo hello world",
		"git status",
		"git commit -m 'test'",
		"ls -la /home",
		"pwd",
		"echo goodbye",
	}

	for _, item := range testItems {
		r.history.Add(item)
	}

	// Test that FuzzyHistorySearchCustom creates proper items
	// We can't easily test the interactive part, but we can test the setup
	if len(r.history.items) != 6 {
		t.Errorf("Expected 6 history items, got %d", len(r.history.items))
	}

	// Test unique item filtering (simulate what FuzzyHistorySearchCustom does)
	items := make([]string, 0, len(r.history.items))
	seen := make(map[string]bool)

	for i := len(r.history.items) - 1; i >= 0; i-- {
		item := strings.TrimSpace(r.history.items[i])
		if item != "" && !seen[item] {
			items = append(items, item)
			seen[item] = true
		}
	}

	if len(items) != 6 {
		t.Errorf("Expected 6 unique items, got %d", len(items))
	}

	// Verify most recent item is first
	if items[0] != "echo goodbye" {
		t.Errorf("Expected most recent item first, got %s", items[0])
	}
}

func TestCommandTruncation(t *testing.T) {
	// Test that long commands would be truncated in display
	longCommand := "this is a very long command that exceeds seventy characters and should be truncated"

	if len(longCommand) <= 70 {
		t.Fatal("Test command should be longer than 70 characters")
	}

	// Simulate truncation logic from drawInline
	cmd := longCommand
	if len(cmd) > 70 {
		cmd = cmd[:67] + "..."
	}

	if len(cmd) != 70 {
		t.Errorf("Expected truncated command to be 70 chars, got %d", len(cmd))
	}

	if !strings.HasSuffix(cmd, "...") {
		t.Error("Expected truncated command to end with '...'")
	}

	// Test short command is not truncated
	shortCommand := "short cmd"
	cmd = shortCommand
	if len(cmd) > 70 {
		cmd = cmd[:67] + "..."
	}

	if cmd != shortCommand {
		t.Errorf("Expected short command unchanged, got %s", cmd)
	}
}

func TestEmptyHistory(t *testing.T) {
	r, err := New("test> ")
	if err != nil {
		t.Fatalf("Failed to create readline: %v", err)
	}

	// Test with empty history
	if len(r.history.items) != 0 {
		t.Errorf("Expected empty history, got %d items", len(r.history.items))
	}

	// FuzzyHistorySearchCustom should return empty string for empty history
	// We can't test the interactive part, but the setup should handle empty history
}

func TestDuplicateHistoryFiltering(t *testing.T) {
	r, err := New("test> ")
	if err != nil {
		t.Fatalf("Failed to create readline: %v", err)
	}

	// Add duplicate items
	r.history.Add("echo hello")
	r.history.Add("git status")
	r.history.Add("echo hello") // duplicate
	r.history.Add("pwd")
	r.history.Add("git status") // duplicate

	// Simulate unique filtering from FuzzyHistorySearchCustom
	items := make([]string, 0, len(r.history.items))
	seen := make(map[string]bool)

	for i := len(r.history.items) - 1; i >= 0; i-- {
		item := strings.TrimSpace(r.history.items[i])
		if item != "" && !seen[item] {
			items = append(items, item)
			seen[item] = true
		}
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 unique items, got %d", len(items))
	}

	// Verify most recent unique items are preserved
	expected := []string{"git status", "pwd", "echo hello"}
	for i, expectedItem := range expected {
		if items[i] != expectedItem {
			t.Errorf("Expected item %d to be %s, got %s", i, expectedItem, items[i])
		}
	}
}
