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
	c := &Completion{}
	c.loadCommands()
	return c
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

// CompletionItem represents a completion with its type.
type CompletionItem struct {
	Text string
	Type string // "builtin", "command", "file", "directory"
}

// Complete performs tab completion.
func (c *Completion) Complete(input string, cursor int) ([]CompletionItem, string) {
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

// completeCommand completes command names.
func (c *Completion) completeCommand(prefix string) ([]CompletionItem, string) {
	var matches []CompletionItem
	builtins := []string{"cd", "pwd", "help", "exit", "todo"}
	
	// Add builtin matches
	for _, cmd := range builtins {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, CompletionItem{Text: cmd, Type: "builtin"})
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
				matches = append(matches, CompletionItem{Text: cmd, Type: "command"})
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
		dir = filepath.Dir(prefix)
		filename = filepath.Base(prefix)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, ""
	}

	var matches []CompletionItem
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, filename) {
			if entry.IsDir() {
				matches = append(matches, CompletionItem{Text: name + "/", Type: "directory"})
			} else {
				matches = append(matches, CompletionItem{Text: name, Type: "file"})
			}
		}
	}

	if len(matches) == 1 {
		completion := matches[0].Text[len(filename):]
		return matches, completion
	}

	return matches, c.commonPrefixItems(matches, filename)
}

// commonPrefixItems finds the common prefix of CompletionItems.
func (c *Completion) commonPrefixItems(matches []CompletionItem, current string) string {
	if len(matches) == 0 {
		return ""
	}

	if len(matches) == 1 {
		return matches[0].Text[len(current):]
	}

	// Find common prefix
	prefix := matches[0].Text
	for _, match := range matches[1:] {
		for i := 0; i < len(prefix) && i < len(match.Text); i++ {
			if prefix[i] != match.Text[i] {
				prefix = prefix[:i]
				break
			}
		}
	}

	if len(prefix) > len(current) {
		return prefix[len(current):]
	}

	return ""
}

// commonPrefix finds the common prefix of matches.
func (c *Completion) commonPrefix(matches []string, current string) string {
	if len(matches) == 0 {
		return ""
	}

	if len(matches) == 1 {
		return matches[0][len(current):]
	}

	// Find common prefix
	prefix := matches[0]
	for _, match := range matches[1:] {
		for i := 0; i < len(prefix) && i < len(match); i++ {
			if prefix[i] != match[i] {
				prefix = prefix[:i]
				break
			}
		}
	}

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
