package ast

import (
	"testing"

	"github.com/pqppq/writing-an-interpreter-in-go/monkey/token"
)

func TestString(t *testing.T) {
	statements := []Statement{
		&LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "foo"}, Value: "foo"},
			Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "100"}, Value: "100"},
		},
	}
	program := &Program{Statements: statements}
	expr := "let foo = 100;"

	if program.String() != expr {
		t.Errorf("expected expression `%s`, got %s instead", expr, program.String())
	}
}
