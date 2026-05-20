package token

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	EOL
	// Keywords
	FUNC
	VAR
	PRINT
	INT
	FLOAT
	BOOL
	STRING

	IDENTIFIER
	NUMBER
	INT_LIT
	FLOAT_LIT
	STRING_LITERAL

	// Operators
	ASSIGNMENT
	PLUS
	MINUS
	ASTERISK
	SLASH

	PLUSPLUS
	MINUSMINUS

	L_PAREN
	R_PAREN
	L_BRACE
	R_BRACE
)

type Token struct {
	Token  TokenType
	Lexeme string
	Line   int
}

var Keywords = map[string]TokenType{
	"func":   FUNC,
	"var":    VAR,
	"int":    INT,
	"float":  FLOAT,
	"bool":   BOOL,
	"string": STRING,
}

func (t Token) String() string {
	return fmt.Sprintf("%-4d %-16s %q", t.Line, t.Token, t.Lexeme)
}

func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case EOL:
		return "EOL"
	case FUNC:
		return "FUNC"
	case VAR:
		return "VAR"
	case PRINT:
		return "PRINT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case INT_LIT:
		return "INT_LIT"
	case FLOAT_LIT:
		return "FLOAT_LIT"
	case BOOL:
		return "BOOL"
	case STRING:
		return "STRING"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NUMBER:
		return "NUMBER"
	case STRING_LITERAL:
		return "STRING_LITERAL"
	case ASSIGNMENT:
		return "ASSIGNMENT"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case PLUSPLUS:
		return "PLUSPLUS"
	case MINUSMINUS:
		return "MINUSMINUS"
	case L_PAREN:
		return "L_PAREN"
	case R_PAREN:
		return "R_PAREN"
	case L_BRACE:
		return "L_BRACE"
	case R_BRACE:
		return "R_BRACE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}
