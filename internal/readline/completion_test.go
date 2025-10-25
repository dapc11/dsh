package readline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompletion_Commands(t *testing.T) {
	c := NewCompletion()

	// Test command completion
	matches, _ := c.Complete("ec", 0)

	// Should find "echo" command
	found := false
	for _, match := range matches {
		if match.Text == "echo" && match.Type == "command" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'echo' command in completion")
	}

	// Test builtin completion
	matches, _ = c.Complete("c", 0)
	found = false
	for _, match := range matches {
		if match.Text == "cd" && match.Type == itemTypeBuiltin {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'cd' builtin in completion")
	}
}

func TestCompletion_Files(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "testfile.txt")
	testDir := filepath.Join(tmpDir, "testdir")

	// Create test files
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Change to test directory
	t.Chdir(tmpDir)

	c := NewCompletion()

	// Test file completion after command
	matches, _ := c.Complete("cat test", 0)

	// Should find both file and directory
	foundFile := false
	foundDir := false
	for _, match := range matches {
		if match.Text == "testfile.txt" && match.Type == "file" {
			foundFile = true
		}
		if match.Text == "testdir/" && match.Type == "directory" {
			foundDir = true
		}
	}

	if !foundFile {
		t.Error("Expected to find test file in completion")
	}
	if !foundDir {
		t.Error("Expected to find test directory in completion")
	}
}

func TestCompletion_TrailingSpace(t *testing.T) {
	c := NewCompletion()

	// Test completion with trailing space (should complete files)
	matches, _ := c.Complete("cat ", 0)

	// Should return file completions, not command completions
	if len(matches) > 0 {
		// Check that we're not getting command completions
		for _, match := range matches {
			if match.Type == itemTypeBuiltin || match.Type == itemTypeCommand {
				t.Errorf("Expected file completion, got %s type", match.Type)
			}
		}
	}
}

func TestExpandTilde(t *testing.T) {
	// Set test HOME
	testHome := "/test/home"
	t.Setenv("HOME", testHome)

	tests := []struct {
		input    string
		expected string
	}{
		{"~/file.txt", testHome + "/file.txt"},
		{"~", testHome},
		{"~/dir/file", testHome + "/dir/file"},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
		{"file.txt", "file.txt"},
	}

	for _, test := range tests {
		result := expandTilde(test.input)
		if result != test.expected {
			t.Errorf("expandTilde(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestCompletion_Empty(t *testing.T) {
	c := NewCompletion()

	// Empty input should return nothing
	matches, completion := c.Complete("", 0)

	if len(matches) != 0 {
		t.Errorf("Expected no matches for empty input, got %d", len(matches))
	}

	if completion != "" {
		t.Errorf("Expected empty completion for empty input, got %s", completion)
	}
}
