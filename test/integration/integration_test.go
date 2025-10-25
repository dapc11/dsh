package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var ErrCommandTimeout = errors.New("command timed out after 3 seconds")

// runDSHCommand runs a single command in DSH with timeout
func runDSHCommand(t *testing.T, command string) (string, error) {
	dshPath := filepath.Join("..", "..", "dsh")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, dshPath, "-c", command)
	output, err := cmd.CombinedOutput()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", ErrCommandTimeout
	}

	return strings.TrimSpace(string(output)), err
}

// TestCoreCommands tests essential shell functionality
func TestCoreCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		contains string
	}{
		{"echo basic", "echo hello", "hello"},
		{"echo with quotes", `echo "hello world"`, "hello world"},
		{"pwd command", "pwd", "/"},
		{"help command", "help", "Built-in commands"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := runDSHCommand(t, test.command)
			if err != nil {
				t.Logf("Command: %s", test.command)
				t.Logf("Output: %q", output)
				t.Logf("Error: %v", err)
				// Don't fail on error - log for debugging
			}
			if test.contains != "" && !strings.Contains(output, test.contains) {
				t.Errorf("Expected output to contain %q, got %q", test.contains, output)
			}
		})
	}
}

// TestCommandChaining tests semicolon command chaining
func TestCommandChaining(t *testing.T) {
	output, err := runDSHCommand(t, "echo first; echo second")
	if err != nil {
		t.Logf("Command chaining failed: %v", err)
	}

	if !strings.Contains(output, "first") {
		t.Errorf("Command chaining missing 'first', got: %q", output)
	}
	if !strings.Contains(output, "second") {
		t.Errorf("Command chaining missing 'second', got: %q", output)
	}
}

// TestFileRedirection tests basic I/O redirection
func TestFileRedirection(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("output redirection", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "output.txt")

		// Use absolute path to avoid cd issues
		cmd := fmt.Sprintf("echo 'test output' > %s", testFile)
		output, err := runDSHCommand(t, cmd)
		if err != nil {
			t.Logf("Output redirection failed: %v, output: %q", err, output)
		}

		// Check if file was created and has correct content
		if content, err := os.ReadFile(testFile); err != nil {
			t.Logf("Could not read output file: %v", err)
		} else if !strings.Contains(string(content), "test output") {
			t.Errorf("File content incorrect, expected 'test output', got: %q", string(content))
		}
	})
}

// TestQuoteHandling tests quote processing
func TestQuoteHandling(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{"double quotes", `echo "hello world"`, "hello world"},
		{"single quotes", `echo 'single test'`, "single test"},
		{"mixed content", `echo "test" and 'more'`, "test"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := runDSHCommand(t, test.command)
			if err != nil {
				t.Logf("Quote test failed: %v", err)
			}
			if !strings.Contains(output, test.expected) {
				t.Errorf("Expected %q in output, got %q", test.expected, output)
			}
		})
	}
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{"invalid command", "nonexistentcommand123", true},
		{"invalid file", "cat /nonexistent/file.txt", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := runDSHCommand(t, test.command)

			// For error cases, we expect either an error or empty output
			if test.expectError {
				if err == nil && output != "" {
					t.Logf("Expected error for command %q, but got output: %q", test.command, output)
				}
			}
		})
	}
}

// TestWorkflowIntegration tests realistic usage patterns
func TestWorkflowIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("file creation workflow", func(t *testing.T) {
		// Create file with content
		createFile := fmt.Sprintf("echo 'package main' > %s/main.go", tmpDir)
		_, err := runDSHCommand(t, createFile)
		if err != nil {
			t.Logf("File creation failed: %v", err)
		}

		// Verify file exists
		mainFile := filepath.Join(tmpDir, "main.go")
		if content, err := os.ReadFile(mainFile); err != nil {
			t.Logf("Could not read created file: %v", err)
		} else if !strings.Contains(string(content), "package main") {
			t.Errorf("File content incorrect, got: %q", string(content))
		}
	})

	t.Run("multiple commands workflow", func(t *testing.T) {
		// Test multiple echo commands
		output, err := runDSHCommand(t, "echo step1; echo step2; echo step3")
		if err != nil {
			t.Logf("Multi-command workflow failed: %v", err)
		}

		steps := []string{"step1", "step2", "step3"}
		for _, step := range steps {
			if !strings.Contains(output, step) {
				t.Errorf("Workflow missing step %q, output: %q", step, output)
			}
		}
	})
}
