package rendering

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"dsh/internal/completion"
)

// CaptureStdout captures stdout during function execution.
func CaptureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// TestActualRenderingOutput tests the real rendering output.
func TestActualRenderingOutput(t *testing.T) {
	colorProvider := &MockColorProvider{}
	terminalProvider := &MockTerminalProvider{width: 80, height: 24}

	renderer := completion.NewRenderer(colorProvider, terminalProvider)
	menu := completion.NewMenu()

	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
		{Text: "ls", Type: "command"},
		{Text: "file.txt", Type: "file"},
	}

	menu.Show(items, "")

	// Capture actual rendering output
	output := CaptureStdout(func() {
		renderer.Render(menu)
	})

	t.Logf("Actual rendering output: %q", output)

	// Validate ANSI color codes are present
	if !strings.Contains(output, "\033[7m") {
		t.Error("Output should contain reverse video ANSI code for selection")
	}

	if !strings.Contains(output, "\033[32m") {
		t.Error("Output should contain green ANSI code for command")
	}

	if !strings.Contains(output, "\033[0m") {
		t.Error("Output should contain reset ANSI code")
	}

	// Validate text content
	if !strings.Contains(output, "echo") {
		t.Error("Output should contain 'echo' text")
	}

	if !strings.Contains(output, "ls") {
		t.Error("Output should contain 'ls' text")
	}

	// Test that selection overrides type color (echo is selected, so reverse instead of cyan)
	if strings.Contains(output, "\033[36mecho") {
		t.Error("Selected item should not have type color (should be reverse instead)")
	}

	// Test that non-selected items get type colors
	if !strings.Contains(output, "\033[32mls") {
		t.Error("Non-selected command should have green color")
	}
}

// TestColorFormatValidation tests specific color format expectations.
func TestColorFormatValidation(t *testing.T) {
	t.Parallel()
	colorProvider := &MockColorProvider{}

	tests := []struct {
		text     string
		color    string
		expected string
	}{
		{"echo", "reverse", "\033[7mecho\033[0m"},
		{"ls", "green", "\033[32mls\033[0m"},
		{"dir/", "blue", "\033[34mdir/\033[0m"},
		{"help", "cyan", "\033[36mhelp\033[0m"},
	}

	for _, test := range tests {
		t.Run(test.text+"_"+test.color, func(t *testing.T) {
			result := colorProvider.Colorize(test.text, test.color)

			if result != test.expected {
				t.Errorf("Color format mismatch:\nExpected: %q\nGot:      %q", test.expected, result)
			}
		})
	}
}

// TestRenderingRegression tests that rendering behavior doesn't regress.
func TestRenderingRegression(t *testing.T) {
	colorProvider := &MockColorProvider{}
	terminalProvider := &MockTerminalProvider{width: 80, height: 24}

	renderer := completion.NewRenderer(colorProvider, terminalProvider)
	menu := completion.NewMenu()

	// Test with known items
	items := []completion.Item{
		{Text: "echo", Type: "builtin"},
		{Text: "exit", Type: "builtin"},
	}

	menu.Show(items, "")

	// Capture baseline output
	baseline := CaptureStdout(func() {
		renderer.Render(menu)
	})

	// Expected patterns that should always be present
	expectedPatterns := []string{
		"\033[7mecho\033[0m",  // Selected item in reverse
		"\033[36mexit\033[0m", // Builtin in cyan
		"\r\n",                // Line breaks
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(baseline, pattern) {
			t.Errorf("Baseline output missing expected pattern: %q\nOutput: %q", pattern, baseline)
		}
	}

	t.Logf("Baseline rendering validated: %q", baseline)
}
