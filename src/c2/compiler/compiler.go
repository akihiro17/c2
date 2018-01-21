package compiler

import (
	"c2/lexer"
	"c2/parser"
	"io"
)

type Compiler struct {
	l *lexer.Lexer
	p *parser.Parser
}

func New(code string) *Compiler {
	c := &Compiler{}
	c.l = lexer.New(code)
	c.p = parser.New(c.l)

	return c
}

func (c *Compiler) Compile(out io.Writer) {
	program := c.p.ParserProgram()
	program.Compile(out)
}
