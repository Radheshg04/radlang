package interpreter

import "radlang/token"

type Numeric interface {
	~int64 | ~int | ~float64
}

type Value interface {
	valueType() token.TokenType
}

type IntValue struct {
	Val int
}

func (v IntValue) valueType() token.TokenType {
	return token.INT
}

type FloatValue struct {
	Val float64
}

func (v FloatValue) valueType() token.TokenType {
	return token.FLOAT
}

type BoolValue struct {
	Val bool
}

func (v BoolValue) valueType() token.TokenType {
	return token.BOOL
}

type StringValue struct {
	Val string
}

func (v StringValue) valueType() token.TokenType {
	return token.STRING
}
