package parser

import (
	"fmt"
	"radlang/token"
)

// helper functions for parsing

func (p *Parser) peek() (token.Token, error) {
	if p.pos >= len(p.tokens) {
		return token.Token{}, fmt.Errorf("unexpected end of file")
	}
	return p.tokens[p.pos], nil
}

func (p *Parser) peekAt(index int) (token.Token, error) {
	if p.pos+index >= len(p.tokens) {
		return token.Token{}, fmt.Errorf("unexpected end of file")
	}
	return p.tokens[p.pos+index], nil
}

func (p *Parser) consume() (token.Token, error) {
	tok, err := p.peek()
	if err != nil {
		return token.Token{}, err
	}
	p.pos++
	return tok, nil
}

func (p *Parser) expect(tt token.TokenType) token.Token {
	tok, err := p.peek()
	if err != nil || tok.Token != tt {
		panic(fmt.Sprintf("expected %v on line: %d, got %v", tt, tok.Line, tok.Token))
	}
	p.consume()
	return tok
}

func (p *Parser) expectType() (token.Token, bool) {
	tok, _ := p.consume()
	switch tok.Token {
	case token.INT, token.FLOAT, token.BOOL, token.STRING:
		return tok, true
	default:
		return tok, false
	}
}
