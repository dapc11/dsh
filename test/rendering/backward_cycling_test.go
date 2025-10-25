package rendering

import (
	"testing"

	"dsh/internal/completion"
)

// TestBackwardCycling tests basic backward navigation through completion menu
func TestBackwardCycling(t *testing.T) {
	menu := completion.NewMenu()
	items := []completion.Item{
		{Text: "first", Type: "builtin"},
		{Text: "second", Type: "builtin"},
		{Text: "third", Type: "builtin"},
	}
	menu.Show(items, "")

	// Start at first item
	selected, _ := menu.GetSelected()
	if selected.Text != "first" {
		t.Errorf("Expected 'first', got %q", selected.Text)
	}

	// Backward should wrap to last item
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "third" {
		t.Errorf("Backward from first should wrap to 'third', got %q", selected.Text)
	}

	// Backward should go to second item
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "second" {
		t.Errorf("Backward from third should go to 'second', got %q", selected.Text)
	}

	// Backward should go to first item
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "first" {
		t.Errorf("Backward from second should go to 'first', got %q", selected.Text)
	}
}

// TestForwardBackwardCombination tests mixing forward and backward navigation
func TestForwardBackwardCombination(t *testing.T) {
	menu := completion.NewMenu()
	items := []completion.Item{
		{Text: "alpha", Type: "builtin"},
		{Text: "beta", Type: "builtin"},
		{Text: "gamma", Type: "builtin"},
	}
	menu.Show(items, "")

	// Forward: alpha -> beta
	menu.NextItem()
	selected, _ := menu.GetSelected()
	if selected.Text != "beta" {
		t.Errorf("Forward should go to 'beta', got %q", selected.Text)
	}

	// Backward: beta -> alpha
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "alpha" {
		t.Errorf("Backward should go to 'alpha', got %q", selected.Text)
	}

	// Forward twice: alpha -> beta -> gamma
	menu.NextItem()
	menu.NextItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "gamma" {
		t.Errorf("Forward twice should go to 'gamma', got %q", selected.Text)
	}

	// Backward: gamma -> beta
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "beta" {
		t.Errorf("Backward from gamma should go to 'beta', got %q", selected.Text)
	}
}

// TestSingleItemBackward tests backward navigation with only one item
func TestSingleItemBackward(t *testing.T) {
	menu := completion.NewMenu()
	items := []completion.Item{
		{Text: "only", Type: "builtin"},
	}
	menu.Show(items, "")

	// Should stay on same item
	menu.PrevItem()
	selected, _ := menu.GetSelected()
	if selected.Text != "only" {
		t.Errorf("Single item backward should stay on 'only', got %q", selected.Text)
	}

	// Multiple backwards should still stay
	menu.PrevItem()
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "only" {
		t.Errorf("Multiple backwards should stay on 'only', got %q", selected.Text)
	}
}

// TestBackwardRenderingSelection tests that backward navigation updates visual selection
func TestBackwardRenderingSelection(t *testing.T) {
	colorProvider := &MockColorProvider{}
	terminalProvider := &MockTerminalProvider{width: 80, height: 24}
	renderer := completion.NewRenderer(colorProvider, terminalProvider)

	menu := completion.NewMenu()
	items := []completion.Item{
		{Text: "cmd1", Type: "builtin"},
		{Text: "cmd2", Type: "builtin"},
	}
	menu.Show(items, "")

	// Navigate backward to last item
	menu.PrevItem()

	// Render and check selection
	output := CaptureStdout(func() {
		renderer.Render(menu)
	})

	// Should show cmd2 as selected (reverse video)
	if !contains(output, "\033[7mcmd2\033[0m") {
		t.Errorf("Backward navigation should select 'cmd2', output: %q", output)
	}

	// Should show cmd1 as unselected (cyan)
	if !contains(output, "\033[36mcmd1\033[0m") {
		t.Errorf("cmd1 should be unselected after backward nav, output: %q", output)
	}
}

// contains checks if string contains substring (helper function)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
