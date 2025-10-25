package readline

import (
	"testing"
)

func TestKillRing_Add(t *testing.T) {
	t.Parallel()
	kr := NewKillRing()

	// Test adding items
	kr.Add("first")
	kr.Add("second")
	kr.Add("third")

	// Should return most recent item
	if got := kr.Yank(); got != "third" {
		t.Errorf("Yank() = %v, want %v", got, "third")
	}
}

func TestKillRing_Cycle(t *testing.T) {
	t.Parallel()
	kr := NewKillRing()
	kr.Add("first")
	kr.Add("second")
	kr.Add("third")

	// Start with most recent
	if got := kr.Yank(); got != "third" {
		t.Errorf("Yank() = %v, want %v", got, "third")
	}

	// Cycle forward
	if got := kr.Cycle(1); got != "second" {
		t.Errorf("Cycle(1) = %v, want %v", got, "second")
	}

	if got := kr.Cycle(1); got != "first" {
		t.Errorf("Cycle(1) = %v, want %v", got, "first")
	}

	// Should wrap around to beginning
	if got := kr.Cycle(1); got != "third" {
		t.Errorf("Cycle(1) wrap = %v, want %v", got, "third")
	}
}

func TestKillRing_CycleBackward(t *testing.T) {
	t.Parallel()
	kr := NewKillRing()
	kr.Add("first")
	kr.Add("second")
	kr.Add("third")

	kr.Yank() // Start at "third"

	// Cycle backward should wrap to end
	if got := kr.Cycle(-1); got != "first" {
		t.Errorf("Cycle(-1) wrap = %v, want %v", got, "first")
	}

	if got := kr.Cycle(-1); got != "second" {
		t.Errorf("Cycle(-1) = %v, want %v", got, "second")
	}
}

func TestKillRing_LastYank(t *testing.T) {
	t.Parallel()
	kr := NewKillRing()

	// Initially no yank
	if got := kr.GetLastYank(); got != 0 {
		t.Errorf("GetLastYank() = %v, want %v", got, 0)
	}

	// Set last yank
	kr.SetLastYank(5)
	if got := kr.GetLastYank(); got != 5 {
		t.Errorf("GetLastYank() = %v, want %v", got, 5)
	}

	// Reset yank
	kr.ResetYank()
	if got := kr.GetLastYank(); got != 0 {
		t.Errorf("GetLastYank() after reset = %v, want %v", got, 0)
	}
}

func TestKillRing_EmptyRing(t *testing.T) {
	t.Parallel()
	kr := NewKillRing()

	// Empty ring should return empty string
	if got := kr.Yank(); got != "" {
		t.Errorf("Yank() on empty ring = %v, want empty string", got)
	}

	if got := kr.Cycle(1); got != "" {
		t.Errorf("Cycle(1) on empty ring = %v, want empty string", got)
	}
}

func TestKillRing_SingleItem(t *testing.T) {
	t.Parallel()
	kr := NewKillRing()
	kr.Add("only")

	if got := kr.Yank(); got != "only" {
		t.Errorf("Yank() = %v, want %v", got, "only")
	}

	// Cycling with single item should return empty
	if got := kr.Cycle(1); got != "" {
		t.Errorf("Cycle(1) with single item = %v, want empty string", got)
	}
}
