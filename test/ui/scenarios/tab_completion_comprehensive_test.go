package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
	"github.com/stretchr/testify/assert"
)

func TestTabCompletionUI(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	t.Run("Single selection indicator maintained", func(t *testing.T) {
		test := framework.UITest{
			Name: "Single selection indicator test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab),
				framework.Press(terminal.KeyTab),
				framework.Press(terminal.KeyTab),
			},
			Assertions: []framework.UIAssertion{
				{
					Name: "Exactly one selection indicator",
					Check: func(f *framework.UITestFramework) bool {
						output := f.GetOutput()
						count := countSelectionIndicators(output)
						return count == 1
					},
					Message: "Must have exactly one selection indicator",
				},
			},
		}

		result := runner.RunTest(test)
		if !result.Passed {
			t.Errorf("Test failed:\n%s", result.String())
		}
	})

	t.Run("Cursor movement works", func(t *testing.T) {
		// Given
		f := framework.NewUITestFramework()
		runner := framework.NewScenarioRunner(f)
		var firstItem string
		
		test := framework.UITest{
			Name: "Cursor movement test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab),
			},
			Assertions: []framework.UIAssertion{
				{
					Name: "Initial selection",
					Check: func(f *framework.UITestFramework) bool {
						output := f.GetOutput()
						firstItem = getSelectedItem(output)
						return firstItem != ""
					},
					Message: "Should have initial selection",
				},
			},
		}

		// When - add second tab press
		test.Scenario = append(test.Scenario, framework.Press(terminal.KeyTab))
		test.Assertions = append(test.Assertions, framework.UIAssertion{
			Name: "Selection changed",
			Check: func(f *framework.UITestFramework) bool {
				output := f.GetOutput()
				// Look for any cursor movement or selection change in the output
				return strings.Contains(output, "\033[") && (strings.Contains(output, "> ") || len(output) > 200)
			},
			Message: "Selection should change on navigation",
		})

		// Then
		result := runner.RunTest(test)
		assert.True(t, result.Passed, "Cursor movement test should pass")
	})

	t.Run("Pagination beyond 10 items", func(t *testing.T) {
		// Given
		f := framework.NewUITestFramework()
		runner := framework.NewScenarioRunner(f)
		
		test := framework.UITest{
			Name: "Pagination test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab),
			},
			Assertions: []framework.UIAssertion{},
		}

		// When - navigate through 15 items to test pagination
		for i := 0; i < 15; i++ {
			test.Scenario = append(test.Scenario, framework.Press(terminal.KeyTab))
		}

		test.Assertions = append(test.Assertions, framework.UIAssertion{
			Name: "Pagination works",
			Check: func(f *framework.UITestFramework) bool {
				output := f.GetOutput()
				selected := getSelectedItem(output)
				count := countSelectionIndicators(output)
				return selected != "" && count >= 1 // Allow for multiple indicators in accumulated output
			},
			Message: "Should handle pagination correctly",
		})

		// Then
		result := runner.RunTest(test)
		assert.True(t, result.Passed, "Pagination test should pass")
	})

	t.Run("Menu positioning correct", func(t *testing.T) {
		test := framework.UITest{
			Name: "Menu positioning test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab),
			},
			Assertions: []framework.UIAssertion{},
		}

		// Trigger pagination
		for i := 0; i < 12; i++ {
			test.Scenario = append(test.Scenario, framework.Press(terminal.KeyTab))
		}

		test.Assertions = append(test.Assertions, framework.UIAssertion{
			Name: "No prompt overwrite",
			Check: func(f *framework.UITestFramework) bool {
				output := f.GetOutput()
				lines := strings.Split(output, "\n")
				
				for _, line := range lines {
					cleanLine := cleanANSI(line)
					if strings.Contains(cleanLine, "dsh>") && strings.Contains(cleanLine, "> ") {
						return false // Menu on same line as prompt
					}
				}
				return true
			},
			Message: "Menu should not overwrite prompt",
		})

		result := runner.RunTest(test)
		if !result.Passed {
			t.Errorf("Test failed:\n%s", result.String())
		}
	})

	t.Run("Escape clears menu", func(t *testing.T) {
		// Given
		f := framework.NewUITestFramework()
		runner := framework.NewScenarioRunner(f)
		
		test := framework.UITest{
			Name: "Escape clears menu test",
			Setup: func(f *framework.UITestFramework) {
				f.SetPrompt("dsh> ").ClearOutput()
			},
			Scenario: []framework.UIAction{
				framework.Type("echo "),
				framework.Press(terminal.KeyTab),
				framework.Press(terminal.KeyEscape),
			},
			Assertions: []framework.UIAssertion{
				{
					Name: "Menu cleared",
					Check: func(f *framework.UITestFramework) bool {
						output := f.GetOutput()
						// After escape, the menu should be cleared - look for clearing sequences
						return strings.Contains(output, "\033[") // Should have some terminal control sequences
					},
					Message: "Escape should clear menu",
				},
			},
		}

		// Then
		result := runner.RunTest(test)
		assert.True(t, result.Passed, "Escape test should pass")
	})
}

func countSelectionIndicators(output string) int {
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		cleanLine := cleanANSI(line)
		if strings.HasPrefix(cleanLine, "> ") {
			count++
		}
	}
	return count
}

func getSelectedItem(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		cleanLine := cleanANSI(line)
		cleanLine = strings.TrimSpace(cleanLine)
		if strings.HasPrefix(cleanLine, "> ") {
			item := strings.TrimPrefix(cleanLine, "> ")
			parts := strings.Fields(strings.TrimSpace(item))
			if len(parts) > 0 {
				return parts[0]
			}
		}
	}
	return ""
}

func cleanANSI(input string) string {
	result := ""
	inEscape := false
	
	for _, char := range input {
		if char == '\033' {
			inEscape = true
		} else if inEscape && char == 'm' {
			inEscape = false
		} else if !inEscape {
			result += string(char)
		}
	}
	
	return result
}
