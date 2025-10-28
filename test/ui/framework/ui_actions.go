package framework

import (
	"dsh/internal/terminal"
)

// ActionType defines the type of UI action
type ActionType int

const (
	ActionKeyPress ActionType = iota
	ActionType_
	ActionWait
	ActionValidate
	ActionSetupKeys
)

// UIAction represents a single UI interaction
type UIAction struct {
	Type     ActionType
	Data     interface{}
	Expected string
}

// UITest represents a complete UI test scenario
type UITest struct {
	Name       string
	Setup      func(*UITestFramework)
	Scenario   []UIAction
	Assertions []UIAssertion
	Cleanup    func(*UITestFramework)
}

// UIAssertion represents a UI state assertion
type UIAssertion struct {
	Name    string
	Check   func(*UITestFramework) bool
	Message string
}

// Action builders for fluent API

// Type creates a typing action
func Type(text string) UIAction {
	return UIAction{
		Type: ActionType_,
		Data: text,
	}
}

// Press creates a key press action
func Press(key terminal.Key) UIAction {
	return UIAction{
		Type: ActionKeyPress,
		Data: terminal.KeyEvent{Key: key},
	}
}

// PressCtrl creates a Ctrl+key action
func PressCtrl(key rune) UIAction {
	var ctrlKey terminal.Key
	switch key {
	case 'a':
		ctrlKey = terminal.KeyCtrlA
	case 'b':
		ctrlKey = terminal.KeyCtrlB
	case 'c':
		ctrlKey = terminal.KeyCtrlC
	case 'd':
		ctrlKey = terminal.KeyCtrlD
	case 'e':
		ctrlKey = terminal.KeyCtrlE
	case 'f':
		ctrlKey = terminal.KeyCtrlF
	case 'k':
		ctrlKey = terminal.KeyCtrlK
	case 'l':
		ctrlKey = terminal.KeyCtrlL
	case 'n':
		ctrlKey = terminal.KeyCtrlN
	case 'p':
		ctrlKey = terminal.KeyCtrlP
	case 'r':
		ctrlKey = terminal.KeyCtrlR
	case 'u':
		ctrlKey = terminal.KeyCtrlU
	case 'w':
		ctrlKey = terminal.KeyCtrlW
	case 'y':
		ctrlKey = terminal.KeyCtrlY
	case 'z':
		ctrlKey = terminal.KeyCtrlZ
	default:
		ctrlKey = terminal.KeyCtrlA + terminal.Key(key-'a') // fallback
	}

	return UIAction{
		Type: ActionKeyPress,
		Data: terminal.KeyEvent{Key: ctrlKey},
	}
}

// Expect creates an output expectation
func Expect(output string) UIAction {
	return UIAction{
		Type:     ActionValidate,
		Expected: output,
	}
}
