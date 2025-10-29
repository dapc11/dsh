package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestBackspaceNoTrailingSpace(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace should not leave trailing spaces in terminal output",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("cat"),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Final rendered line should not have trailing space",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Full output: %q", output)

					// Parse the output to find the final rendered state
					// Split by \r to get individual line states
					lines := strings.Split(output, "\r")
					if len(lines) == 0 {
						return false
					}

					// Get the last line (final state)
					finalLine := lines[len(lines)-1]
					t.Logf("Final line: %q", finalLine)

					// Remove ANSI escape sequences to get the visible text
					visibleText := removeANSISequences(finalLine)
					t.Logf("Visible text: %q", visibleText)

					// The visible text should be "dsh> ca" without trailing space
					expected := "dsh> ca"
					if visibleText != expected {
						t.Logf("Expected %q, got %q", expected, visibleText)
						return false
					}

					return true
				},
				Message: "Final rendered line should be 'dsh> ca' without trailing space",
			},
			{
				Name: "Buffer state should be correct",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "ca"
				},
				Message: "Buffer should contain 'ca'",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// removeANSISequences removes ANSI escape sequences from a string
func removeANSISequences(s string) string {
	result := ""
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			i++ // skip the '['
			continue
		}

		if inEscape {
			// Skip until we find a letter (end of escape sequence)
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}

		result += string(s[i])
	}

	return result
}

func TestBackspaceMultipleChars(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Multiple backspaces should not accumulate trailing spaces",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello"),
			framework.Press(terminal.KeyBackspace),
			framework.Press(terminal.KeyBackspace),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should be 'he'",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after multiple backspaces: %q", buffer)
					return buffer == "he"
				},
				Message: "Buffer should be 'he' after three backspaces",
			},
			{
				Name: "No trailing spaces in final output",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					lines := strings.Split(output, "\r")
					if len(lines) == 0 {
						return false
					}

					finalLine := lines[len(lines)-1]
					visibleText := removeANSISequences(finalLine)
					t.Logf("Final visible text: %q", visibleText)

					expected := "dsh> he"
					return visibleText == expected
				},
				Message: "Final output should be 'dsh> he' without trailing spaces",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
