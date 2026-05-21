package lexer

import "radlang/token"

// helper functions for lexer package

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

func (l *lexer) get_lexeme(tok *token.TokenType) string {
	end := min(l.readPos, len(l.input))
	if tok != nil && *tok == token.STRING_LIT {
		return l.input[l.pos+1 : end-1]
	}
	return l.input[l.pos:end]
}

func (l *lexer) emit(tok token.TokenType) {
	l.tokenStream = append(l.tokenStream, token.Token{Lexeme: l.get_lexeme(&tok), Token: tok, Line: l.line})
	if tok == token.EOF {
		return
	}
	l.advance()
}

func (l *lexer) match(expected byte) bool {
	if l.peek() != expected {
		return false
	}
	l.readPos++
	return true
}

func (l *lexer) isNumberTerminator(b byte) bool {
	switch b {
	case ' ', '\t', '\n', 0, ')', ',', '+', '-', '*', '/', '=', '<', '>', '!', '{', '}', '(':
		return true
	}
	return false
}

func (l *lexer) isTokenTerminator(b byte) bool {
	switch b {
	case ' ', '\t', '\n', 0:
		return true
	}
	return false
}
