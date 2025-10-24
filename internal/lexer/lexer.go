// Package lexer provides tokenization of shell input with quote and escape handling.
package lexer

import (
	"errors"
	"strings"
)

// TokenType represents different types of shell tokens.
type TokenType int

const (
	// Word represents a command word or argument.
	Word TokenType = iota
	// Pipe represents the pipe operator |.
	Pipe
	// RedirectOut represents the output redirection operator >.
	RedirectOut
	// RedirectIn represents the input redirection operator <.
	RedirectIn
	// RedirectAppend represents the append redirection operator >>.
	RedirectAppend
	// Background represents the background operator &.
	Background
	// Semicolon represents the command separator ;.
	Semicolon
	// EOF represents end of file.
	EOF
)

// Token represents a lexical token with its type and value.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer tokenizes shell input.
type Lexer struct {
	input    string
	position int
	current  rune
}

var (
	// ErrUnexpectedEOF indicates unexpected end of file in quoted string.
	ErrUnexpectedEOF = errors.New("unexpected EOF in quoted string")
	// ErrUnterminatedString indicates an unterminated quoted string.
	ErrUnterminatedString = errors.New("unterminated quoted string")
)

// New creates a new lexer for the given input.
func New(input string) *Lexer {
	lexer := &Lexer{
		input:    input,
		position: 0,
		current:  0,
	}
	lexer.readChar()

	return lexer
}

// NextToken returns the next token from the input.
func (lexer *Lexer) NextToken() Token {
	lexer.skipWhitespace()

	switch lexer.current {
	case 0:
		return Token{Type: EOF, Value: ""}
	case '#':
		lexer.skipComment()

		return lexer.NextToken()
	case '|':
		lexer.readChar()

		return Token{Type: Pipe, Value: "|"}
	case ';':
		lexer.readChar()

		return Token{Type: Semicolon, Value: ";"}
	case '&':
		lexer.readChar()

		return Token{Type: Background, Value: "&"}
	case '>':
		if lexer.peekChar() == '>' {
			lexer.readChar()
			lexer.readChar()

			return Token{Type: RedirectAppend, Value: ">>"}
		}
		lexer.readChar()

		return Token{Type: RedirectOut, Value: ">"}
	case '<':
		lexer.readChar()

		return Token{Type: RedirectIn, Value: "<"}
	default:
		word := lexer.readWord()

		return Token{Type: Word, Value: word}
	}
}

func (lexer *Lexer) readChar() {
	if lexer.position >= len(lexer.input) {
		lexer.current = 0 // EOF
	} else {
		lexer.current = rune(lexer.input[lexer.position])
	}
	lexer.position++
}

func (lexer *Lexer) peekChar() rune {
	if lexer.position >= len(lexer.input) {
		return 0
	}

	return rune(lexer.input[lexer.position])
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.current == ' ' || lexer.current == '\t' || lexer.current == '\n' || lexer.current == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) skipComment() {
	for lexer.current != 0 && lexer.current != '\n' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readQuotedString(quote rune) (string, error) {
	var result strings.Builder
	lexer.readChar() // skip opening quote

	for lexer.current != 0 && lexer.current != quote {
		if quote == '"' && lexer.current == '\\' {
			lexer.readChar()
			if lexer.current == 0 {
				return "", ErrUnexpectedEOF
			}

			result.WriteRune(lexer.handleEscapeSequence())
		} else {
			result.WriteRune(lexer.current)
		}
		lexer.readChar()
	}

	if lexer.current != quote {
		return "", ErrUnterminatedString
	}
	lexer.readChar() // skip closing quote

	return result.String(), nil
}

func (lexer *Lexer) handleEscapeSequence() rune {
	switch lexer.current {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case '\\':
		return '\\'
	case '"':
		return '"'
	case '$':
		return '$'
	default:
		// For unknown escape sequences, return both backslash and character
		return lexer.current
	}
}

func (lexer *Lexer) readWord() string {
	var result strings.Builder

	for lexer.current != 0 && !isWhitespace(lexer.current) && !isSpecialChar(lexer.current) {
		switch lexer.current {
		case '\'', '"':
			quoted, err := lexer.readQuotedString(lexer.current)
			if err != nil {
				result.WriteRune(lexer.current)
				lexer.readChar()
			} else {
				result.WriteString(quoted)
			}
		case '\\':
			lexer.readChar()
			if lexer.current != 0 {
				result.WriteRune(lexer.current)
				lexer.readChar()
			}
		default:
			result.WriteRune(lexer.current)
			lexer.readChar()
		}
	}

	return result.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isSpecialChar(ch rune) bool {
	return ch == '|' || ch == '>' || ch == '<' || ch == ';' || ch == '&'
}
