package test

import (
	"os"
	"path/filepath"
	"testing"

	"dsh/internal/readline"
)

func TestPathCompletionWithTrailingSlash(t *testing.T) {
	completion := readline.NewCompletion()

	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "dsh_completion_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files and directories
	testFiles := []string{"file1.txt", "file2.log", "another.txt"}
	testDirs := []string{"dir1", "dir2", "subdir"}

	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	for _, dir := range testDirs {
		err := os.Mkdir(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test dir %s: %v", dir, err)
		}
	}

	// Test completion with trailing slash - should show all files
	input := "cat " + tmpDir + "/"
	matches, _ := completion.Complete(input, len(input))

	expectedCount := len(testFiles) + len(testDirs)
	if len(matches) != expectedCount {
		t.Errorf("Expected %d matches for trailing slash, got %d", expectedCount, len(matches))
	}

	// Verify we have both files and directories
	hasFile := false
	hasDir := false
	for _, match := range matches {
		if match.Type == "file" {
			hasFile = true
		}
		if match.Type == "directory" {
			hasDir = true
		}
	}

	if !hasFile {
		t.Error("Expected at least one file in matches")
	}
	if !hasDir {
		t.Error("Expected at least one directory in matches")
	}
}

func TestPartialPathCompletion(t *testing.T) {
	completion := readline.NewCompletion()

	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "dsh_completion_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files with different prefixes
	testFiles := []string{"apple.txt", "application.log", "banana.txt", "cherry.txt"}

	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	// Test partial completion - should only show files starting with "app"
	input := "cat " + filepath.Join(tmpDir, "app")
	matches, commonPrefix := completion.Complete(input, len(input))

	// Should match "apple.txt" and "application.log"
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'app' prefix, got %d", len(matches))
	}

	// Check that matches contain the expected files
	foundApple := false
	foundApplication := false
	for _, match := range matches {
		if filepath.Base(match.Text) == "apple.txt" {
			foundApple = true
		}
		if filepath.Base(match.Text) == "application.log" {
			foundApplication = true
		}
	}

	if !foundApple {
		t.Error("Expected to find apple.txt in matches")
	}
	if !foundApplication {
		t.Error("Expected to find application.log in matches")
	}

	// Common prefix should be "l" (common part of "apple" and "application")
	if commonPrefix != "l" {
		t.Errorf("Expected common prefix 'l', got '%s'", commonPrefix)
	}
}

func TestRootPathCompletion(t *testing.T) {
	completion := readline.NewCompletion()

	// Test completion in root directory with partial path
	input := "cat /h"
	matches, commonPrefix := completion.Complete(input, len(input))

	// Should find /home if it exists
	if len(matches) == 0 {
		t.Skip("No matches for /h - this is expected if /home doesn't exist")
	}

	// If we have matches, they should all start with /h
	for _, match := range matches {
		if !filepath.HasPrefix(match.Text, "/h") {
			t.Errorf("Match %s doesn't start with /h", match.Text)
		}
	}

	// If there's a common prefix, it should extend beyond "h"
	if commonPrefix != "" && len(commonPrefix) == 0 {
		t.Error("Common prefix should not be empty if matches exist")
	}
}

func TestCommandCompletion(t *testing.T) {
	completion := readline.NewCompletion()

	// Test command completion
	input := "ec"
	matches, commonPrefix := completion.Complete(input, len(input))

	// Should find "echo" command
	foundEcho := false
	for _, match := range matches {
		if match.Text == "echo" {
			foundEcho = true
			if match.Type != "builtin" && match.Type != "command" {
				t.Errorf("Expected echo to be builtin or command, got %s", match.Type)
			}
		}
	}

	if !foundEcho {
		t.Error("Expected to find 'echo' command in completion")
	}

	// Common prefix should be "ho" (from "echo")
	if commonPrefix != "ho" {
		t.Errorf("Expected common prefix 'ho', got '%s'", commonPrefix)
	}
}

func TestEmptyDirectoryCompletion(t *testing.T) {
	completion := readline.NewCompletion()

	// Create empty temporary directory
	tmpDir, err := os.MkdirTemp("", "dsh_empty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test completion in empty directory
	input := "cat " + tmpDir + "/"
	matches, commonPrefix := completion.Complete(input, len(input))

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches in empty directory, got %d", len(matches))
	}

	if commonPrefix != "" {
		t.Errorf("Expected empty common prefix, got '%s'", commonPrefix)
	}
}

func TestHiddenFileCompletion(t *testing.T) {
	completion := readline.NewCompletion()

	// Create temporary directory with hidden files
	tmpDir, err := os.MkdirTemp("", "dsh_hidden_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create hidden and regular files
	files := []string{".hidden1", ".hidden2", "visible1", "visible2"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	// Test completion without dot - should not show hidden files
	input := "cat " + tmpDir + "/"
	matches, _ := completion.Complete(input, len(input))

	hiddenCount := 0
	for _, match := range matches {
		if filepath.Base(match.Text)[0] == '.' {
			hiddenCount++
		}
	}

	if hiddenCount > 0 {
		t.Errorf("Expected 0 hidden files without dot prefix, got %d", hiddenCount)
	}

	// Test completion with dot - should show hidden files
	input = "cat " + filepath.Join(tmpDir, ".h")
	matches, _ = completion.Complete(input, len(input))

	hiddenCount = 0
	for _, match := range matches {
		if filepath.Base(match.Text)[0] == '.' {
			hiddenCount++
		}
	}

	if hiddenCount != 2 {
		t.Errorf("Expected 2 hidden files with dot prefix, got %d", hiddenCount)
	}
}
