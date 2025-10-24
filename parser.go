package main

import (
	"fmt"
	"os"
)

type Command struct {
	Args        []string
	InputFile   string
	OutputFile  string
	AppendMode  bool
	Background  bool
}

type Pipeline struct {
	Commands []*Command
}

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseCommandLine() ([]*Pipeline, error) {
	var pipelines []*Pipeline
	
	for p.currentToken.Type != EOF {
		pipeline, err := p.parsePipeline()
		if err != nil {
			return nil, err
		}
		if pipeline != nil {
			pipelines = append(pipelines, pipeline)
		}
		
		if p.currentToken.Type == SEMICOLON {
			p.nextToken()
		}
	}
	
	return pipelines, nil
}

func (p *Parser) parsePipeline() (*Pipeline, error) {
	pipeline := &Pipeline{}
	
	cmd, err := p.parseCommand()
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, nil
	}
	
	pipeline.Commands = append(pipeline.Commands, cmd)
	
	for p.currentToken.Type == PIPE {
		p.nextToken()
		cmd, err := p.parseCommand()
		if err != nil {
			return nil, err
		}
		if cmd == nil {
			return nil, fmt.Errorf("expected command after pipe")
		}
		pipeline.Commands = append(pipeline.Commands, cmd)
	}
	
	return pipeline, nil
}

func (p *Parser) parseCommand() (*Command, error) {
	if p.currentToken.Type != WORD {
		return nil, nil
	}
	
	cmd := &Command{}
	
	for p.currentToken.Type == WORD || p.currentToken.Type == REDIRECT_OUT || 
		p.currentToken.Type == REDIRECT_IN || p.currentToken.Type == REDIRECT_APPEND {
		
		switch p.currentToken.Type {
		case WORD:
			cmd.Args = append(cmd.Args, p.currentToken.Value)
			p.nextToken()
		case REDIRECT_OUT:
			p.nextToken()
			if p.currentToken.Type != WORD {
				return nil, fmt.Errorf("expected filename after >")
			}
			cmd.OutputFile = p.currentToken.Value
			p.nextToken()
		case REDIRECT_APPEND:
			p.nextToken()
			if p.currentToken.Type != WORD {
				return nil, fmt.Errorf("expected filename after >>")
			}
			cmd.OutputFile = p.currentToken.Value
			cmd.AppendMode = true
			p.nextToken()
		case REDIRECT_IN:
			p.nextToken()
			if p.currentToken.Type != WORD {
				return nil, fmt.Errorf("expected filename after <")
			}
			cmd.InputFile = p.currentToken.Value
			p.nextToken()
		}
	}
	
	if p.currentToken.Type == BACKGROUND {
		cmd.Background = true
		p.nextToken()
	}
	
	if len(cmd.Args) == 0 {
		return nil, nil
	}
	
	return cmd, nil
}

func (p *Parser) expandVariables(arg string) string {
	// Simple variable expansion - just handle $HOME for now
	if arg == "$HOME" {
		if home := os.Getenv("HOME"); home != "" {
			return home
		}
	}
	return arg
}
