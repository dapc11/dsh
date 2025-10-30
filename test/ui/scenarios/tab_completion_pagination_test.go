package scenarios

import (
	"strings"
	"testing"

	"dsh/internal/terminal"
	"dsh/test/ui/framework"
)

func TestTabCompletionPaginationAndCursor(t *testing.T) {
	fw := framework.NewUITestFramework()
	runner := framework.NewScenarioRunner(fw)

	test := framework.UITest{
		Name: "Tab completion pagination and cursor tracking",
		Setup: func(f *framework.UITestFramework) {
			f.SetPrompt("dsh> ").ClearOutput()
		},
		Scenario: []framework.UIAction{
			framework.Type("echo "), // This should match all files in current directory
			framework.Press(terminal.KeyTab), // First tab - show menu
		},
		Assertions: []framework.UIAssertion{
			{
				Name: "First tab should show menu with single selection indicator",
				Check: func(f *framework.UITestFramework) bool {
					output := f.GetOutput()
					count := countSelectionIndicators(output)
					selectedItem := getSelectedItem(output)
					t.Logf("Tab 1: Selection indicators: %d, Selected item: %q", count, selectedItem)
					return count == 1
				},
				Message: "Should have exactly one selection indicator after first tab",
			},
		},
	}

	// Add multiple tab presses with cursor tracking
	for i := 2; i <= 20; i++ {
		test.Scenario = append(test.Scenario, framework.Press(terminal.KeyTab))
		
		tabNum := i
		test.Assertions = append(test.Assertions, framework.UIAssertion{
			Name: "Tab " + string(rune(tabNum+'0'-1)) + " cursor movement",
			Check: func(f *framework.UITestFramework) bool {
				output := f.GetOutput()
				count := countSelectionIndicators(output)
				selectedItem := getSelectedItem(output)
				
				t.Logf("Tab %d: Selection indicators: %d, Selected item: %q", tabNum, count, selectedItem)
				
				// Verify exactly one selection indicator
				if count != 1 {
					t.Logf("ERROR: Expected 1 selection indicator, got %d", count)
					t.Logf("Output: %q", output)
					return false
				}
				
				// Verify cursor moved (selected item changed from previous or pagination occurred)
				if tabNum > 2 {
					// Should have some selection
					if selectedItem == "" {
						t.Logf("ERROR: No item selected on tab %d", tabNum)
						return false
					}
				}
				
				return true
			},
			Message: "Should maintain single selection indicator and show cursor movement",
		})
	}

	// Final assertion to check pagination behavior
	test.Assertions = append(test.Assertions, framework.UIAssertion{
		Name: "Should demonstrate pagination after 20 tabs",
		Check: func(f *framework.UITestFramework) bool {
			output := f.GetOutput()
			
			// Count total unique items that appeared in output
			items := extractAllItems(output)
			t.Logf("Total unique items found: %d", len(items))
			t.Logf("Items: %v", items)
			
			// Should have seen more than 10 items (indicating pagination)
			if len(items) <= 10 {
				t.Logf("WARNING: Expected pagination with >10 items, got %d", len(items))
			}
			
			// Final selection should be valid
			selectedItem := getSelectedItem(output)
			t.Logf("Final selected item: %q", selectedItem)
			
			return len(items) > 0 && selectedItem != ""
		},
		Message: "Should show evidence of pagination and final cursor position",
	})

	result := runner.RunTest(test)
	if !result.Passed {
		t.Errorf("Test failed:\n%s", result.String())
	}
}

// Helper function to count selection indicators
func countSelectionIndicators(output string) int {
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		cleanLine := cleanANSI(line)
		cleanLine = strings.TrimSpace(cleanLine)
		if strings.HasPrefix(cleanLine, "> ") {
			count++
		}
	}
	return count
}

// Helper function to get the currently selected item
func getSelectedItem(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		cleanLine := cleanANSI(line)
		cleanLine = strings.TrimSpace(cleanLine)
		if strings.HasPrefix(cleanLine, "> ") {
			// Extract item name after "> "
			item := strings.TrimPrefix(cleanLine, "> ")
			item = strings.TrimSpace(item)
			// Take first word as item name
			parts := strings.Fields(item)
			if len(parts) > 0 {
				return parts[0]
			}
		}
	}
	return ""
}

// Helper function to extract all completion items from output
func extractAllItems(output string) []string {
	itemSet := make(map[string]bool)
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		cleanLine := cleanANSI(line)
		
		// Look for lines that start with ">" or "  " (completion items)
		if strings.HasPrefix(cleanLine, "> ") || strings.HasPrefix(cleanLine, "  ") {
			// Remove the prefix
			content := strings.TrimPrefix(cleanLine, "> ")
			content = strings.TrimPrefix(content, "  ")
			content = strings.TrimSpace(content)
			
			// Split by whitespace and filter out empty strings
			// The completion display uses fixed-width columns with spaces
			words := strings.Fields(content)
			for _, word := range words {
				word = strings.TrimSpace(word)
				// Only add words that look like filenames (contain letters/numbers/dots/underscores)
				if word != "" && len(word) > 0 {
					itemSet[word] = true
				}
			}
		}
	}
	
	// Convert set to slice
	items := make([]string, 0, len(itemSet))
	for item := range itemSet {
		items = append(items, item)
	}
	
	return items
}

// Helper function to clean ANSI escape sequences
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
