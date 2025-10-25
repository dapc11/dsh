package executor

import (
	"testing"
	"dsh/internal/parser"
)

func TestExecutor_SimpleCommand(t *testing.T) {
	cmd := &parser.Command{
		Args: []string{"echo", "hello", "world"},
	}

	if !ExecuteCommand(cmd) {
		t.Error("Expected command to succeed")
	}
}

func TestExecutor_NonExistentCommand(t *testing.T) {
	cmd := &parser.Command{
		Args: []string{"nonexistentcommand12345"},
	}

	if ExecuteCommand(cmd) {
		t.Error("Expected command to fail for non-existent command")
	}
}

func TestExecutor_EmptyCommand(t *testing.T) {
	cmd := &parser.Command{
		Args: []string{},
	}

	if !ExecuteCommand(cmd) {
		t.Error("Expected empty command to succeed")
	}
}

func TestExecutor_BuiltinCommand(t *testing.T) {
	cmd := &parser.Command{
		Args: []string{"pwd"},
	}

	if !ExecuteCommand(cmd) {
		t.Error("Expected builtin command to succeed")
	}
}
