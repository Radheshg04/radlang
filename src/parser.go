package main

import (
	"fmt"
)

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

func (p *Parser) expect(tt TokenType) Token {
	tok, err := p.peek()
	if err != nil || tok.Token != tt {
		panic(fmt.Sprintf("expected %v on line: %d, got %v", tt, tok.line, tok.Token))
	}
	p.consume()
	return tok
}

func (p *Parser) expectType() (Token, bool) {
	tok, _ := p.consume()
	switch tok.Token {
	case INT, FLOAT, BOOL, STRING:
		return tok, true
	default:
		return tok, false
	}
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
	Type  TokenType
}

func (*Number_lit) exprNode() {}

// Parsing Logic
func (p *Parser) parseProgram() (*Program, error) {
	program := &Program{}
	for {
		if val, err := p.peek(); err == nil && val.Token == EOF {
			p.consume()
			break
		}
		if val, err := p.peek(); err == nil && val.Token == EOL {
			p.consume()
			continue
		}
		if val, err := p.peek(); err == nil && val.Token == FUNC {
			// skip func keyword and move to ident + block
			p.consume()

			// consume ident
			token := p.expect(IDENTIFIER)
			func_ident := token.lexeme

			// expect no args (mk 1)
			p.expect(L_PAREN)
			p.expect(R_PAREN)

			// consume block
			func_body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			program.Functions = append(program.Functions,
				&Func_Decl{
					Name: func_ident,
					Body: func_body,
				},
			)
		} else {
			val, _ := p.peek()
			return nil, fmt.Errorf("unexpected token %v on line: %d", val.Token, val.line)
		}
	}
	return program, nil
}

func (p *Parser) parseBlock() (*Block, error) {
	block := &Block{}
	p.expect(L_BRACE)
	for {
		if val, err := p.peek(); err == nil && val.Token == EOL {
			// skip eol
			p.consume()
			continue
		}
		if val, err := p.peek(); err == nil && val.Token == EOF {
			return nil, fmt.Errorf("unclosed block, unexpected EOF")
		}
		if val, err := p.peek(); err == nil && val.Token == R_BRACE {
			// consume the r_brace and return (block parsed)
			p.consume()
			return block, nil
		}
		for {
			if val, err := p.peek(); err == nil && val.Token == R_BRACE {
				// break on r_brace
				break
			}
			// skip eols
			if val, err := p.peek(); err == nil && val.Token == EOL {
				p.consume()
				continue
			}
			statement, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			block.Statement_Group = append(block.Statement_Group, statement)
		}
	}
}

func (p *Parser) parseStatement() (Statement, error) {

	for {
		// catch unexplected EOF
		if val, err := p.peek(); err == nil && val.Token == EOF {
			return nil, fmt.Errorf("unexpected EOF on line: %v", val.line)
		}

		// DECL_STMT
		if val, err := p.peek(); err == nil && val.Token == VAR {
			decl_stmt := &Decl_stmt{}
			// consume "var"
			p.consume()

			ident := p.expect(IDENTIFIER)
			ident_type, ok := p.expectType()
			if !ok {
				return nil, fmt.Errorf(
					"unexpected %v on line %v, expected one of INT, FLOAT, BOOL, STRING",
					ident_type,
					ident.line,
				)

			}

			p.expect(EOL)

			decl_stmt.Name = ident.lexeme
			decl_stmt.Type = ident_type.Token

			return decl_stmt, nil
		}

		// ASSIGN_STMT
		if val, err := p.peekAt(1); err == nil && val.Token == ASSIGNMENT {
			assign_stmt := &Assign_stmt{}

			ident := p.expect(IDENTIFIER)

			p.expect(ASSIGNMENT)

			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			p.expect(EOL)

			assign_stmt.Target = ident.lexeme
			assign_stmt.Value = expr

			return assign_stmt, nil
		}

		// EXPR_STMT
		if val, err := p.peekAt(1); err == nil {
			switch val.Token {
			case L_PAREN:
				ident, err := p.consume()
				if err != nil {
					return nil, err
				}
				p.expect(L_PAREN)
				p.expect(R_PAREN)

				p.expect(EOL)

				return &Expr_stmt{
					&Call_expr{
						Name: ident.lexeme,
					}}, nil

			case PLUSPLUS, MINUSMINUS:
				ident, err := p.consume()
				if err != nil {
					return nil, err
				}
				op, err := p.consume()
				if err != nil {
					return nil, err
				}

				p.expect(EOL)

				return &Update_stmt{
					Target: ident.lexeme,
					Op:     op.Token,
				}, nil

			default:
				expr_stmt := &Expr_stmt{}

				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}

				expr_stmt.Expr = expr

				return expr_stmt, nil
			}
		}

	}
}

func (p *Parser) parseExpression() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	var op TokenType
	var right Expression
	for {
		if val, err := p.peek(); err == nil && (val.Token == PLUS || val.Token == MINUS) {
			tok, err := p.consume()
			if err != nil {
				return nil, err
			}
			op = tok.Token
			right, err = p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = &Binary_expr{Left: left, Op: op, Right: right}
		} else {
			return left, nil
		}
	}
}

func (p *Parser) parseTerm() (Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	var op TokenType
	var right Expression
	for {
		if val, err := p.peek(); err == nil && (val.Token == ASTERISK || val.Token == SLASH) {
			tok, err := p.consume()
			if err != nil {
				return nil, err
			}
			op = tok.Token
			right, err = p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = &Binary_expr{Left: left, Op: op, Right: right}
		} else {
			return left, nil
		}
	}
}

func (p *Parser) parseFactor() (Expression, error) {
	tok, err := p.peek()
	if err != nil {
		return nil, err
	}
	switch tok.Token {
	case INT:
		p.consume()
		return &Number_lit{Value: tok.lexeme, Type: INT}, nil
	case FLOAT:
		p.consume()
		return &Number_lit{Value: tok.lexeme, Type: FLOAT}, nil
	case IDENTIFIER:
		if next, err := p.peekAt(1); err == nil && next.Token == L_PAREN {
			p.consume()
			p.expect(L_PAREN)
			p.expect(R_PAREN)
			return &Call_expr{Name: tok.lexeme}, nil
		}
		p.consume()
		return &Identifier_expr{Name: tok.lexeme}, nil
	case L_PAREN:
		p.consume()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.expect(R_PAREN)
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token %v in expression", tok.Token)

	}
}

func Parse(TokenStream []Token) (*Program, error) {
	parser := Parser{tokens: TokenStream, pos: 0}
	return parser.parseProgram()
}
