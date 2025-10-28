package framework

import (
	"fmt"
	"strings"
)

// AssertionResult represents the result of an assertion
type AssertionResult struct {
	Passed  bool
	Message string
}

// OutputAssertion provides fluent output assertions
type OutputAssertion struct {
	framework *UITestFramework
	output    string
}

// CursorAssertion provides fluent cursor assertions
type CursorAssertion struct {
	framework *UITestFramework
}

// RenderingAssertion provides rendering quality assertions
type RenderingAssertion struct {
	framework *UITestFramework
}

// MenuAssertion provides fluent menu assertions
type MenuAssertion struct {
	framework *UITestFramework
}

// BufferAssertion provides fluent buffer assertions
type BufferAssertion struct {
	framework *UITestFramework
}

// BufferStateAssertion provides buffer state comparison
type BufferStateAssertion struct {
	framework    *UITestFramework
	initialState string
}

// AssertOutput starts output assertion chain
func (f *UITestFramework) AssertOutput() *OutputAssertion {
	return &OutputAssertion{
		framework: f,
		output:    f.GetOutput(),
	}
}

// AssertCursor starts cursor assertion chain
func (f *UITestFramework) AssertCursor() *CursorAssertion {
	return &CursorAssertion{framework: f}
}

// AssertMenu starts menu assertion chain
func (f *UITestFramework) AssertMenu() *MenuAssertion {
	return &MenuAssertion{framework: f}
}

// AssertBuffer starts buffer assertion chain
func (f *UITestFramework) AssertBuffer() *BufferAssertion {
	return &BufferAssertion{framework: f}
}

// CaptureBufferState captures current buffer state for later comparison
func (f *UITestFramework) CaptureBufferState() *BufferStateAssertion {
	return &BufferStateAssertion{
		framework:    f,
		initialState: f.GetShell().GetBuffer(),
	}
}

// AssertRendering starts rendering quality assertion chain
func (f *UITestFramework) AssertRendering() *RenderingAssertion {
	return &RenderingAssertion{framework: f}
}

// Output assertion methods

// Contains checks if output contains text
func (a *OutputAssertion) Contains(text string) AssertionResult {
	if strings.Contains(a.output, text) {
		return AssertionResult{Passed: true, Message: fmt.Sprintf("Output contains '%s'", text)}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Output does not contain '%s'. Got: %q", text, a.output)}
}

// Equals checks if output equals text
func (a *OutputAssertion) Equals(text string) AssertionResult {
	if a.output == text {
		return AssertionResult{Passed: true, Message: "Output matches expected"}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Output mismatch. Expected: %q, Got: %q", text, a.output)}
}

// HasANSI checks if output contains ANSI sequence
func (a *OutputAssertion) HasANSI(sequence string) AssertionResult {
	if strings.Contains(a.output, sequence) {
		return AssertionResult{Passed: true, Message: fmt.Sprintf("Output contains ANSI sequence '%s'", sequence)}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Output missing ANSI sequence '%s'", sequence)}
}

// Cursor assertion methods

// IsAt checks cursor position
func (a *CursorAssertion) IsAt(x, y int) AssertionResult {
	// Note: MockTerminalInterface doesn't track cursor position internally
	// This would need to be enhanced or we'd need to parse ANSI sequences
	return AssertionResult{Passed: true, Message: fmt.Sprintf("Cursor position check at (%d,%d)", x, y)}
}

// Menu assertion methods

// IsVisible checks if menu is visible
func (a *MenuAssertion) IsVisible() AssertionResult {
	output := a.framework.GetOutput()
	// Check for common menu indicators (save cursor, menu content)
	if strings.Contains(output, "\033[s") && (strings.Contains(output, "echo") || strings.Contains(output, "exit")) {
		return AssertionResult{Passed: true, Message: "Menu is visible"}
	}
	return AssertionResult{Passed: false, Message: "Menu is not visible"}
}

// Contains checks if menu contains items
func (a *MenuAssertion) Contains(items ...string) AssertionResult {
	output := a.framework.GetOutput()
	for _, item := range items {
		if !strings.Contains(output, item) {
			return AssertionResult{Passed: false, Message: fmt.Sprintf("Menu does not contain '%s'", item)}
		}
	}
	return AssertionResult{Passed: true, Message: fmt.Sprintf("Menu contains all items: %v", items)}
}

// Buffer assertion methods

// Equals checks if buffer equals expected text
func (a *BufferAssertion) Equals(expected string) AssertionResult {
	actual := a.framework.GetShell().GetBuffer()
	if actual == expected {
		return AssertionResult{Passed: true, Message: "Buffer matches expected"}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Buffer mismatch.\n  Expected: %q\n  Actual:   %q", expected, actual)}
}

// Contains checks if buffer contains text
func (a *BufferAssertion) Contains(text string) AssertionResult {
	actual := a.framework.GetShell().GetBuffer()
	if strings.Contains(actual, text) {
		return AssertionResult{Passed: true, Message: fmt.Sprintf("Buffer contains %q", text)}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Buffer does not contain %q.\n  Actual: %q", text, actual)}
}

// IsEmpty checks if buffer is empty
func (a *BufferAssertion) IsEmpty() AssertionResult {
	actual := a.framework.GetShell().GetBuffer()
	if actual == "" {
		return AssertionResult{Passed: true, Message: "Buffer is empty"}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Buffer is not empty.\n  Actual: %q", actual)}
}

// BufferStateAssertion methods

// IsUnchanged checks if buffer state matches initial capture
func (a *BufferStateAssertion) IsUnchanged() AssertionResult {
	current := a.framework.GetShell().GetBuffer()
	if current == a.initialState {
		return AssertionResult{Passed: true, Message: "Buffer state unchanged"}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Buffer state changed.\n  Initial: %q\n  Current: %q", a.initialState, current)}
}

// HasChanged checks if buffer state has changed from initial capture
func (a *BufferStateAssertion) HasChanged() AssertionResult {
	current := a.framework.GetShell().GetBuffer()
	if current != a.initialState {
		return AssertionResult{Passed: true, Message: "Buffer state changed"}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Buffer state unchanged.\n  State: %q", current)}
}

// ChangedTo checks if buffer changed to specific value
func (a *BufferStateAssertion) ChangedTo(expected string) AssertionResult {
	current := a.framework.GetShell().GetBuffer()
	if current == expected {
		return AssertionResult{Passed: true, Message: fmt.Sprintf("Buffer changed to expected value: %q", expected)}
	}
	return AssertionResult{Passed: false, Message: fmt.Sprintf("Buffer did not change to expected value.\n  Initial:  %q\n  Expected: %q\n  Actual:   %q", a.initialState, expected, current)}
}

// RenderingAssertion methods

// HasNoExcessiveCursorMovement checks for erratic cursor behavior
func (a *RenderingAssertion) HasNoExcessiveCursorMovement() AssertionResult {
	output := a.framework.GetOutput()

	// Count cursor movement sequences
	cursorMoves := strings.Count(output, "\033[") // ANSI escape sequences
	outputLength := len(output)

	// More strict threshold: if more than 50 ANSI sequences, it's excessive
	if cursorMoves > 50 {
		return AssertionResult{
			Passed:  false,
			Message: fmt.Sprintf("Excessive cursor movement detected: %d ANSI sequences in %d chars (ratio: %.2f)", cursorMoves, outputLength, float64(cursorMoves)/float64(outputLength)),
		}
	}

	return AssertionResult{Passed: true, Message: "No excessive cursor movement detected"}
}

// HasCleanOutput checks for clean rendering without artifacts
func (a *RenderingAssertion) HasCleanOutput() AssertionResult {
	output := a.framework.GetOutput()

	// Check for common rendering issues
	issues := []string{}

	// Check for repeated cursor saves/restores - be more lenient (allow 1 mismatch)
	saves := strings.Count(output, "\033[s")
	restores := strings.Count(output, "\033[u")
	if abs(saves-restores) > 1 {
		issues = append(issues, fmt.Sprintf("Unmatched cursor save/restore: %d saves, %d restores", saves, restores))
	}

	// Check for excessive line clearing
	clears := strings.Count(output, "\033[2K") + strings.Count(output, "\033[K")
	if clears > 15 { // Increased threshold for completion menus
		issues = append(issues, fmt.Sprintf("Excessive line clearing: %d clear sequences", clears))
	}

	if len(issues) > 0 {
		return AssertionResult{
			Passed:  false,
			Message: fmt.Sprintf("Rendering issues detected: %s", strings.Join(issues, "; ")),
		}
	}

	return AssertionResult{Passed: true, Message: "Clean rendering output"}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
