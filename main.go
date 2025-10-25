// Package main implements Daniel's Shell (dsh) - a minimal POSIX-compatible shell.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"

	"dsh/internal/executor"
	"dsh/internal/lexer"
	"dsh/internal/parser"
	"dsh/internal/readline"
)

func main() {
	var commandFlag = flag.String("c", "", "execute command and exit")
	flag.Parse()

	// If -c flag is provided, execute command and exit
	if *commandFlag != "" {
		success := processCommandLine(*commandFlag)
		if !success {
			os.Exit(1)
		}
		// Exit with last command's exit status
		os.Exit(executor.GetLastExitStatus())
	}

	// Interactive mode
	if isatty.IsTerminal(os.Stdin.Fd()) {
		rl, err := readline.New("dsh> ")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "dsh: failed to initialize readline: %v\n", err)
			os.Exit(1)
		}

		for {
			line, err := rl.ReadLine()
			if err != nil {
				if errors.Is(err, readline.ErrEOF) {
					// Ctrl+D pressed on empty line - exit gracefully
					break
				}
				_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)

				break
			}

			if line == "" {
				continue
			}

			if !processCommandLine(line) {
				// Command returned false (likely exit command)
				os.Exit(0)
			}
		}
	} else {
		// Non-interactive mode - read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if !processCommandLine(line) {
				// Command returned false (likely exit command)
				os.Exit(0)
			}
		}
	}

	// Exit with last command's exit status
	os.Exit(executor.GetLastExitStatus())
}

func processCommandLine(line string) bool {
	l := lexer.New(line)
	p := parser.New(l)

	pipelines, err := p.ParseCommandLine()
	if err != nil {
		if errors.Is(err, parser.ErrEmptyPipeline) {
			// Skip empty pipelines, continue parsing
			return true
		}
		_, _ = fmt.Fprintf(os.Stderr, "dsh: %v\n", err)

		return true
	}

	for _, pipeline := range pipelines {
		if !executor.ExecutePipeline(pipeline) {
			return false
		}
	}

	return true
}
