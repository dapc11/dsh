// Package executor handles command execution with I/O redirection and process management.
package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"dsh/internal/builtins"
	"dsh/internal/parser"
)

var (
	lastExitStatus int
	exitStatusMu   sync.RWMutex
)

// GetLastExitStatus returns the exit status of the last executed command.
func GetLastExitStatus() int {
	exitStatusMu.RLock()
	defer exitStatusMu.RUnlock()
	return lastExitStatus
}

// setExitStatus sets the exit status from a command execution.
func setExitStatus(err error) {
	exitStatusMu.Lock()
	defer exitStatusMu.Unlock()

	if err == nil {
		lastExitStatus = 0
		return
	}

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			lastExitStatus = status.ExitStatus()
			return
		}
	}

	// Default to 1 for other errors
	lastExitStatus = 1
}

// ExecuteCommand executes a single command.
func ExecuteCommand(cmd *parser.Command) bool {
	if len(cmd.Args) == 0 {
		setExitStatus(nil) // Empty command succeeds
		return true
	}

	// Handle built-in commands
	if builtins.IsBuiltin(cmd.Args[0]) {
		success := builtins.ExecuteBuiltin(cmd.Args)
		if success {
			setExitStatus(nil) // Builtin succeeded
		} else {
			lastExitStatus = 1 // Builtin failed (or exit command)
		}
		return success
	}

	// Execute external command
	return executeExternal(cmd)
}

// ExecutePipeline executes a pipeline of commands.
func ExecutePipeline(pipeline *parser.Pipeline) bool {
	if len(pipeline.Commands) == 1 {
		return ExecuteCommand(pipeline.Commands[0])
	}

	return executeMultiCommandPipeline()
}

func executeExternal(cmd *parser.Command) bool {
	ctx := context.Background()
	execCmd := exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...) //nolint:gosec

	// Handle I/O redirection
	if cmd.InputFile != "" {
		file, err := os.Open(cmd.InputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)

			return true
		}
		defer func() { _ = file.Close() }()

		execCmd.Stdin = file
	} else {
		execCmd.Stdin = os.Stdin
	}

	if cmd.OutputFile != "" {
		file, err := openOutputFile(cmd.OutputFile, cmd.AppendMode)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)

			return true
		}
		defer func() { _ = file.Close() }()

		execCmd.Stdout = file
	} else {
		execCmd.Stdout = os.Stdout
	}

	execCmd.Stderr = os.Stderr

	if cmd.Background {
		startBackgroundProcess(execCmd)

		return true
	}

	return runForegroundProcess(execCmd)
}

func openOutputFile(filename string, appendMode bool) (*os.File, error) {
	if appendMode {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600) //nolint:gosec
		if err != nil {
			return nil, fmt.Errorf("failed to open file for append %s: %w", filename, err)
		}

		return file, nil
	}

	file, err := os.Create(filename) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filename, err)
	}

	return file, nil
}

func startBackgroundProcess(execCmd *exec.Cmd) {
	err := execCmd.Start()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)

		return
	}

	_, _ = fmt.Fprintf(os.Stdout, "[%d]\n", execCmd.Process.Pid)
}

func runForegroundProcess(execCmd *exec.Cmd) bool {
	err := execCmd.Run()
	setExitStatus(err)

	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			// Not an exit error, print the error message
			_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
		}
		// Command failed, but continue processing (don't exit shell)
		return true
	}

	return true
}

func executeMultiCommandPipeline() bool {
	_, _ = fmt.Fprintf(os.Stderr, "dsh: multi-command pipelines not yet implemented\n")

	return true
}
