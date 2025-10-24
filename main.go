package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		_, _ = fmt.Fprint(os.Stdout, "dsh> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		if !processCommandLine(line) {
			break
		}
	}
}

func processCommandLine(line string) bool {
	lexer := NewLexer(line)
	parser := NewParser(lexer)

	pipelines, err := parser.ParseCommandLine()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)

		return true
	}

	for _, pipeline := range pipelines {
		if !executePipeline(pipeline) {
			return false
		}
	}

	return true
}

func executePipeline(pipeline *Pipeline) bool {
	if len(pipeline.Commands) == 1 {
		return executeCommand(pipeline.Commands[0])
	}

	return executeMultiCommandPipeline()
}

func executeCommand(cmd *Command) bool {
	if len(cmd.Args) == 0 {
		return true
	}

	// Handle built-in commands
	switch cmd.Args[0] {
	case "exit":
		return false
	case "cd":
		return handleCD(cmd.Args)
	case "help":
		_, _ = fmt.Fprintln(os.Stdout, "dsh - Daniel's Shell")
		_, _ = fmt.Fprintln(os.Stdout, "Built-in commands: cd, exit, help")
		_, _ = fmt.Fprintln(os.Stdout, "Features: quotes, pipes, I/O redirection")

		return true
	case "pwd":
		return handlePWD()
	default:
		return executeExternal(cmd)
	}
}

func handlePWD() bool {
	pwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: pwd: %v\n", err)

		return true
	}

	_, _ = fmt.Fprintln(os.Stdout, pwd)

	return true
}

func handleCD(args []string) bool {
	var target string
	if len(args) < 2 {
		target = os.Getenv("HOME")
		if target == "" {
			_, _ = fmt.Fprintf(os.Stderr, "dsh: cd: HOME not set\n")

			return true
		}
	} else {
		target = args[1]
	}

	err := os.Chdir(target)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dsh: cd: %v\n", err)
	}

	return true
}

func executeExternal(cmd *Command) bool {
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
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				// Exit with the same status as the child process
				os.Exit(status.ExitStatus())
			}
		}
		_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
	}

	return true
}

func executeMultiCommandPipeline() bool {
	_, _ = fmt.Fprintf(os.Stderr, "dsh: multi-command pipelines not yet implemented\n")

	return true
}
