package semantic

import (
	"radlang/parser"
)

type Scope struct {
	SymbolTable map[string]Symbol
	Parent      *Scope
}

type SemanticCtx struct {
	Scope       *Scope
	CurrentFunc *FuncSymbol
	LoopDepth   int
	Diagnostics []Diagnostic
}

type Symbol interface {
	symbol()
}

type VarSymbol struct {
	Type     ValueType
	Value    interface{}
	Declared bool
}

func (*VarSymbol) symbol() {}

type FuncSymbol struct {
	Params    map[string]ValueType
	Returns   []ValueType
	Decl      *parser.Func_Decl
	isBuiltin bool
}

func (*FuncSymbol) symbol() {}
