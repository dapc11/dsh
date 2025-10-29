package readline

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"dsh/internal/terminal"
)

// Terminal handles raw terminal operations with foundation.
type Terminal struct {
	fd            int
	original      syscall.Termios
	termInterface *terminal.Interface
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
		fd:            fd,
		original:      original,
		termInterface: terminal.NewInterface(),
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

	err := setTermios(t.fd, &raw)
	if err == nil {
		_ = t.termInterface.EnableRawMode()
	}
	return err
}

// Restore restores original terminal mode.
func (t *Terminal) Restore() error {
	_ = t.termInterface.DisableRawMode()
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

// GetTerminalSize returns terminal width and height.
func (t *Terminal) GetTerminalSize() (int, int) {
	var ws struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws))) //nolint:gosec // Required for terminal size detection
	if errno != 0 {
		return 80, 24 // Default fallback
	}

	return int(ws.Col), int(ws.Row)
}

// Interface returns the terminal interface for advanced operations.
func (t *Terminal) Interface() *terminal.Interface {
	return t.termInterface
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
