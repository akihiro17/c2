package lexer

import (
	"testing"
	"../token"
)

func TestNextToken(t *testing.T) {
	input := "{}();"

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLOM, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken();
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%g, got=%g", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%g, got=%g", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestSimpleMain(t *testing.T) {
	input := `
int main(){
  return 2;
}
        `

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.INT, "int"},
		{token.IDENT, "main"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT_LITERAL, "2"},
		{token.SEMICOLOM, ";"},
		{token.RBRACE, "}"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken();
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%g, got=%g", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%g, got=%g", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
