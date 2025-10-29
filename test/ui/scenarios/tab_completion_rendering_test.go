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
					
					// Look for in-place cursor positioning updates
					// Should see cursor moves to update selection indicators
					cursorMoves := strings.Count(output, "\x1b[2;")
					t.Logf("Cursor positioning sequences: %d", cursorMoves)
					
					// Should have cursor moves for updating selection (old + new position)
					if cursorMoves < 2 {
						t.Errorf("Expected at least 2 cursor moves for selection update, got %d", cursorMoves)
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
