// Package parser converts tokens into command structures for execution.
package parser

import (
	"errors"

	"dsh/internal/lexer"
)

var (
	// ErrExpectedCommandAfterPipe indicates missing command after pipe operator.
	ErrExpectedCommandAfterPipe = errors.New("expected command after pipe")
	// ErrExpectedFilenameAfterOut indicates missing filename after > operator.
	ErrExpectedFilenameAfterOut = errors.New("expected filename after >")
	// ErrExpectedFilenameAfterAppend indicates missing filename after >> operator.
	ErrExpectedFilenameAfterAppend = errors.New("expected filename after >>")
	// ErrExpectedFilenameAfterIn indicates missing filename after < operator.
	ErrExpectedFilenameAfterIn = errors.New("expected filename after <")
	// ErrNoCommand indicates no command was found in input.
	ErrNoCommand = errors.New("no command found")
	// ErrEmptyPipeline indicates an empty pipeline.
	ErrEmptyPipeline = errors.New("empty pipeline")
	// ErrNoTokens indicates no tokens to parse.
	ErrNoTokens = errors.New("no tokens to parse")
)

// Command represents a single command with its arguments and redirections.
type Command struct {
	Args       []string
	InputFile  string
	OutputFile string
	AppendMode bool
	Background bool
}

// Pipeline represents a sequence of commands connected by pipes.
type Pipeline struct {
	Commands []*Command
}

// Parser parses tokens into command structures.
type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
}

// New creates a new parser with the given lexer.
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:        l,
		currentToken: lexer.Token{Type: lexer.EOF, Value: ""},
		peekToken:    lexer.Token{Type: lexer.EOF, Value: ""},
	}
	parser.nextToken()
	parser.nextToken()

	return parser
}

// ParseCommandLine parses a complete command line into pipelines.
func (parser *Parser) ParseCommandLine() ([]*Pipeline, error) {
	var pipelines []*Pipeline

	for parser.currentToken.Type != lexer.EOF {
		pipeline, err := parser.parsePipeline()
		if err != nil {
			if errors.Is(err, ErrEmptyPipeline) {
				// Skip empty pipelines, continue parsing
				if parser.currentToken.Type == lexer.Semicolon {
					parser.nextToken()
				}

				continue
			}

			return nil, err
		}

		if pipeline != nil {
			pipelines = append(pipelines, pipeline)
		}

		if parser.currentToken.Type == lexer.Semicolon {
			parser.nextToken()
		}
	}

	return pipelines, nil
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) parsePipeline() (*Pipeline, error) {
	pipeline := &Pipeline{
		Commands: []*Command{},
	}

	cmd, err := parser.parseCommand()
	if err != nil {
		if errors.Is(err, ErrNoTokens) {
			return nil, ErrEmptyPipeline // Use sentinel error instead of nil
		}

		return nil, err
	}

	if cmd == nil {
		return nil, ErrEmptyPipeline // Use sentinel error instead of nil
	}

	pipeline.Commands = append(pipeline.Commands, cmd)

	for parser.currentToken.Type == lexer.Pipe {
		parser.nextToken()

		cmd, err := parser.parseCommand()
		if err != nil {
			if errors.Is(err, ErrNoTokens) {
				return nil, ErrExpectedCommandAfterPipe
			}

			return nil, err
		}

		if cmd == nil {
			return nil, ErrExpectedCommandAfterPipe
		}

		pipeline.Commands = append(pipeline.Commands, cmd)
	}

	return pipeline, nil
}

func (parser *Parser) parseCommand() (*Command, error) {
	if parser.currentToken.Type != lexer.Word {
		return nil, ErrNoTokens
	}

	cmd := &Command{
		Args:       []string{},
		InputFile:  "",
		OutputFile: "",
		AppendMode: false,
		Background: false,
	}

	err := parser.processCommandTokens(cmd)
	if err != nil {
		return nil, err
	}

	if parser.currentToken.Type == lexer.Background {
		cmd.Background = true
		parser.nextToken()
	}

	if len(cmd.Args) == 0 {
		return nil, ErrNoCommand
	}

	return cmd, nil
}

func (parser *Parser) processCommandTokens(cmd *Command) error {
	for parser.isCommandToken() {
		switch parser.currentToken.Type {
		case lexer.Word:
			cmd.Args = append(cmd.Args, expandTilde(parser.currentToken.Value))
			parser.nextToken()
		case lexer.RedirectOut:
			err := parser.handleOutputRedirect(cmd, false)
			if err != nil {
				return err
			}
		case lexer.RedirectAppend:
			err := parser.handleOutputRedirect(cmd, true)
			if err != nil {
				return err
			}
		case lexer.RedirectIn:
			err := parser.handleInputRedirect(cmd)
			if err != nil {
				return err
			}
		case lexer.Pipe, lexer.Background, lexer.Semicolon, lexer.EOF:
			return nil
		}
	}

	return nil
}

func (parser *Parser) isCommandToken() bool {
	return parser.currentToken.Type == lexer.Word ||
		parser.currentToken.Type == lexer.RedirectOut ||
		parser.currentToken.Type == lexer.RedirectIn ||
		parser.currentToken.Type == lexer.RedirectAppend
}

func (parser *Parser) handleOutputRedirect(cmd *Command, appendMode bool) error {
	parser.nextToken()
	if parser.currentToken.Type != lexer.Word {
		if appendMode {
			return ErrExpectedFilenameAfterAppend
		}

		return ErrExpectedFilenameAfterOut
	}

	cmd.OutputFile = expandTilde(parser.currentToken.Value)
	cmd.AppendMode = appendMode
	parser.nextToken()

	return nil
}

func (parser *Parser) handleInputRedirect(cmd *Command) error {
	parser.nextToken()
	if parser.currentToken.Type != lexer.Word {
		return ErrExpectedFilenameAfterIn
	}

	cmd.InputFile = expandTilde(parser.currentToken.Value)
	parser.nextToken()

	return nil
}
