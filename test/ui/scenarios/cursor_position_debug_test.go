package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestCursorPositionCalculation(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Debug cursor position calculation after backspace",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("cat"),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Analyze cursor positioning",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					buffer := f.GetShell().GetBuffer()

					t.Logf("Buffer: %q (length: %d)", buffer, len(buffer))
					t.Logf("Prompt: %q (length: %d)", "dsh> ", len("dsh> "))
					t.Logf("Full output: %q", output)

					// Expected calculation:
					// promptLen = len("dsh> ") = 5
					// cursor = 2 (after "ca")
					// totalPos = 5 + 2 + 1 = 8
					// So cursor should be at column 8 (1-based)

					promptLen := len("dsh> ")
					cursorPos := len(buffer) // cursor should be at end after backspace
					expectedPos := promptLen + cursorPos + 1

					t.Logf("Prompt length: %d", promptLen)
					t.Logf("Cursor position: %d", cursorPos)
					t.Logf("Expected ANSI position: %d", expectedPos)

					// Look for the cursor positioning command in output
					expectedANSI := "\033[" + string(rune('0'+expectedPos)) + "G"
					if expectedPos >= 10 {
						// Handle two-digit positions
						expectedANSI = "\033[" + string(rune('0'+expectedPos/10)) + string(rune('0'+expectedPos%10)) + "G"
					}

					t.Logf("Expected ANSI sequence: %q", expectedANSI)

					// Check if this sequence appears in the output
					hasExpectedPos := strings.Contains(output, expectedANSI)
					t.Logf("Found expected position: %v", hasExpectedPos)

					// Also check what actual position sequences are in the output
					ansiPositions := extractANSIPositions(output)
					t.Logf("All ANSI position sequences found: %v", ansiPositions)

					return true // Always pass, this is just for debugging
				},
				Message: "Debug cursor positioning",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// extractANSIPositions finds all ANSI cursor positioning sequences in a string
func extractANSIPositions(s string) []string {
	var positions []string

	for i := 0; i < len(s)-2; i++ {
		if s[i] == '\033' && s[i+1] == '[' {
			// Found start of ANSI sequence
			start := i
			j := i + 2

			// Find the end (letter that terminates the sequence)
			for j < len(s) && !((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z')) {
				j++
			}

			if j < len(s) && s[j] == 'G' {
				// This is a cursor position sequence
				positions = append(positions, s[start:j+1])
			}

			i = j // Skip past this sequence
		}
	}

	return positions
}

func TestManualBackspaceSequence(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Step by step backspace analysis",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("c"),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "After typing 'c'",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					buffer := f.GetShell().GetBuffer()

					t.Logf("After 'c' - Buffer: %q, Output: %q", buffer, output)
					positions := extractANSIPositions(output)
					t.Logf("Cursor positions after 'c': %v", positions)

					return true
				},
				Message: "Debug after 'c'",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}

	// Now test adding 'a'
	test2 := framework.UITest{
		Name: "After adding 'a'",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("ca"),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "After typing 'ca'",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					buffer := f.GetShell().GetBuffer()

					t.Logf("After 'ca' - Buffer: %q, Output: %q", buffer, output)
					positions := extractANSIPositions(output)
					t.Logf("Cursor positions after 'ca': %v", positions)

					return true
				},
				Message: "Debug after 'ca'",
			},
		},
	}

	result2 := runner.RunTest(test2)
	if !result2.Passed {
		t.Errorf("Test2 failed:\n%s", result2.String())
	}
}
