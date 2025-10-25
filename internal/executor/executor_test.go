package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dsh/internal/parser"
)

func TestExecutor_SimpleCommand(t *testing.T) {
	cmd := &parser.Command{
		Args: []string{"echo", "hello", "world"},
	}

	success := ExecuteCommand(cmd)
	if !success {
		t.Error("Expected command to succeed")
	}
}

func TestExecutor_NonExistentCommand(t *testing.T) {
	

	cmd := &parser.Command{
		Args: []string{"nonexistentcommand12345"},
	}

	success := ExecuteCommand(cmd)
	if err == nil {
		t.Error("Expected error for non-existent command")
	}
}

func TestExecutor_OutputRedirection(t *testing.T) {
	
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	cmd := &parser.Command{
		Args:       []string{"echo", "hello", "world"},
		OutputFile: outputFile,
	}

	success := ExecuteCommand(cmd)
	if !success {
		t.Error("Expected command to succeed")
	}

	// Check if file was created and contains expected content
	content, err := os.ReadFile(outputFile)
	if !success {
		t.Errorf("Failed to read output file: %v", err)
	}

	expected := "hello world"
	if !strings.Contains(string(content), expected) {
		t.Errorf("Output file content: expected to contain '%s', got '%s'", expected, string(content))
	}
}

func TestExecutor_InputRedirection(t *testing.T) {
	
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	// Create input file
	inputContent := "test input content"
	err := os.WriteFile(inputFile, []byte(inputContent), 0644)
	if !success {
		t.Fatalf("Failed to create input file: %v", err)
	}

	cmd := &parser.Command{
		Args:      []string{"cat"},
		InputFile: inputFile,
	}

	err = ExecuteCommand(cmd)
	if !success {
		t.Error("Expected command to succeed")
	}
}

func TestExecutor_AppendRedirection(t *testing.T) {
	
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	// Create initial file content
	initialContent := "initial content\n"
	err := os.WriteFile(outputFile, []byte(initialContent), 0644)
	if !success {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	cmd := &parser.Command{
		Args:         []string{"echo", "appended content"},
		OutputFile:   outputFile,
		AppendMode: true,
	}

	err = ExecuteCommand(cmd)
	if !success {
		t.Error("Expected command to succeed")
	}

	// Check if content was appended
	content, err := os.ReadFile(outputFile)
	if !success {
		t.Errorf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "initial content") {
		t.Error("Original content not found")
	}
	if !strings.Contains(contentStr, "appended content") {
		t.Error("Appended content not found")
	}
}

func TestExecutor_BackgroundCommand(t *testing.T) {
	

	cmd := &parser.Command{
		Args:       []string{"sleep", "0.1"},
		Background: true,
	}

	start := time.Now()
	success := ExecuteCommand(cmd)
	duration := time.Since(start)

	if !success {
		t.Error("Expected command to succeed")
	}

	// Background command should return immediately
	if duration > 50*time.Millisecond {
		t.Errorf("Background command took too long: %v", duration)
	}
}

func TestExecutor_EmptyCommand(t *testing.T) {
	

	cmd := &parser.Command{
		Args: []string{},
	}

	success := ExecuteCommand(cmd)
	if err == nil {
		t.Error("Expected error for empty command")
	}
}

func TestExecutor_CommandWithArguments(t *testing.T) {
	

	cmd := &parser.Command{
		Args: []string{"echo", "-n", "no newline"},
	}

	success := ExecuteCommand(cmd)
	if !success {
		t.Error("Expected command to succeed")
	}
}

func TestExecutor_InvalidOutputFile(t *testing.T) {
	

	cmd := &parser.Command{
		Args:       []string{"echo", "hello"},
		OutputFile: "/invalid/path/that/does/not/exist/file.txt",
	}

	success := ExecuteCommand(cmd)
	if err == nil {
		t.Error("Expected error for invalid output file path")
	}
}

func TestExecutor_InvalidInputFile(t *testing.T) {
	

	cmd := &parser.Command{
		Args:      []string{"cat"},
		InputFile: "/nonexistent/file.txt",
	}

	success := ExecuteCommand(cmd)
	if err == nil {
		t.Error("Expected error for non-existent input file")
	}
}
