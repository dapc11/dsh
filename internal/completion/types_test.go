package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenu_NewMenu(t *testing.T) {
	// Given/When
	menu := NewMenu()

	// Then
	require.NotNil(t, menu)
	assert.Equal(t, 10, menu.maxRows)
	assert.False(t, menu.IsDisplayed())
}

func TestMenu_ShowHide(t *testing.T) {
	// Given
	menu := NewMenu()
	items := []Item{
		{Text: "item1", Type: "command"},
		{Text: "item2", Type: "file"},
	}

	// When
	menu.Show(items, "test")

	// Then
	assert.True(t, menu.IsDisplayed())
	assert.Len(t, menu.items, 2)
	assert.Equal(t, "test", menu.GetBase())

	// When
	menu.Hide()

	// Then
	assert.False(t, menu.IsDisplayed())
	assert.False(t, menu.HasItems())
}

func TestMenu_Navigation(t *testing.T) {
	// Given
	menu := NewMenu()
	items := []Item{
		{Text: "item1", Type: "command"},
		{Text: "item2", Type: "file"},
		{Text: "item3", Type: "directory"},
	}
	menu.Show(items, "test")

	// When/Then - initial selection
	item, ok := menu.GetSelected()
	require.True(t, ok)
	assert.Equal(t, "item1", item.Text)

	// When
	menu.NextItem()

	// Then
	item, ok = menu.GetSelected()
	require.True(t, ok)
	assert.Equal(t, "item2", item.Text)

	// When
	menu.PrevItem()

	// Then
	item, ok = menu.GetSelected()
	require.True(t, ok)
	assert.Equal(t, "item1", item.Text)
}

func TestMenu_NavigationWrapAround(t *testing.T) {
	// Given
	menu := NewMenu()
	items := []Item{
		{Text: "item1", Type: "command"},
		{Text: "item2", Type: "file"},
	}
	menu.Show(items, "test")

	// When - go to last item
	menu.NextItem()

	// Then
	item, _ := menu.GetSelected()
	assert.Equal(t, "item2", item.Text)

	// When - should wrap to first item
	menu.NextItem()

	// Then
	item, _ = menu.GetSelected()
	assert.Equal(t, "item1", item.Text)

	// When - should wrap to last item
	menu.PrevItem()

	// Then
	item, _ = menu.GetSelected()
	assert.Equal(t, "item2", item.Text)
}

func TestMenu_EmptyMenu(t *testing.T) {
	// Given
	menu := NewMenu()

	// When/Then
	assert.False(t, menu.HasItems())

	_, ok := menu.GetSelected()
	assert.False(t, ok)

	// When - navigation on empty menu should not panic
	menu.NextItem()
	menu.PrevItem()

	// Then - should still be empty
	assert.False(t, menu.HasItems())
}

func TestMenu_SingleItem(t *testing.T) {
	// Given
	menu := NewMenu()
	items := []Item{
		{Text: "only", Type: "command"},
	}
	menu.Show(items, "test")

	// When/Then - initial selection
	item, ok := menu.GetSelected()
	require.True(t, ok)
	assert.Equal(t, "only", item.Text)

	// When - navigation should stay on same item
	menu.NextItem()

	// Then
	item, _ = menu.GetSelected()
	assert.Equal(t, "only", item.Text)

	// When
	menu.PrevItem()

	// Then
	item, _ = menu.GetSelected()
	assert.Equal(t, "only", item.Text)
}
