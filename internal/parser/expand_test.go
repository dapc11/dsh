package parser

import (
	"os/user"
	"testing"
)

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

func TestExpandTildeUser(t *testing.T) {
	// Get current user for testing
	currentUser, err := user.Current()
	if err != nil {
		t.Skip("Cannot get current user, skipping user expansion test")
	}

	username := currentUser.Username
	homeDir := currentUser.HomeDir

	tests := []struct {
		input    string
		expected string
	}{
		{"~" + username, homeDir},
		{"~" + username + "/file.txt", homeDir + "/file.txt"},
		{"~nonexistentuser", "~nonexistentuser"}, // Should remain unchanged
	}

	for _, test := range tests {
		result := expandTilde(test.input)
		if result != test.expected {
			t.Errorf("expandTilde(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}
