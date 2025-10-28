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

func TestLineEditingCtrlE(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+E moves cursor to end",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
			framework.Press(terminal.KeyCtrlA), // Move to beginning
			framework.Press(terminal.KeyCtrlE), // Move to end
			framework.Type("!"),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Text should be appended at end",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("hello world!")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Ctrl+E should move cursor to end of line",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlY(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+Y yanks killed text",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
			framework.Press(terminal.KeyCtrlA),
			framework.Press(terminal.KeyCtrlK), // Kill entire line
			framework.Press(terminal.KeyCtrlY), // Yank it back
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Text should be yanked back",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("hello world")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Ctrl+Y should yank killed text back",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingArrowKeys(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Arrow keys move cursor",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello"),
			framework.Press(terminal.KeyArrowLeft),
			framework.Press(terminal.KeyArrowLeft),
			framework.Type(" world"),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Text should be inserted at cursor position",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("hel worldlo")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Arrow keys should move cursor for text insertion",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlW(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+W kills word backward",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world test"),
			framework.Press(terminal.KeyCtrlW), // Kill "test"
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Last word should be killed",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("hello world ")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Ctrl+W should kill word backward",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlD(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+D deletes character at cursor",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello"),
			framework.Press(terminal.KeyArrowLeft),
			framework.Press(terminal.KeyArrowLeft),
			framework.Press(terminal.KeyCtrlD), // Delete 'l'
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Character at cursor should be deleted",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("helo")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Ctrl+D should delete character at cursor",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestLineEditingCtrlBF(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+B/F move cursor left/right",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello"),
			framework.Press(terminal.KeyCtrlB), // Move left
			framework.Press(terminal.KeyCtrlB), // Move left
			framework.Type("X"),
			framework.Press(terminal.KeyCtrlF), // Move right
			framework.Type("Y"),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Text should be inserted at correct positions",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertBuffer().Equals("helXlYo")
					if !result.Passed {
						t.Error(result.Message)
					}
					return result.Passed
				},
				Message: "Ctrl+B/F should move cursor left/right",
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
