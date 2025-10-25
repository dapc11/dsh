package tests

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestShell_BasicCommands tests basic shell functionality end-to-end.
func TestShell_BasicCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "echo command",
			input:    "echo hello world\nexit\n",
			expected: "hello world",
		},
		{
			name:     "pwd command",
			input:    "pwd\nexit\n",
			expected: "/", // Will contain current directory
		},
		{
			name:     "help command",
			input:    "help\nexit\n",
			expected: "Built-in commands:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runShellCommand(tt.input)
			if err != nil {
				t.Fatalf("Shell execution failed: %v", err)
			}

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expected, output)
			}
		})
	}
}

// TestShell_FileRedirection tests I/O redirection functionality.
func TestShell_FileRedirection(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.txt")

	// Create input file
	err := os.WriteFile(inputFile, []byte("test content\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test output redirection
	input := "echo hello world > " + outputFile + "\nexit\n"
	_, err = runShellCommand(input)
	if err != nil {
		t.Fatalf("Shell execution failed: %v", err)
	}

	// Check output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "hello world") {
		t.Errorf("Output file should contain 'hello world', got: %s", string(content))
	}

	// Test input redirection
	input = "cat < " + inputFile + "\nexit\n"
	output, err := runShellCommand(input)
	if err != nil {
		t.Fatalf("Shell execution failed: %v", err)
	}

	if !strings.Contains(output, "test content") {
		t.Errorf("Expected 'test content' in output, got: %s", output)
	}
}

// TestShell_BackgroundCommands tests background process execution.
func TestShell_BackgroundCommands(t *testing.T) {
	input := "sleep 0.1 &\necho done\nexit\n"

	start := time.Now()
	output, err := runShellCommand(input)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Shell execution failed: %v", err)
	}

	// Should complete quickly (background process)
	if duration > 500*time.Millisecond {
		t.Errorf("Background command took too long: %v", duration)
	}

	if !strings.Contains(output, "done") {
		t.Errorf("Expected 'done' in output, got: %s", output)
	}
}

// TestShell_MultipleCommands tests command chaining with semicolons.
func TestShell_MultipleCommands(t *testing.T) {
	input := "echo first; echo second; echo third\nexit\n"

	output, err := runShellCommand(input)
	if err != nil {
		t.Fatalf("Shell execution failed: %v", err)
	}

	expected := []string{"first", "second", "third"}
	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected '%s' in output, got: %s", exp, output)
		}
	}
}

// TestShell_ErrorHandling tests error conditions.
func TestShell_ErrorHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "non-existent command",
			input: "nonexistentcommand12345\nexit\n",
		},
		{
			name:  "invalid redirection",
			input: "echo hello > /invalid/path/file.txt\nexit\n",
		},
		{
			name:  "cd to non-existent directory",
			input: "cd /nonexistent/directory\nexit\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runShellCommand(tt.input)
			// Shell should handle errors gracefully, not crash
			if err != nil && strings.Contains(err.Error(), "signal: killed") {
				t.Errorf("Shell crashed on error case: %v", err)
			}
		})
	}
}

// runShellCommand executes the shell with given input and returns output.
func runShellCommand(input string) (string, error) {
	ctx := context.Background()

	// Build the shell first
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", "dsh_test", ".")
	buildCmd.Dir = ".."
	err := buildCmd.Run()
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove("../dsh_test") }()

	// Run the shell with input
	cmd := exec.CommandContext(ctx, "./dsh_test")
	cmd.Dir = ".."
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Combine stdout and stderr for analysis
	output := stdout.String() + stderr.String()

	return output, err
}
