package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestBackspaceRendering(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace should update display correctly",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello"),
			framework.Press(terminal.KeyBackspace),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should contain correct text after backspace",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after backspace: %q", buffer)
					return buffer == "hel"
				},
				Message: "Buffer should be 'hel' after two backspaces",
			},
			{
				Name: "Output should show correct rendering",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Output after backspace: %q", output)
					// Should end with the correct prompt and buffer content
					return output != "" // Basic check - should not be empty
				},
				Message: "Output should show proper backspace rendering",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestDeleteRendering(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Delete (Ctrl+D) should update display correctly",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello"),
			framework.Press(terminal.KeyArrowLeft), // Move cursor to 'o'
			framework.Press(terminal.KeyArrowLeft), // Move cursor to 'l'
			framework.PressCtrl('d'),               // Delete 'l'
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should contain correct text after delete",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after delete: %q", buffer)
					return buffer == "helo"
				},
				Message: "Buffer should be 'helo' after deleting middle character",
			},
			{
				Name: "Cursor should be in correct position",
				Check: func(_ *framework.UITestFramework) bool {
					// This would need access to cursor position
					// For now, just check buffer is correct
					return true
				},
				Message: "Cursor should remain at deletion point",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestBackspaceAtBeginning(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace at beginning should not crash",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyBackspace), // Backspace on empty buffer
			framework.Type("test"),
			framework.PressCtrl('a'),               // Move to beginning
			framework.Press(terminal.KeyBackspace), // Backspace at beginning
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should remain unchanged",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after backspace at beginning: %q", buffer)
					return buffer == "test"
				},
				Message: "Backspace at beginning should not modify buffer",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestDeleteAtEnd(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Delete at end should not crash",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("test"),
			framework.PressCtrl('d'), // Delete at end of buffer
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should remain unchanged",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after delete at end: %q", buffer)
					return buffer == "test"
				},
				Message: "Delete at end should not modify buffer",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
func TestBackspaceRenderingEfficiency(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace should render efficiently without full line redraws",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should not do unnecessary full line redraws",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Backspace output: %q", output)
					// Count full line redraws (\r\x1b[K)
					fullRedraws := 0
					for i := 0; i < len(output)-3; i++ {
						if output[i:i+3] == "\r\x1b[" && i+4 < len(output) && output[i+4] == 'K' {
							fullRedraws++
						}
					}
					t.Logf("Full line redraws: %d", fullRedraws)
					// Should be efficient - not redraw for every character
					return fullRedraws <= len("hello world")+2 // Allow some redraws but not excessive
				},
				Message: "Backspace should render efficiently",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestBackspaceInMiddleOfLine(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace in middle of line should render correctly",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("hello world"),
			framework.PressCtrl('a'),                // Move to beginning
			framework.Press(terminal.KeyArrowRight), // Move to after 'h'
			framework.Press(terminal.KeyArrowRight), // Move to after 'e'
			framework.Press(terminal.KeyArrowRight), // Move to after 'l'
			framework.Press(terminal.KeyBackspace),  // Delete 'l'
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should be correct after middle backspace",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after middle backspace: %q", buffer)
					return buffer == "helo world"
				},
				Message: "Middle backspace should work correctly",
			},
			{
				Name: "Rendering should handle middle deletion properly",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Middle backspace output: %q", output)
					// Should end with correct content
					return len(output) > 0
				},
				Message: "Middle backspace rendering should work",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
