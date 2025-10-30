package completion

import (
	"fmt"
	"strings"
	"testing"

	"dsh/internal/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockTerminal for testing renderer
type MockTerminal struct {
	output  strings.Builder
	width   int
	height  int
	cursorX int
	cursorY int
	savedX  int
	savedY  int
}

func NewMockTerminal() *MockTerminal {
	return &MockTerminal{
		width:  80,
		height: 24,
	}
}

func (m *MockTerminal) Write(data []byte) (int, error) {
	return m.output.Write(data)
}

func (m *MockTerminal) WriteString(s string) (int, error) {
	m.output.WriteString(s)
	return len(s), nil
}

func (m *MockTerminal) Printf(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...))
}

func (m *MockTerminal) Size() (int, int) {
	return m.width, m.height
}

func (m *MockTerminal) MoveCursor(x, y int) {
	m.cursorX, m.cursorY = x, y
	m.output.WriteString(fmt.Sprintf("\033[%d;%dH", y+1, x+1))
}

func (m *MockTerminal) ClearLine() {
	m.output.WriteString("\033[2K")
}

func (m *MockTerminal) ClearToEnd() {
	m.output.WriteString("\033[K")
}

func (m *MockTerminal) ClearFromCursor() {
	m.output.WriteString("\033[0J")
}

func (m *MockTerminal) SaveCursor() {
	m.savedX, m.savedY = m.cursorX, m.cursorY
	m.output.WriteString("\033[s")
}

func (m *MockTerminal) RestoreCursor() {
	m.cursorX, m.cursorY = m.savedX, m.savedY
	m.output.WriteString("\033[u")
}

func (m *MockTerminal) HideCursor() {
	m.output.WriteString("\033[?25l")
}

func (m *MockTerminal) ShowCursor() {
	m.output.WriteString("\033[?25h")
}

func (m *MockTerminal) ReadKey() (terminal.KeyEvent, error) {
	return terminal.KeyEvent{}, nil
}

func (m *MockTerminal) Colorize(text string, color terminal.Color) string {
	switch color {
	case terminal.ColorRed:
		return "\033[31m" + text + "\033[0m"
	case terminal.ColorGreen:
		return "\033[32m" + text + "\033[0m"
	case terminal.ColorBlue:
		return "\033[34m" + text + "\033[0m"
	case terminal.ColorCyan:
		return "\033[36m" + text + "\033[0m"
	case terminal.ColorBrightBlack:
		return "\033[90m" + text + "\033[0m"
	default:
		return text
	}
}

func (m *MockTerminal) StyleText(text string, style terminal.Style) string {
	if style.Reverse {
		return "\033[7m" + text + "\033[0m"
	}
	return text
}

func (m *MockTerminal) EnableRawMode() error {
	return nil
}

func (m *MockTerminal) DisableRawMode() error {
	return nil
}

func (m *MockTerminal) IsRawMode() bool {
	return false
}

func (m *MockTerminal) Cleanup() {
	// No-op for mock
}

func (m *MockTerminal) GetOutput() string {
	return m.output.String()
}

func (m *MockTerminal) Reset() {
	m.output.Reset()
	m.cursorX, m.cursorY = 0, 0
	m.savedX, m.savedY = 0, 0
}

func TestNewRenderer(t *testing.T) {
	// Given
	term := NewMockTerminal()

	// When
	renderer := NewRenderer(term)

	// Then
	require.NotNil(t, renderer)
}

func TestRender_WithItems(t *testing.T) {
	// Given
	term := NewMockTerminal()
	renderer := NewRenderer(term)
	menu := NewMenu()
	items := []Item{
		{Text: "echo", Type: "builtin"},
		{Text: "ls", Type: "command"},
		{Text: "dir/", Type: "directory"},
	}
	menu.Show(items, "")

	// When
	renderer.Render(menu)

	// Then
	output := term.GetOutput()
	require.Contains(t, output, "\033[s") // save cursor
	require.Contains(t, output, "echo")
	require.Contains(t, output, "ls")
	require.Contains(t, output, "\033[32m") // green for command
	require.Contains(t, output, "\033[34m") // blue for directory
	require.Contains(t, output, "\033[7m")  // reverse for selection
}

func TestRender_BuiltinColorWhenNotSelected(t *testing.T) {
	// Given
	term := NewMockTerminal()
	renderer := NewRenderer(term)
	menu := NewMenu()
	items := []Item{
		{Text: "echo", Type: "builtin"},
		{Text: "ls", Type: "command"},
	}
	menu.Show(items, "")
	menu.NextItem() // Move to second item

	// When
	renderer.Render(menu)

	// Then
	output := term.GetOutput()
	require.Contains(t, output, "\033[36m") // cyan for builtin when not selected
}

func TestRender_WithPagination(t *testing.T) {
	// Given
	term := NewMockTerminal()
	term.width = 20
	term.height = 8
	renderer := NewRenderer(term)
	menu := NewMenu()
	items := make([]Item, 20)
	for i := range items {
		items[i] = Item{Text: fmt.Sprintf("item%d", i), Type: "command"}
	}
	menu.Show(items, "")

	// When
	renderer.Render(menu)

	// Then
	output := term.GetOutput()
	require.Contains(t, output, "Page")
}

func TestRender_EmptyMenu(t *testing.T) {
	// Given
	term := NewMockTerminal()
	renderer := NewRenderer(term)
	menu := NewMenu()

	// When
	renderer.Render(menu)

	// Then
	output := term.GetOutput()
	assert.Empty(t, output)
}

func TestRender_SmallTerminal(t *testing.T) {
	// Given
	term := NewMockTerminal()
	term.width = 5
	term.height = 2
	renderer := NewRenderer(term)
	menu := NewMenu()
	items := []Item{{Text: "verylongitemname", Type: "command"}}
	menu.Show(items, "")

	// When
	renderer.Render(menu)

	// Then
	output := term.GetOutput()
	require.Contains(t, output, "verylongitemname")
}

func TestClear_DisplayedMenu(t *testing.T) {
	// Given
	term := NewMockTerminal()
	renderer := NewRenderer(term)
	menu := NewMenu()
	items := []Item{{Text: "echo", Type: "builtin"}}
	menu.Show(items, "")
	menu.linesDrawn = 3

	// When
	renderer.Clear(menu)

	// Then
	output := term.GetOutput()
	require.Contains(t, output, "\033[0J") // clear from cursor
	require.Contains(t, output, "\033[u")  // restore cursor
}

func TestClear_NotDisplayedMenu(t *testing.T) {
	// Given
	term := NewMockTerminal()
	renderer := NewRenderer(term)
	menu := NewMenu()
	menu.displayed = false

	// When
	renderer.Clear(menu)

	// Then
	output := term.GetOutput()
	assert.Empty(t, output)
}

func TestRenderer_Pagination(t *testing.T) {
	term := NewMockTerminal()
	term.width = 20 // Small width to force pagination
	term.height = 8 // Small height
	renderer := NewRenderer(term)

	menu := NewMenu()
	items := make([]Item, 20)
	for i := range items {
		items[i] = Item{Text: fmt.Sprintf("item%d", i), Type: "command"}
	}
	menu.Show(items, "")

	renderer.Render(menu)
	output := term.GetOutput()

	// Should show pagination info
	if !strings.Contains(output, "Page") {
		t.Error("Should show pagination info")
	}
}

func TestRenderer_ClearNotDisplayed(t *testing.T) {
	term := NewMockTerminal()
	renderer := NewRenderer(term)

	menu := NewMenu()
	menu.displayed = false

	renderer.Clear(menu)
	output := term.GetOutput()

	// Should not output anything for non-displayed menu
	if output != "" {
		t.Errorf("Expected no output for non-displayed menu, got: %q", output)
	}
}

func TestRenderer_LayoutEdgeCases(t *testing.T) {
	term := NewMockTerminal()
	term.width = 5  // Very small width
	term.height = 2 // Very small height
	renderer := NewRenderer(term)

	menu := NewMenu()
	items := []Item{
		{Text: "verylongitemname", Type: "command"},
	}
	menu.Show(items, "")

	renderer.Render(menu)
	// Should not crash with small terminal
}
