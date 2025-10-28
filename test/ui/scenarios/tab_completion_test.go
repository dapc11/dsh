package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestTabCompletionBasic(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Basic tab completion shows menu",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Menu should be visible",
				Check: func(f *framework.UITestFramework) bool {
					return f.AssertMenu().IsVisible().Passed
				},
				Message: "Tab completion menu should appear",
			},
			{
				Name: "Menu should contain echo and exit",
				Check: func(f *framework.UITestFramework) bool {
					return f.AssertMenu().Contains("echo", "exit").Passed
				},
				Message: "Menu should show available completions",
			},
		},
	}

	result := runner.RunTest(test)

	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	} else {
		t.Logf("Test passed:\n%s", result.String())
	}
}

func TestTabCompletionNavigation(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion menu navigation",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyArrowDown),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Menu should be visible after navigation",
				Check: func(f *framework.UITestFramework) bool {
					return f.AssertMenu().IsVisible().Passed
				},
				Message: "Menu should remain visible during navigation",
			},
		},
	}

	result := runner.RunTest(test)

	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	} else {
		t.Logf("Test passed:\n%s", result.String())
	}
}
