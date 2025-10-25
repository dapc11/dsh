package parser

import (
	"os"
	"os/user"
	"strings"
)

// expandTilde expands tilde (~) in paths according to POSIX rules.
func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	// Handle ~ alone or ~/...
	if path == "~" || strings.HasPrefix(path, "~/") {
		home := os.Getenv("HOME")
		if home == "" {
			return path // Return unchanged if HOME not set
		}
		return strings.Replace(path, "~", home, 1)
	}

	// Handle ~username or ~username/...
	slashIndex := strings.Index(path, "/")
	var username string
	if slashIndex == -1 {
		username = path[1:] // Everything after ~
	} else {
		username = path[1:slashIndex] // Between ~ and /
	}

	u, err := user.Lookup(username)
	if err != nil {
		return path // Return unchanged if user not found
	}

	return strings.Replace(path, "~"+username, u.HomeDir, 1)
}
