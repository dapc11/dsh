package parser

import (
	"testing"

	"dsh/internal/lexer"
)

func TestParser_SimpleCommand(t *testing.T) {
	input := "echo hello world"
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	pipeline := commands[0]
	if len(pipeline.Commands) != 1 {
		t.Errorf("Expected 1 command in pipeline, got %d", len(pipeline.Commands))
	}
	
	cmd := pipeline.Commands[0]
	if len(cmd.Args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(cmd.Args))
	}

	expected := []string{"echo", "hello", "world"}
	for i, arg := range cmd.Args {
		if arg != expected[i] {
			t.Errorf("Arg %d: expected '%s', got '%s'", i, expected[i], arg)
		}
	}
}

func TestParser_Redirection(t *testing.T) {
	input := "cat < input.txt > output.txt"
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	pipeline := commands[0]; cmd := pipeline.Commands[0]
	if cmd.InputFile != "input.txt" {
		t.Errorf("Expected input file 'input.txt', got '%s'", cmd.InputFile)
	}
	if cmd.OutputFile != "output.txt" {
		t.Errorf("Expected output file 'output.txt', got '%s'", cmd.OutputFile)
	}
}

func TestParser_Background(t *testing.T) {
	input := "sleep 5 &"
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	pipeline := commands[0]; cmd := pipeline.Commands[0]
	if !cmd.Background {
		t.Error("Expected background command")
	}
}

func TestParser_MultipleCommands(t *testing.T) {
	input := "echo hello; echo world"
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	if commands[0].Commands[0].Args[1] != "hello" {
		t.Errorf("First command: expected 'hello', got '%s'", commands[0].Commands[0].Args[1])
	}
	if commands[1].Commands[0].Args[1] != "world" {
		t.Errorf("Second command: expected 'world', got '%s'", commands[1].Commands[0].Args[1])
	}
}

func TestParser_AppendRedirection(t *testing.T) {
	input := "echo hello >> output.txt"
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	pipeline := commands[0]; cmd := pipeline.Commands[0]
	if cmd.OutputFile != "output.txt" {
		t.Errorf("Expected output file 'output.txt', got '%s'", cmd.OutputFile)
	}
	if !cmd.AppendMode {
		t.Error("Expected append output mode")
	}
}

func TestParser_EmptyInput(t *testing.T) {
	input := ""
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 0 {
		t.Errorf("Expected 0 commands for empty input, got %d", len(commands))
	}
}

func TestParser_QuotedArguments(t *testing.T) {
	input := `echo "hello world" 'single quotes'`
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	pipeline := commands[0]; cmd := pipeline.Commands[0]
	if len(cmd.Args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(cmd.Args))
	}

	if cmd.Args[1] != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", cmd.Args[1])
	}
	if cmd.Args[2] != "single quotes" {
		t.Errorf("Expected 'single quotes', got '%s'", cmd.Args[2])
	}
}

func TestParser_ComplexCommand(t *testing.T) {
	input := "cat < input.txt | grep pattern > output.txt &"
	l := lexer.New(input)
	p := New(l)

	commands, err := p.ParseCommandLine()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// This tests the parser's ability to handle complex syntax
	// The exact behavior depends on implementation
	if len(commands) == 0 {
		t.Error("Expected at least one command")
	}
}
