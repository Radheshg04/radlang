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
	slotCounter *int
}

type Symbol interface {
	symbol()
}

type VarSymbol struct {
	Slot     int
	Type     ValueType
	Declared bool
}

func (*VarSymbol) symbol() {}

type FuncSymbol struct {
	ID        int
	Params    map[string]ValueType
	Returns   []ValueType
	Decl      *parser.Func_Decl
	isBuiltin bool
	Slots     int
}

func (*FuncSymbol) symbol() {}
