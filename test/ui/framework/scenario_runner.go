package framework

import (
	"fmt"
	"strings"

	"dsh/internal/terminal"
)

// ScenarioRunner executes UI test scenarios
type ScenarioRunner struct {
	framework *UITestFramework
}

// NewScenarioRunner creates a new scenario runner
func NewScenarioRunner(framework *UITestFramework) *ScenarioRunner {
	return &ScenarioRunner{framework: framework}
}

// RunTest executes a complete UI test
func (r *ScenarioRunner) RunTest(test UITest) TestResult {
	result := TestResult{
		Name:    test.Name,
		Passed:  true,
		Results: make([]AssertionResult, 0),
	}

	// Setup
	if test.Setup != nil {
		test.Setup(r.framework)
	}

	// Execute scenario actions
	for i, action := range test.Scenario {
		actionResult := r.executeAction(action)
		if !actionResult.Passed {
			result.Passed = false
			result.Results = append(result.Results, AssertionResult{
				Passed:  false,
				Message: fmt.Sprintf("Action %d failed: %s", i, actionResult.Message),
			})
			break
		}
	}

	// Run assertions
	for _, assertion := range test.Assertions {
		assertResult := AssertionResult{
			Passed:  assertion.Check(r.framework),
			Message: assertion.Message,
		}
		result.Results = append(result.Results, assertResult)
		if !assertResult.Passed {
			result.Passed = false
		}
	}

	// Cleanup
	if test.Cleanup != nil {
		test.Cleanup(r.framework)
	}

	return result
}

// executeAction executes a single UI action
func (r *ScenarioRunner) executeAction(action UIAction) AssertionResult {
	switch action.Type {
	case ActionType_:
		return r.executeType(action.Data.(string))
	case ActionKeyPress:
		return r.executeKeyPress(action.Data.(terminal.KeyEvent))
	case ActionValidate:
		return r.executeValidation(action.Expected)
	case ActionSetupKeys:
		return r.executeSetupKeys(action.Data.([]terminal.KeyEvent))
	default:
		return AssertionResult{Passed: false, Message: "Unknown action type"}
	}
}

// executeType simulates typing text
func (r *ScenarioRunner) executeType(text string) AssertionResult {
	for _, char := range text {
		keyEvent := terminal.KeyEvent{
			Key:  terminal.KeyNone,
			Rune: char,
		}
		r.framework.recorder.RecordKey(keyEvent)
		// Process the key through the real readline
		r.framework.readline.ProcessKey(keyEvent)
	}
	return AssertionResult{Passed: true, Message: fmt.Sprintf("Typed: %s", text)}
}

// executeKeyPress simulates pressing a key
func (r *ScenarioRunner) executeKeyPress(keyEvent terminal.KeyEvent) AssertionResult {
	r.framework.recorder.RecordKey(keyEvent)
	// Process the key through the real readline
	shouldContinue := r.framework.readline.ProcessKey(keyEvent)

	// If ProcessKey returns false, simulate what the main readline loop does
	if !shouldContinue {
		// This means the line is complete (Enter was pressed)
		_, _ = r.framework.mockTerm.WriteString("\r\n")
		_, _ = r.framework.mockTerm.WriteString(r.framework.readline.GetPrompt())
		// Reset the buffer for the next line
		r.framework.readline.SetBuffer("")
	}

	return AssertionResult{Passed: true, Message: fmt.Sprintf("Pressed key: %v", keyEvent.Key)}
}

// executeValidation validates expected output
func (r *ScenarioRunner) executeValidation(expected string) AssertionResult {
	output := r.framework.GetOutput()
	if strings.Contains(output, expected) {
		return AssertionResult{Passed: true, Message: fmt.Sprintf("Found expected output: %s", expected)}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Expected output not found: %s", expected)}
}

// executeSetupKeys queues keys for interactive components
func (r *ScenarioRunner) executeSetupKeys(keys []terminal.KeyEvent) AssertionResult {
	mockTerm := r.framework.GetMockTerminal()
	mockTerm.QueueKeys(keys...)
	return AssertionResult{Passed: true, Message: fmt.Sprintf("Queued %d keys for interactive input", len(keys))}
}

// TestResult represents the result of a UI test
type TestResult struct {
	Name    string
	Passed  bool
	Results []AssertionResult
}

// String returns a string representation of the test result
func (r TestResult) String() string {
	status := "PASS"
	if !r.Passed {
		status = "FAIL"
	}

	var details strings.Builder
	for _, result := range r.Results {
		if result.Passed {
			_, _ = details.WriteString(fmt.Sprintf("  ✓ %s\n", result.Message))
		} else {
			_, _ = details.WriteString(fmt.Sprintf("  ✗ %s\n", result.Message))
		}
	}

	return fmt.Sprintf("%s: %s\n%s", status, r.Name, details.String())
}
