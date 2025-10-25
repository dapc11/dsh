package executor

import (
	"dsh/internal/parser"
	"testing"
)

func TestExecutor_SimpleCommand(t *testing.T) {
	t.Parallel()
	cmd := &parser.Command{
		Args: []string{"echo", "hello", "world"},
	}

	if !ExecuteCommand(cmd) {
		t.Error("Expected command to succeed")
	}
}

func TestExecutor_NonExistentCommand(t *testing.T) {
	t.Parallel()
	cmd := &parser.Command{
		Args: []string{"nonexistentcommand12345"},
	}

	// Command should continue processing (return true) but set exit status
	if !ExecuteCommand(cmd) {
		t.Error("Expected command to continue processing even for non-existent command")
	}

	// Check that exit status was set correctly
	if GetLastExitStatus() == 0 {
		t.Error("Expected non-zero exit status for non-existent command")
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
