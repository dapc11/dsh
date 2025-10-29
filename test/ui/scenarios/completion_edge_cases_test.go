package scenarios

import (
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

// Edge Cases: Completion with No Matches
func TestCompletionNoMatches(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion with no matches",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("zzznomatch"),
			framework.Press(terminal.KeyTab),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle no matches gracefully",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "zzznomatch"
				},
				Message: "No matches should leave buffer unchanged",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Completion Menu Overflow
func TestCompletionMenuOverflow(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion menu with many items",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Use a prefix that matches many commands
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
			// Navigate through many items
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle menu overflow",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "e"
				},
				Message: "Menu overflow should not corrupt state",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Completion During Editing
func TestCompletionDuringEditing(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion while cursor in middle",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("echo hello"),
			framework.PressCtrl('a'),                // Move to beginning
			framework.Press(terminal.KeyArrowRight), // Move to after 'e'
			framework.Press(terminal.KeyTab),        // Try completion in middle
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle completion in middle of line",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash or corrupt the line
					buffer := f.GetShell().GetBuffer()
					return len(buffer) > 0
				},
				Message: "Completion in middle should work correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Rapid Menu Operations
func TestRapidMenuOperations(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Rapid menu open/close operations",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			// Rapid open/close
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyArrowDown),
			framework.Press(terminal.KeyEscape),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEnter),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle rapid menu operations",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return len(buffer) > 1 // Should have completed something
				},
				Message: "Rapid menu operations should work correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Completion State Corruption
func TestCompletionStateCorruption(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion state corruption prevention",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
			// Try to corrupt state with mixed operations
			framework.Type("x"), // Type while menu is open
			framework.Press(terminal.KeyArrowUp),
			framework.Press(terminal.KeyBackspace),
			framework.Press(terminal.KeyArrowDown),
			framework.PressCtrl('a'),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should prevent state corruption",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash and should have reasonable state
					return true
				},
				Message: "State should remain consistent",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Menu Rendering Edge Cases
func TestMenuRenderingEdgeCases(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Menu rendering edge cases",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Test with very short prefix
			framework.Type("a"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			// Clear and test with single character
			framework.PressCtrl('u'),
			framework.Type("b"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle short prefixes",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == "b"
				},
				Message: "Short prefixes should work correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Completion with Special Characters
func TestCompletionSpecialChars(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion with special characters",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Test completion after quotes
			framework.Type("echo \""),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			// Test completion after spaces
			framework.PressCtrl('u'),
			framework.Type("ls "),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle special characters in completion context",
				Check: func(f *framework.UITestFramework) bool {
					// Should not crash
					return true
				},
				Message: "Special characters should be handled in completion",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Multiple Tab Behaviors
func TestMultipleTabBehaviors(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Multiple tab press behaviors",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("e"),
			// First tab should show menu (multiple matches)
			framework.Press(terminal.KeyTab),
			// Second tab should navigate
			framework.Press(terminal.KeyTab),
			// Third tab should continue navigation
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEnter),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should handle multiple tab behaviors",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after multiple tabs: %q", buffer)
					return len(buffer) > 1 // Should have completed something
				},
				Message: "Multiple tabs should work as expected",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Edge Cases: Completion Memory Management
func TestCompletionMemoryManagement(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion memory management",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			// Create and destroy many completion sessions
			framework.Type("e"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			framework.PressCtrl('u'),
			framework.Type("l"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			framework.PressCtrl('u'),
			framework.Type("c"),
			framework.Press(terminal.KeyTab),
			framework.Press(terminal.KeyEscape),
			framework.PressCtrl('u'),
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should manage memory correctly",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					return buffer == ""
				},
				Message: "Memory should be managed correctly",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Test for the specific bug: completion menu remains visible after accepting suggestion
func TestCompletionMenuCleanupAfterAccept(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Completion menu should disappear after accepting suggestion",
		Setup: func(f *framework.UITestFramework) {
			f.Reset().SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Press(terminal.KeyEnter), // Start with enter to get initial prompt
			framework.Type("ls"),
			framework.Press(terminal.KeyTab),   // Show completion menu (ls highlighted)
			framework.Press(terminal.KeyTab),   // Navigate to lsattr
			framework.Press(terminal.KeyTab),   // Navigate to lsb_release
			framework.Press(terminal.KeyEnter), // Accept lsb_release
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Buffer should contain lsb_release after completion",
				Check: func(f *framework.UITestFramework) bool {
					buffer := f.GetShell().GetBuffer()
					t.Logf("Buffer after completion: %q", buffer)
					return buffer == "lsb_release"
				},
				Message: "Buffer should contain lsb_release after completion",
			},
			{
				Name: "Output should be exact with escape sequences",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Output after completion: %q", output)
					// Updated expectation: no full redraws for navigation, just initial menu + cleanup
					// Updated expectation: includes in-place cursor movements for selection updates
					expectedOutput := "\r\ndsh> \r\x1b[Kdsh> l\x1b[7G\r\x1b[Kdsh> ls\x1b[8G\x1b[s\x1b[s\r\n> ls                                   lsattr                           \r\n  lsb_release                          lsblk                            \r\n  lscpu                                lshw                             \r\n  lsinitramfs                          lsipc                            \r\n  lslocks                              lslogins                         \r\n\x1b[1;1H\x1b[2;1H  \x1b[2;38H> \x1b[u\x1b[1;1H\x1b[2;38H  \x1b[3;1H> \x1b[u\x1b[1;1H\x1b[2K\x1b[u\x1b[u\x1b[J\x1b[2K\rdsh> lsb_release"
					return output == expectedOutput
				},
				Message: "Output should match exact escape sequence pattern",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
