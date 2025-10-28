package terminal

import (
	"bufio"
	"io"
)

// Key represents special keys.
type Key int

const (
	KeyNone Key = iota
	KeyEnter
	KeyTab
	KeyBackspace
	KeyDelete
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyEscape
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlK
	KeyCtrlL
	KeyCtrlN
	KeyCtrlP
	KeyCtrlR
	KeyCtrlU
	KeyCtrlW
	KeyCtrlY
	KeyCtrlZ
)

// KeyEvent represents a key press.
type KeyEvent struct {
	Key  Key
	Rune rune
	Alt  bool
	Ctrl bool
}

// InputReader handles terminal input.
type InputReader struct {
	reader *bufio.Reader
}

// NewInputReader creates an input reader.
func NewInputReader(r io.Reader) *InputReader {
	return &InputReader{
		reader: bufio.NewReader(r),
	}
}

// ReadKey reads a single key event.
func (r *InputReader) ReadKey() (KeyEvent, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return KeyEvent{}, err
	}

	// Handle escape sequences
	if b == 27 { // ESC
		return r.readEscapeSequence()
	}

	// Handle control characters
	if b < 32 {
		return KeyEvent{Key: r.ctrlKey(b)}, nil
	}

	// Handle printable characters
	return KeyEvent{Rune: rune(b)}, nil
}

// readEscapeSequence reads ANSI escape sequences.
func (r *InputReader) readEscapeSequence() (KeyEvent, error) {
	// Peek next byte
	next, err := r.reader.Peek(1)
	if err != nil || len(next) == 0 {
		return KeyEvent{Key: KeyEscape}, nil
	}

	if next[0] == '[' {
		r.reader.ReadByte() // consume '['
		return r.readCSISequence()
	}

	// Alt + key
	b, err := r.reader.ReadByte()
	if err != nil {
		return KeyEvent{Key: KeyEscape}, nil
	}

	return KeyEvent{Rune: rune(b), Alt: true}, nil
}

// readCSISequence reads Control Sequence Introducer sequences.
func (r *InputReader) readCSISequence() (KeyEvent, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return KeyEvent{}, err
	}

	switch b {
	case 'A':
		return KeyEvent{Key: KeyArrowUp}, nil
	case 'B':
		return KeyEvent{Key: KeyArrowDown}, nil
	case 'C':
		return KeyEvent{Key: KeyArrowRight}, nil
	case 'D':
		return KeyEvent{Key: KeyArrowLeft}, nil
	case 'H':
		return KeyEvent{Key: KeyHome}, nil
	case 'F':
		return KeyEvent{Key: KeyEnd}, nil
	default:
		return KeyEvent{Rune: rune(b)}, nil
	}
}

// ctrlKey maps control characters to keys.
func (r *InputReader) ctrlKey(b byte) Key {
	switch b {
	case 1:
		return KeyCtrlA
	case 2:
		return KeyCtrlB
	case 3:
		return KeyCtrlC
	case 4:
		return KeyCtrlD
	case 5:
		return KeyCtrlE
	case 6:
		return KeyCtrlF
	case 8, 127:
		return KeyBackspace
	case 9:
		return KeyTab
	case 10, 13:
		return KeyEnter
	case 11:
		return KeyCtrlK
	case 12:
		return KeyCtrlL
	case 14:
		return KeyCtrlN
	case 16:
		return KeyCtrlP
	case 18:
		return KeyCtrlR
	case 21:
		return KeyCtrlU
	case 23:
		return KeyCtrlW
	case 25:
		return KeyCtrlY
	case 26:
		return KeyCtrlZ
	default:
		return KeyNone
	}
}
