package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestTabCompletionExcessiveRendering(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion should update highlight without excessive rendering",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("ls "), // Space after ls to trigger file completion
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should not render completion menu multiple times",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Full output: %q", output)

					// Count how many times we see files
					testFileCount := strings.Count(output, "backspace_cursor_test.go")
					t.Logf("backspace_cursor_test.go appears %d times", testFileCount)

					// With proper behavior: file should appear once per menu state
					// First tab: shows menu (1 time)
					// Second tab: redraws menu with new selection (1 more time)
					// Total: 2 times is acceptable for navigation
					// But more than 3 would indicate excessive rendering
					if testFileCount > 3 {
						t.Errorf("File appears %d times - excessive rendering", testFileCount)
						return false
					}

					return true
				},
				Message: "Should not render excessively",
			},
			{
				Name: "Should show in-place selection updates",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()

					// Look for cursor restore sequences (new positioning method)
					cursorRestores := strings.Count(output, "\x1b[u")
					downMoves := strings.Count(output, "\x1b[1B") + strings.Count(output, "\x1b[2B")
					t.Logf("Up movements: %d, Down movements: %d", 0, downMoves)

					// Should have cursor restores for updating selection
					if cursorRestores < 2 {
						t.Errorf("Expected at least 2 cursor restores for selection update, got %d", cursorRestores)
						return false
					}

					return true
				},
				Message: "Should show in-place selection updates",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed: %s", result.String())
	}
}
