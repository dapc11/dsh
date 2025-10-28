package readline

import (
	"fmt"
	"strings"
	"testing"

	"dsh/internal/terminal"
)

// MockTerminal for testing buffer management.
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

func (m *MockTerminal) Colorize(text string, _ terminal.Color) string {
	return text
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

func (m *MockTerminal) GetCursorPosition() (int, int) {
	return m.cursorX, m.cursorY
}

func (m *MockTerminal) SetCursorPosition(x, y int) {
	m.cursorX, m.cursorY = x, y
}

func (m *MockTerminal) Reset() {
	m.output.Reset()
	m.cursorX, m.cursorY = 0, 0
	m.savedX, m.savedY = 0, 0
}

func TestBufferManager_Creation(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	if bm == nil {
		t.Fatal("BufferManager should not be nil")
	}

	if bm.terminal == nil {
		t.Error("BufferManager should store terminal reference")
	}
}

func TestBufferManager_SaveRestore(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Save current state
	bm.SaveState()

	// Verify save cursor was called
	output := term.GetOutput()
	if !strings.Contains(output, "\033[s") {
		t.Error("Save cursor sequence not found in output")
	}

	term.Reset()

	// Restore state
	bm.RestoreState()

	// Verify restore cursor was called
	output = term.GetOutput()
	if !strings.Contains(output, "\033[u") {
		t.Error("Restore cursor sequence not found in output")
	}
}

func TestBufferManager_TemporaryBuffer(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Create temporary buffer
	buffer := bm.CreateTemporaryBuffer("test-menu")

	if buffer == nil {
		t.Fatal("Temporary buffer should not be nil")
	}

	if buffer.ID != "test-menu" {
		t.Errorf("Expected buffer ID 'test-menu', got '%s'", buffer.ID)
	}

	// Write to buffer
	buffer.WriteString("Hello World")

	// Render buffer
	bm.RenderBuffer(buffer)

	output := term.GetOutput()
	if !strings.Contains(output, "Hello World") {
		t.Error("Buffer content not found in output")
	}
}

func TestBufferManager_ClearBuffer(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Create and render buffer
	buffer := bm.CreateTemporaryBuffer("test-clear")
	buffer.WriteString("Content to clear")
	buffer.linesUsed = 3 // Simulate 3 lines used
	bm.RenderBuffer(buffer)

	term.Reset()

	// Clear buffer
	bm.ClearBuffer(buffer)

	output := term.GetOutput()
	// Should contain clear line sequences (one for each line used + 1)
	clearLineCount := strings.Count(output, "\033[2K")
	expectedClearLines := 4 // linesUsed + 1
	if clearLineCount != expectedClearLines {
		t.Errorf("Expected %d clear line sequences, got %d", expectedClearLines, clearLineCount)
	}
}

func TestBufferManager_MultipleBuffers(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Create multiple buffers
	buffer1 := bm.CreateTemporaryBuffer("menu1")
	buffer2 := bm.CreateTemporaryBuffer("menu2")

	if buffer1.ID == buffer2.ID {
		t.Error("Buffer IDs should be unique")
	}

	// Both should be tracked
	if len(bm.activeBuffers) != 2 {
		t.Errorf("Expected 2 active buffers, got %d", len(bm.activeBuffers))
	}
}

func TestBufferManager_CleanupAll(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Create multiple buffers
	buffer1 := bm.CreateTemporaryBuffer("menu1")
	buffer2 := bm.CreateTemporaryBuffer("menu2")
	buffer1.linesUsed = 2
	buffer2.linesUsed = 3

	// Cleanup all
	bm.CleanupAll()

	// All buffers should be cleared
	if len(bm.activeBuffers) != 0 {
		t.Errorf("Expected 0 active buffers after cleanup, got %d", len(bm.activeBuffers))
	}

	output := term.GetOutput()
	// Should contain clear line sequences for both buffers
	clearLineCount := strings.Count(output, "\033[2K")
	expectedClearLines := (2 + 1) + (3 + 1) // buffer1 + buffer2 lines
	if clearLineCount != expectedClearLines {
		t.Errorf("Expected %d clear line sequences, got %d", expectedClearLines, clearLineCount)
	}
}

func TestBufferManager_CursorPositionTracking(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Set mock terminal to specific position
	term.SetCursorPosition(10, 5)

	// Create buffer - should capture current cursor position
	buffer := bm.CreateTemporaryBufferAtCursor("positioned-menu")

	if buffer.startX != 10 {
		t.Errorf("Expected startX to be 10, got %d", buffer.startX)
	}

	if buffer.startY != 5 {
		t.Errorf("Expected startY to be 5, got %d", buffer.startY)
	}
}

func TestBufferManager_RestoreToSpecificPosition(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Create buffer at specific position
	buffer := bm.CreateTemporaryBufferAtCursor("positioned-menu")
	buffer.startX = 15
	buffer.startY = 8
	buffer.linesUsed = 3

	// Clear buffer - should clear precisely at that position
	bm.ClearBuffer(buffer)

	output := term.GetOutput()

	// Should move cursor to buffer start position
	expectedMove := fmt.Sprintf("\033[%d;%dH", buffer.startY+1, buffer.startX+1)
	if !strings.Contains(output, expectedMove) {
		t.Errorf("Expected cursor move to %s, but not found in output: %q", expectedMove, output)
	}

	// Should clear only the lines used by this buffer
	clearLineCount := strings.Count(output, "\033[2K")
	expectedClearLines := buffer.linesUsed + 1
	if clearLineCount != expectedClearLines {
		t.Errorf("Expected %d clear line sequences, got %d", expectedClearLines, clearLineCount)
	}
}

func TestBufferManager_MultipleBufferPositions(t *testing.T) {
	term := NewMockTerminal()
	bm := NewBufferManager(term)

	// Create buffers at different positions
	buffer1 := bm.CreateTemporaryBufferAtCursor("menu1")
	buffer1.startX = 5
	buffer1.startY = 10

	buffer2 := bm.CreateTemporaryBufferAtCursor("menu2")
	buffer2.startX = 20
	buffer2.startY = 15

	// Verify both buffers have different positions
	if buffer1.startX == buffer2.startX {
		t.Error("Buffers should have different X positions")
	}

	if buffer1.startY == buffer2.startY {
		t.Error("Buffers should have different Y positions")
	}
}
