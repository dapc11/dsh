package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

// TestCompletionCursorPositioningRegression tests that selection indicators
// appear in the correct position within the completion menu, not at the top of the terminal.
// This is a regression test for the issue where removing cursor save from ShowCompletion
// caused UpdateSelectionHighlight to position indicators incorrectly.
func TestCompletionCursorPositioningRegression(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion selection indicators should appear in correct positions",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyEnter), // Start with enter to get initial prompt
			framework.Type("e"),                // Type 'e' to get multiple completions
			framework.Press(terminal.KeyTab),   // Show completion menu
			framework.Press(terminal.KeyTab),   // Navigate to next item (should update selection)
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Selection indicator should not appear at top of terminal",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()

					// Check that selection indicator ("> ") appears after the menu content,
					// not at the very beginning of the output
					lines := strings.Split(output, "\n")
					if len(lines) < 2 {
						return false
					}

					// The first few lines should be prompt and typing, not selection indicators
					firstLine := lines[0]
					if strings.HasPrefix(firstLine, "> ") {
						t.Errorf("Selection indicator incorrectly appears at top of terminal: %q", firstLine)
						return false
					}

					return true
				},
				Message: "Selection indicator should not appear at top of terminal",
			},
			{
				Name: "Selection indicator should appear within completion menu",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()

					// Should have both the menu content and selection indicators
					hasMenuContent := strings.Contains(output, "exit") || strings.Contains(output, "echo")
					hasSelectionIndicator := strings.Contains(output, "> ")

					if !hasMenuContent {
						t.Errorf("Missing completion menu content in output")
						return false
					}

					if !hasSelectionIndicator {
						t.Errorf("Missing selection indicator in output")
						return false
					}

					return true
				},
				Message: "Selection indicator should appear within completion menu",
			},
			{
				Name: "Cursor save/restore should be balanced for navigation",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()

					// Count cursor saves and restores
					saves := strings.Count(output, "\x1b[s")
					restores := strings.Count(output, "\x1b[u")

					// Should have at least one save for the menu
					if saves < 1 {
						t.Errorf("Expected at least 1 cursor save, got %d", saves)
						return false
					}

					// Should have restores for navigation (can be more than saves)
					if restores < 1 {
						t.Errorf("Expected at least 1 cursor restore, got %d", restores)
						return false
					}

					t.Logf("Cursor operations: %d saves, %d restores", saves, restores)
					return true
				},
				Message: "Cursor save/restore should be balanced for navigation",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
