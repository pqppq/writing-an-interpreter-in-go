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
		{"true != false", "!=", true, false},
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
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)

		if actual := program.String(); actual != tt.expected {
			t.Errorf("expected %q, got %q instead", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expextedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
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

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("expected *ast.Boolan, got %T instead", stmt.Expression)
		}
		if boolean.Value != tt.expextedBoolean {
			t.Fatalf("expected %t, got %t instead", tt.expextedBoolean, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := getProgram(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.IfExpressionStatement, got %T instead", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expected *ast.IfExpression, got %T instead", stmt.Expression)
	}
	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		t.Fatalf("expected condition to be x < y, got %s instead", expr.Condition.String())
		return
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Fatalf("expected 1 consequence statement, got %d instead", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !testIdentifier(t, consequence.Expression, "x") {
		t.Fatalf("expected consequence to be x, got %s instead", consequence.Expression.String())
		return
	}
	if expr.Alternative != nil {
		t.Errorf("expected nil alternative, got %T instead", expr.Alternative)

	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got %T instead", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expected *ast.IfExpression, got %T instead", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements does not contain 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got %T instead", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expected *ast.FunctionLiteral, got %T instead", stmt.Expression)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d instead", len(fn.Parameters))
	}

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(fn.Body.Statements))
	}

	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got %T instead", fn.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fn := stmt.Expression.(*ast.FunctionLiteral)

		if len(fn.Parameters) != len(tt.expectedParams) {
			t.Errorf("expected %d parameters, got %d instead", len(tt.expectedParams), len(fn.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, fn.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d instead", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.CallExpression, got %T instead", stmt.Expression)
	}

	expr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got %T instead", stmt.Expression)
	}

	if !testIdentifier(t, expr.Function, "add") {
		return
	}

	if len(expr.Arguments) != 3 {
		fmt.Println(expr.Arguments)
		t.Fatalf("expected 3 arguments, got %d instead", len(expr.Arguments))
	}

	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expr.Arguments[2], 4, "+", 5)
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
	case bool:
		return testBooleanLiteral(t, exp, v)
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

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}
