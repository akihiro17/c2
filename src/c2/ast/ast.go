package ast

import (
	"bytes"
	"c2/token"
	"io"
)

type Node interface {
	TokenLiteral() string
	String() string
	Compile(out io.Writer)
}

type Function interface {
	Node
	functionNode()
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
	Func Function
}

func (p *Program) TokenLiteral() string {
	return p.Func.TokenLiteral()
}

func (p *Program) String() string {
	return p.Func.String()
}

func (p *Program) Compile(out io.Writer) {
	p.Func.Compile(out)
}

type SimpleFunction struct {
	Token token.Token
	Name  *Identifier
	Value Statement
}

func (sf *SimpleFunction) functionNode() {}
func (sf *SimpleFunction) TokenLiteral() string {
	return sf.Token.Literal
}
func (sf *SimpleFunction) String() string {
	var out bytes.Buffer

	out.WriteString(sf.Token.Literal + " ")
	out.WriteString(sf.Name.Value)
	out.WriteString("()")
	out.WriteString("{")
	out.WriteString(sf.Value.String())
	out.WriteString("}")

	return out.String()
}
func (sf *SimpleFunction) Compile(out io.Writer) {
	out.Write([]byte(".globl " + sf.Name.Value))
	out.Write([]byte("\n"))
	out.Write([]byte(sf.Name.Value + ":"))
	out.Write([]byte("\n"))
	sf.Value.Compile(out)
}

type Identifier struct {
	Token token.Token
	Value string
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
func (r *ReturnStatement) String() string {
	return r.Token.Literal + " " + r.Value.String() + ";"
}
func (sf *ReturnStatement) Compile(out io.Writer) {
	sf.Value.Compile(out)
	out.Write([]byte("ret"))
	out.Write([]byte("\n"))
}

type IntegerLiteral struct {
	Token token.Token
	Value string
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return i.Value
}
func (i *IntegerLiteral) Compile(out io.Writer) {
	out.Write([]byte("movq $" + i.Value + ", " + "%rax"))
	out.Write([]byte("\n"))
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {
	return
}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
func (pe *PrefixExpression) Compile(out io.Writer) {
	pe.Right.Compile(out)
	switch pe.Operator {
	case "-":
		out.Write([]byte("neg %rax"))
	case "~":
		out.Write([]byte("not %rax"))
	case "!":
		out.Write([]byte("cmpq $0, %rax\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("sete %al"))
	default:
	}
	out.Write([]byte("\n"))
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (pe *InfixExpression) expressionNode() {
}
func (pe *InfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(" " + pe.Operator + " ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
func (pe *InfixExpression) Compile(out io.Writer) {
	pe.Left.Compile(out)
	out.Write([]byte("pushq %rax\n"))
	pe.Right.Compile(out)
	out.Write([]byte("popq %rcx\n"))
	out.Write([]byte("addq %rcx, %rax\n"))
}
