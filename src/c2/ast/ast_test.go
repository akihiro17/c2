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
		"movl $5, %eax",
		"neg %eax",
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
		"movl $5, %eax",
		"cmpl $0, %eax",
		"movl $0, %eax",
		"sete %al",
		"",
	}
	for i, line := range strings.Split(out.String(), "\n") {
		if expected[i] != line {
			t.Errorf("expected = %q. got = %q", expected[i], line)
		}
	}
}
