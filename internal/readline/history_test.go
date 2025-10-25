package readline

import (
	"fmt"
	"testing"
)

func TestHistory_Add(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	h.Add("first command")
	h.Add("second command")
	h.Add("third command")

	if len(h.items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(h.items))
	}

	// Most recent should be last (appended)
	if h.items[len(h.items)-1] != "third command" {
		t.Errorf("Expected 'third command' last, got '%s'", h.items[len(h.items)-1])
	}
}

func TestHistory_AddEmpty(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	h.Add("")
	h.Add("   ")

	if len(h.items) != 0 {
		t.Errorf("Expected 0 items for empty strings, got %d", len(h.items))
	}
}

func TestHistory_AddDuplicate(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	h.Add("command")
	h.Add("other")
	h.Add("command") // Duplicate

	if len(h.items) != 2 {
		t.Errorf("Expected 2 items (duplicate removed), got %d", len(h.items))
	}

	// Most recent "command" should be last
	if h.items[len(h.items)-1] != "command" {
		t.Errorf("Expected 'command' last, got '%s'", h.items[len(h.items)-1])
	}
}

func TestHistory_Navigation(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	h.Add("first")
	h.Add("second")
	h.Add("third")

	// Test previous (going back in history)
	if prev := h.Previous(); prev != "third" {
		t.Errorf("Expected 'third', got '%s'", prev)
	}

	if prev := h.Previous(); prev != "second" {
		t.Errorf("Expected 'second', got '%s'", prev)
	}

	if prev := h.Previous(); prev != "first" {
		t.Errorf("Expected 'first', got '%s'", prev)
	}

	// Should stay at oldest
	if prev := h.Previous(); prev != "first" {
		t.Errorf("Expected 'first' (at end), got '%s'", prev)
	}

	// Test next (going forward in history)
	if next := h.Next(); next != "second" {
		t.Errorf("Expected 'second', got '%s'", next)
	}

	if next := h.Next(); next != "third" {
		t.Errorf("Expected 'third', got '%s'", next)
	}

	// Should return empty at end
	if next := h.Next(); next != "" {
		t.Errorf("Expected empty string at end, got '%s'", next)
	}
}

func TestHistory_ResetPosition(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	h.Add("first")
	h.Add("second")

	// Navigate in history
	h.Previous()
	h.Previous()

	// Reset position
	h.ResetPosition()

	// Should start from beginning again
	if prev := h.Previous(); prev != "second" {
		t.Errorf("Expected 'second' after reset, got '%s'", prev)
	}
}

func TestHistory_EmptyHistory(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	// Navigation on empty history should return empty strings
	if prev := h.Previous(); prev != "" {
		t.Errorf("Expected empty string for empty history, got '%s'", prev)
	}

	if next := h.Next(); next != "" {
		t.Errorf("Expected empty string for empty history, got '%s'", next)
	}
}

func TestHistory_MaxSize(t *testing.T) {
	h := &History{
		items:   make([]string, 0, 1000),
		pos:     0,
		maxSize: 1000,
	}

	// Add more than max size
	for i := 0; i < 1200; i++ {
		h.Add(fmt.Sprintf("command%d", i))
	}

	// Should be limited to max size
	if len(h.items) > 1000 {
		t.Errorf("Expected max 1000 items, got %d", len(h.items))
	}

	// Most recent should be last
	if h.items[len(h.items)-1] != "command1199" {
		t.Errorf("Expected 'command1199' last, got '%s'", h.items[len(h.items)-1])
	}
}
