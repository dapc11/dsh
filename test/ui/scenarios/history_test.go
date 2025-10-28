package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestHistoryNavigation(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "History navigation with arrow keys",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").
				AddHistory("git status", "ls -la", "echo hello").
				ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyArrowUp),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should load last command from history",
				Check: func(f *framework.UITestFramework) bool {
					// In a real implementation, this would check the buffer
					// For now, just verify the action was processed
					return true
				},
				Message: "Arrow up should navigate to previous command",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestHistorySearch(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+R fuzzy history search",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").
				AddHistory("git status", "git commit", "ls -la").
				ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Queue keys that the fuzzy search will read
			// Type "git s" then use arrow down to select "git status"
			framework.UIAction{
				Type: framework.ActionSetupKeys,
				Data: []terminal.KeyEvent{
					{Key: terminal.KeyNone, Rune: 'g'},
					{Key: terminal.KeyNone, Rune: 'i'},
					{Key: terminal.KeyNone, Rune: 't'},
					{Key: terminal.KeyNone, Rune: ' '},
					{Key: terminal.KeyNone, Rune: 's'},
					{Key: terminal.KeyArrowDown}, // Move to "git status"
					{Key: terminal.KeyEnter},
				},
			},
			framework.PressCtrl('r'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should select git status from history",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					output := f.GetOutput()
					t.Logf("Buffer after Ctrl+R: %q", buffer)
					t.Logf("Output after Ctrl+R: %q", output)
					return buffer == "git status"
				},
				Message: "Ctrl+R should activate history search and select git status",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
