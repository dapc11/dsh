package rendering

import (
	"fmt"

	"dsh/internal/terminal"
)

// MockTerminalInterface implements terminal.TerminalInterface for testing.
type MockTerminalInterface struct {
	width   int
	height  int
	output  string
	rawMode bool
}

// NewMockTerminalInterface creates a mock terminal interface.
func NewMockTerminalInterface(width, height int) *MockTerminalInterface {
	return &MockTerminalInterface{
		width:   width,
		height:  height,
		output:  "",
		rawMode: false,
	}
}

// Terminal methods
func (m *MockTerminalInterface) Write(data []byte) (int, error) {
	m.output += string(data)
	return len(data), nil
}

func (m *MockTerminalInterface) WriteString(s string) (int, error) {
	m.output += s
	return len(s), nil
}

func (m *MockTerminalInterface) Printf(format string, args ...interface{}) (int, error) {
	formatted := fmt.Sprintf(format, args...)
	return m.WriteString(formatted)
}

func (m *MockTerminalInterface) Size() (width, height int) {
	return m.width, m.height
}

func (m *MockTerminalInterface) MoveCursor(x, y int) {
	m.output += fmt.Sprintf("\033[%d;%dH", y+1, x+1)
}

func (m *MockTerminalInterface) ClearLine() {
	m.output += "\033[2K"
}

func (m *MockTerminalInterface) ClearToEnd() {
	m.output += "\033[K"
}

func (m *MockTerminalInterface) ClearFromCursor() {
	m.output += "\033[0J"
}

func (m *MockTerminalInterface) SaveCursor() {
	m.output += "\033[s"
}

func (m *MockTerminalInterface) RestoreCursor() {
	m.output += "\033[u"
}

func (m *MockTerminalInterface) HideCursor() {
	m.output += "\033[?25l"
}

func (m *MockTerminalInterface) ShowCursor() {
	m.output += "\033[?25h"
}

// ColorManager methods
func (m *MockTerminalInterface) Colorize(text string, color terminal.Color) string {
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

func (m *MockTerminalInterface) StyleText(text string, style terminal.Style) string {
	if style.Reverse {
		return "\033[7m" + text + "\033[0m"
	}
	return text
}

// InputReader methods
func (m *MockTerminalInterface) ReadKey() (terminal.KeyEvent, error) {
	return terminal.KeyEvent{}, nil
}

// Interface methods
func (m *MockTerminalInterface) EnableRawMode() error {
	m.rawMode = true
	return nil
}

func (m *MockTerminalInterface) DisableRawMode() error {
	m.rawMode = false
	return nil
}

func (m *MockTerminalInterface) IsRawMode() bool {
	return m.rawMode
}

func (m *MockTerminalInterface) Cleanup() {
	m.ShowCursor()
	_ = m.DisableRawMode()
}

// Test helper methods
func (m *MockTerminalInterface) GetOutput() string {
	return m.output
}

func (m *MockTerminalInterface) ClearOutput() {
	m.output = ""
}
