package semantic

import (
	"radlang/parser"
	"radlang/token"
)

type Scope struct {
	SymbolTable map[string]Symbol
	Parent      *Scope
}

type SemanticCtx struct {
	Scope       *Scope
	Diagnostics []Diagnostic
}

type Symbol interface {
	symbol()
}

type VarSymbol struct {
	Type     token.TokenType
	Value    interface{}
	Declared bool
}

func (*VarSymbol) symbol()

type FuncSymbol struct {
	Params  map[string]token.TokenType
	Returns []token.TokenType
	Decl    *parser.Func_Decl
}

func (*FuncSymbol) symbol()
