package semantic

import (
	"fmt"
	"radlang/parser"
	"radlang/token"
)

func resolve(scope *Scope, name string) (Symbol, error) {
	symbol, ok := scope.SymbolTable[name]
	if ok {
		return symbol, nil
	}
	if scope.Parent != nil {
		return resolve(scope.Parent, name)
	}
	return nil, fmt.Errorf("%s could not be resolved", name)
}

func resolveParams(Params []*parser.Param) map[string]token.TokenType {
	if Params == nil {
		return nil
	}
	paramTable := make(map[string]token.TokenType)
	for _, param := range Params {
		paramTable[param.Name] = param.Type
	}
	return paramTable
}

func newChildCtx(ctx *SemanticCtx) *SemanticCtx {
	return &SemanticCtx{
		Scope:       &Scope{SymbolTable: make(map[string]Symbol), Parent: ctx.Scope},
		CurrentFunc: ctx.CurrentFunc,
		LoopDepth:   ctx.LoopDepth,
		Diagnostics: []Diagnostic{},
	}
}

func symbolExistAs[T any](ctx *SemanticCtx, target string) bool {
	sym, ok := ctx.Scope.SymbolTable[target]
	if !ok {
		return false
	}
	_, ok = sym.(T)

	return ok
}
