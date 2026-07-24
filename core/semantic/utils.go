package semantic

import (
	"fmt"
	"radlang/parser"
	"radlang/token"
)

func Resolve(scope *Scope, name string) (Symbol, error) {
	symbol, ok := scope.SymbolTable[name]
	if ok {
		return symbol, nil
	}
	if scope.Parent != nil {
		return Resolve(scope.Parent, name)
	}
	return nil, fmt.Errorf("%s could not be resolved", name)
}

func resolveParams(Params []*parser.Param) map[string]ValueType {
	if Params == nil {
		return nil
	}
	paramTable := make(map[string]ValueType)
	for _, param := range Params {
		paramTable[param.Name] = resolveType(param.Type)
	}
	return paramTable
}

func resolveTypes(Types []token.TokenType) []ValueType {
	out := make([]ValueType, len(Types))
	for i, t := range Types {
		out[i] = resolveType(t)
	}
	return out
}

func newChildCtx(ctx *SemanticCtx) *SemanticCtx {
	return &SemanticCtx{
		Scope:       &Scope{SymbolTable: make(map[string]Symbol), Parent: ctx.Scope},
		CurrentFunc: ctx.CurrentFunc,
		LoopDepth:   ctx.LoopDepth,
		Diagnostics: []Diagnostic{},
		slotCounter: ctx.slotCounter,
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

func resolveType(tok token.TokenType) ValueType {
	switch tok {
	case token.INT:
		return IntType
	case token.FLOAT:
		return FloatType
	case token.BOOL:
		return BoolType
	case token.STRING:
		return StringType
	case token.ERROR_LITERAL:
		return ErrorType
	default:
		return InvalidType
	}
}

func getNextSlot(ctx *SemanticCtx) int {
	nextSlot := *ctx.slotCounter
	*ctx.slotCounter++
	return nextSlot
}
