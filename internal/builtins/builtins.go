// Package builtins implements built-in shell commands like cd, pwd, help, and exit.
package builtins

import (
	"fmt"
	"os"
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
}

// IsBuiltin checks if a command is a built-in.
func IsBuiltin(name string) bool {
	_, exists := builtinCommands[name]

	return exists
}

// ExecuteBuiltin executes a built-in command.
func ExecuteBuiltin(args []string) bool {
	if len(args) == 0 {
		return true
	}

	if fn, exists := builtinCommands[args[0]]; exists {
		return fn(args)
	}

	return true
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
	_, _ = fmt.Fprintln(os.Stdout, "Built-in commands: cd, exit, help, pwd")
	_, _ = fmt.Fprintln(os.Stdout, "Features: quotes, pipes, I/O redirection, emacs-like editing")

	return true
}
