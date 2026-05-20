package interpreter

import (
	"fmt"
	"radlang/parser"
	"radlang/token"
	"strconv"
)

var env = make(map[string]Value)

func exec(statement parser.Statement) error {
	switch s := statement.(type) {
	case *parser.Decl_stmt:
		switch s.Type {
		case token.INT:
			env[s.Name] = IntValue{}
		case token.FLOAT:
			env[s.Name] = FloatValue{}
		case token.BOOL:
			env[s.Name] = BoolValue{}
		case token.STRING:
			env[s.Name] = StringValue{}
		default:
			return fmt.Errorf("Undefined dtype for decl stmt")
		}
		return nil
	case *parser.Assign_stmt:
		val, err := eval(s.Value)
		if err != nil {
			return err
		}
		env[s.Target] = val
		return nil

	case *parser.Expr_stmt:
		_, err := eval(s.Expr)
		if err != nil {
			return err
		}
		return nil

	case *parser.Update_stmt:
		delta := 1
		if s.Op == token.MINUSMINUS {
			delta = -1
		}
		switch v := env[s.Target].(type) {
		case IntValue:
			env[s.Target] = IntValue{Val: v.Val + delta}
			return nil
		case FloatValue:
			env[s.Target] = FloatValue{Val: v.Val + float64(delta)}
			return nil
		default:
			return fmt.Errorf("Couldnt update stmt")
		}
	}
	return fmt.Errorf("statement interpretation failed")
}

func eval(expression parser.Expression) (Value, error) {
	switch e := expression.(type) {
	case *parser.Number_lit:
		switch e.Type {
		case token.INT:
			val, err := strconv.Atoi(e.Value)
			if err != nil {
				return nil, err
			}
			return IntValue{Val: val}, nil
		case token.FLOAT:
			val, err := strconv.ParseFloat(e.Value, 64)
			if err != nil {
				return nil, err
			}
			return FloatValue{Val: val}, nil
		}

	case *parser.Call_expr:
		if e.Name == "print" {
			fmt.Println(env)
		}
		return nil, nil
	case *parser.Identifier_expr:
		return env[e.Name], nil
	case *parser.Binary_expr:
		left, err := eval(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := eval(e.Right)
		if err != nil {
			return nil, err
		}
		switch l := left.(type) {
		case IntValue:
			r := right.(IntValue)
			val := performOp(l.Val, r.Val, e.Op)
			return IntValue{Val: val}, nil
		case FloatValue:
			r := right.(FloatValue)
			val := performOp(l.Val, r.Val, e.Op)
			return FloatValue{Val: val}, nil
		}
	}
	return nil, fmt.Errorf("")
}

// use generic func to perform bin op
func performOp[T Numeric](l T, r T, op token.TokenType) T {
	switch op {
	case token.PLUS:
		return l + r
	case token.MINUS:
		return l - r
	case token.ASTERISK:
		return l * r
	case token.SLASH:
		return l / r
	}
	return 0
}

func Interpret(program *parser.Program) error {
	for _, function := range program.Functions {
		if function.Name == "main" {
			for _, stmt := range function.Body.Statement_Group {
				err := exec(stmt)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
