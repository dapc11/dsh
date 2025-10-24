package main

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	WORD TokenType = iota
	PIPE
	REDIRECT_OUT
	REDIRECT_IN
	REDIRECT_APPEND
	BACKGROUND
	SEMICOLON
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input    string
	position int
	current  rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.current = 0 // EOF
	} else {
		l.current = rune(l.input[l.position])
	}
	l.position++
}

func (l *Lexer) peekChar() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readQuotedString(quote rune) (string, error) {
	var result strings.Builder
	l.readChar() // skip opening quote
	
	for l.current != 0 && l.current != quote {
		if quote == '"' && l.current == '\\' {
			// Handle escape sequences in double quotes
			l.readChar()
			if l.current == 0 {
				return "", fmt.Errorf("unexpected EOF in quoted string")
			}
			switch l.current {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case 'r':
				result.WriteRune('\r')
			case '\\':
				result.WriteRune('\\')
			case '"':
				result.WriteRune('"')
			case '$':
				result.WriteRune('$')
			default:
				result.WriteRune('\\')
				result.WriteRune(l.current)
			}
		} else {
			result.WriteRune(l.current)
		}
		l.readChar()
	}
	
	if l.current != quote {
		return "", fmt.Errorf("unterminated quoted string")
	}
	l.readChar() // skip closing quote
	
	return result.String(), nil
}

func (l *Lexer) readWord() string {
	var result strings.Builder
	
	for l.current != 0 && !isWhitespace(l.current) && !isSpecialChar(l.current) {
		if l.current == '\'' || l.current == '"' {
			quoted, err := l.readQuotedString(l.current)
			if err != nil {
				// For now, just include the quote as literal
				result.WriteRune(l.current)
				l.readChar()
			} else {
				result.WriteString(quoted)
			}
		} else if l.current == '\\' {
			l.readChar()
			if l.current != 0 {
				result.WriteRune(l.current)
				l.readChar()
			}
		} else {
			result.WriteRune(l.current)
			l.readChar()
		}
	}
	
	return result.String()
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	
	switch l.current {
	case 0:
		return Token{Type: EOF, Value: ""}
	case '#':
		l.skipComment()
		return l.NextToken()
	case '|':
		l.readChar()
		return Token{Type: PIPE, Value: "|"}
	case ';':
		l.readChar()
		return Token{Type: SEMICOLON, Value: ";"}
	case '&':
		l.readChar()
		return Token{Type: BACKGROUND, Value: "&"}
	case '>':
		if l.peekChar() == '>' {
			l.readChar()
			l.readChar()
			return Token{Type: REDIRECT_APPEND, Value: ">>"}
		}
		l.readChar()
		return Token{Type: REDIRECT_OUT, Value: ">"}
	case '<':
		l.readChar()
		return Token{Type: REDIRECT_IN, Value: "<"}
	default:
		word := l.readWord()
		return Token{Type: WORD, Value: word}
	}
}

func (l *Lexer) skipComment() {
	for l.current != 0 && l.current != '\n' {
		l.readChar()
	}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isSpecialChar(ch rune) bool {
	return ch == '|' || ch == '>' || ch == '<' || ch == ';' || ch == '&'
}
