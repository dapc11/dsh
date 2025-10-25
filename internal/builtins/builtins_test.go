package builtins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteBuiltin_Help(t *testing.T) {
	result := ExecuteBuiltin([]string{"help"})
	if !result {
		t.Error("help command should return true")
	}
}

func TestExecuteBuiltin_Exit(t *testing.T) {
	result := ExecuteBuiltin([]string{"exit"})
	if result {
		t.Error("exit command should return false")
	}
}

func TestExecuteBuiltin_Pwd(t *testing.T) {
	result := ExecuteBuiltin([]string{"pwd"})
	if !result {
		t.Error("pwd command should return true")
	}
}

func TestExecuteBuiltin_Cd(t *testing.T) {
	// Save current directory and restore after test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	// Test cd to home
	result := ExecuteBuiltin([]string{"cd"})
	if !result {
		t.Error("cd command should return true")
	}

	// Test cd to specific directory
	tmpDir := t.TempDir()
	result = ExecuteBuiltin([]string{"cd", tmpDir})
	if !result {
		t.Error("cd to valid directory should return true")
	}

	// Verify we're in the right directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if currentDir != tmpDir {
		t.Errorf("cd failed: current dir = %s, want %s", currentDir, tmpDir)
	}

	// Test cd to non-existent directory
	result = ExecuteBuiltin([]string{"cd", "/nonexistent/directory"})
	if result {
		t.Error("cd to invalid directory should return false")
	}
}

func TestExecuteBuiltin_Todo(t *testing.T) {
	// Test todo list (empty)
	result := ExecuteBuiltin([]string{"todo"})
	if !result {
		t.Error("todo list command should return true")
	}

	// Test adding todo
	result = ExecuteBuiltin([]string{"todo", "test", "item"})
	if !result {
		t.Error("todo add command should return true")
	}
}

func TestExecuteBuiltin_InvalidCommand(t *testing.T) {
	result := ExecuteBuiltin([]string{"nonexistent"})
	if result {
		t.Error("invalid command should return false")
	}
}

func TestExecuteBuiltin_EmptyArgs(t *testing.T) {
	result := ExecuteBuiltin([]string{})
	if result {
		t.Error("empty args should return false")
	}
}

func TestIsBuiltin(t *testing.T) {
	builtins := []string{"cd", "pwd", "help", "exit", "todo"}

	for _, cmd := range builtins {
		if !IsBuiltin(cmd) {
			t.Errorf("IsBuiltin(%s) = false, want true", cmd)
		}
	}

	nonBuiltins := []string{"ls", "echo", "cat", "grep"}
	for _, cmd := range nonBuiltins {
		if IsBuiltin(cmd) {
			t.Errorf("IsBuiltin(%s) = true, want false", cmd)
		}
	}
}

func TestAddTodo(t *testing.T) {
	// Create temporary home directory
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	err := addTodo("test todo item")
	if err != nil {
		t.Errorf("addTodo failed: %v", err)
	}

	// Check if file was created
	todoFile := filepath.Join(tmpDir, ".dsh_todos")
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		t.Error("Todo file was not created")
	}

	// Check file contents
	content, err := os.ReadFile(todoFile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "test todo item") {
		t.Error("Todo item not found in file")
	}
}

func TestListTodos(t *testing.T) {
	// Create temporary home directory
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Test with no todos file
	todos := loadTodos()
	if todos != nil {
		t.Error("Expected nil for non-existent todos file")
	}

	// Create todos file
	todoFile := filepath.Join(tmpDir, ".dsh_todos")
	content := "[2023-01-01 12:00:00] First todo\n[2023-01-02 13:00:00] Second todo\n"
	err := os.WriteFile(todoFile, []byte(content), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Test reading todos
	todos = loadTodos()
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	if !strings.Contains(todos[0], "First todo") {
		t.Error("First todo not found")
	}
}
