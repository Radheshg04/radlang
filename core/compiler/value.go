package compiler

import "radlang/semantic"

type Value interface {
	valueType() semantic.ValueType
	Add(rhs Value) Value
	Sub(rhs Value) Value
	Mul(rhs Value) Value
	Div(rhs Value) Value
}

type IntValue struct {
	Val int64
}

func (v IntValue) valueType() semantic.ValueType {
	return semantic.IntType
}

func (v IntValue) Add(rhs Value) Value {
	return IntValue{Val: v.Val + rhs.(IntValue).Val}
}
func (v IntValue) Sub(rhs Value) Value {
	return IntValue{Val: v.Val - rhs.(IntValue).Val}
}
func (v IntValue) Mul(rhs Value) Value {
	return IntValue{Val: v.Val * rhs.(IntValue).Val}
}
func (v IntValue) Div(rhs Value) Value {
	return IntValue{Val: v.Val / rhs.(IntValue).Val}
}

type FloatValue struct {
	Val float64
}

func (v FloatValue) valueType() semantic.ValueType {
	return semantic.FloatType
}

func (v FloatValue) Add(rhs Value) Value {
	return FloatValue{Val: v.Val + rhs.(FloatValue).Val}
}
func (v FloatValue) Sub(rhs Value) Value {
	return FloatValue{Val: v.Val - rhs.(FloatValue).Val}
}
func (v FloatValue) Mul(rhs Value) Value {
	return FloatValue{Val: v.Val * rhs.(FloatValue).Val}
}
func (v FloatValue) Div(rhs Value) Value {
	return FloatValue{Val: v.Val / rhs.(FloatValue).Val}
}

type BoolValue struct {
	Val bool
}

func (v BoolValue) valueType() semantic.ValueType {
	return semantic.BoolType
}

func (v BoolValue) Add(rhs Value) Value {
	return BoolValue{}
}
func (v BoolValue) Sub(rhs Value) Value {
	return BoolValue{}
}
func (v BoolValue) Mul(rhs Value) Value {
	return BoolValue{}
}
func (v BoolValue) Div(rhs Value) Value {
	return BoolValue{}
}

type StringValue struct {
	Val string
}

func (v StringValue) valueType() semantic.ValueType {
	return semantic.StringType
}

func (v StringValue) Add(rhs Value) Value {
	return StringValue{Val: v.Val + rhs.(StringValue).Val}
}
func (v StringValue) Sub(rhs Value) Value {
	return StringValue{}
}
func (v StringValue) Mul(rhs Value) Value {
	return StringValue{}
}
func (v StringValue) Div(rhs Value) Value {
	return StringValue{}
}
