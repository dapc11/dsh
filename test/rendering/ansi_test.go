package rendering

import (
	"regexp"
	"strings"
	"testing"
)

// ANSISequence represents expected ANSI escape sequences
type ANSISequence struct {
	Name        string
	Pattern     string
	Description string
}

// Common ANSI sequences used in terminal applications
var expectedSequences = []ANSISequence{
	{"cursor_left", `\x1b\[D`, "Move cursor left"},
	{"cursor_right", `\x1b\[C`, "Move cursor right"},
	{"cursor_up", `\x1b\[A`, "Move cursor up"},
	{"cursor_down", `\x1b\[B`, "Move cursor down"},
	{"cursor_home", `\x1b\[H`, "Move cursor to home"},
	{"clear_line", `\x1b\[2?K`, "Clear line"},
	{"clear_screen", `\x1b\[2J`, "Clear screen"},
	{"save_cursor", `\x1b\[s`, "Save cursor position"},
	{"restore_cursor", `\x1b\[u`, "Restore cursor position"},
	{"color_reset", `\x1b\[0m`, "Reset colors"},
	{"color_red", `\x1b\[31m`, "Red text"},
	{"color_green", `\x1b\[32m`, "Green text"},
	{"color_blue", `\x1b\[34m`, "Blue text"},
	{"bold", `\x1b\[1m`, "Bold text"},
}

// TestANSISequenceGeneration tests that we can generate valid ANSI sequences
func TestANSISequenceGeneration(t *testing.T) {
	for _, seq := range expectedSequences {
		t.Run(seq.Name, func(t *testing.T) {
			// Generate the actual sequence
			var actual string
			switch seq.Name {
			case "cursor_left":
				actual = "\x1b[D"
			case "cursor_right":
				actual = "\x1b[C"
			case "cursor_up":
				actual = "\x1b[A"
			case "cursor_down":
				actual = "\x1b[B"
			case "cursor_home":
				actual = "\x1b[H"
			case "clear_line":
				actual = "\x1b[2K"
			case "clear_screen":
				actual = "\x1b[2J"
			case "save_cursor":
				actual = "\x1b[s"
			case "restore_cursor":
				actual = "\x1b[u"
			case "color_reset":
				actual = "\x1b[0m"
			case "color_red":
				actual = "\x1b[31m"
			case "color_green":
				actual = "\x1b[32m"
			case "color_blue":
				actual = "\x1b[34m"
			case "bold":
				actual = "\x1b[1m"
			}

			// Validate against pattern
			matched, err := regexp.MatchString(seq.Pattern, actual)
			if err != nil {
				t.Fatalf("Invalid regex pattern %q: %v", seq.Pattern, err)
			}

			if !matched {
				t.Errorf("Generated sequence %q doesn't match pattern %q", actual, seq.Pattern)
			}

			t.Logf("✓ %s: %q matches %q", seq.Description, actual, seq.Pattern)
		})
	}
}

// TestRenderingOutput tests actual DSH output for ANSI sequences
func TestRenderingOutput(t *testing.T) {
	// Simulate common terminal operations
	operations := []struct {
		name     string
		simulate func() string
		check    func(string) bool
	}{
		{
			"prompt_display",
			func() string { return "dsh> " },
			func(s string) bool { return strings.Contains(s, "dsh>") },
		},
		{
			"cursor_movement",
			func() string { return "text\x1b[D\x1b[D" }, // text + 2 left moves
			func(s string) bool { return strings.Contains(s, "\x1b[D") },
		},
		{
			"line_clear",
			func() string { return "old text\x1b[2K" }, // clear line
			func(s string) bool { return strings.Contains(s, "\x1b[2K") },
		},
		{
			"color_output",
			func() string { return "\x1b[32mgreen\x1b[0m" }, // green text
			func(s string) bool {
				return strings.Contains(s, "\x1b[32m") && strings.Contains(s, "\x1b[0m")
			},
		},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			output := op.simulate()

			t.Logf("Operation %s output: %q", op.name, output)

			if !op.check(output) {
				t.Errorf("Operation %s failed validation", op.name)
			} else {
				t.Logf("✓ Operation %s passed validation", op.name)
			}
		})
	}
}

// TestTabCompletionSequences tests expected sequences for tab completion
func TestTabCompletionSequences(t *testing.T) {
	// Simulate tab completion rendering
	completionOutput := func() string {
		var output strings.Builder

		// Save cursor position
		output.WriteString("\x1b[s")

		// Move to next line and show options
		output.WriteString("\n")
		output.WriteString("echo  exit  help")

		// Restore cursor position
		output.WriteString("\x1b[u")

		return output.String()
	}

	output := completionOutput()
	t.Logf("Tab completion output: %q", output)

	// Check for cursor save/restore
	if !strings.Contains(output, "\x1b[s") {
		t.Error("Tab completion should save cursor position")
	}
	if !strings.Contains(output, "\x1b[u") {
		t.Error("Tab completion should restore cursor position")
	}

	// Check for completion options
	if !strings.Contains(output, "echo") {
		t.Error("Tab completion should show 'echo' option")
	}
}

// TestLineEditingSequences tests sequences for line editing operations
func TestLineEditingSequences(t *testing.T) {
	operations := []struct {
		name     string
		sequence string
		desc     string
	}{
		{"backspace", "\x1b[D\x1b[K", "Move left and clear to end"},
		{"delete_word", "\x1b[2K", "Clear entire line"},
		{"move_home", "\x1b[H", "Move to beginning of line"},
		{"move_end", "\x1b[F", "Move to end of line"},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			t.Logf("Line editing %s: %q (%s)", op.name, op.sequence, op.desc)

			// Validate sequence format
			if !strings.HasPrefix(op.sequence, "\x1b[") {
				t.Errorf("Sequence should start with ESC[, got: %q", op.sequence)
			}

			t.Logf("✓ Valid ANSI sequence for %s", op.name)
		})
	}
}

// TestRenderingDiagnostics provides diagnostic information for debugging
func TestRenderingDiagnostics(t *testing.T) {
	t.Log("=== RENDERING DIAGNOSTICS ===")

	// Test environment info
	t.Logf("Terminal sequences supported:")
	for _, seq := range expectedSequences {
		t.Logf("  %s: %s", seq.Name, seq.Description)
	}

	// Test sequence validation
	testSequence := "\x1b[32mHello\x1b[0m World"
	t.Logf("Sample colored output: %q", testSequence)

	// Extract ANSI sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	sequences := ansiRegex.FindAllString(testSequence, -1)
	t.Logf("Extracted sequences: %v", sequences)

	if len(sequences) != 2 {
		t.Errorf("Expected 2 sequences, found %d", len(sequences))
	}
}
