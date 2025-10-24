// Package readline provides emacs-like line editing functionality.
package readline

import (
	"errors"
	"fmt"
)

var (
	// ErrEOF indicates end of file condition.
	ErrEOF = errors.New("EOF")
)

// Readline provides emacs-like line editing functionality.
type Readline struct {
	prompt       string
	terminal     *Terminal
	history      *History
	killRing     *KillRing
	buffer       []rune
	cursor       int
	suggestion   string
	searchPrefix string
	browseMode   bool
	color        *Color
}

// New creates a new readline instance.
func New(prompt string) (*Readline, error) {
	terminal, err := NewTerminal()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize terminal: %w", err)
	}

	return &Readline{
		prompt:       prompt,
		terminal:     terminal,
		history:      NewHistory(),
		killRing:     NewKillRing(),
		buffer:       make([]rune, 0, 256),
		cursor:       0,
		suggestion:   "",
		searchPrefix: "",
		browseMode:   false,
		color:        NewColor(),
	}, nil
}

// ReadLine reads a line with emacs-like editing.
func (r *Readline) ReadLine() (string, error) {
	err := r.terminal.SetRawMode()
	if err != nil {
		return "", fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer func() { _ = r.terminal.Restore() }()

	r.buffer = r.buffer[:0]
	r.cursor = 0
	r.history.ResetPosition()

	r.displayPrompt()

	for {
		ch, err := r.readChar()
		if err != nil {
			return "", fmt.Errorf("failed to read character: %w", err)
		}

		if r.handleKey(ch) {
			continue
		}

		// Check for EOF case
		if ch == KeyCtrlD && len(r.buffer) == 0 {
			return "", ErrEOF
		}

		// Return completed line
		r.moveCursorToEnd()
		_, _ = fmt.Print("\r\n") //nolint:forbidigo
		line := string(r.buffer)
		if line != "" {
			r.history.Add(line)
		}

		return line, nil
	}
}

func (r *Readline) readChar() (byte, error) {
	buf := make([]byte, 1)
	_, err := r.terminal.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read character: %w", err)
	}

	return buf[0], nil
}

func (r *Readline) displayPrompt() {
	_, _ = fmt.Print(r.prompt) //nolint:forbidigo
}
