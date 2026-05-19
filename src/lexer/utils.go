package lexer

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

func (l *lexer) get_lexeme(offset int) string {
	return l.input[l.pos : l.readPos+offset]
}
