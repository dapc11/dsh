package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	ErrEOF = errors.New("EOF")
)

// Readline provides emacs-like line editing functionality.
type Readline struct {
	buffer    []rune
	cursor    int
	history   []string
	histPos   int
	prompt    string
	terminal  *Terminal
	killRing  []string // Store multiple killed items
	killIndex int      // Current position in kill ring
	lastYank  int      // Length of last yanked text
}

// Terminal handles raw terminal operations.
type Terminal struct {
	fd       int
	original syscall.Termios
}

// Key codes for special keys.
const (
	KeyCtrlA     = 1
	KeyCtrlB     = 2
	KeyCtrlC     = 3
	KeyCtrlD     = 4
	KeyCtrlE     = 5
	KeyCtrlF     = 6
	KeyCtrlK     = 11
	KeyCtrlL     = 12
	KeyCtrlN     = 14
	KeyCtrlP     = 16
	KeyCtrlU     = 21
	KeyCtrlW     = 23
	KeyCtrlY     = 25
	KeyCtrlZ     = 26
	KeyEscape    = 27
	KeyBackspace = 127
	KeyDelete    = 127
)

// NewReadline creates a new readline instance.
func NewReadline(prompt string) (*Readline, error) {
	terminal, err := NewTerminal()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize terminal: %w", err)
	}

	return &Readline{
		buffer:    make([]rune, 0, 256),
		cursor:    0,
		history:   make([]string, 0, 100),
		histPos:   0,
		prompt:    prompt,
		terminal:  terminal,
		killRing:  make([]string, 0, 10),
		killIndex: 0,
		lastYank:  0,
	}, nil
}

// NewTerminal initializes terminal for raw mode.
func NewTerminal() (*Terminal, error) {
	fd := int(os.Stdin.Fd())

	var original syscall.Termios
	err := getTermios(fd, &original)
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal attributes: %w", err)
	}

	return &Terminal{
		fd:       fd,
		original: original,
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
	r.histPos = len(r.history)

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
			r.addToHistory(line)
		}

		return line, nil
	}
}

func (r *Readline) handleKey(ch byte) bool {
	switch ch {
	case '\r', '\n':
		return false // Signal completion
	case KeyCtrlA:
		r.resetYank()
		r.moveCursorToStart()
	case KeyCtrlE:
		r.resetYank()
		r.moveCursorToEnd()
	case KeyCtrlB:
		r.resetYank()
		r.moveCursorLeft()
	case KeyCtrlC:
		r.resetYank()
		r.clearLine()
		_, _ = fmt.Print("^C\r\n") //nolint:forbidigo
		r.displayPrompt()
	case KeyCtrlF:
		r.resetYank()
		r.moveCursorRight()
	case KeyCtrlD:
		r.resetYank()
		if len(r.buffer) == 0 {
			return false // Will be handled as EOF in caller
		}
		r.deleteChar()
	case KeyCtrlK:
		r.resetYank()
		r.killToEnd()
	case KeyCtrlU:
		r.resetYank()
		r.killLine()
	case KeyCtrlW:
		r.resetYank()
		r.killWordBackward()
	case KeyCtrlY:
		r.yank()
	case KeyCtrlL:
		r.resetYank()
		r.clearScreen()
	case KeyCtrlP:
		r.resetYank()
		r.historyPrevious()
	case KeyCtrlN:
		r.resetYank()
		r.historyNext()
	case KeyCtrlZ:
		r.resetYank()
		r.clearLine()
		_, _ = fmt.Print("^Z\r\n") //nolint:forbidigo
		r.displayPrompt()
	case KeyBackspace:
		r.resetYank()
		r.backspace()
	case KeyEscape:
		err := r.handleEscapeSequence()
		if err != nil {
			return true // Continue on error
		}
	default:
		if ch >= 32 && ch < 127 {
			r.resetYank()
			r.insertChar(rune(ch))
		}
	}

	return true // Continue reading
}

func (r *Readline) readChar() (byte, error) {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read character: %w", err)
	}

	return buf[0], nil
}

func (r *Readline) handleEscapeSequence() error {
	ch1, err := r.readChar()
	if err != nil {
		return err
	}

	switch ch1 {
	case '[':
		ch2, err := r.readChar()
		if err != nil {
			return err
		}

		switch ch2 {
		case 'A': // Up arrow
			r.historyPrevious()
		case 'B': // Down arrow
			r.historyNext()
		case 'C': // Right arrow
			r.moveCursorRight()
		case 'D': // Left arrow
			r.moveCursorLeft()
		case '1':
			ch3, err := r.readChar()
			if err == nil && ch3 == ';' {
				ch4, err := r.readChar()
				if err == nil && ch4 == '5' {
					ch5, err := r.readChar()
					if err == nil {
						switch ch5 {
						case 'C': // Ctrl+Right
							r.moveWordForward()
						case 'D': // Ctrl+Left
							r.moveWordBackward()
						}
					}
				}
			}
		}
	case 'y': // Alt+Y (forward)
		r.yankPop()
	case 'Y': // Alt+Shift+Y (backward)
		r.yankCycle(-1)
	case 'd': // Alt+D (delete word forward)
		r.resetYank()
		r.killWordForward()
	}

	return nil
}

func (r *Readline) displayPrompt() {
	_, _ = fmt.Print(r.prompt) //nolint:forbidigo
}

func (r *Readline) redraw() {
	_, _ = fmt.Print("\r" + strings.Repeat(" ", len(r.prompt)+len(r.buffer)+10) + "\r") //nolint:forbidigo
	_, _ = fmt.Print(r.prompt + string(r.buffer))                                       //nolint:forbidigo
	r.setCursorPosition()
}

func (r *Readline) setCursorPosition() {
	if r.cursor < len(r.buffer) {
		_, _ = fmt.Printf("\r%s%s\r%s%s", //nolint:forbidigo
			r.prompt,
			string(r.buffer),
			r.prompt,
			string(r.buffer[:r.cursor]))
	}
}

func (r *Readline) insertChar(ch rune) {
	if r.cursor == len(r.buffer) {
		r.buffer = append(r.buffer, ch)
		_, _ = fmt.Printf("%c", ch) //nolint:forbidigo
	} else {
		r.buffer = append(r.buffer[:r.cursor+1], r.buffer[r.cursor:]...)
		r.buffer[r.cursor] = ch
		r.redraw()
	}
	r.cursor++
}

func (r *Readline) backspace() {
	if r.cursor > 0 {
		r.buffer = append(r.buffer[:r.cursor-1], r.buffer[r.cursor:]...)
		r.cursor--
		r.redraw()
	}
}

func (r *Readline) deleteChar() {
	if r.cursor < len(r.buffer) {
		r.buffer = append(r.buffer[:r.cursor], r.buffer[r.cursor+1:]...)
		r.redraw()
	}
}

func (r *Readline) moveCursorLeft() {
	if r.cursor > 0 {
		r.cursor--
		_, _ = fmt.Print("\b") //nolint:forbidigo
	}
}

func (r *Readline) moveCursorRight() {
	if r.cursor < len(r.buffer) {
		_, _ = fmt.Printf("%c", r.buffer[r.cursor]) //nolint:forbidigo
		r.cursor++
	}
}

func (r *Readline) moveCursorToStart() {
	for r.cursor > 0 {
		r.moveCursorLeft()
	}
}

func (r *Readline) moveCursorToEnd() {
	for r.cursor < len(r.buffer) {
		r.moveCursorRight()
	}
}

func (r *Readline) clearLine() {
	r.buffer = r.buffer[:0]
	r.cursor = 0
}

func (r *Readline) resetYank() {
	r.lastYank = 0
}

func (r *Readline) addToKillRing(text string) {
	if text == "" {
		return
	}

	// Add to front of kill ring
	r.killRing = append([]string{text}, r.killRing...)

	// Limit kill ring size
	if len(r.killRing) > 10 {
		r.killRing = r.killRing[:10]
	}

	r.killIndex = 0
}

func (r *Readline) killToEnd() {
	if r.cursor < len(r.buffer) {
		killed := string(r.buffer[r.cursor:])
		r.addToKillRing(killed)
		r.buffer = r.buffer[:r.cursor]
		_, _ = fmt.Print("\033[K") //nolint:forbidigo // Clear to end of line
	}
}

func (r *Readline) killLine() {
	killed := string(r.buffer)
	r.addToKillRing(killed)
	r.buffer = r.buffer[:0]
	r.cursor = 0
	r.redraw()
}

func (r *Readline) killWordBackward() {
	if r.cursor == 0 {
		return
	}

	start := r.cursor - 1
	for start > 0 && r.buffer[start] == ' ' {
		start--
	}
	for start > 0 && r.buffer[start] != ' ' {
		start--
	}
	if r.buffer[start] == ' ' {
		start++
	}

	killed := string(r.buffer[start:r.cursor])
	r.addToKillRing(killed)
	r.buffer = append(r.buffer[:start], r.buffer[r.cursor:]...)
	r.cursor = start
	r.redraw()
}

func (r *Readline) yank() {
	if len(r.killRing) == 0 {
		return
	}

	r.killIndex = 0
	yankText := []rune(r.killRing[r.killIndex])
	r.lastYank = len(yankText)

	if r.cursor >= len(r.buffer) {
		// At end of buffer
		r.buffer = append(r.buffer, yankText...)
		_, _ = fmt.Print(r.killRing[r.killIndex]) //nolint:forbidigo
	} else {
		// In middle of buffer - make space and insert
		newBuffer := make([]rune, len(r.buffer)+len(yankText))
		copy(newBuffer, r.buffer[:r.cursor])
		copy(newBuffer[r.cursor:], yankText)
		copy(newBuffer[r.cursor+len(yankText):], r.buffer[r.cursor:])
		r.buffer = newBuffer
		r.redraw()
	}
	r.cursor += len(yankText)
}

func (r *Readline) yankPop() {
	r.yankCycle(1) // Forward direction
}

func (r *Readline) yankCycle(direction int) {
	if len(r.killRing) <= 1 || r.lastYank == 0 {
		return
	}

	// Calculate start position of previously yanked text
	startPos := r.cursor - r.lastYank
	if startPos < 0 {
		startPos = 0
	}

	// Remove the previously yanked text
	r.buffer = append(r.buffer[:startPos], r.buffer[r.cursor:]...)
	r.cursor = startPos

	// Cycle through kill ring (with proper wrap-around)
	newIndex := r.killIndex + direction
	if newIndex < 0 {
		newIndex = len(r.killRing) - 1
	} else if newIndex >= len(r.killRing) {
		newIndex = 0
	}
	r.killIndex = newIndex

	yankText := []rune(r.killRing[r.killIndex])
	r.lastYank = len(yankText)

	// Insert new text at cursor position
	if r.cursor >= len(r.buffer) {
		// At end of buffer
		r.buffer = append(r.buffer, yankText...)
	} else {
		// In middle of buffer - make space and insert
		newBuffer := make([]rune, len(r.buffer)+len(yankText))
		copy(newBuffer, r.buffer[:r.cursor])
		copy(newBuffer[r.cursor:], yankText)
		copy(newBuffer[r.cursor+len(yankText):], r.buffer[r.cursor:])
		r.buffer = newBuffer
	}

	r.cursor += len(yankText)
	r.redraw()
}

func (r *Readline) killWordForward() {
	if r.cursor >= len(r.buffer) {
		return
	}

	end := r.cursor
	// Skip whitespace
	for end < len(r.buffer) && r.buffer[end] == ' ' {
		end++
	}
	// Skip word characters
	for end < len(r.buffer) && r.buffer[end] != ' ' {
		end++
	}

	if end > r.cursor {
		killed := string(r.buffer[r.cursor:end])
		r.addToKillRing(killed)
		r.buffer = append(r.buffer[:r.cursor], r.buffer[end:]...)
		r.redraw()
	}
}

func (r *Readline) moveWordForward() {
	for r.cursor < len(r.buffer) && r.buffer[r.cursor] != ' ' {
		r.moveCursorRight()
	}
	for r.cursor < len(r.buffer) && r.buffer[r.cursor] == ' ' {
		r.moveCursorRight()
	}
}

func (r *Readline) moveWordBackward() {
	if r.cursor == 0 {
		return
	}

	r.moveCursorLeft()
	for r.cursor > 0 && r.buffer[r.cursor] == ' ' {
		r.moveCursorLeft()
	}
	for r.cursor > 0 && r.buffer[r.cursor] != ' ' {
		r.moveCursorLeft()
	}
	if r.cursor > 0 && r.buffer[r.cursor] == ' ' {
		r.moveCursorRight()
	}
}

func (r *Readline) clearScreen() {
	_, _ = fmt.Print("\033[2J\033[H") //nolint:forbidigo // Clear screen and move to top
	r.displayPrompt()
	_, _ = fmt.Print(string(r.buffer)) //nolint:forbidigo
	r.setCursorPosition()
}

func (r *Readline) addToHistory(line string) {
	if len(r.history) == 0 || r.history[len(r.history)-1] != line {
		r.history = append(r.history, line)
		if len(r.history) > 100 {
			r.history = r.history[1:]
		}
	}
}

func (r *Readline) historyPrevious() {
	if r.histPos > 0 {
		r.histPos--
		r.setBufferFromHistory()
	}
}

func (r *Readline) historyNext() {
	if r.histPos < len(r.history) {
		r.histPos++
		r.setBufferFromHistory()
	}
}

func (r *Readline) setBufferFromHistory() {
	if r.histPos < len(r.history) {
		r.buffer = []rune(r.history[r.histPos])
	} else {
		r.buffer = r.buffer[:0]
	}
	r.cursor = len(r.buffer)
	r.redraw()
}

// SetRawMode enables raw terminal mode.
func (t *Terminal) SetRawMode() error {
	raw := t.original
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Cflag |= syscall.CS8
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	return setTermios(t.fd, &raw)
}

// Restore restores original terminal mode.
func (t *Terminal) Restore() error {
	return setTermios(t.fd, &t.original)
}

func getTermios(fd int, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(termios))) //nolint:gosec
	if errno != 0 {
		return errno
	}

	return nil
}

func setTermios(fd int, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(termios))) //nolint:gosec
	if errno != 0 {
		return errno
	}

	return nil
}
