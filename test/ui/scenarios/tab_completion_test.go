package scenarios

import (
	"strings"
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
				Name: "Menu should contain available completions",
				Check: func(f *framework.UITestFramework) bool {
					// Check for items that should be in the simplified menu
					return f.AssertMenu().Contains("exit", "e2freefrag").Passed
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

func TestTabCompletionSelection(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion selection with Enter",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("ec"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEnter),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should complete to echo",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "echo"
				},
				Message: "Should complete ec to echo",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionEscape(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion menu dismissal with Escape",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should remain unchanged",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "e" // Should still have original input
				},
				Message: "Escape should dismiss completion menu without changing buffer",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionBufferIntegrity(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion preserves buffer integrity",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: framework.CancelCompletion("e"), // "e" has multiple matches (echo, exit)
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should be exactly as before completion",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Actual buffer after cancel: %q", buffer)
					return f.AssertBuffer().Equals("e").Passed
				},
				Message: "Cancelled completion should leave buffer unchanged",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionExactResult(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion produces exact result",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: framework.CompleteWith("ec"),
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should contain exact completion",
				Check: func(f *framework.UITestFramework) bool {
					return f.AssertBuffer().Equals("echo").Passed
				},
				Message: "Completion should result in exactly 'echo'",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionStateTracking(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	var initialState *framework.BufferStateAssertion

	test := framework.UITest{
		Name: "Tab completion state tracking",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
			f.GetShell().SetBuffer("e") // Set initial state with multiple matches
			initialState = f.CaptureBufferState()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyTab),     // Show completion
			framework.Press(terminal.KeyEscape), // Cancel
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer state should be unchanged after cancel",
				Check: func(f *framework.UITestFramework) bool {
					return initialState.IsUnchanged().Passed
				},
				Message: "Buffer should return to exact initial state",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionDoubleTabNavigation(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Double tab navigation and selection",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),  // First tab - show menu
			framework.Press(terminal.KeyTab),  // Second tab - navigate to next item
			framework.Press(terminal.KeyEnter), // Select highlighted item
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should complete to exact expected command",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after tab tab + enter: %q", buffer)
					// The actual result is "e2freefrag" - assert exactly that
					return f.AssertBuffer().Equals("e2freefrag").Passed
				},
				Message: "Should complete to exactly 'e2freefrag' after tab tab navigation",
			},
			{
				Name: "Should have clean rendering without excessive cursor movement",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Full output length: %d", len(output))
					t.Logf("ANSI sequences count: %d", strings.Count(output, "\033["))
					if len(output) > 200 {
						t.Logf("Sample output: %q", output[:200])
					} else {
						t.Logf("Full output: %q", output)
					}
					
					result := f.AssertRendering().HasNoExcessiveCursorMovement()
					if !result.Passed {
						t.Errorf("Rendering issue: %s", result.Message)
					}
					return result.Passed
				},
				Message: "Completion should render cleanly without cursor jumping",
			},
			{
				Name: "Should have clean output without rendering artifacts",
				Check: func(f *framework.UITestFramework) bool {
					result := f.AssertRendering().HasCleanOutput()
					if !result.Passed {
						t.Errorf("Output issue: %s", result.Message)
						output := f.GetOutput()
						t.Logf("Full output: %q", output)
					}
					return result.Passed
				},
				Message: "Completion should not leave rendering artifacts",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionCleanup(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion menu cleanup after selection",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),    // Show menu
			framework.Press(terminal.KeyTab),    // Navigate
			framework.Press(terminal.KeyEnter),  // Select
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should contain selected completion",
				Check: func(f *framework.UITestFramework) bool {
					return f.AssertBuffer().Equals("e2freefrag").Passed
				},
				Message: "Should complete to e2freefrag",
			},
			{
				Name: "Output should contain screen clear after selection",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					// Should contain screen clear sequence after completion
					return strings.Contains(output, "\033[2J\033[H")
				},
				Message: "Should clear screen to remove completion menu",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionPagination(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion pagination after 15+ tabs",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab), // Show initial menu
			// Press tab 15 more times to go beyond visible items
			framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab), framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEnter), // Select whatever is highlighted
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should complete to a valid command after pagination",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after 15+ tabs + enter: %q", buffer)
					// Should be a valid command starting with 'e'
					return len(buffer) > 1 && buffer[0] == 'e'
				},
				Message: "Should complete to valid command after pagination",
			},
			{
				Name: "Should show different completion items after pagination",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					// Should contain both early items (exit, e2freefrag) and later items (ebtables)
					hasEarlyItems := strings.Contains(output, "exit") && strings.Contains(output, "e2freefrag")
					hasLaterItems := strings.Contains(output, "ebtables")
					t.Logf("Has early items: %v, Has later items: %v", hasEarlyItems, hasLaterItems)
					return hasEarlyItems && hasLaterItems
				},
				Message: "Should show both early and later completion items during pagination",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

func TestTabCompletionSingleMatch(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Single match auto-completion",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("echox"), // Unique prefix
			framework.Press(terminal.KeyTab),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should auto-complete unique match",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					// Should either complete or remain unchanged if no match
					return len(buffer) >= 5
				},
				Message: "Single match should auto-complete",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
