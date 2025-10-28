package readline

import (
	"strings"

	"dsh/internal/terminal"
)

// CursorTracker interface for terminals that can track cursor position.
type CursorTracker interface {
	GetCursorPosition() (int, int)
}

// TemporaryBuffer represents a temporary screen buffer for menus, completions, etc.
type TemporaryBuffer struct {
	ID        string
	content   strings.Builder
	linesUsed int
	startX    int
	startY    int
	active    bool
}

// WriteString writes content to the buffer.
func (tb *TemporaryBuffer) WriteString(s string) {
	tb.content.WriteString(s)
	// Count newlines to track lines used
	tb.linesUsed += strings.Count(s, "\n")
}

// GetContent returns the buffer content.
func (tb *TemporaryBuffer) GetContent() string {
	return tb.content.String()
}

// Clear clears the buffer content.
func (tb *TemporaryBuffer) Clear() {
	tb.content.Reset()
	tb.linesUsed = 0
}

// BufferManager manages temporary screen buffers and cursor state.
type BufferManager struct {
	terminal      terminal.TerminalInterface
	activeBuffers map[string]*TemporaryBuffer
	savedCursor   bool
}

// NewBufferManager creates a new buffer manager.
func NewBufferManager(term terminal.TerminalInterface) *BufferManager {
	return &BufferManager{
		terminal:      term,
		activeBuffers: make(map[string]*TemporaryBuffer),
		savedCursor:   false,
	}
}

// SaveState saves the current cursor position and terminal state.
func (bm *BufferManager) SaveState() {
	if !bm.savedCursor {
		bm.terminal.SaveCursor()
		bm.savedCursor = true
	}
}

// RestoreState restores the saved cursor position and terminal state.
func (bm *BufferManager) RestoreState() {
	if bm.savedCursor {
		bm.terminal.RestoreCursor()
		bm.savedCursor = false
	}
}

// CreateTemporaryBuffer creates a new temporary buffer.
func (bm *BufferManager) CreateTemporaryBuffer(id string) *TemporaryBuffer {
	buffer := &TemporaryBuffer{
		ID:     id,
		active: true,
		startX: 0, // Default position
		startY: 0, // Default position
	}

	bm.activeBuffers[id] = buffer
	return buffer
}

// CreateTemporaryBufferAtCursor creates a new temporary buffer at the current cursor position.
func (bm *BufferManager) CreateTemporaryBufferAtCursor(id string) *TemporaryBuffer {
	buffer := &TemporaryBuffer{
		ID:     id,
		active: true,
		startX: 0,
		startY: 0,
	}

	// If terminal supports cursor tracking, get current position
	if tracker, ok := bm.terminal.(CursorTracker); ok {
		buffer.startX, buffer.startY = tracker.GetCursorPosition()
	}

	bm.activeBuffers[id] = buffer
	return buffer
}

// CreateTemporaryBufferAfterCursor creates a buffer positioned after the current cursor.
func (bm *BufferManager) CreateTemporaryBufferAfterCursor(id string) *TemporaryBuffer {
	buffer := &TemporaryBuffer{
		ID:     id,
		active: true,
		startX: 0, // Always start at beginning of line
	}

	// Position buffer on the line after current cursor
	if tracker, ok := bm.terminal.(CursorTracker); ok {
		_, currentY := tracker.GetCursorPosition()
		buffer.startY = currentY + 1
	} else {
		buffer.startY = 1 // Default to line 1 if no tracking
	}

	bm.activeBuffers[id] = buffer
	return buffer
}

// RenderBuffer renders a temporary buffer to the terminal.
func (bm *BufferManager) RenderBuffer(buffer *TemporaryBuffer) {
	if buffer == nil || !buffer.active {
		return
	}

	// Move to buffer start position (don't clear - just write over)
	bm.terminal.MoveCursor(buffer.startX, buffer.startY)

	content := buffer.GetContent()
	if content != "" {
		bm.terminal.WriteString(content)
	}
}

// ClearBuffer clears a temporary buffer from the screen.
func (bm *BufferManager) ClearBuffer(buffer *TemporaryBuffer) {
	if buffer == nil || !buffer.active {
		return
	}

	// Move cursor to buffer start position
	bm.terminal.MoveCursor(buffer.startX, buffer.startY)

	// Clear only the lines used by this buffer
	for i := 0; i <= buffer.linesUsed; i++ {
		bm.terminal.ClearLine()
		if i < buffer.linesUsed {
			bm.terminal.MoveCursor(0, buffer.startY+i+1)
		}
	}

	// Mark buffer as inactive
	buffer.active = false
	delete(bm.activeBuffers, buffer.ID)
}

// CleanupAll clears all active buffers and restores state.
func (bm *BufferManager) CleanupAll() {
	// Clear all active buffers
	for _, buffer := range bm.activeBuffers {
		if buffer.active {
			bm.ClearBuffer(buffer)
		}
	}

	// Clear the map
	bm.activeBuffers = make(map[string]*TemporaryBuffer)

	// Restore cursor state
	bm.RestoreState()
}

// GetBuffer retrieves a buffer by ID.
func (bm *BufferManager) GetBuffer(id string) *TemporaryBuffer {
	return bm.activeBuffers[id]
}

// HasActiveBuffers returns true if there are any active buffers.
func (bm *BufferManager) HasActiveBuffers() bool {
	return len(bm.activeBuffers) > 0
}

// ClearAt clears text at specific position (for selective redraw)
func (bm *BufferManager) ClearAt(line, col int, text string) error {
	bm.terminal.MoveCursor(col, line)
	// Clear the text by overwriting with spaces
	spaces := strings.Repeat(" ", len(text))
	bm.terminal.WriteString(spaces)
	return nil
}

// DrawAt draws text at specific position with optional highlighting
func (bm *BufferManager) DrawAt(line, col int, text string, highlight bool) error {
	bm.terminal.MoveCursor(col, line)
	if highlight {
		bm.terminal.WriteString("\033[7m" + text + "\033[0m") // Reverse video
	} else {
		bm.terminal.WriteString(text)
	}
	return nil
}
