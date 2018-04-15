package parser

import (
	"c2/ast"
	"c2/lexer"
	"c2/token"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_OR_EQ: LESSGREATER,
	token.GT_OR_EQ: LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.ASSIGN:   ASSIGN,
}

const (
	_ int = iota
	LOWEST
	ASSIGN
	OR
	AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.INT_LITERAL, p.ParseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.ParsePrefixExpression)
	p.registerPrefix(token.BITWISE_COMPLEMENT, p.ParsePrefixExpression)
	p.registerPrefix(token.LOGICAL_NEGATION, p.ParsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.ParseGroupedExpression)
	p.registerPrefix(token.IDENT, p.ParseIdentifier)

	// infix
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.ParseInfixExpression)
	p.registerInfix(token.MINUS, p.ParseInfixExpression)
	p.registerInfix(token.ASTERISK, p.ParseInfixExpression)
	p.registerInfix(token.SLASH, p.ParseInfixExpression)
	p.registerInfix(token.GT, p.ParseInfixExpression)
	p.registerInfix(token.LT, p.ParseInfixExpression)
	p.registerInfix(token.GT_OR_EQ, p.ParseInfixExpression)
	p.registerInfix(token.LT_OR_EQ, p.ParseInfixExpression)
	p.registerInfix(token.EQ, p.ParseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.ParseInfixExpression)
	p.registerInfix(token.AND, p.ParseInfixExpression)
	p.registerInfix(token.OR, p.ParseInfixExpression)
	p.registerInfix(token.ASSIGN, p.ParseInfixExpression)
	return p
}

func (p *Parser) peekPrecedence() int {
	fmt.Println("peektoken:", p.peekToken)
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectCur(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.curError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Func = p.ParseFunction()

	return program
}

func (p *Parser) ParseFunction() ast.Function {
	sf := &ast.SimpleFunction{Token: p.curToken}

	p.nextToken()

	sf.Name = &ast.FunctionName{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	for p.curToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			sf.Statements = append(sf.Statements, stmt)
			p.nextToken()
		}
	}

	// TODO: invalid function
	if !p.expectCur(token.EOF) {
		return nil
	}

	return sf
}

func (p *Parser) ParseStatement() ast.Statement {
	fmt.Println(p.curToken)
	switch p.curToken.Type {
	case token.RETURN:
		returnStatement := &ast.ReturnStatement{Token: p.curToken}

		p.nextToken()
		returnStatement.Value = p.ParseExpression(LOWEST)

		// parseExpressionでセミコロンまで進む場合がある
		if p.expectCur(token.SEMICOLOM) {
			p.nextToken()
			return returnStatement
		}

		if !p.expectPeek(token.SEMICOLOM) {
			panic("semicolom")
		} else {
			p.nextToken()
		}

		return returnStatement
	case token.INT:
		fmt.Println("parse int")
		intAssignmentStatement := &ast.IntAssignmentStatement{Token: p.curToken}
		p.nextToken()

		intAssignmentStatement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(token.ASSIGN) {
			intAssignmentStatement.Value = nil
		} else {
			p.nextToken()
			fmt.Println("assign with value", p.curToken)
			intAssignmentStatement.Value = p.ParseExpression(LOWEST)
			fmt.Println("int assignment right value", intAssignmentStatement.Value)
		}

		if !p.expectPeek(token.SEMICOLOM) {
			return nil
		}

		return intAssignmentStatement
	// case token.IDENT:
	// 	fmt.Println("parse ident")
	// 	// 予約語チェック
	// 	if (p.curToken.Literal == "RETURN") {
	// 		panic("return")
	// 	}

	// 	identifierStatement := &ast.IdentifierStatement{Token: token.Token{Type: token.INT, Literal: "INT"}}
	// 	identifierStatement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 	if !p.expectPeek(token.ASSIGN) {
	// 		panic("need assign")
	// 	} else {
	// 		p.nextToken()
	// 		identifierStatement.Value = p.ParseExpression(LOWEST)
	// 		fmt.Println("int assignment right value")
	// 	}

	// 	if !p.expectPeek(token.SEMICOLOM) {
	// 		return nil
	// 	}

	// 	return identifierStatement
	default:
		exp := p.ParseExpression(LOWEST)

		expStatement := &ast.ExpressionStatement{Value: exp, Token: p.curToken}

		fmt.Println("exp:", exp)

		return expStatement
		// fmt.Println("not supported", p.curToken)
		// panic("not supported")

	}
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	fmt.Println("start parse expression")
	fmt.Println(p.curToken)
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	fmt.Println("before prefix", p.curToken)
	leftExp := prefix()
	fmt.Println("after prefix", p.curToken, "left:", leftExp)

	fmt.Println(precedence, p.peekPrecedence())
	for !p.peekTokenIs(token.SEMICOLOM) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return nil
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	fmt.Println("final left:", leftExp)

	return leftExp
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	_, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("does not Parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = p.curToken.Literal
	return lit

}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.ParseExpression(PREFIX)

	return expression
}

func (p *Parser) ParseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.ParseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) ParseIdentifier() ast.Expression {
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// p.nextToken()

	return identifier
}

func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	if expression.Operator == "=" {
		expression.Right = p.ParseExpression(precedence - 1)
	} else {
		expression.Right = p.ParseExpression(precedence)
	}

	return expression
}
