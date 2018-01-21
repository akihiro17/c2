package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGALE"
	EOF     = "EOF"

	LBRACE = "{"
	RBRACE = "{"

	LPAREN = "("
	RPAREN = ")"

	SEMICOLOM = ";"

	INT = "int"

	RETURN = "return"

	IDENT       = "IDENT"
	INT_LITERAL = "INT_LITERAL"

	BITWISE_COMPLEMENT = "~"
	LOGICAL_NEGATION   = "!"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	GT       = ">"
	LT       = "<"
	EQ       = "=="
	NOT_EQ   = "!="
)

var keywords = map[string]TokenType{
	"int":    INT,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
