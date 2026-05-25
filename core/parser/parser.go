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
		if val := p.peek(); val.Token == token.EOF {
			p.consume()
			break
		}
		if val := p.peek(); val.Token == token.EOL {
			p.consume()
			continue
		}
		if val := p.peek(); val.Token == token.FUNC {
			// skip func keyword and move to ident + block
			p.consume()

			// parse signature
			func_sig, err := p.parseFuncSignature()
			if err != nil {
				return nil, err
			}

			func_body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			program.Functions = append(program.Functions, &Func_Decl{
				Signature: func_sig,
				Body:      func_body,
			})
		} else {
			val := p.peek()
			return nil, fmt.Errorf("unexpected token %v on line: %d", val.Token, val.Line)
		}
	}
	return program, nil
}

func (p *Parser) parseFuncSignature() (*Func_Signature, error) {
	func_sig := &Func_Signature{}
	name := p.expect(token.IDENTIFIER)

	// extract params
	p.expect(token.L_PAREN)
	var func_params []*Param
	for val := p.peek(); val.Token != token.R_PAREN; {
		paramName := p.expect(token.IDENTIFIER)

		paramType := p.expectType()
		param := &Param{
			Name: paramName.Lexeme,
			Type: paramType.Token,
		}
		func_params = append(func_params, param)
		if val := p.peek(); val.Token == token.COMMA {
			p.consume()
		} else {
			break
		}
	}
	p.expect(token.R_PAREN)

	// parse return vals
	ret := p.peek()

	var returns []token.TokenType
	var err error
	// this scales badly, would need some type checker for custom types
	switch ret.Token {
	// single return type
	case token.INT, token.FLOAT, token.BOOL, token.STRING:
		val := p.expectType()
		returns = append(returns, val.Token)

	// multiple return types
	case token.L_PAREN:
		p.consume()
		returns, err = parseMany(p, token.COMMA, token.R_PAREN, func() (token.TokenType, error) {
			return p.expectType().Token, nil
		})
		if err != nil {
			return nil, err
		}

	// no return
	case token.L_BRACE:
		break
	default:
		return nil, fmt.Errorf("Unexpected token ")
	}

	func_sig.Name = name.Lexeme
	func_sig.Params = func_params
	func_sig.Returns = returns

	return func_sig, nil
}

func (p *Parser) parseBlock() (*Block, error) {
	block := &Block{}
	p.expect(token.L_BRACE)
	for {
		if val := p.peek(); val.Token == token.EOL {
			// skip eol
			p.consume()
			continue
		}
		if val := p.peek(); val.Token == token.EOF {
			return nil, fmt.Errorf("unclosed block, unexpected EOF")
		}
		if val := p.peek(); val.Token == token.R_BRACE {
			// consume the r_brace and return (block parsed)
			p.consume()
			return block, nil
		}
		for {
			if val := p.peek(); val.Token == token.R_BRACE {
				// break on r_brace
				break
			}
			// skip eols
			if val := p.peek(); val.Token == token.EOL {
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
		if val := p.peek(); val.Token == token.EOF {
			return nil, fmt.Errorf("unexpected EOF on line: %v", val.Line)
		}

		// JUMP_STMT
		if val := p.peek(); val.Token == token.BREAK || val.Token == token.CONTINUE {
			p.consume()
			return &Jump_stmt{Type: val.Token}, nil
		}

		// RETURN_STMT
		if val := p.peek(); val.Token == token.RETURN {
			p.consume()
			var exprList []Expression
			for p.peek().Token != token.EOL {
				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				exprList = append(exprList, expr)
				if p.peek().Token == token.COMMA {
					p.consume()
					continue
				}
				break
			}
			return &Return_stmt{Returns: exprList}, nil
		}

		// LOOP_STMT
		if val := p.peek(); val.Token == token.FOR {
			p.consume()
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			block, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			return &Loop_stmt{Expression: expr, Loop_block: block}, nil
		}

		// parse decl stmt
		if val := p.peek(); val.Token == token.VAR {
			decl_stmt := &Decl_stmt{}
			// consume "var"
			p.consume()
			idents, err := parseMany(p, token.COMMA, token.INVALID,
				func() (string, error) {
					return p.expect(token.IDENTIFIER).Lexeme, nil
				})
			if err != nil {
				return nil, err
			}
			ident_type := p.expectType()
			p.expect(token.EOL)

			decl_stmt.Name = idents
			decl_stmt.Type = ident_type.Token

			return decl_stmt, nil
		}

		// ASSIGN_STMT
		if p.containsAny(token.ASSIGNMENT, token.WALRUS) {
			assign_stmt := &Assign_stmt{}
			targets, err := parseMany(p, token.COMMA, token.INVALID,
				func() (string, error) {
					return p.expect(token.IDENTIFIER).Lexeme, nil
				})
			if err != nil {
				return nil, err
			}

			op := p.expectAny(token.WALRUS, token.ASSIGNMENT)

			values, err := parseMany(p, token.COMMA, token.INVALID,
				func() (Expression, error) {
					return p.parseExpression()
				})
			if err != nil {
				return nil, err
			}
			p.expect(token.EOL)

			assign_stmt.Targets = targets
			assign_stmt.Op = op.Token
			assign_stmt.Values = values

			return assign_stmt, nil
		}

		// CONTROL_STMT
		if p.peek().Token == token.IF {
			return p.parseControlStatement()
		}

		// EXPR_STMT

		expr_stmt := &Expr_stmt{}

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		expr_stmt.Expression = expr

		return expr_stmt, nil

	}
}

func (p *Parser) parseControlStatement() (*Control_stmt, error) {
	p.expect(token.IF)
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	ifBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	// if p.peek().Token == token.EOL {
	// 	p.consume()
	// 	return &Control_stmt{Expression: expr, IfBlock: ifBlock}, nil
	// }
	if p.peek().Token == token.ELSE {
		p.consume()
		if p.peek().Token == token.IF {
			elseStatement, err := p.parseControlStatement()
			if err != nil {
				return nil, err
			}
			return &Control_stmt{Expression: expr, IfBlock: ifBlock, ElseStmt: elseStatement}, nil
		}
		if p.peek().Token == token.L_BRACE {
			elseBlock, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			return &Control_stmt{Expression: expr, IfBlock: ifBlock, ElseBlock: elseBlock}, nil

		}
	}
	return &Control_stmt{Expression: expr, IfBlock: ifBlock}, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	left, err := p.parseArithmetic()
	if err != nil {
		return nil, err
	}
	for {
		switch op := p.peek().Token; op {
		case token.GT, token.GTE, token.LT, token.LTE, token.EQ, token.NEQ:
			// consume op
			p.consume()
			// parse right
			right, err := p.parseArithmetic()
			if err != nil {
				return nil, err
			}
			return &Binary_expr{
				Left:  left,
				Op:    op,
				Right: right,
			}, nil

		default:
			return left, nil
		}

	}

}

func (p *Parser) parseArithmetic() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	var op token.TokenType
	var right Expression
	for {
		if val := p.peek(); val.Token == token.PLUS || val.Token == token.MINUS {
			tok := p.consume()
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
		if val := p.peek(); val.Token == token.ASTERISK || val.Token == token.SLASH {
			tok := p.consume()
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
	tok := p.peek()
	switch tok.Token {

	// parse literals
	case token.INT_LIT:
		p.consume()
		return &Lit_val{Value: tok.Lexeme, Type: token.INT}, nil
	case token.FLOAT_LIT:
		p.consume()
		return &Lit_val{Value: tok.Lexeme, Type: token.FLOAT}, nil
	case token.STRING_LIT:
		p.consume()
		return &Lit_val{Value: tok.Lexeme, Type: token.STRING}, nil
	case token.BOOL_LIT:
		p.consume()
		return &Lit_val{Value: tok.Lexeme, Type: token.BOOL}, nil
	case token.ERR:
		p.consume()
		p.expect(token.L_PAREN)
		str := p.expect(token.STRING_LIT)
		p.expect(token.R_PAREN)
		return &Lit_val{Value: str.Lexeme, Type: token.ERR}, nil

	// parse postfix op, fn call and identifier expr
	case token.IDENTIFIER:
		// fn call
		if next := p.peekAt(1); next.Token == token.L_PAREN {
			p.consume()
			args, err := p.parseArgs()
			if err != nil {
				return nil, err
			}
			return &Call_expr{Name: tok.Lexeme, Args: args}, nil
		}

		// consume ident
		p.consume()
		identExpr := &Identifier_expr{Name: tok.Lexeme}

		// parse postfix op
		if next := p.peek(); next.Token == token.PLUSPLUS || next.Token == token.MINUSMINUS {
			p.consume()
			return &Postfix_expr{
				Target: identExpr,
				Op:     next.Token,
			}, nil
		}
		return identExpr, nil

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

func (p *Parser) parseArgs() ([]Expression, error) {
	p.expect(token.L_PAREN)
	if p.peek().Token == token.R_PAREN {
		p.consume()
		return nil, nil
	}
	return parseMany(p, token.COMMA, token.R_PAREN, p.parseExpression)
}

// Parsing Logic
func Parse(TokenStream []token.Token) (program *Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	parser := Parser{tokens: TokenStream, pos: 0}
	return parser.parseProgram()
}
