package ast

import (
	"bytes"
	"c2/token"
	_ "fmt"
	"strings"
	"testing"
)

func TestCompilePrefixLogical(t *testing.T) {
	prefix := &PrefixExpression{
		Token:    token.Token{Type: token.MINUS, Literal: "-"},
		Operator: "-",
		Right:    &IntegerLiteral{Token: token.Token{Type: token.INT_LITERAL, Literal: "5"}, Value: "5"},
	}
	out := new(bytes.Buffer)
	prefix.Compile(out)

	expected := [3]string{
		"movq $5, %rax",
		"neg %rax",
		"",
	}
	for i, line := range strings.Split(out.String(), "\n") {
		if expected[i] != line {
			t.Errorf("expected = %q. got = %q", expected[i], line)
		}
	}
}

func TestCompilePrefixLogicalNegation(t *testing.T) {
	prefix := &PrefixExpression{
		Token:    token.Token{Type: token.LOGICAL_NEGATION, Literal: "!"},
		Operator: "!",
		Right:    &IntegerLiteral{Token: token.Token{Type: token.INT_LITERAL, Literal: "5"}, Value: "5"},
	}

	out := new(bytes.Buffer)
	prefix.Compile(out)

	expected := [5]string{
		"movq $5, %rax",
		"cmpq $0, %rax",
		"movq $0, %rax",
		"sete %al",
		"",
	}
	for i, line := range strings.Split(out.String(), "\n") {
		if expected[i] != line {
			t.Errorf("expected = %q. got = %q", expected[i], line)
		}
	}
}

func TestCompilePlusOperator(t *testing.T) {
	prefix := &InfixExpression{
		Token:    token.Token{Type: token.PLUS, Literal: "+"},
		Operator: "+",
		Right:    &IntegerLiteral{Token: token.Token{Type: token.INT_LITERAL, Literal: "5"}, Value: "5"},
		Left:     &IntegerLiteral{Token: token.Token{Type: token.INT_LITERAL, Literal: "5"}, Value: "5"},
	}

	out := new(bytes.Buffer)
	prefix.Compile(out)

	expected := [6]string{
		"movq $5, %rax",
		"pushq %rax",
		"movq $5, %rax",
		"popq %rcx",
		"addq %rcx, %rax",
		"",
	}
	for i, line := range strings.Split(out.String(), "\n") {
		if expected[i] != line {
			t.Errorf("expected = %q. got = %q", expected[i], line)
		}
	}
}
