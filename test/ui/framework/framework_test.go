package framework

import (
	"testing"

	"dsh/internal/terminal"
)

func TestUITestFramework_Creation(t *testing.T) {
	fw := NewUITestFramework()

	if fw == nil {
		t.Fatal("Framework should not be nil")
	}

	if fw.mockTerm == nil {
		t.Fatal("MockTerminal should not be nil")
	}

	if fw.recorder == nil {
		t.Fatal("InteractionRecorder should not be nil")
	}

	if fw.shell == nil {
		t.Fatal("TestShell should not be nil")
	}
}

func TestUITestFramework_SetPrompt(t *testing.T) {
	fw := NewUITestFramework()
	result := fw.SetPrompt("test> ")

	// Basic validation that SetPrompt returns the framework (fluent interface)
	if result != fw {
		t.Fatal("SetPrompt should return the same framework instance")
	}
}

func TestUITestFramework_AddHistory(t *testing.T) {
	fw := NewUITestFramework()
	fw.AddHistory("git status", "ls -la", "echo hello")

	if len(fw.history) != 3 {
		t.Errorf("Expected 3 history items, got %d", len(fw.history))
	}

	expected := []string{"git status", "ls -la", "echo hello"}
	for i, cmd := range expected {
		if fw.history[i] != cmd {
			t.Errorf("History[%d]: expected %q, got %q", i, cmd, fw.history[i])
		}
	}
}

func TestInteractionRecorder_RecordKey(t *testing.T) {
	recorder := NewInteractionRecorder()

	keyEvent := terminal.KeyEvent{
		Key:  terminal.KeyTab,
		Rune: 0,
	}

	recorder.RecordKey(keyEvent)

	keystrokes := recorder.GetKeystrokes()
	if len(keystrokes) != 1 {
		t.Errorf("Expected 1 keystroke, got %d", len(keystrokes))
	}

	if keystrokes[0].Key != terminal.KeyTab {
		t.Errorf("Expected KeyTab, got %v", keystrokes[0].Key)
	}
}

func TestOutputAssertion_Contains(t *testing.T) {
	fw := NewUITestFramework()

	// Simulate some output
	fw.mockTerm.WriteString("hello world")

	result := fw.AssertOutput().Contains("hello")
	if !result.Passed {
		t.Errorf("Assertion should pass: %s", result.Message)
	}

	result = fw.AssertOutput().Contains("missing")
	if result.Passed {
		t.Errorf("Assertion should fail: %s", result.Message)
	}
}

func TestMenuAssertion_IsVisible(t *testing.T) {
	fw := NewUITestFramework()

	// Simulate menu output (save cursor + menu content)
	fw.mockTerm.WriteString("\033[s\necho  exit  help")

	result := fw.AssertMenu().IsVisible()
	if !result.Passed {
		t.Errorf("Menu should be visible: %s", result.Message)
	}
}

func TestScenarioRunner_ExecuteType(t *testing.T) {
	fw := NewUITestFramework()
	runner := NewScenarioRunner(fw)

	result := runner.executeType("hello")
	if !result.Passed {
		t.Errorf("Type action should succeed: %s", result.Message)
	}

	// Check that keystrokes were recorded
	keystrokes := fw.recorder.GetKeystrokes()
	if len(keystrokes) != 5 { // "hello" = 5 characters
		t.Errorf("Expected 5 keystrokes, got %d", len(keystrokes))
	}
}
