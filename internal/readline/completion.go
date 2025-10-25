package readline

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Completion handles tab completion for commands and files.
type Completion struct {
	commands []string
}

// NewCompletion creates a new completion instance.
func NewCompletion() *Completion {
	c := &Completion{
		commands: make([]string, 0),
	}
	c.loadCommands()
	return c
}

// CompletionItem represents a completion with its type.
type CompletionItem struct {
	Text string
	Type string // "builtin", "command", "file", "directory"
}

// Complete performs tab completion.
func (c *Completion) Complete(input string, _ int) ([]CompletionItem, string) {
	if input == "" {
		return nil, ""
	}

	// Split input into words
	words := strings.Fields(input)
	if len(words) == 0 {
		return nil, ""
	}

	// Determine what we're completing
	if len(words) == 1 && !strings.HasSuffix(input, " ") {
		// Completing command name
		return c.completeCommand(words[0])
	}

	// Completing file/directory name
	lastWord := ""
	if len(words) > 0 {
		if strings.HasSuffix(input, " ") {
			lastWord = ""
		} else {
			lastWord = words[len(words)-1]
		}
	}

	return c.completeFile(lastWord)
}

// loadCommands loads available commands from PATH and builtins.
func (c *Completion) loadCommands() {
	var commands []string

	// Add builtin commands
	builtinNames := []string{"cd", "pwd", "help", "exit", "todo"}
	commands = append(commands, builtinNames...)

	// Add commands from PATH
	pathCommands := c.getPathCommands()
	commands = append(commands, pathCommands...)

	// Sort and deduplicate
	sort.Strings(commands)
	c.commands = c.deduplicate(commands)
}

// getPathCommands gets executable commands from PATH.
func (c *Completion) getPathCommands() []string {
	var commands []string
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return commands
	}

	paths := strings.Split(pathEnv, ":")
	for _, path := range paths {
		if path == "" {
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				info, err := entry.Info()
				if err == nil && info.Mode()&0111 != 0 { // Executable
					commands = append(commands, entry.Name())
				}
			}
		}
	}

	return commands
}

// completeCommand completes command names.
func (c *Completion) completeCommand(prefix string) ([]CompletionItem, string) {
	var matches []CompletionItem
	builtins := []string{"cd", "pwd", "help", "exit", "todo"}

	// Add builtin matches
	for _, cmd := range builtins {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, CompletionItem{Text: cmd, Type: itemTypeBuiltin})
		}
	}

	// Add command matches
	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd, prefix) {
			// Skip if already added as builtin
			isBuiltin := false
			for _, builtin := range builtins {
				if cmd == builtin {
					isBuiltin = true
					break
				}
			}
			if !isBuiltin {
				matches = append(matches, CompletionItem{Text: cmd, Type: itemTypeCommand})
			}
		}
	}

	if len(matches) == 1 {
		return matches, matches[0].Text[len(prefix):]
	}

	return matches, c.commonPrefixItems(matches, prefix)
}

// completeFile completes file and directory names.
func (c *Completion) completeFile(prefix string) ([]CompletionItem, string) {
	dir := "."
	filename := prefix

	if strings.Contains(prefix, "/") {
		// Handle trailing slash case - show all files in directory
		if strings.HasSuffix(prefix, "/") {
			dir = prefix
			filename = ""
		} else {
			// Partial path - get directory and filename to match
			dir = filepath.Dir(prefix)
			filename = filepath.Base(prefix)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, ""
	}

	var matches []CompletionItem
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files unless explicitly requested
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(filename, ".") {
			continue
		}

		if strings.HasPrefix(name, filename) {
			// For partial paths, we need to return the full path from the original prefix
			var displayText string
			if dir == "." {
				displayText = name
			} else {
				// Replace the filename part with the matched name
				displayText = filepath.Join(filepath.Dir(prefix), name)
			}

			if entry.IsDir() {
				matches = append(matches, CompletionItem{Text: displayText + "/", Type: itemTypeDirectory})
			} else {
				matches = append(matches, CompletionItem{Text: displayText, Type: "file"})
			}
		}
	}

	if len(matches) == 1 {
		// Calculate completion from the original prefix
		completion := matches[0].Text[len(prefix):]
		return matches, completion
	}

	return matches, c.commonPrefixItems(matches, prefix)
}

// commonPrefixItems finds the common prefix of CompletionItems.
func (c *Completion) commonPrefixItems(matches []CompletionItem, current string) string {
	if len(matches) == 0 {
		return ""
	}

	if len(matches) == 1 {
		return matches[0].Text[len(current):]
	}

	// Find common prefix among all matches
	prefix := matches[0].Text
	for _, match := range matches[1:] {
		for i := 0; i < len(prefix) && i < len(match.Text); i++ {
			if prefix[i] != match.Text[i] {
				prefix = prefix[:i]
				break
			}
		}
	}

	// Return the part that extends beyond the current input
	if len(prefix) > len(current) {
		return prefix[len(current):]
	}

	return ""
}

// deduplicate removes duplicate strings.
func (c *Completion) deduplicate(strs []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range strs {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}
