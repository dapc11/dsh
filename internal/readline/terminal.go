package readline

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	TIOCGWINSZ = 0x5413
)

// Terminal handles raw terminal operations.
type Terminal struct {
	fd       int
	original syscall.Termios
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

// Read reads from terminal.
func (t *Terminal) Read(buf []byte) (int, error) {
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return n, fmt.Errorf("terminal read error: %w", err)
	}

	return n, nil
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

// GetTerminalSize returns terminal width and height.
func (t *Terminal) GetTerminalSize() (int, int) {
	var ws struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return 80, 24 // Default fallback
	}

	return int(ws.Col), int(ws.Row)
}
