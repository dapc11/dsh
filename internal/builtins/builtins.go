// Package builtins implements built-in shell commands like cd, pwd, help, and exit.
package builtins

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BuiltinCommand represents a built-in shell command.
type BuiltinCommand struct {
	Name string
	Func func([]string) bool
}

// builtinCommands maps command names to their implementations.
var builtinCommands = map[string]func([]string) bool{ //nolint:gochecknoglobals // Required for builtin command registry
	"cd":   handleCD,
	"pwd":  handlePWD,
	"help": handleHelp,
	"exit": handleExit,
	"todo": handleTodo,
}

// IsBuiltin checks if a command is a built-in.
func IsBuiltin(name string) bool {
	_, exists := builtinCommands[name]

	return exists
}

// ExecuteBuiltin executes a built-in command.
func ExecuteBuiltin(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if fn, exists := builtinCommands[args[0]]; exists {
		return fn(args)
	}

	return false
}

func handleExit(_ []string) bool {
	return false
}

func handleCD(args []string) bool {
	var target string
	if len(args) < 2 {
		target = os.Getenv("HOME")
		if target == "" {
			_, _ = fmt.Fprintf(os.Stderr, "dsh: cd: HOME not set\n")

			return false
		}
	} else {
		target = args[1]
	}

	err := os.Chdir(target)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: cd: %v\n", err)

		return false
	}

	return true
}

func handlePWD(_ []string) bool { //nolint:unparam // Always returns true as pwd command always succeeds
	pwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: pwd: %v\n", err)

		return true
	}

	_, _ = fmt.Fprintln(os.Stdout, pwd)

	return true
}

func handleHelp(_ []string) bool {
	_, _ = fmt.Fprintln(os.Stdout, "dsh - Daniel's Shell")
	_, _ = fmt.Fprintln(os.Stdout, "Built-in commands: cd, exit, help, pwd, todo")
	_, _ = fmt.Fprintln(os.Stdout, "Features: quotes, pipes, I/O redirection, emacs-like editing, history, autosuggestions")

	return true
}

func handleTodo(args []string) bool {
	if len(args) < 2 {
		// List todos
		todos := loadTodos()
		if len(todos) == 0 {
			_, _ = fmt.Fprintln(os.Stdout, "No todos found.")
		} else {
			_, _ = fmt.Fprintln(os.Stdout, "DSH Todo List:")
			for i, todo := range todos {
				_, _ = fmt.Fprintf(os.Stdout, "%d. %s\n", i+1, todo)
			}
		}
		return true
	}

	// Add new todo
	todoText := strings.Join(args[1:], " ")
	err := addTodo(todoText)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: todo: %v\n", err)
		return false
	}

	_, _ = fmt.Fprintf(os.Stdout, "Added todo: %s\n", todoText)
	return true
}

func loadTodos() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	todoFile := filepath.Join(homeDir, ".dsh_todos")
	file, err := os.Open(todoFile) //nolint:gosec // User's home directory is safe
	if err != nil {
		return nil // File doesn't exist yet
	}
	defer func() {
		_ = file.Close()
	}()

	var todos []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			todos = append(todos, line)
		}
	}

	return todos
}

func addTodo(todo string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	todoFile := filepath.Join(homeDir, ".dsh_todos")
	file, err := os.OpenFile(todoFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600) //nolint:gosec // User's home directory is safe
	if err != nil {
		return fmt.Errorf("failed to open todo file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = fmt.Fprintf(file, "[%s] %s\n", timestamp, todo)
	if err != nil {
		return fmt.Errorf("failed to write todo: %w", err)
	}
	return nil
}
