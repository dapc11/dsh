package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

// TestBackspaceDoesNotInsertDEL ensures backspace removes characters instead of inserting DEL
func TestBackspaceDoesNotInsertDEL(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace must remove character, not insert DEL (\\x7f)",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("cat"),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer must not contain DEL character",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()

					// Critical: buffer must be "ca", not "cat\x7f"
					if buffer != "ca" {
						t.Errorf("REGRESSION: Expected 'ca', got %q", buffer)
						return false
					}

					// Ensure no DEL character exists
					for _, r := range buffer {
						if r == 127 {
							t.Errorf("REGRESSION: Buffer contains DEL character")
							return false
						}
					}
					return true
				},
				Message: "Backspace must actually remove characters",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Fatalf("CRITICAL REGRESSION: %s", result.String())
	}
}

// TestBackspaceSequentialRemoval tests multiple backspaces
func TestBackspaceSequentialRemoval(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Sequential backspaces must remove characters properly",
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
				Name: "Buffer must be 'hel' after two backspaces",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					if buffer != "hel" {
						t.Errorf("Expected 'hel', got %q", buffer)
						return false
					}

					// Ensure no control characters
					for _, r := range buffer {
						if r < 32 || r == 127 {
							t.Errorf("Buffer contains control character: %d", r)
							return false
						}
					}
					return true
				},
				Message: "Sequential backspaces must work correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed: %s", result.String())
	}
}

// TestBackspaceEmptyBuffer ensures backspace on empty buffer doesn't crash
func TestBackspaceEmptyBuffer(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Backspace on empty buffer must not crash",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyBackspace),
			framework.Press(terminal.KeyBackspace),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer must remain empty",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == ""
				},
				Message: "Empty buffer backspace must be safe",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed: %s", result.String())
	}
}
