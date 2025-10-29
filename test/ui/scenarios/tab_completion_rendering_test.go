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
		Name: "Tab completion should not render menu multiple times (CURRENTLY FAILING)",
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
					
					// Count how many times we see the menu rendering pattern
					// Look for the pattern where files are listed multiple times
					testFileCount := strings.Count(output, "backspace_cursor_test.go")
					t.Logf("backspace_cursor_test.go appears %d times", testFileCount)
					
					// Count menu clear/redraw sequences
					menuClearCount := strings.Count(output, "\x1b[1;1H\x1b[u\x1b[J")
					t.Logf("Menu clear sequences: %d", menuClearCount)
					
					// EXPECTED TO FAIL: Currently renders menu multiple times
					// Each tab should update the menu in place, not render it again
					// Currently: file appears 2 times (should be 1)
					if testFileCount > 1 {
						t.Errorf("REGRESSION: File appears %d times - should appear only once with in-place updates", testFileCount)
						return false
					}
					
					// Should not have menu clear sequences for navigation - should update in place
					// Currently: 1 clear sequence (should be 0 for navigation)
					if menuClearCount > 0 {
						t.Errorf("REGRESSION: %d menu clear sequences - should update selection in place", menuClearCount)
						return false
					}
					
					return true
				},
				Message: "Tab completion menu is being rendered excessively",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("EXPECTED FAILURE - documents excessive rendering bug: %s", result.String())
	}
}
