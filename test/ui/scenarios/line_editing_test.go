package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestLineEditingBasic(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Basic line editing operations",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should contain typed text",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer content: %q", buffer)
					return buffer == "hello world"
				},
				Message: "Typing should work correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlA(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+A moves cursor to beginning",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello "),
			framework.Press(terminal.KeyCtrlA),
			framework.Type("world "),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Text should be inserted at beginning",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after Ctrl+A and typing: %q", buffer)
					return buffer == "world hello "
				},
				Message: "Ctrl+A should move cursor to beginning",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlK(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+K kills to end of line",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
			framework.Press(terminal.KeyCtrlA),
			framework.Press(terminal.KeyCtrlK),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should be empty after kill",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Ctrl+K should kill entire line from beginning",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlU(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Delete operations",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
			framework.Press(terminal.KeyCtrlU), // Clear line
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should be empty after clear",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after Ctrl+U: %q", buffer)
					return buffer == ""
				},
				Message: "Ctrl+U should clear the line",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
