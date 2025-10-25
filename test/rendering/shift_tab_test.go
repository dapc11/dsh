package rendering

import (
	"strings"
	"testing"

	"dsh/internal/completion"
)

// TestShiftTabBackwardCycling tests that shift-tab cycles backward through completions
func TestShiftTabBackwardCycling(t *testing.T) {
	colorProvider := &MockColorProvider{}
	terminalProvider := &MockTerminalProvider{width: 80, height: 24}

	renderer := completion.NewRenderer(colorProvider, terminalProvider)
	menu := completion.NewMenu()

	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
		{Text: "exit", Type: "builtin"},
		{Text: "help", Type: "builtin"},
	}

	menu.Show(items, "")

	// Initial state: "echo" should be selected (index 0)
	selected, _ := menu.GetSelected()
	if selected.Text != "echo" {
		t.Errorf("Initial selection should be 'echo', got %q", selected.Text)
	}

	// Test forward navigation (tab) - should go to "exit" (index 1)
	menu.NextItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "exit" {
		t.Errorf("After NextItem() should be 'exit', got %q", selected.Text)
	}

	// Test backward navigation (shift-tab) - should go back to "echo" (index 0)
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "echo" {
		t.Errorf("After PrevItem() should be 'echo', got %q", selected.Text)
	}

	// Test backward from first item - should wrap to last item "help" (index 2)
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "help" {
		t.Errorf("PrevItem() from first should wrap to 'help', got %q", selected.Text)
	}

	// Validate rendering shows correct selection
	output := CaptureStdout(func() {
		renderer.Render(menu)
	})

	// Should show "help" as selected (reverse video)
	if !strings.Contains(output, "\033[7mhelp\033[0m") {
		t.Errorf("Output should show 'help' as selected, got: %q", output)
	}

	// Should show "echo" and "exit" as unselected (cyan)
	if !strings.Contains(output, "\033[36mecho\033[0m") {
		t.Errorf("Output should show 'echo' as unselected builtin, got: %q", output)
	}

	if !strings.Contains(output, "\033[36mexit\033[0m") {
		t.Errorf("Output should show 'exit' as unselected builtin, got: %q", output)
	}
}
