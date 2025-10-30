package readline

import (
	"testing"

	"dsh/test/rendering"
)

func TestCompletionPagination(t *testing.T) {
	mockTerm := &rendering.MockTerminalInterface{}
	renderer := NewCompletionRenderer(mockTerm)

	// Create 25 test items (spans 3 pages)
	items := make([]CompletionItem, 25)
	for i := 0; i < 25; i++ {
		items[i] = CompletionItem{
			Text: "item" + string(rune('A'+i)),
			Type: "file",
		}
	}

	t.Run("ShowCompletion displays first page correctly", func(t *testing.T) {
		mockTerm.ClearOutput()
		renderer.ShowCompletion(items, 0)

		output := mockTerm.GetOutput()
		if !contains(output, "itemA") {
			t.Error("Should show first item")
		}
		if !contains(output, "itemJ") {
			t.Error("Should show 10th item")
		}
		if contains(output, "itemK") {
			t.Error("Should not show 11th item on first page")
		}
	})

	t.Run("Same page navigation uses efficient update", func(t *testing.T) {
		mockTerm.ClearOutput()
		renderer.ShowCompletion(items, 0)
		mockTerm.ClearOutput()
		
		renderer.UpdateSelectionHighlight(0, 5) // Same page
		
		output := mockTerm.GetOutput()
		if !contains(output, "\033[u") {
			t.Error("Should use cursor positioning for same page")
		}
		if contains(output, "itemA") {
			t.Error("Should not re-render items for same page")
		}
	})

	t.Run("Cross page navigation re-renders menu", func(t *testing.T) {
		mockTerm.ClearOutput()
		renderer.ShowCompletion(items, 0)
		mockTerm.ClearOutput()
		
		renderer.UpdateSelectionHighlight(5, 15) // Page 0 to page 1
		
		output := mockTerm.GetOutput()
		if !contains(output, "itemK") {
			t.Error("Should show items from page 1")
		}
		if !contains(output, "\033[J") {
			t.Error("Should clear menu area")
		}
	})

	t.Run("Page calculation works correctly", func(t *testing.T) {
		testCases := []struct {
			selected int
			page     int
		}{
			{0, 0}, {9, 0}, {10, 1}, {19, 1}, {20, 2}, {24, 2},
		}

		for _, tc := range testCases {
			page := tc.selected / renderer.itemsPerPage
			if page != tc.page {
				t.Errorf("Item %d: expected page %d, got %d", tc.selected, tc.page, page)
			}
		}
	})

	t.Run("Third page displays correctly", func(t *testing.T) {
		mockTerm.ClearOutput()
		renderer.ShowCompletion(items, 22) // Item on page 2
		
		output := mockTerm.GetOutput()
		if !contains(output, "itemU") {
			t.Error("Should show 21st item")
		}
		if !contains(output, "itemY") {
			t.Error("Should show 25th item")
		}
		if contains(output, "itemA") {
			t.Error("Should not show first page items")
		}
	})
}

func TestCompletionMenuNavigation(t *testing.T) {
	mockTerm := &rendering.MockTerminalInterface{}
	menu := NewCompletionMenu(mockTerm)

	items := []CompletionItem{
		{Text: "apple", Type: "file"},
		{Text: "banana", Type: "file"},
		{Text: "cherry", Type: "file"},
	}

	t.Run("Show initializes menu correctly", func(t *testing.T) {
		menu.Show(items)
		
		if !menu.IsActive() {
			t.Error("Menu should be active after Show")
		}
		if menu.selected != 0 {
			t.Error("Should start with first item selected")
		}
		if len(menu.items) != 3 {
			t.Error("Should store all items")
		}
	})

	t.Run("Next advances selection", func(t *testing.T) {
		menu.Show(items)
		menu.Next()
		
		if menu.selected != 1 {
			t.Errorf("Expected selection 1, got %d", menu.selected)
		}
	})

	t.Run("Next wraps at end", func(t *testing.T) {
		menu.Show(items)
		menu.selected = 2 // Last item
		menu.Next()
		
		if menu.selected != 0 {
			t.Errorf("Expected wrap to 0, got %d", menu.selected)
		}
	})

	t.Run("Prev moves backward", func(t *testing.T) {
		menu.Show(items)
		menu.selected = 1
		menu.Prev()
		
		if menu.selected != 0 {
			t.Errorf("Expected selection 0, got %d", menu.selected)
		}
	})

	t.Run("Prev wraps at beginning", func(t *testing.T) {
		menu.Show(items)
		menu.Prev()
		
		if menu.selected != 2 {
			t.Errorf("Expected wrap to 2, got %d", menu.selected)
		}
	})

	t.Run("GetSelected returns correct item", func(t *testing.T) {
		menu.Show(items)
		menu.selected = 1
		
		item, ok := menu.GetSelected()
		if !ok {
			t.Error("Should return selected item")
		}
		if item.Text != "banana" {
			t.Errorf("Expected banana, got %s", item.Text)
		}
	})

	t.Run("Hide deactivates menu", func(t *testing.T) {
		menu.Show(items)
		menu.Hide()
		
		if menu.IsActive() {
			t.Error("Menu should be inactive after Hide")
		}
	})
}

func TestCompletionRenderer(t *testing.T) {
	mockTerm := &rendering.MockTerminalInterface{}
	renderer := NewCompletionRenderer(mockTerm)

	t.Run("NewCompletionRenderer initializes correctly", func(t *testing.T) {
		if renderer.itemsPerPage != 10 {
			t.Error("Should default to 10 items per page")
		}
		if renderer.currentPage != 0 {
			t.Error("Should start on page 0")
		}
		if renderer.active {
			t.Error("Should start inactive")
		}
	})

	t.Run("UpdateSelection handles same page", func(t *testing.T) {
		items := []CompletionItem{
			{Text: "item1", Type: "file"},
			{Text: "item2", Type: "file"},
		}
		mockTerm.ClearOutput()
		renderer.ShowCompletion(items, 0)
		mockTerm.ClearOutput()
		
		renderer.UpdateSelection(items, 0, 1)
		
		output := mockTerm.GetOutput()
		// UpdateSelection uses video buffer, not cursor positioning
		if len(output) == 0 {
			t.Error("Should produce some output for selection update")
		}
	})

	t.Run("UpdateSelection handles page change", func(t *testing.T) {
		items := make([]CompletionItem, 15)
		for i := 0; i < 15; i++ {
			items[i] = CompletionItem{Text: "item" + string(rune('A'+i)), Type: "file"}
		}
		
		mockTerm.ClearOutput()
		renderer.ShowCompletion(items, 0)
		mockTerm.ClearOutput()
		
		renderer.UpdateSelection(items, 5, 12) // Cross page boundary
		
		output := mockTerm.GetOutput()
		// UpdateSelection uses video buffer approach, not clearing
		if len(output) == 0 {
			t.Error("Should produce output for page change")
		}
	})
}

func contains(output, substr string) bool {
	for i := 0; i <= len(output)-len(substr); i++ {
		if output[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
