package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("dsh> ")
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
		fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
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
	
	// Handle multi-command pipelines
	return executeMultiCommandPipeline(pipeline)
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
		fmt.Println("dsh - Daniel's Shell")
		fmt.Println("Built-in commands: cd, exit, help")
		fmt.Println("Features: quotes, pipes, I/O redirection")
		return true
	case "pwd":
		if pwd, err := os.Getwd(); err == nil {
			fmt.Println(pwd)
		} else {
			fmt.Fprintf(os.Stderr, "dsh: pwd: %v\n", err)
		}
		return true
	default:
		return executeExternal(cmd)
	}
}

func handleCD(args []string) bool {
	var target string
	if len(args) < 2 {
		target = os.Getenv("HOME")
		if target == "" {
			fmt.Fprintf(os.Stderr, "dsh: cd: HOME not set\n")
			return true
		}
	} else {
		target = args[1]
	}
	
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, "dsh: cd: %v\n", err)
	}
	return true
}

func executeExternal(cmd *Command) bool {
	execCmd := exec.Command(cmd.Args[0], cmd.Args[1:]...)
	
	// Handle I/O redirection
	if cmd.InputFile != "" {
		file, err := os.Open(cmd.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
			return true
		}
		defer file.Close()
		execCmd.Stdin = file
	} else {
		execCmd.Stdin = os.Stdin
	}
	
	if cmd.OutputFile != "" {
		var file *os.File
		var err error
		if cmd.AppendMode {
			file, err = os.OpenFile(cmd.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		} else {
			file, err = os.Create(cmd.OutputFile)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
			return true
		}
		defer file.Close()
		execCmd.Stdout = file
	} else {
		execCmd.Stdout = os.Stdout
	}
	
	execCmd.Stderr = os.Stderr
	
	if cmd.Background {
		if err := execCmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
		} else {
			fmt.Printf("[%d]\n", execCmd.Process.Pid)
		}
		return true
	}
	
	if err := execCmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		fmt.Fprintf(os.Stderr, "dsh: %v\n", err)
	}
	return true
}

func executeMultiCommandPipeline(pipeline *Pipeline) bool {
	// Simple pipe implementation for now
	fmt.Fprintf(os.Stderr, "dsh: multi-command pipelines not yet implemented\n")
	return true
}
