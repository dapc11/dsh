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
				Name: "Should show initial selection but not redraw for navigation",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					
					// Look for selection indicators - should only see initial selection
					firstSelection := strings.Count(output, "> backspace_cursor_test.go")
					secondSelection := strings.Count(output, "> backspace_delete_test.go")
					t.Logf("First item selected: %d times", firstSelection)
					t.Logf("Second item selected: %d times", secondSelection)
					
					// With no redraws for navigation, only initial selection is visible
					if firstSelection != 1 || secondSelection != 0 {
						t.Errorf("Expected only initial selection visible, got first=%d, second=%d", firstSelection, secondSelection)
						return false
					}
					
					return true
				},
				Message: "Should show only initial selection without navigation redraws",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed: %s", result.String())
	}
}
