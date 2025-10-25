package integration

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestExitStatusHandling tests that DSH properly propagates command exit status.
func TestExitStatusHandling(t *testing.T) {
	dshPath := filepath.Join("..", "..", "dsh")

	tests := []struct {
		name           string
		command        string
		expectedStatus int
	}{
		{"successful command", "true", 0},
		{"failing command", "false", 1},
		{"nonexistent command", "nonexistentcommand123", 1},
		{"echo command", "echo hello", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			cmd := exec.CommandContext(ctx, dshPath, "-c", test.command)
			err := cmd.Run()

			var actualStatus int
			if err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) {
					actualStatus = exitError.ExitCode()
				} else {
					t.Fatalf("Unexpected error type: %v", err)
				}
			} else {
				actualStatus = 0
			}

			if actualStatus != test.expectedStatus {
				t.Errorf("Command %q: expected exit status %d, got %d",
					test.command, test.expectedStatus, actualStatus)
			}
		})
	}
}

// TestBuiltinExitStatus tests exit status for builtin commands.
func TestBuiltinExitStatus(t *testing.T) {
	dshPath := filepath.Join("..", "..", "dsh")

	tests := []struct {
		name           string
		command        string
		expectedStatus int
	}{
		{"pwd builtin", "pwd", 0},
		{"help builtin", "help", 0},
		{"echo builtin", "echo test", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			cmd := exec.CommandContext(ctx, dshPath, "-c", test.command)
			err := cmd.Run()

			var actualStatus int
			if err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) {
					actualStatus = exitError.ExitCode()
				} else {
					t.Fatalf("Unexpected error type: %v", err)
				}
			} else {
				actualStatus = 0
			}

			if actualStatus != test.expectedStatus {
				t.Errorf("Builtin %q: expected exit status %d, got %d",
					test.command, test.expectedStatus, actualStatus)
			}
		})
	}
}

// TestCommandChainExitStatus tests exit status with command chaining.
func TestCommandChainExitStatus(t *testing.T) {
	dshPath := filepath.Join("..", "..", "dsh")

	tests := []struct {
		name           string
		command        string
		expectedStatus int
		description    string
	}{
		{
			"success then success",
			"true; true",
			0,
			"Last command succeeds",
		},
		{
			"success then failure",
			"true; false",
			1,
			"Last command fails",
		},
		{
			"failure then success",
			"false; true",
			0,
			"Last command succeeds",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			cmd := exec.CommandContext(ctx, dshPath, "-c", test.command)
			err := cmd.Run()

			var actualStatus int
			if err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) {
					actualStatus = exitError.ExitCode()
				} else {
					t.Fatalf("Unexpected error type: %v", err)
				}
			} else {
				actualStatus = 0
			}

			if actualStatus != test.expectedStatus {
				t.Errorf("Command chain %q: expected exit status %d, got %d (%s)",
					test.command, test.expectedStatus, actualStatus, test.description)
			}
		})
	}
}
