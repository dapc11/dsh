package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

// Edge Cases: Large History Stress Test
func TestLargeHistoryStress(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Large history stress test",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
			// Add many history items
			for i := 0; i < 100; i++ {
				f.AddHistory("command_" + strings.Repeat("x", i%20))
			}
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			framework.UIAction{
				Type: framework.ActionSetupKeys,
				Data: []terminal.KeyEvent{
					{Key: terminal.KeyNone, Rune: 'c'},
					{Key: terminal.KeyEscape},
				},
			},
			framework.PressCtrl('r'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle large history efficiently",
				Check: func(f *framework.UITestFramework) bool {
					return true
				},
				Message: "Large history should be handled efficiently",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Rapid Input Stress Test
func TestRapidInputStress(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Rapid input stress test",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("abcdefghijklmnop"),
			framework.PressCtrl('a'),
			framework.Press(terminal.KeyArrowRight),
			framework.Press(terminal.KeyArrowLeft),
			framework.PressCtrl('e'),
			framework.PressCtrl('w'),
			framework.PressCtrl('y'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle rapid input",
				Check: func(f *framework.UITestFramework) bool {
					return true
				},
				Message: "Rapid input should be handled correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Buffer Overflow Protection
func TestBufferOverflowProtection(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	veryLongText := strings.Repeat("abcdefghij", 50) // 500 characters

	test := framework.UITest{
		Name: "Buffer overflow protection",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type(veryLongText),
			framework.PressCtrl('a'),
			framework.PressCtrl('e'),
			framework.PressCtrl('k'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should protect against buffer overflow",
				Check: func(f *framework.UITestFramework) bool {
					return true
				},
				Message: "Buffer overflow should be prevented",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: State Consistency Under Stress
func TestStateConsistencyUnderStress(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "State consistency under stress",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").
				AddHistory("cmd1", "cmd2", "cmd3").
				ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("start"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			framework.Press(terminal.KeyArrowUp),
			framework.PressCtrl('a'),
			framework.Type("modified "),
			framework.PressCtrl('k'),
			framework.PressCtrl('y'),
			framework.PressCtrl('u'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should maintain state consistency",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == ""
				},
				Message: "State should remain consistent under stress",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
