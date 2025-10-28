// Package readline provides emacs-like line editing functionality.
package readline

import (
	"errors"
	"fmt"

	"dsh/internal/terminal"
)

var (
	// ErrEOF indicates end of file condition.
	ErrEOF = errors.New("EOF")
)

// Readline provides emacs-like line editing functionality.
type Readline struct {
	prompt         string
	terminal       *terminal.Interface
	rawTerminal    *Terminal
	history        *History
	killRing       *KillRing
	buffer         []rune
	cursor         int
	suggestion     string
	searchPrefix   string
	browseMode     bool
	completion     *Completion
	completionMenu *CompletionMenu
	bufferManager  *BufferManager
}

// New creates a new readline instance.
func New(prompt string) (*Readline, error) {
	termInterface := terminal.NewInterface()

	rawTerminal, err := NewTerminal()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize terminal: %w", err)
	}

	bufferManager := NewBufferManager(termInterface)

	return &Readline{
		prompt:         prompt,
		terminal:       termInterface,
		rawTerminal:    rawTerminal,
		history:        NewHistory(),
		killRing:       NewKillRing(),
		buffer:         make([]rune, 0, 256),
		cursor:         0,
		suggestion:     "",
		searchPrefix:   "",
		browseMode:     false,
		completion:     NewCompletion(),
		completionMenu: NewCompletionMenu(termInterface),
		bufferManager:  bufferManager,
	}, nil
}

// ReadLine reads a line with emacs-like editing.
func (r *Readline) ReadLine() (string, error) {
	err := r.rawTerminal.SetRawMode()
	if err != nil {
		return "", fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer func() {
		// Clean up any temporary buffers before restoring terminal
		r.bufferManager.CleanupAll()
		_ = r.rawTerminal.Restore()
	}()

	r.buffer = r.buffer[:0]
	r.cursor = 0
	r.history.ResetPosition()

	r.displayPrompt()

	for {
		keyEvent, err := r.terminal.ReadKey()
		if err != nil {
			return "", fmt.Errorf("failed to read key: %w", err)
		}

		if r.handleKeyEvent(keyEvent) {
			continue
		}

		// Check for EOF case
		if keyEvent.Key == terminal.KeyCtrlD && len(r.buffer) == 0 {
			return "", ErrEOF
		}

		// Return completed line
		r.moveCursorToEnd()
		r.terminal.WriteString("\r\n")
		line := string(r.buffer)
		if line != "" {
			r.history.Add(line)
		}

		return line, nil
	}
}

func (r *Readline) displayPrompt() {
	if r.terminal != nil {
		r.terminal.WriteString(r.prompt)
	}
}
