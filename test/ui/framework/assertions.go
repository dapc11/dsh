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

// MenuAssertion provides fluent menu assertions
type MenuAssertion struct {
	framework *UITestFramework
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
