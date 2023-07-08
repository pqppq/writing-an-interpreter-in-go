package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pqppq/writing-an-interpreter-in-go/monkey/token"
)

type Node interface {
	// returns the literal value of the token associated with
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// let <identifier> = <expression>;
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // <identifier>
	Value Expression  // <expression>
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	expr := fmt.Sprintf("let %s = ", ls.Name.String())

	if ls.Value != nil {
		expr += ls.Value.String()
	}
	expr += ";"

	out.WriteString(expr)
	return out.String()
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// return <expression>;
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	expr := fmt.Sprintf("return ")

	if rs.ReturnValue != nil {
		expr += rs.ReturnValue.String()
	}

	expr += ";"
	return out.String()
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer
	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}
	return out.String()
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (oe *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", oe.Left.String(), oe.Operator, oe.Right.String())
}
func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// if (<condition>) <consequence> else <alternative>
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) String() string {
	expr := fmt.Sprintf("if%s %s", ie.Condition.String(), ie.Consequence.String())
	if ie.Alternative != nil {
		expr += fmt.Sprintf("else %s", ie.Alternative.String())
	}

	return expr
}
func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// fn <parameters> <block statement>
// <parameters> = <parameter one>, <parameter two>, ...
type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	expr := fmt.Sprintf("fn(%s) %s", strings.Join(params, ", "), fl.Body.String())

	return expr
}
func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
