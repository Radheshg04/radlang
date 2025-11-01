// usr/bin/go run $0 $@ ; exit
package main

import "unicode"

type TokenType string

type Token struct {
	Token  TokenType
	lexeme string
	line   int
}

type lexer struct {
	input   string // input string to the lexer
	pos     int    // position of pointer in input
	readPos int    // next Token
	ch      byte   // current character
}

const (
	// Keywords
	FUNC   TokenType = "func"
	VAR    TokenType = "var"
	INT    TokenType = "type"
	FLOAT  TokenType = "type"
	BOOL   TokenType = "type"
	STRING TokenType = "type"
	PRINT  TokenType = "print"

	IDENTIFIER     TokenType = "identifier"
	NUMBER         TokenType = "number"
	STRING_LITERAL TokenType = "string_literal"

	ILLEGAL TokenType = "illegal"
	EOF     TokenType = "eof"
	EOL     TokenType = "eol"

	// Operators
	ASSIGNMENT TokenType = "="
	PLUS       TokenType = "+"
	MINUS      TokenType = "-"
	ASTERISK   TokenType = "*"
	SLASH      TokenType = "/"

	PLUSPLUS   TokenType = "++"
	MINUSMINUS TokenType = "--"

	L_PAREN TokenType = "("
	R_PAREN TokenType = ")"
	L_BRACE TokenType = "{"
	R_BRACE TokenType = "}"
)

var keywords = map[string]TokenType{
	"func":   FUNC,
	"var":    VAR,
	"int":    INT,
	"float":  FLOAT,
	"bool":   BOOL,
	"string": STRING,
}

func (l *lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *lexer) advance() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *lexer) isDigit(readPos int) bool {
	ch := l.input[readPos]
	return ch >= '0' && ch <= '9'
}

func (l *lexer) get_lexeme(offset int) string {
	return l.input[l.pos : l.readPos+offset]
}

func Lex(code string) []Token {
	var L lexer
	var tokenStream []Token
	line := 1
	L.input = code
	L.pos, L.readPos = 0, 1
	L.ch = L.input[L.pos]
	for {
		// fmt.Println(L.pos, L.readPos, L.ch)
		if L.pos >= len(code) || L.ch == 0 {
			tokenStream = append(tokenStream, Token{lexeme: "", Token: EOF, line: line})
			break
		}
		L.ch = L.input[L.pos]
		if L.ch == ' ' || L.ch == '\t' {
			L.advance()
			continue
		}
		if L.ch == '\n' {
			tokenStream = append(tokenStream, Token{lexeme: "\\n", Token: EOL, line: line})
			line++
			L.advance()
			continue
		}

		if L.ch == '"' {
			L.pos++
			for {
				if L.input[L.readPos] == '"' {
					break
				}
				L.readPos++
			}
			// make a token of type string literal which starts at L.pos and ends at L.readPos
			// TODO: Add support for escape sequences
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: STRING_LITERAL, line: line})
			L.readPos++
			L.advance()
			continue
		}
		switch L.ch {
		case '=':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: ASSIGNMENT, line: line})
			L.advance()
			continue
		case '+':
			if L.peek() == '+' {
				tokenStream = append(tokenStream, Token{lexeme: "++", Token: PLUSPLUS, line: line})
				L.readPos++
				L.advance()
			} else {
				tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: PLUS, line: line})
				L.advance()
				continue
			}

		case '-':
			if L.peek() == '-' {
				tokenStream = append(tokenStream, Token{lexeme: "--", Token: MINUSMINUS, line: line})
				L.readPos++
				L.advance()
			} else {
				tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: MINUS, line: line})
				L.advance()
				continue
			}
		case '*':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: ASTERISK, line: line})
			L.advance()
			continue
		case '/':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: SLASH, line: line})
			L.advance()
			continue
		case '(':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: L_PAREN, line: line})
			L.advance()
			continue
		case ')':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: R_PAREN, line: line})
			L.advance()
			continue
		case '{':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: L_BRACE, line: line})
			L.advance()
			continue
		case '}':
			tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: R_BRACE, line: line})
			L.advance()
			continue
		}
		// get identifiers/keywords
		// BUG: x++ being treated as an identifier, it should be lexed as x identifier and ++ INCREMENT
		if (L.ch >= 'a' && L.ch <= 'z') || (L.ch >= 'A' && L.ch <= 'Z') {
			for {
				next := L.peek()
				if L.readPos >= len(code) || unicode.IsDigit(rune(next)) || !unicode.IsLetter(rune(next)) {
					token, exists := keywords[L.get_lexeme(0)]
					if exists {
						tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: token, line: line})
						L.advance()
						break

					} else {
						tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: IDENTIFIER, line: line})
						L.advance()
						break

					}
				} else {
					L.readPos++
				}
			}
			continue
		}
		if L.isDigit(L.pos) {
			for {
				if !L.isDigit(L.readPos) {
					tokenStream = append(tokenStream, Token{lexeme: L.get_lexeme(0), Token: NUMBER, line: line})
					L.advance()
					break
				} else {
					L.readPos++
				}
			}
			continue
		}
		// TODO: Work on catchinbg illegal tokens

	}
	return tokenStream
}
