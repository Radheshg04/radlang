// usr/bin/go run $0 $@ ; exit
package lexer

import (
	"unicode"

	"radlang/token"
)

type lexer struct {
	input       string // input string to the lexer
	pos         int    // position of pointer in input
	readPos     int    // next Token
	ch          byte   // current character
	line        int
	tokenStream []token.Token
}

func Lex(code string) []token.Token {
	var l lexer
	l.line = 1
	l.input = code
	l.readPos = 1
	// guard against empty input
	if len(code) == 0 {
		l.emit(token.EOF)
		return l.tokenStream
	}
	l.ch = l.input[l.pos]
	for {
		// EOF Handling
		if l.pos >= len(code) {
			l.emit(token.EOF)
			break
		}

		// Skip whitespaces and tabs
		if l.ch == ' ' || l.ch == '\t' {
			l.advance()
			continue
		}

		// Tokenize newlines
		if l.ch == '\n' {
			l.emit(token.EOL)
			l.line++
			continue
		}

		// String handing
		if l.ch == '"' {
			for {
				if l.peek() == '\n' || l.readPos >= len(l.input) {
					l.emit(token.ILLEGAL)
					break
				}
				if l.match('"') {
					// make a token of type string literal which starts at l.pos and ends at l.readPos
					// TODO: Add support for escape sequences
					l.emit(token.STRING_LIT)
					break
				}
				l.readPos++
			}
			continue
		}

		switch l.ch {
		case ',':
			l.emit(token.COMMA)
			continue

		case ':':
			if l.match('=') {
				l.emit(token.WALRUS)
				continue
			}
			l.emit(token.ILLEGAL)
			continue

		case '=':
			if l.match('=') {
				l.emit(token.EQ)
			} else {
				l.emit(token.ASSIGNMENT)
			}
			continue

		case '!':
			if l.match('=') {
				l.emit(token.NEQ)
				continue
			}
			l.emit(token.ILLEGAL)
			continue

		case '>':
			if l.match('=') {
				l.emit(token.GTE)
			} else {
				l.emit(token.GT)
			}
			continue

		case '<':
			if l.match('=') {
				l.emit(token.LTE)
			} else {
				l.emit(token.LT)
			}
			continue

		case '+':
			if l.match('+') {
				l.emit(token.PLUSPLUS)
			} else {
				l.emit(token.PLUS)
			}
			continue

		case '-':
			if l.match('-') {
				l.emit(token.MINUSMINUS)
			} else {
				l.emit(token.MINUS)
			}
			continue

		case '*':
			l.emit(token.ASTERISK)
			continue

		case '/':
			l.emit(token.SLASH)
			continue

		case '(':
			l.emit(token.L_PAREN)
			continue

		case ')':
			l.emit(token.R_PAREN)
			continue

		case '{':
			l.emit(token.L_BRACE)
			continue

		case '}':
			l.emit(token.R_BRACE)
			continue
		}

		// get identifiers/keywords
		if (l.ch >= 'a' && l.ch <= 'z') || (l.ch >= 'A' && l.ch <= 'Z') || l.ch == '_' {
			for {
				next := l.peek()
				if l.readPos >= len(code) || (!unicode.IsDigit(rune(next)) && !unicode.IsLetter(rune(next)) && next != '_') {
					tok, exists := token.Keywords[l.get_lexeme(nil)]
					if exists {
						l.emit(tok)
						break

					} else {
						l.emit(token.IDENTIFIER)
						break
					}
				} else {
					l.readPos++
				}
			}
			continue
		}

		// Number detection
		if l.isDigit(l.pos) {
			for {
				if l.isNumberTerminator(l.peek()) {
					l.emit(token.INT_LIT)
					break
				}
				if l.match('.') {
					for {
						if l.isNumberTerminator(l.peek()) {
							l.emit(token.FLOAT_LIT)
							break
						}
						if !l.isDigit(l.readPos) {
							for {
								if l.isTokenTerminator(l.peek()) {
									l.emit(token.ILLEGAL)
									break
								}
								l.readPos++
							}
							break
						}
						l.readPos++
					}
					break
				}
				if !l.isDigit(l.readPos) {
					for {
						if l.isTokenTerminator(l.peek()) {
							l.emit(token.ILLEGAL)
							break
						} else {
							l.readPos++
						}
					}
					break
				}
				l.readPos++
			}
			continue
		}

		// Catch Illegal tokens
		for {
			if l.isTokenTerminator(l.peek()) {
				l.emit(token.ILLEGAL)
				break
			} else {
				l.readPos++
			}
		}
	}
	return l.tokenStream
}
