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
	executor := New()
	
	cmd := &parser.Command{
		Args: []string{"echo", "hello", "world"},
	}

	err := executor.Execute(cmd)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}

func TestExecutor_NonExistentCommand(t *testing.T) {
	executor := New()
	
	cmd := &parser.Command{
		Args: []string{"nonexistentcommand12345"},
	}

	err := executor.Execute(cmd)
	if err == nil {
		t.Error("Expected error for non-existent command")
	}
}

func TestExecutor_OutputRedirection(t *testing.T) {
	executor := New()
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")
	
	cmd := &parser.Command{
		Args:       []string{"echo", "hello", "world"},
		OutputFile: outputFile,
	}

	err := executor.Execute(cmd)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	// Check if file was created and contains expected content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
	}

	expected := "hello world"
	if !strings.Contains(string(content), expected) {
		t.Errorf("Output file content: expected to contain '%s', got '%s'", expected, string(content))
	}
}

func TestExecutor_InputRedirection(t *testing.T) {
	executor := New()
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	
	// Create input file
	inputContent := "test input content"
	err := os.WriteFile(inputFile, []byte(inputContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	
	cmd := &parser.Command{
		Args:      []string{"cat"},
		InputFile: inputFile,
	}

	err = executor.Execute(cmd)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}

func TestExecutor_AppendRedirection(t *testing.T) {
	executor := New()
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")
	
	// Create initial file content
	initialContent := "initial content\n"
	err := os.WriteFile(outputFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}
	
	cmd := &parser.Command{
		Args:         []string{"echo", "appended content"},
		OutputFile:   outputFile,
		AppendOutput: true,
	}

	err = executor.Execute(cmd)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	// Check if content was appended
	content, err := os.ReadFile(outputFile)
	if err != nil {
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
	executor := New()
	
	cmd := &parser.Command{
		Args:       []string{"sleep", "0.1"},
		Background: true,
	}

	start := time.Now()
	err := executor.Execute(cmd)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	// Background command should return immediately
	if duration > 50*time.Millisecond {
		t.Errorf("Background command took too long: %v", duration)
	}
}

func TestExecutor_EmptyCommand(t *testing.T) {
	executor := New()
	
	cmd := &parser.Command{
		Args: []string{},
	}

	err := executor.Execute(cmd)
	if err == nil {
		t.Error("Expected error for empty command")
	}
}

func TestExecutor_CommandWithArguments(t *testing.T) {
	executor := New()
	
	cmd := &parser.Command{
		Args: []string{"echo", "-n", "no newline"},
	}

	err := executor.Execute(cmd)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}

func TestExecutor_InvalidOutputFile(t *testing.T) {
	executor := New()
	
	cmd := &parser.Command{
		Args:       []string{"echo", "hello"},
		OutputFile: "/invalid/path/that/does/not/exist/file.txt",
	}

	err := executor.Execute(cmd)
	if err == nil {
		t.Error("Expected error for invalid output file path")
	}
}

func TestExecutor_InvalidInputFile(t *testing.T) {
	executor := New()
	
	cmd := &parser.Command{
		Args:      []string{"cat"},
		InputFile: "/nonexistent/file.txt",
	}

	err := executor.Execute(cmd)
	if err == nil {
		t.Error("Expected error for non-existent input file")
	}
}
