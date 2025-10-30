package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
	"github.com/stretchr/testify/assert"
)

func TestTabCompletionRegression(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	t.Run("Single selection indicator throughout navigation", func(t *testing.T) {
		test := framework.UITest{
			Name: "Single selection indicator regression test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab), // Show menu
				framework.Press(terminal.KeyTab), // Navigate
				framework.Press(terminal.KeyTab), // Navigate
				framework.Press(terminal.KeyTab), // Navigate
			},
			Assertions: []framework.UIAssertion{
				{
					Name: "Should have exactly one selection indicator",
					Check: func(f *framework.UITestFramework) bool {
						output := f.GetOutput()
						count := countSelectionIndicators(output)
						if count != 1 {
							t.Logf("Expected 1 selection indicator, got %d", count)
							t.Logf("Output: %q", output)
						}
						return count == 1
					},
					Message: "Must maintain exactly one selection indicator",
				},
			},
		}

		result := runner.RunTest(test)
		if !result.Passed {
			t.Errorf("Test failed:\n%s", result.String())
		}
	})

	t.Run("Cursor movement through multiple items", func(t *testing.T) {
		// Given
		f := framework.NewUITestFramework()
		runner := framework.NewScenarioRunner(f)

		test := framework.UITest{
			Name: "Cursor movement regression test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab), // Show menu
			},
			Assertions: []framework.UIAssertion{
				{
					Name: "Should show initial selection",
					Check: func(f *framework.UITestFramework) bool {
						output := f.GetOutput()
						selected := getSelectedItem(output)
						return selected != ""
					},
					Message: "Should have initial selection",
				},
			},
		}

		// When - add navigation steps
		for i := 0; i < 5; i++ {
			test.Scenario = append(test.Scenario, framework.Press(terminal.KeyTab))
		}

		test.Assertions = append(test.Assertions, framework.UIAssertion{
			Name: "Navigation occurred",
			Check: func(f *framework.UITestFramework) bool {
				output := f.GetOutput()
				// Look for cursor movement sequences or multiple selection indicators
				return strings.Contains(output, "\033[") && len(output) > 300
			},
			Message: "Cursor should move through items",
		})

		// Then
		result := runner.RunTest(test)
		assert.True(t, result.Passed, "Cursor movement regression test should pass")
	})

	t.Run("Menu positioning after pagination", func(t *testing.T) {
		test := framework.UITest{
			Name: "Menu positioning regression test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab), // Show menu
			},
			Assertions: []framework.UIAssertion{},
		}

		// Navigate to trigger pagination
		for i := 0; i < 12; i++ {
			test.Scenario = append(test.Scenario, framework.Press(terminal.KeyTab))
		}

		test.Assertions = append(test.Assertions, framework.UIAssertion{
			Name: "Menu should not overwrite prompt after pagination",
			Check: func(f *framework.UITestFramework) bool {
				output := f.GetOutput()

				// Check that prompt and menu are properly separated
				lines := strings.Split(output, "\n")
				for _, line := range lines {
					cleanLine := cleanANSI(line)
					// Menu items should not appear on same line as prompt
					if strings.Contains(cleanLine, "dsh>") && strings.Contains(cleanLine, "> ") {
						t.Logf("Menu item on same line as prompt: %q", cleanLine)
						return false
					}
				}

				return true
			},
			Message: "Menu should not overwrite prompt line",
		})

		result := runner.RunTest(test)
		if !result.Passed {
			t.Errorf("Test failed:\n%s", result.String())
		}
	})
}
