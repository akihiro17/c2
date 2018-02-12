package ast

import (
	"bytes"
	"c2/token"
	"fmt"
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
	fmt.Println("program string")
	fmt.Println(p.Func)
	return p.Func.String()
}

func (p *Program) Compile(out io.Writer) {
	p.Func.Compile(out)
}

type SimpleFunction struct {
	Token      token.Token
	Name       *Identifier
	Statements []Statement
}

func (sf *SimpleFunction) functionNode() {}
func (sf *SimpleFunction) TokenLiteral() string {
	return sf.Token.Literal
}
func (sf *SimpleFunction) String() string {
	fmt.Println("simple func")
	var out bytes.Buffer

	out.WriteString(sf.Token.Literal + " ")
	out.WriteString(sf.Name.Value)
	out.WriteString("()")
	out.WriteString("{")
	for _, stmt := range sf.Statements {
		out.WriteString(stmt.String())
	}
	out.WriteString("}")

	return out.String()
}
func (sf *SimpleFunction) Compile(out io.Writer) {
	out.Write([]byte(".globl " + sf.Name.Value))
	out.Write([]byte("\n"))
	out.Write([]byte(sf.Name.Value + ":"))
	out.Write([]byte("\n"))
	for _, stmt := range sf.Statements {
		stmt.Compile(out)
	}
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}
func (i *Identifier) Compile(out io.Writer) {
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

type IntAssignmentStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (r *IntAssignmentStatement) statementNode() {}
func (r *IntAssignmentStatement) TokenLiteral() string {
	return r.Token.Literal
}
func (r *IntAssignmentStatement) String() string {
	return r.Token.Literal + " " + r.Value.String() + ";"
}
func (sf *IntAssignmentStatement) Compile(out io.Writer) {
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
	switch pe.Operator {
	case "+":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("addq %rcx, %rax\n"))
	case "*":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("imulq %rcx\n"))
	case "-":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("subq %rax, %rcx\n"))
		out.Write([]byte("movq %rcx, %rax\n"))
	case "/":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rcx\n"))
		out.Write([]byte("popq %rax\n"))
		out.Write([]byte("movq $0, %rdx\n"))
		out.Write([]byte("idivq %rcx\n"))
	case ">":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rbx\n"))
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("cmpq %rbx, %rcx\n"))
		out.Write([]byte("setg %al\n"))
	case ">=":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rbx\n"))
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("cmpq %rbx, %rcx\n"))
		out.Write([]byte("setge %al\n"))
	case "<":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rbx\n"))
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("cmpq %rbx, %rcx\n"))
		out.Write([]byte("setl %al\n"))
	case "<=":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rbx\n"))
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("cmpq %rbx, %rcx\n"))
		out.Write([]byte("setle %al\n"))
	case "==":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rbx\n"))
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("cmpq %rbx, %rcx\n"))
		out.Write([]byte("sete %al\n"))
	case "!=":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("movq %rax, %rbx\n"))
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("cmpq %rbx, %rcx\n"))
		out.Write([]byte("setne %al\n"))
	case "||":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("popq %rcx\n"))
		out.Write([]byte("orq %rcx, %rax\n"))
		out.Write([]byte("movq $0, %rax\n"))
		out.Write([]byte("setne %al\n"))
	case "&&":
		pe.Left.Compile(out)
		out.Write([]byte("pushq %rax\n"))
		pe.Right.Compile(out)
		out.Write([]byte("popq %rcx\n"))
		// set cl = 1 if e1 != 0
		out.Write([]byte("cmpq $0, %rcx\n"))
		out.Write([]byte("setne %cl\n"))
		// set al = 1 if e1 != 0
		out.Write([]byte("cmpq $0, %rax\n"))
		out.Write([]byte("setne %al\n"))
		// compute al & cl
		// store it in al
		out.Write([]byte("andb %cl, %al\n"))
	default:
	}
}
