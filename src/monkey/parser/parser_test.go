package parser

import (
	"fmt"
	"testing"

	"github.com/pqppq/writing-an-interpreter-in-go/monkey/ast"
	"github.com/pqppq/writing-an-interpreter-in-go/monkey/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;

		let foo = 100;
	`
	program := getProgram(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d instead", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnsStatements(t *testing.T) {
	input := `
		return -1;
		return 10;
		return 100;
	`
	program := getProgram(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d instead", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected *ast.ReturnStatement, got %T instead", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("expected returnStmt.TokenLiteral() to be 'return', got %q instead", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got %T instead", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got %T instead", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("expected ident.Value to be 'foobar', got %q instead", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("expected ident.TokenLiteral() to be 'foobar', got %q instead", ident.Value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got %T instead", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T instead", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("expected literal.Value to be 5, got %d instead", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("expected literal.TokenLiteral() to be '5', got %q instead", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got %T instead", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expected *ast.PrefixExpression, got %T instead", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("expected operator to be '%s', got '%s' instead", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		operator   string
		leftValue  any
		rightValue any
	}{
		{"5 + 5;", "+", 5, 5},
		{"5 - 5;", "-", 5, 5},
		{"5 * 5;", "*", 5, 5},
		{"5 / 5;", "/", 5, 5},
		{"5 > 5;", ">", 5, 5},
		{"5 < 5;", "<", 5, 5},
		{"5 == 5;", "==", 5, 5},
		{"5 != 5;", "!=", 5, 5},
		{"foobar + barfoo;", "+", "foobar", "barfoo"},
		{"foobar - barfoo;", "-", "foobar", "barfoo"},
		{"foobar * barfoo;", "*", "foobar", "barfoo"},
		{"foobar / barfoo;", "/", "foobar", "barfoo"},
		{"foobar > barfoo;", ">", "foobar", "barfoo"},
		{"foobar < barfoo;", "<", "foobar", "barfoo"},
		{"foobar == barfoo;", "==", "foobar", "barfoo"},
		{"foobar != barfoo;", "!=", "foobar", "barfoo"},
		{"true == true", "==", true, true},
		{"true != false", "!=", false, false},
		{"false == false", "==", false, false},
	}

	for _, test := range tests {
		program := getProgram(t, test.input)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got %T instead", program.Statements[0])
		}

		if !ok {
			t.Fatalf("expected *ast.InfixExpression, got %T instead", stmt.Expression)
		}
		if !testInfixExpression(t, stmt.Expression, test.leftValue,
			test.operator, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "(a + (b + c))"},
		{"a + b - c", "(a + (b - c))"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(a + ((b * c) + ((d / e) - f)))"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)

		if actual := program.String(); actual != tt.expected {
			t.Errorf("expected %q, got %q instead", tt.expected, actual)
		}
	}
}

func getProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	return program
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("expected literal 'let', got %s instead", s.TokenLiteral())
		return false
	}

	stmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected *ast.LetStatement, got %T instead", s)
		return false
	}

	if stmt.Name.Value != name {
		t.Errorf("expected stmt.Name.Value '%s', got '%s' instead", name, stmt.Name.Value)
		return false
	}

	if stmt.Name.TokenLiteral() != name {
		t.Errorf("expected stmt.Name.TokenLiteral() '%s', got '%s' instead", name, stmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intg, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expected *ast.IntegerLiteral, got %T instead", il)
		return false
	}
	if intg.Value != value {
		t.Errorf("expected intg.Value to be %d, got %d instead", value, intg.Value)
		return false
	}
	if intg.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("expected intg.TokenLiteral() to be '%d', got '%s' instead", value, intg.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expected *ast.Identifier, got %T instead", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("expected ident.Value to be %s, got %s instead", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("expected ident.TokenLiteral() to be %s, got %s instead", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled, got %T instead", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expected *ast.InfixExpression, got %T instead", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("expected opExp.Operator to be '%s', got '%s' instead", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
