package parser

import (
	"fmt"
	"radlang/token"
)

type Node interface{}

type Parser struct {
	tokens []token.Token
	pos    int
}

func (p *Parser) parseProgram() (*Program, error) {
	program := &Program{}
	for {
		if val, err := p.peek(); err == nil && val.Token == token.EOF {
			p.consume()
			break
		}
		if val, err := p.peek(); err == nil && val.Token == token.EOL {
			p.consume()
			continue
		}
		if val, err := p.peek(); err == nil && val.Token == token.FUNC {
			// skip func keyword and move to ident + block
			p.consume()

			// consume ident
			tok := p.expect(token.IDENTIFIER)
			func_ident := tok.Lexeme

			// expect no args (mk 1)
			p.expect(token.L_PAREN)
			p.expect(token.R_PAREN)

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
			return nil, fmt.Errorf("unexpected token %v on line: %d", val.Token, val.Line)
		}
	}
	return program, nil
}

func (p *Parser) parseBlock() (*Block, error) {
	block := &Block{}
	p.expect(token.L_BRACE)
	for {
		if val, err := p.peek(); err == nil && val.Token == token.EOL {
			// skip eol
			p.consume()
			continue
		}
		if val, err := p.peek(); err == nil && val.Token == token.EOF {
			return nil, fmt.Errorf("unclosed block, unexpected EOF")
		}
		if val, err := p.peek(); err == nil && val.Token == token.R_BRACE {
			// consume the r_brace and return (block parsed)
			p.consume()
			return block, nil
		}
		for {
			if val, err := p.peek(); err == nil && val.Token == token.R_BRACE {
				// break on r_brace
				break
			}
			// skip eols
			if val, err := p.peek(); err == nil && val.Token == token.EOL {
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
		if val, err := p.peek(); err == nil && val.Token == token.EOF {
			return nil, fmt.Errorf("unexpected EOF on line: %v", val.Line)
		}

		// DECL_STMT
		if val, err := p.peek(); err == nil && val.Token == token.VAR {
			decl_stmt := &Decl_stmt{}
			// consume "var"
			p.consume()

			ident := p.expect(token.IDENTIFIER)
			ident_type, ok := p.expectType()
			if !ok {
				return nil, fmt.Errorf(
					"unexpected %v on line %v, expected one of INT, FLOAT, BOOL, STRING",
					ident_type,
					ident.Line,
				)

			}

			p.expect(token.EOL)

			decl_stmt.Name = ident.Lexeme
			decl_stmt.Type = ident_type.Token

			return decl_stmt, nil
		}

		// ASSIGN_STMT
		if val, err := p.peekAt(1); err == nil && val.Token == token.ASSIGNMENT {
			assign_stmt := &Assign_stmt{}

			ident := p.expect(token.IDENTIFIER)

			p.expect(token.ASSIGNMENT)

			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			p.expect(token.EOL)

			assign_stmt.Target = ident.Lexeme
			assign_stmt.Value = expr

			return assign_stmt, nil
		}

		// EXPR_STMT
		if val, err := p.peekAt(1); err == nil {
			switch val.Token {
			case token.L_PAREN:
				ident, err := p.consume()
				if err != nil {
					return nil, err
				}
				p.expect(token.L_PAREN)
				p.expect(token.R_PAREN)

				p.expect(token.EOL)

				return &Expr_stmt{
					&Call_expr{
						Name: ident.Lexeme,
					}}, nil

			case token.PLUSPLUS, token.MINUSMINUS:
				ident, err := p.consume()
				if err != nil {
					return nil, err
				}
				op, err := p.consume()
				if err != nil {
					return nil, err
				}

				p.expect(token.EOL)

				return &Update_stmt{
					Target: ident.Lexeme,
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
	var op token.TokenType
	var right Expression
	for {
		if val, err := p.peek(); err == nil && (val.Token == token.PLUS || val.Token == token.MINUS) {
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
	var op token.TokenType
	var right Expression
	for {
		if val, err := p.peek(); err == nil && (val.Token == token.ASTERISK || val.Token == token.SLASH) {
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
	case token.INT_LIT:
		p.consume()
		return &Number_lit{Value: tok.Lexeme, Type: token.INT}, nil
	case token.FLOAT_LIT:
		p.consume()
		return &Number_lit{Value: tok.Lexeme, Type: token.FLOAT}, nil
	case token.IDENTIFIER:
		if next, err := p.peekAt(1); err == nil && next.Token == token.L_PAREN {
			p.consume()
			p.expect(token.L_PAREN)
			p.expect(token.R_PAREN)
			return &Call_expr{Name: tok.Lexeme}, nil
		}
		p.consume()
		return &Identifier_expr{Name: tok.Lexeme}, nil
	case token.L_PAREN:
		p.consume()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.expect(token.R_PAREN)
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token %v in expression", tok.Token)

	}
}

// Parsing Logic
func Parse(TokenStream []token.Token) (*Program, error) {
	parser := Parser{tokens: TokenStream, pos: 0}
	return parser.parseProgram()
}
