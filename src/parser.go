package main

import "fmt"

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) peek() (Token, error) {
	if p.pos >= len(p.tokens) {
		return Token{}, fmt.Errorf("unexpected end of file")
	}
	return p.tokens[p.pos], nil
}

func (p *Parser) peekAt(index int) (Token, error) {
	if p.pos+index >= len(p.tokens) {
		return Token{}, fmt.Errorf("unexpected end of file")
	}
	return p.tokens[p.pos+index], nil
}

func (p *Parser) consume() (Token, error) {
	tok, err := p.peek()
	if err != nil {
		return Token{}, err
	}
	p.pos++
	return tok, nil
}

func (p *Parser) expect(tt TokenType) bool {
	tok, err := p.peek()
	if err != nil || tok.Token != tt {
		panic(fmt.Sprintf("expected %v got %v", tt, tok.Token))
	}
	return tt == tok.Token
}

type Node interface{}

type Program struct {
	Functions []*Func_Decl
}

// Add returntype and params
type Func_Decl struct {
	Name string
	Body *Block
}

type Block struct {
	Statement_Group []Statement
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

type Decl_stmt struct {
	Name string
	Type TokenType
}

func (*Decl_stmt) stmtNode() {}

type Assign_stmt struct {
	Target string
	Value  Expression
}

func (*Assign_stmt) stmtNode() {}

// Handles x++ and x--
type Update_stmt struct {
	Target string
	Op     TokenType
}

func (*Update_stmt) stmtNode() {}

type Expr_stmt struct {
	Expr Expression
}

func (*Expr_stmt) stmtNode() {}

type Identifier_expr struct {
	Name string
}

func (*Identifier_expr) exprNode() {}

type Binary_expr struct {
	Left  Expression
	Op    TokenType
	Right Expression
}

func (*Binary_expr) exprNode() {}

// Function Call expr
type Call_expr struct {
	Name string
}

func (*Call_expr) exprNode() {}

type Number_lit struct {
	Value string
}

func (*Number_lit) exprNode() {}

// func (*Expression) exprNode() {}
