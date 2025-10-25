package rendering

import (
	"testing"

	"dsh/internal/completion"
)

// MockColorProvider provides test color formatting using real ANSI codes.
type MockColorProvider struct{}

func (m *MockColorProvider) Colorize(text, colorName string) string {
	var colorCode string
	switch colorName {
	case "red":
		colorCode = "\033[31m"
	case "green":
		colorCode = "\033[32m"
	case "yellow":
		colorCode = "\033[33m"
	case "blue":
		colorCode = "\033[34m"
	case "magenta":
		colorCode = "\033[35m"
	case "cyan":
		colorCode = "\033[36m"
	case "gray":
		colorCode = "\033[90m"
	case "reverse":
		colorCode = "\033[7m"
	default:
		return text // Unknown color, return as-is
	}

	return colorCode + text + "\033[0m"
}

// MockTerminalProvider provides test terminal size.
type MockTerminalProvider struct {
	width  int
	height int
}

func (m *MockTerminalProvider) GetTerminalSize() (int, int) {
	return m.width, m.height
}

// TestCompletionMenuDisplay tests basic menu display.
func TestCompletionMenuDisplay(t *testing.T) {
	colorProvider := &MockColorProvider{}
	terminalProvider := &MockTerminalProvider{width: 80, height: 24}

	renderer := completion.NewRenderer(colorProvider, terminalProvider)
	menu := completion.NewMenu()

	// Show menu with items
	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
		{Text: "exit", Type: "builtin"},
		{Text: "help", Type: "builtin"},
	}

	menu.Show(items, "e")

	// Test that menu is displayed
	if !menu.IsDisplayed() {
		t.Error("Menu should be displayed after Show()")
	}

	if !menu.HasItems() {
		t.Error("Menu should have items after Show()")
	}

	// Test selected item
	selected, ok := menu.GetSelected()
	if !ok {
		t.Error("Should have selected item")
	}
	if selected.Text != "echo" {
		t.Errorf("Expected first item 'echo', got %q", selected.Text)
	}

	// Test base string
	if menu.GetBase() != "e" {
		t.Errorf("Expected base 'e', got %q", menu.GetBase())
	}

	// Test that renderer can render the menu
	renderer.Render(menu)
}

// TestCompletionNavigation tests menu navigation.
func TestCompletionNavigation(t *testing.T) {
	menu := completion.NewMenu()

	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
		{Text: "exit", Type: "builtin"},
		{Text: "help", Type: "builtin"},
	}

	menu.Show(items, "")

	// Test initial selection
	selected, _ := menu.GetSelected()
	if selected.Text != "echo" {
		t.Errorf("Initial selection should be 'echo', got %q", selected.Text)
	}

	// Test next navigation
	menu.NextItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "exit" {
		t.Errorf("After NextItem() should be 'exit', got %q", selected.Text)
	}

	// Test previous navigation
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "echo" {
		t.Errorf("After PrevItem() should be 'echo', got %q", selected.Text)
	}

	// Test wrap-around (previous from first item)
	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "help" {
		t.Errorf("PrevItem() from first should wrap to 'help', got %q", selected.Text)
	}
}

// TestCompletionHideShow tests menu hide/show behavior.
func TestCompletionHideShow(t *testing.T) {
	menu := completion.NewMenu()

	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
	}

	// Initially not displayed
	if menu.IsDisplayed() {
		t.Error("New menu should not be displayed")
	}

	// Show menu
	menu.Show(items, "e")
	if !menu.IsDisplayed() {
		t.Error("Menu should be displayed after Show()")
	}

	// Hide menu
	menu.Hide()
	if menu.IsDisplayed() {
		t.Error("Menu should not be displayed after Hide()")
	}

	// Should have no items after hide
	if menu.HasItems() {
		t.Error("Menu should have no items after Hide()")
	}
}

// TestCompletionRendering tests actual rendering behavior.
func TestCompletionRendering(_ *testing.T) {
	colorProvider := &MockColorProvider{}
	terminalProvider := &MockTerminalProvider{width: 80, height: 24}

	renderer := completion.NewRenderer(colorProvider, terminalProvider)
	menu := completion.NewMenu()

	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
		{Text: "ls", Type: "command"},
		{Text: "file.txt", Type: "file"},
		{Text: "dir/", Type: "directory"},
	}

	menu.Show(items, "")

	// Test rendering
	renderer.Render(menu)

	// Test clearing
	renderer.Clear(menu)
}

// TestCompletionTypeHandling tests different completion types.
func TestCompletionTypeHandling(t *testing.T) {
	menu := completion.NewMenu()

	// Test different item types
	testCases := []struct {
		name  string
		items []completion.Item
	}{
		{
			"builtin_commands",
			[]completion.Item{
				{Text: "echo", Type: "builtin"},
				{Text: "exit", Type: "builtin"},
				{Text: "help", Type: "builtin"},
			},
		},
		{
			"system_commands",
			[]completion.Item{
				{Text: "ls", Type: "command"},
				{Text: "grep", Type: "command"},
				{Text: "cat", Type: "command"},
			},
		},
		{
			"files_and_dirs",
			[]completion.Item{
				{Text: "file.txt", Type: "file"},
				{Text: "script.sh", Type: "file"},
				{Text: "directory/", Type: "directory"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			menu.Show(tc.items, "")

			// Verify all items are present by navigating through them
			for i, expectedItem := range tc.items {
				// Reset to first item, then navigate to target
				menu.Show(tc.items, "")
				for range i {
					menu.NextItem()
				}

				selected, ok := menu.GetSelected()
				if !ok {
					t.Errorf("Should have selected item at index %d", i)
					continue
				}

				if selected.Text != expectedItem.Text {
					t.Errorf("Expected item %q, got %q", expectedItem.Text, selected.Text)
				}

				if selected.Type != expectedItem.Type {
					t.Errorf("Expected type %q, got %q", expectedItem.Type, selected.Type)
				}
			}
		})
	}
}

// TestCompletionEdgeCases tests edge cases and error conditions.
func TestCompletionEdgeCases(t *testing.T) {
	menu := completion.NewMenu()

	// Test empty menu
	menu.Show([]completion.Item{}, "")
	if menu.HasItems() {
		t.Error("Empty menu should not have items")
	}

	_, ok := menu.GetSelected()
	if ok {
		t.Error("Empty menu should not have selected item")
	}

	// Test navigation on empty menu (should not crash)
	menu.NextItem()
	menu.PrevItem()

	// Test single item menu
	menu.Show([]completion.Item{{Text: "single", Type: "builtin"}}, "")

	// Navigation should stay on same item
	menu.NextItem()
	selected, _ := menu.GetSelected()
	if selected.Text != "single" {
		t.Error("Single item navigation should stay on same item")
	}

	menu.PrevItem()
	selected, _ = menu.GetSelected()
	if selected.Text != "single" {
		t.Error("Single item navigation should stay on same item")
	}
}
