package lexer

import (
	"testing"
)

func TestLexer_SimpleCommand(t *testing.T) {
	input := "echo hello"
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	expected := []TokenType{Word, Word, EOF}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token.Type != expected[i] {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], token.Type)
		}
	}

	if tokens[0].Value != "echo" {
		t.Errorf("First token value: expected 'echo', got '%s'", tokens[0].Value)
	}
	if tokens[1].Value != "hello" {
		t.Errorf("Second token value: expected 'hello', got '%s'", tokens[1].Value)
	}
}

func TestLexer_Redirection(t *testing.T) {
	input := "cat < input.txt > output.txt"
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	expected := []TokenType{Word, RedirectIn, Word, RedirectOut, Word, EOF}
	for i, token := range tokens {
		if i < len(expected) && token.Type != expected[i] {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], token.Type)
		}
	}
}

func TestLexer_Quotes(t *testing.T) {
	input := `echo "hello world" 'single quotes'`
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	if len(tokens) < 3 {
		t.Fatal("Expected at least 3 tokens")
	}

	// Check quoted strings are handled properly
	if tokens[1].Value != "hello world" {
		t.Errorf("Double quoted string: expected 'hello world', got '%s'", tokens[1].Value)
	}
	if tokens[2].Value != "single quotes" {
		t.Errorf("Single quoted string: expected 'single quotes', got '%s'", tokens[2].Value)
	}
}

func TestLexer_Background(t *testing.T) {
	input := "sleep 5 &"
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	expected := []TokenType{Word, Word, Background, EOF}
	for i, token := range tokens {
		if i < len(expected) && token.Type != expected[i] {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], token.Type)
		}
	}
}

func TestLexer_Semicolon(t *testing.T) {
	input := "echo hello; echo world"
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	expected := []TokenType{Word, Word, Semicolon, Word, Word, EOF}
	for i, token := range tokens {
		if i < len(expected) && token.Type != expected[i] {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], token.Type)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	input := "echo hello # this is a comment"
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	// Comments should be ignored
	expected := []TokenType{Word, Word, EOF}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens (comment ignored), got %d", len(expected), len(tokens))
	}
}

func TestLexer_EmptyInput(t *testing.T) {
	input := ""
	lexer := New(input)

	token := lexer.NextToken()
	if token.Type != EOF {
		t.Errorf("Expected EOF for empty input, got %v", token.Type)
	}
}

func TestLexer_WhitespaceOnly(t *testing.T) {
	input := "   \t  \n  "
	lexer := New(input)

	token := lexer.NextToken()
	if token.Type != EOF {
		t.Errorf("Expected EOF for whitespace-only input, got %v", token.Type)
	}
}

func TestLexer_AppendRedirection(t *testing.T) {
	input := "echo hello >> output.txt"
	lexer := New(input)

	tokens := []Token{}
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	expected := []TokenType{Word, Word, RedirectAppend, Word, EOF}
	for i, token := range tokens {
		if i < len(expected) && token.Type != expected[i] {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], token.Type)
		}
	}
}
