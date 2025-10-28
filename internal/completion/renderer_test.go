package completion

import (
	"fmt"
	"strings"
	"testing"

	"dsh/internal/terminal"
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

func TestRenderer_Creation(t *testing.T) {
	term := NewMockTerminal()
	renderer := NewRenderer(term)

	if renderer == nil {
		t.Fatal("Renderer should not be nil")
	}
}

func TestRenderer_RenderMenu(t *testing.T) {
	term := NewMockTerminal()
	renderer := NewRenderer(term)

	menu := NewMenu()
	items := []Item{
		{Text: "echo", Type: "builtin"},
		{Text: "ls", Type: "command"},
		{Text: "dir/", Type: "directory"},
	}
	menu.Show(items, "")

	renderer.Render(menu)

	output := term.GetOutput()

	// Debug: print the actual output
	t.Logf("Actual output: %q", output)

	// Should contain save cursor
	if !strings.Contains(output, "\033[s") {
		t.Error("Output should contain save cursor sequence")
	}

	// Should contain the items
	if !strings.Contains(output, "echo") {
		t.Error("Output should contain 'echo'")
	}

	if !strings.Contains(output, "ls") {
		t.Error("Output should contain 'ls'")
	}

	// Should contain colors for non-selected items
	if !strings.Contains(output, "\033[32m") { // Green for command
		t.Error("Output should contain green color for command")
	}

	if !strings.Contains(output, "\033[34m") { // Blue for directory
		t.Error("Output should contain blue color for directory")
	}

	// Should contain selection (reverse video) for first item
	if !strings.Contains(output, "\033[7m") {
		t.Error("Output should contain reverse video for selection")
	}

	// The selected item (echo) should have reverse video, not cyan
	// Let's test with a different selection to see cyan
	menu.NextItem() // Move to second item
	term.Reset()
	renderer.Render(menu)

	output2 := term.GetOutput()
	if !strings.Contains(output2, "\033[36m") { // Cyan for builtin when not selected
		t.Error("Output should contain cyan color for builtin when not selected")
	}
}

func TestRenderer_ClearMenu(t *testing.T) {
	term := NewMockTerminal()
	renderer := NewRenderer(term)

	menu := NewMenu()
	items := []Item{
		{Text: "echo", Type: "builtin"},
	}
	menu.Show(items, "")
	menu.linesDrawn = 3 // Simulate rendered lines

	renderer.Clear(menu)

	output := term.GetOutput()

	// Should contain clear sequence
	if !strings.Contains(output, "\033[0J") {
		t.Error("Output should contain clear from cursor sequence")
	}

	// Should contain restore cursor
	if !strings.Contains(output, "\033[u") {
		t.Error("Output should contain restore cursor sequence")
	}
}

func TestRenderer_EmptyMenu(t *testing.T) {
	term := NewMockTerminal()
	renderer := NewRenderer(term)

	menu := NewMenu()
	// Don't show menu or add items

	renderer.Render(menu)

	output := term.GetOutput()

	// Should not render anything for empty menu
	if output != "" {
		t.Errorf("Expected empty output for empty menu, got: %q", output)
	}
}

func TestRenderer_SingleItem(t *testing.T) {
	term := NewMockTerminal()
	renderer := NewRenderer(term)

	menu := NewMenu()
	items := []Item{
		{Text: "single", Type: "command"},
	}
	menu.Show(items, "")

	renderer.Render(menu)

	output := term.GetOutput()

	// Should contain the single item
	if !strings.Contains(output, "single") {
		t.Error("Output should contain 'single'")
	}

	// Should be selected (reverse video)
	if !strings.Contains(output, "\033[7m") {
		t.Error("Single item should be selected")
	}
}
