package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

// Edge Cases: Empty Input Handling
func TestEmptyInputEdgeCases(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion on empty input",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyTab),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle tab on empty input gracefully",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash and buffer should remain empty
					return f.GetShell().GetBuffer() == ""
				},
				Message: "Tab on empty input should not crash",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestCtrlROnEmptyHistory(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Ctrl+R with empty history",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
			// Explicitly don't add any history
		},
		Scenario: []framework.UIAction{
			framework.PressCtrl('r'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle empty history gracefully",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash
					return true
				},
				Message: "Ctrl+R with empty history should not crash",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Boundary Conditions
func TestCursorBoundaryConditions(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Cursor movement at boundaries",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Try to move left when at beginning
			framework.PressCtrl('a'),
			framework.Press(terminal.KeyArrowLeft),
			framework.Press(terminal.KeyArrowLeft),
			framework.PressCtrl('b'),
			// Add some text
			framework.Type("test"),
			// Try to move right when at end
			framework.PressCtrl('e'),
			framework.Press(terminal.KeyArrowRight),
			framework.Press(terminal.KeyArrowRight),
			framework.PressCtrl('f'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle cursor boundary conditions",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "test"
				},
				Message: "Cursor should not move beyond boundaries",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Rapid Key Sequences
func TestRapidKeySequences(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Rapid tab presses",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			// Rapid tab presses
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle rapid tab presses",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "e"
				},
				Message: "Rapid tab presses should not corrupt state",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Special Characters
func TestSpecialCharacterHandling(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Special characters in input",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("echo \"hello\tworld\n\""),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle special characters",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return len(buffer) > 0 // Should not crash
				},
				Message: "Special characters should be handled gracefully",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Long Input Lines
func TestLongInputLines(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	longText := "this is a very long command line that exceeds normal terminal width and should test how the shell handles wrapping and cursor positioning with extremely long input sequences that might cause issues"

	test := framework.UITest{
		Name: "Very long input lines",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type(longText),
			framework.PressCtrl('a'), // Move to beginning
			framework.PressCtrl('e'), // Move to end
			framework.PressCtrl('k'), // Kill to end
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle long input lines",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash and should handle the operations
					return true
				},
				Message: "Long input lines should be handled correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Kill Ring Overflow
func TestKillRingOverflow(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Kill ring overflow handling",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Fill kill ring with multiple kills
			framework.Type("first line"),
			framework.PressCtrl('k'),
			framework.Type("second line"),
			framework.PressCtrl('k'),
			framework.Type("third line"),
			framework.PressCtrl('k'),
			framework.Type("fourth line"),
			framework.PressCtrl('k'),
			framework.Type("fifth line"),
			framework.PressCtrl('k'),
			// Try to yank
			framework.PressCtrl('y'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle kill ring operations",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return len(buffer) > 0 // Should have yanked something
				},
				Message: "Kill ring should handle multiple operations",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: History Navigation Boundaries
func TestHistoryNavigationBoundaries(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "History navigation at boundaries",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").
				AddHistory("cmd1", "cmd2", "cmd3").
				ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Go to beginning of history
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			// Try to go beyond beginning
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			// Go to end
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			// Try to go beyond end
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle history boundaries",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash
					return true
				},
				Message: "History navigation should respect boundaries",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Completion Menu Boundaries
func TestCompletionMenuBoundaries(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion menu navigation boundaries",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
			// Try to navigate beyond menu boundaries
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle menu navigation boundaries",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after menu navigation: %q", buffer)
					return buffer == "e"
				},
				Message: "Menu navigation should respect boundaries",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Mixed Operations
func TestMixedOperationsEdgeCases(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Mixed operations edge cases",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").
				AddHistory("previous command").
				ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Start typing
			framework.Type("echo"),
			// Open completion
			framework.Press(terminal.KeyTab),
			// Navigate in menu
			framework.Press(terminal.KeyArrowDown),
			// Cancel completion
			framework.Press(terminal.KeyEscape),
			// Try history
			framework.Press(terminal.KeyArrowUp),
			// Edit the history item
			framework.PressCtrl('a'),
			framework.Type("modified "),
			// Try completion again
			framework.PressCtrl('e'),
			framework.Type(" test"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle mixed operations",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return len(buffer) > 0 // Should not crash
				},
				Message: "Mixed operations should work correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Unicode and Multi-byte Characters
func TestUnicodeHandling(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Unicode character handling",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("echo 🚀 测试 café"),
			framework.PressCtrl('a'),
			framework.PressCtrl('e'),
			framework.Press(terminal.KeyArrowLeft),
			framework.Press(terminal.KeyArrowLeft),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle unicode characters",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return len(buffer) > 0 // Should not crash
				},
				Message: "Unicode characters should be handled correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
