package parser

import (
	"fmt"
	"radlang/token"
)

// helper functions for parsing
// panics from helpers are caught at parse boundary

func (p *Parser) peek() token.Token {
	if p.pos >= len(p.tokens) {
		panic(fmt.Errorf("unexpected end of file"))
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekAt(index int) token.Token {
	if p.pos+index >= len(p.tokens) {
		panic(fmt.Errorf("unexpected end of file"))
	}
	return p.tokens[p.pos+index]
}

func (p *Parser) consume() token.Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) expect(tt token.TokenType) token.Token {
	tok := p.peek()
	if tok.Token != tt {
		panic(fmt.Errorf("expected %v on line: %d, got %v", tt, tok.Line, tok.Token))
	}
	p.consume()
	return tok
}

func (p *Parser) expectType() token.Token {
	tok := p.peek()
	switch tok.Token {
	case token.INT, token.FLOAT, token.BOOL, token.STRING, token.ERR:
		p.consume()
		return tok
	default:
		panic(fmt.Errorf("expected type on line: %d, got %v", tok.Line, tok.Token))
	}
}

func (p *Parser) expectAny(tt ...token.TokenType) token.Token {
	found := false
	tok := p.peek()
	for _, token := range tt {
		if tok.Token == token {
			found = true
		}
	}
	if !found {
		panic(fmt.Errorf("expected any of %v on line: %d, got %v", tt, tok.Line, tok.Token))
	}
	p.consume()
	return tok
}

func (p *Parser) containsAny(tokens ...token.TokenType) bool {
	for i := 0; ; i++ {
		val := p.peekAt(i).Token
		if val == token.EOL || val == token.EOF {
			break
		}
		for _, token := range tokens {
			if val == token {
				return true
			}
		}
	}
	return false
}

func parseMany[T any](
	p *Parser,
	delim token.TokenType,
	eos token.TokenType,
	parse func() (T, error),
) ([]T, error) {
	if p == nil {
		return nil, fmt.Errorf("nil parser in parseMany")
	}
	var items []T

	// parse first item
	first, err := parse()
	if err != nil {
		return nil, err
	}
	items = append(items, first)

	// parse all remaining items
	for p.peek().Token == delim {
		p.consume()
		next, err := parse()
		if err != nil {
			return nil, err
		}
		items = append(items, next)
	}
	// expect end of sequence
	if eos != token.INVALID {
		p.expect(eos)
	}
	return items, nil
}
