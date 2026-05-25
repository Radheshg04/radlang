package token

import (
	"fmt"
)

type TokenType int

const (
	INVALID TokenType = -1
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
	ERR
	TYPE
	STRUCT
	INTERFACE

	IDENTIFIER
	INT_LIT
	FLOAT_LIT
	BOOL_LIT
	STRING_LIT
	ERROR_LITERAL
	COMMA

	// Operators
	ASSIGNMENT
	WALRUS
	PLUS
	MINUS
	ASTERISK
	SLASH

	// CMPR_OP
	EQ  // ==
	NEQ // !=
	GT  // >
	GTE // >=
	LT  // <
	LTE // <=

	TRUE
	FALSE

	PLUSPLUS
	MINUSMINUS

	IF
	ELSE
	RETURN
	FOR
	BREAK
	CONTINUE

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
	"func":      FUNC,
	"var":       VAR,
	"int":       INT,
	"float":     FLOAT,
	"bool":      BOOL,
	"string":    STRING,
	"err":       ERR,
	"type":      TYPE,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"for":       FOR,
	"break":     BREAK,
	"continue":  CONTINUE,
	"true":      TRUE,
	"false":     FALSE,
}

func (t Token) String() string {
	return fmt.Sprintf("%-4d %-16s %q", t.Line, t.Token, t.Lexeme)
}

func (t TokenType) String() string {
	switch t {
	case INVALID:
		return "INVALID"
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
	case ERR:
		return "ERR"
	case TYPE:
		return "TYPE"
	case STRUCT:
		return "STRUCT"
	case INTERFACE:
		return "INTERFACE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING_LIT:
		return "STRING_LIT"
	case BOOL_LIT:
		return "BOOL_LIT"
	case ERROR_LITERAL:
		return "ERROR_LITERAL"
	case COMMA:
		return "COMMA"
	case ASSIGNMENT:
		return "ASSIGNMENT"
	case WALRUS:
		return "WALRUS"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case GT:
		return "GT"
	case GTE:
		return "GTE"
	case LT:
		return "LT"
	case LTE:
		return "LTE"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case PLUSPLUS:
		return "PLUSPLUS"
	case MINUSMINUS:
		return "MINUSMINUS"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case FOR:
		return "FOR"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
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
