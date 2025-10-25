package completion

import (
	"testing"
)

func TestMenu_NewMenu(t *testing.T) {
	menu := NewMenu()
	
	if menu == nil {
		t.Error("NewMenu returned nil")
	}
	
	if menu.maxRows != 10 {
		t.Errorf("Expected maxRows = 10, got %d", menu.maxRows)
	}
	
	if menu.IsDisplayed() {
		t.Error("New menu should not be displayed")
	}
}

func TestMenu_ShowHide(t *testing.T) {
	menu := NewMenu()
	items := []Item{
		{Text: "item1", Type: "command"},
		{Text: "item2", Type: "file"},
	}
	
	menu.Show(items, "test")
	
	if !menu.IsDisplayed() {
		t.Error("Menu should be displayed after Show")
	}
	
	if len(menu.items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(menu.items))
	}
	
	if menu.GetBase() != "test" {
		t.Errorf("Expected base 'test', got '%s'", menu.GetBase())
	}
	
	menu.Hide()
	
	if menu.IsDisplayed() {
		t.Error("Menu should not be displayed after Hide")
	}
	
	if menu.HasItems() {
		t.Error("Menu should not have items after Hide")
	}
}

func TestMenu_Navigation(t *testing.T) {
	menu := NewMenu()
	items := []Item{
		{Text: "item1", Type: "command"},
		{Text: "item2", Type: "file"},
		{Text: "item3", Type: "directory"},
	}
	
	menu.Show(items, "test")
	
	// Test initial selection
	item, ok := menu.GetSelected()
	if !ok {
		t.Error("Should have selected item")
	}
	if item.Text != "item1" {
		t.Errorf("Expected 'item1', got '%s'", item.Text)
	}
	
	// Test next item
	menu.NextItem()
	item, ok = menu.GetSelected()
	if !ok {
		t.Error("Should have selected item")
	}
	if item.Text != "item2" {
		t.Errorf("Expected 'item2', got '%s'", item.Text)
	}
	
	// Test previous item
	menu.PrevItem()
	item, ok = menu.GetSelected()
	if !ok {
		t.Error("Should have selected item")
	}
	if item.Text != "item1" {
		t.Errorf("Expected 'item1', got '%s'", item.Text)
	}
}

func TestMenu_NavigationWrapAround(t *testing.T) {
	menu := NewMenu()
	items := []Item{
		{Text: "item1", Type: "command"},
		{Text: "item2", Type: "file"},
	}
	
	menu.Show(items, "test")
	
	// Go to last item
	menu.NextItem()
	item, _ := menu.GetSelected()
	if item.Text != "item2" {
		t.Errorf("Expected 'item2', got '%s'", item.Text)
	}
	
	// Should wrap to first item
	menu.NextItem()
	item, _ = menu.GetSelected()
	if item.Text != "item1" {
		t.Errorf("Expected wrap to 'item1', got '%s'", item.Text)
	}
	
	// Should wrap to last item
	menu.PrevItem()
	item, _ = menu.GetSelected()
	if item.Text != "item2" {
		t.Errorf("Expected wrap to 'item2', got '%s'", item.Text)
	}
}

func TestMenu_EmptyMenu(t *testing.T) {
	menu := NewMenu()
	
	// Test empty menu
	if menu.HasItems() {
		t.Error("Empty menu should not have items")
	}
	
	_, ok := menu.GetSelected()
	if ok {
		t.Error("Empty menu should not have selected item")
	}
	
	// Navigation on empty menu should not panic
	menu.NextItem()
	menu.PrevItem()
}

func TestMenu_SingleItem(t *testing.T) {
	menu := NewMenu()
	items := []Item{
		{Text: "only", Type: "command"},
	}
	
	menu.Show(items, "test")
	
	item, ok := menu.GetSelected()
	if !ok {
		t.Error("Should have selected item")
	}
	if item.Text != "only" {
		t.Errorf("Expected 'only', got '%s'", item.Text)
	}
	
	// Navigation should stay on same item
	menu.NextItem()
	item, _ = menu.GetSelected()
	if item.Text != "only" {
		t.Errorf("Expected 'only' after NextItem, got '%s'", item.Text)
	}
	
	menu.PrevItem()
	item, _ = menu.GetSelected()
	if item.Text != "only" {
		t.Errorf("Expected 'only' after PrevItem, got '%s'", item.Text)
	}
}
