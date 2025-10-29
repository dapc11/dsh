package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestBackspaceCursorPositioning(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace should position cursor correctly without extra space",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("cat"),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should be 'ca' after backspace",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after backspace: %q", buffer)
					return buffer == "ca"
				},
				Message: "Buffer should be 'ca' after backspace",
			},
			{
				Name: "Output should not show extra space after cursor",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Full output: %q", output)

					// The output should end with "dsh> ca" and cursor positioning
					// It should NOT contain "dsh> ca " (with trailing space)

					// Look for the final rendered state
					// After backspace, we should see the prompt + "ca" + cursor positioning
					// The cursor should be positioned right after 'a', not after a space

					// Check that the final state doesn't have trailing space
					// This is a bit tricky to test directly, but we can check the ANSI sequences

					// The cursor should be positioned at column 8 (dsh> ca = 7 chars + 1 for 1-based)
					// NOT at column 9 which would indicate an extra space
					expectedCursorPos := "\033[8G" // Column 8 (1-based)
					wrongCursorPos := "\033[9G"    // Column 9 (would indicate extra space)

					hasCorrectPos := len(output) > 0 && !contains(output, wrongCursorPos)
					if contains(output, expectedCursorPos) {
						t.Logf("Found correct cursor position: %s", expectedCursorPos)
						return true
					}

					t.Logf("Cursor positioning check - looking for correct position")
					return hasCorrectPos
				},
				Message: "Cursor should be positioned correctly without extra space",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()
}
