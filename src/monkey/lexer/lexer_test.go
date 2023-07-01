package lexer

import (
	"testing"

	"github.com/pqppq/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`
	cases := []struct {
		name            string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{"ASSIGN", token.ASSIGN, "="},
		{"PLUS", token.PLUS, "+"},
		{"LPAREN", token.LPAREN, "("},
		{"RPAREN", token.RPAREN, ")"},
		{"LBRACE", token.LBRACE, "{"},
		{"RBRACE", token.RBRACE, "}"},
		{"COMMA", token.COMMA, ","},
		{"SEMICOLON", token.SEMICOLON, ";"},
		{"EOF", token.EOF, ""},
	}

	l := New(input)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tok := l.NextToken()
			if tok.Type != tc.expectedType {
				t.Fatalf("Expected TokenType %q, got %q instead\n", tc.expectedType, tok.Type)
			}
			if tok.Literal != tc.expectedLiteral {
				t.Fatalf("Expected Literal %q, got %q instead\n", tc.expectedLiteral, tok.Literal)
			}
		})
	}
}
