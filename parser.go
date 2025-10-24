package main

import (
	"errors"
)

var (
	ErrExpectedCommandAfterPipe    = errors.New("expected command after pipe")
	ErrExpectedFilenameAfterOut    = errors.New("expected filename after >")
	ErrExpectedFilenameAfterAppend = errors.New("expected filename after >>")
	ErrExpectedFilenameAfterIn     = errors.New("expected filename after <")
	ErrNoCommand                   = errors.New("no command found")
	ErrEmptyPipeline               = errors.New("empty pipeline")
	ErrNoTokens                    = errors.New("no tokens to parse")
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
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

// NewParser creates a new parser with the given lexer.
func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer:        lexer,
		currentToken: Token{Type: EOF, Value: ""},
		peekToken:    Token{Type: EOF, Value: ""},
	}
	parser.nextToken()
	parser.nextToken()

	return parser
}

// ParseCommandLine parses a complete command line into pipelines.
func (parser *Parser) ParseCommandLine() ([]*Pipeline, error) {
	var pipelines []*Pipeline

	for parser.currentToken.Type != EOF {
		pipeline, err := parser.parsePipeline()
		if err != nil {
			if errors.Is(err, ErrEmptyPipeline) {
				// Skip empty pipelines, continue parsing
				if parser.currentToken.Type == Semicolon {
					parser.nextToken()
				}

				continue
			}

			return nil, err
		}

		if pipeline != nil {
			pipelines = append(pipelines, pipeline)
		}

		if parser.currentToken.Type == Semicolon {
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

	for parser.currentToken.Type == Pipe {
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
	if parser.currentToken.Type != Word {
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

	if parser.currentToken.Type == Background {
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
		case Word:
			cmd.Args = append(cmd.Args, parser.currentToken.Value)
			parser.nextToken()
		case RedirectOut:
			err := parser.handleOutputRedirect(cmd, false)
			if err != nil {
				return err
			}
		case RedirectAppend:
			err := parser.handleOutputRedirect(cmd, true)
			if err != nil {
				return err
			}
		case RedirectIn:
			err := parser.handleInputRedirect(cmd)
			if err != nil {
				return err
			}
		case Pipe, Background, Semicolon, EOF:
			return nil
		}
	}

	return nil
}

func (parser *Parser) isCommandToken() bool {
	return parser.currentToken.Type == Word ||
		parser.currentToken.Type == RedirectOut ||
		parser.currentToken.Type == RedirectIn ||
		parser.currentToken.Type == RedirectAppend
}

func (parser *Parser) handleOutputRedirect(cmd *Command, appendMode bool) error {
	parser.nextToken()
	if parser.currentToken.Type != Word {
		if appendMode {
			return ErrExpectedFilenameAfterAppend
		}

		return ErrExpectedFilenameAfterOut
	}

	cmd.OutputFile = parser.currentToken.Value
	cmd.AppendMode = appendMode
	parser.nextToken()

	return nil
}

func (parser *Parser) handleInputRedirect(cmd *Command) error {
	parser.nextToken()
	if parser.currentToken.Type != Word {
		return ErrExpectedFilenameAfterIn
	}

	cmd.InputFile = parser.currentToken.Value
	parser.nextToken()

	return nil
}
