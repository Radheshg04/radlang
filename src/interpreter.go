package main

import (
	"fmt"
	"strconv"
)

type Numeric interface {
	~int64 | ~int | ~float64
}

type Value interface {
	valueType() TokenType
}

type IntValue struct {
	Val int
}

func (v IntValue) valueType() TokenType {
	return INT
}

type FloatValue struct {
	Val float64
}

func (v FloatValue) valueType() TokenType {
	return FLOAT
}

type BoolValue struct {
	Val bool
}

func (v BoolValue) valueType() TokenType {
	return BOOL
}

type StringValue struct {
	Val string
}

func (v StringValue) valueType() TokenType {
	return STRING
}

var env = make(map[string]Value)

func exec(statement Statement) error {
	switch s := statement.(type) {
	case *Decl_stmt:
		switch s.Type {
		case INT:
			env[s.Name] = IntValue{}
		case FLOAT:
			env[s.Name] = FloatValue{}
		case BOOL:
			env[s.Name] = BoolValue{}
		case STRING:
			env[s.Name] = StringValue{}
		default:
			return fmt.Errorf("Undefined dtype for decl stmt")
		}
		return nil
	case *Assign_stmt:
		val, err := eval(s.Value)
		if err != nil {
			return err
		}
		env[s.Target] = val
		return nil

	case *Expr_stmt:
		_, err := eval(s.Expr)
		if err != nil {
			return err
		}
		return nil

	case *Update_stmt:
		delta := 1
		if s.Op == MINUSMINUS {
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

func eval(expression Expression) (Value, error) {
	switch e := expression.(type) {
	case *Number_lit:
		switch e.Type {
		case INT:
			val, err := strconv.Atoi(e.Value)
			if err != nil {
				return nil, err
			}
			return IntValue{Val: val}, nil
		case FLOAT:
			val, err := strconv.ParseFloat(e.Value, 64)
			if err != nil {
				return nil, err
			}
			return FloatValue{Val: val}, nil
		}

	case *Call_expr:
		if e.Name == "print" {
			fmt.Println(env)
		}
		return nil, nil
	case *Identifier_expr:
		return env[e.Name], nil
	case *Binary_expr:
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
func performOp[T Numeric](l T, r T, op TokenType) T {
	switch op {
	case PLUS:
		return l + r
	case MINUS:
		return l - r
	case ASTERISK:
		return l * r
	case SLASH:
		return l / r
	}
	return 0
}

func interpret(program *Program) error {
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
