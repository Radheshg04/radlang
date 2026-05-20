// usr/bin/go run $0 $@ ; exit
package lexer

import (
	"unicode"

	"radlang/token"
)

type lexer struct {
	input   string // input string to the lexer
	pos     int    // position of pointer in input
	readPos int    // next Token
	ch      byte   // current character
}

func Lex(code string) []token.Token {
	var L lexer
	tokenStream := make([]token.Token, 0, len(code))
	line := 1
	L.input = code
	L.pos, L.readPos = 0, 1
	L.ch = L.input[L.pos]
	for {
		// EOF Handling
		if L.pos >= len(code) {
			tokenStream = append(tokenStream, token.Token{Lexeme: "", Token: token.EOF, Line: line})
			break
		}

		// Skip whitespaces and tabs
		if L.ch == ' ' || L.ch == '\t' {
			L.advance()
			continue
		}

		// Tokenize newlines
		if L.ch == '\n' {
			tokenStream = append(tokenStream, token.Token{Lexeme: "", Token: token.EOL, Line: line})
			line++
			L.advance()
			continue
		}

		// String handing
		if L.ch == '"' {
			for {
				if L.peek() == 0 {
					tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.ILLEGAL, Line: line})
					L.advance()
					break
				}
				if L.peek() == '"' {
					// make a token of type string literal which starts at L.pos and ends at L.readPos
					// TODO: Add support for escape sequences
					tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(1), Token: token.STRING_LITERAL, Line: line})
					L.readPos++
					L.advance()
					break
				}
				L.readPos++
			}
			continue
		}

		switch L.ch {
		case '=':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.ASSIGNMENT, Line: line})
			L.advance()
			continue
		case '+':
			if L.peek() == '+' {
				tokenStream = append(tokenStream, token.Token{Lexeme: "++", Token: token.PLUSPLUS, Line: line})
				L.readPos++
				L.advance()
				continue
			} else {
				tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.PLUS, Line: line})
				L.advance()
				continue
			}

		case '-':
			if L.peek() == '-' {
				tokenStream = append(tokenStream, token.Token{Lexeme: "--", Token: token.MINUSMINUS, Line: line})
				L.readPos++
				L.advance()
				continue
			} else {
				tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.MINUS, Line: line})
				L.advance()
				continue
			}
		case '*':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.ASTERISK, Line: line})
			L.advance()
			continue
		case '/':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.SLASH, Line: line})
			L.advance()
			continue
		case '(':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.L_PAREN, Line: line})
			L.advance()
			continue
		case ')':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.R_PAREN, Line: line})
			L.advance()
			continue
		case '{':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.L_BRACE, Line: line})
			L.advance()
			continue
		case '}':
			tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.R_BRACE, Line: line})
			L.advance()
			continue
		}

		// get identifiers/keywords
		if (L.ch >= 'a' && L.ch <= 'z') || (L.ch >= 'A' && L.ch <= 'Z') || L.ch == '_' {
			for {
				next := L.peek()
				if L.readPos >= len(code) || (!unicode.IsDigit(rune(next)) && !unicode.IsLetter(rune(next)) && next != '_') {
					tok, exists := token.Keywords[L.get_lexeme(0)]
					if exists {
						tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: tok, Line: line})
						L.advance()
						break

					} else {
						tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.IDENTIFIER, Line: line})
						L.advance()
						break
					}
				} else {
					L.readPos++
				}
			}
			continue
		}

		// Number detection
		// todo: refactor, remove num and put concrete types (int, float) + support for floats
		if L.isDigit(L.pos) {
			for {
				if L.peek() == ' ' || L.peek() == '\t' || L.peek() == '\n' || L.peek() == 0 {
					tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.INT_LIT, Line: line})
					L.advance()
					break
				}
				if L.peek() == '.' {
					for {
						if L.peek() == ' ' || L.peek() == '\t' || L.peek() == '\n' || L.peek() == 0 {
							tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.FLOAT_LIT, Line: line})
							L.advance()
							break
						} else {
							L.readPos++
						}
					}
					break
				}
				if !L.isDigit(L.readPos) {
					for {
						if L.peek() == ' ' || L.peek() == '\t' || L.peek() == '\n' || L.peek() == 0 {
							tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.ILLEGAL, Line: line})
							L.advance()
							break
						} else {
							L.readPos++
						}
					}
					break
				}
				L.readPos++
			}
			continue
		}

		// Catch Illegal tokens
		for {
			if L.peek() == ' ' || L.peek() == '\t' || L.peek() == '\n' || L.peek() == 0 {
				tokenStream = append(tokenStream, token.Token{Lexeme: L.get_lexeme(0), Token: token.ILLEGAL, Line: line})
				L.advance()
				break
			} else {
				L.readPos++
			}
		}
	}
	return tokenStream
}
