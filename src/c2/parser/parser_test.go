package parser

import (
	"bytes"
	"c2/ast"
	"c2/lexer"
	"c2/token"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	input := `
int main(){
  return 2;
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	fmt.Println(program.String())

	simple, _ := program.Func.(*ast.SimpleFunction)
	fmt.Println(simple.Name.Value)
	ret, _ := simple.Statements[0].(*ast.ReturnStatement)
	fmt.Println(ret.Value)

	if !testIntegerLiteral(t, ret.Value, "2") {
		return
	}

	if program == nil {
		t.Fatalf("returned nil")
	}

	if program.TokenLiteral() != token.INT {
		t.Fatalf("mismatch %g", program.TokenLiteral())
	}
}

func TestCompile(t *testing.T) {
	input := `
int main(){
  int a = 2;
  a = 2;
  return 2;
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()

	out := new(bytes.Buffer)

	program.Compile(out)
	fmt.Println(out.String())
}

func TestParsePrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    string
	}{
		{"!5", "!", "5"},
		{"-15", "-", "15"},
		{"~1", "~", "1"},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		expression := p.ParseExpression(LOWEST)
		prefix := expression.(*ast.PrefixExpression)
		if prefix.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, prefix.Operator)
		}
		if !testIntegerLiteral(t, prefix.Right, tt.value) {
			return
		}
	}

}

func TestParseInfixExpression(t *testing.T) {
	prefixTests := []struct {
		input      string
		leftValue  string
		operator   string
		rightValue string
	}{
		{"5 + 5;", "5", "+", "5"},
		{"5 - 5;", "5", "-", "5"},
		{"5 * 5;", "5", "*", "5"},
		{"5 / 5;", "5", "/", "5"},
		{"5 < 5;", "5", "<", "5"},
		{"5 > 5;", "5", ">", "5"},
		{"5 == 5;", "5", "==", "5"},
		{"5 != 5;", "5", "!=", "5"},
		{"5 <= 5;", "5", "<=", "5"},
		{"5 >= 5;", "5", ">=", "5"},
		{"5 && 5;", "5", "&&", "5"},
		{"5 || 5;", "5", "||", "5"},
	}

	for _, tt := range prefixTests {
		fmt.Println(tt.input)
		l := lexer.New(tt.input)
		p := New(l)
		expression := p.ParseExpression(LOWEST)
		infix := expression.(*ast.InfixExpression)
		if !testIntegerLiteral(t, infix.Left, tt.leftValue) {
			return
		}
		if infix.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, infix.Operator)
		}
		if !testIntegerLiteral(t, infix.Right, tt.rightValue) {
			return
		}
	}

}

func TestOperatorPrecedence(t *testing.T) {
	prefixTests := []struct {
		input    string
		expected string
	}{
		{
			"-4 * 5",
			"((-4) * 5)",
		},
		{
			"!-1",
			"(!(-1))",
		},
		{
			"1 + 2 * 3",
			"(1 + (2 * 3))",
		},
		{
			"1 + 2 / 3",
			"(1 + (2 / 3))",
		},
		{
			"5 > 4 == 5 < 4",
			"((5 > 4) == (5 < 4))",
		},
		{
			"5 > 4 != 5 < 4",
			"((5 > 4) != (5 < 4))",
		},
		{
			"(1 + 2) * 3",
			"((1 + 2) * 3)",
		},
		{
			"5 >= 4 == 5 <= 4",
			"((5 >= 4) == (5 <= 4))",
		},
		{
			"5 >= 4 != 5 <= 4",
			"((5 >= 4) != (5 <= 4))",
		},
		{
			"5 && 4 || 5 && 4",
			"((5 && 4) || (5 && 4))",
		},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		expression := p.ParseExpression(LOWEST)
		fmt.Println(expression.String())
		if expression.String() != tt.expected {
			t.Errorf("expected = %q, got = %q", tt.expected, expression.String())
		}
	}

}

func TestIntAssignment(t *testing.T) {
	input := "int a = 2;"
	l := lexer.New(input)
	p := New(l)

	stmt := p.ParseStatement()
	intAssignmentStatement := stmt.(*ast.IntAssignmentStatement)
	value := intAssignmentStatement.Value.(*ast.IntegerLiteral)
	if value.Value != "2" {
		t.Errorf("expected = %q, got = %q", 2, intAssignmentStatement.Value)
	}
}

func TestIntAssignmentNoInitialization(t *testing.T) {
	input := "int a;"
	l := lexer.New(input)
	p := New(l)

	stmt := p.ParseStatement()
	intAssignmentStatement := stmt.(*ast.IntAssignmentStatement)
	if intAssignmentStatement.Name.Value != "a" {
		t.Errorf("expected = %q, got = %q", "a", intAssignmentStatement.Name.Value)
	}
	if intAssignmentStatement.Value != nil {
		t.Errorf("expected = %q, got = %q", nil, intAssignmentStatement.Value)
	}
}

func TestAssignmentExpression(t *testing.T) {
	assignmentTests := []struct {
		input    string
		expected string
	}{
		{
			"a = 2 + 2;",
			"(a = (2 + 2))",
		},
		{
			"a = 2 + 2 * 2;",
			"(a = (2 + (2 * 2)))",
		},
		{
			"a = -2;",
			"(a = (-2))",
		},
		{
			"a = 2 * (b = 2)",
			"(a = (2 * (b = 2)))",
		},
		{
			"a = b = 2",
			"(a = (b = 2))",
		},
	}

	for _, tt := range assignmentTests {
		fmt.Println("start:", tt.input)
		l := lexer.New(tt.input)
		p := New(l)
		stmt := p.ParseExpression(LOWEST)

		if stmt.String() != tt.expected {
			t.Errorf("expected = %q, got = %q", tt.expected, stmt.String())
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value string) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %s. got=%s", value, integ.Value)
		return false
	}
	return true
}
