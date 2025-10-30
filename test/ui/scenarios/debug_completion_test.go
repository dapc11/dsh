package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestDebugCompletion(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Debug completion items",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("echo "),
			framework.Press(terminal.KeyTab), // First tab - show menu
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "Should show completion menu",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					t.Logf("Full output: %q", output)
					
					// Debug line by line
					lines := strings.Split(output, "\n")
					for i, line := range lines {
						cleanLine := cleanANSI(line)
						t.Logf("Line %d: %q -> clean: %q", i, line, cleanLine)
						
						if strings.HasPrefix(cleanLine, "> ") || strings.HasPrefix(cleanLine, "  ") {
							t.Logf("  -> This is a completion line")
							content := strings.TrimPrefix(cleanLine, "> ")
							content = strings.TrimPrefix(content, "  ")
							content = strings.TrimSpace(content)
							words := strings.Fields(content)
							t.Logf("  -> Words found: %v", words)
						}
					}
					
					// Extract all visible items from the output
					items := extractAllItems(output)
					t.Logf("Found %d completion items: %v", len(items), items)
					
					return len(items) > 0
				},
				Message: "Should find completion items",
			},
		},
	}

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}
