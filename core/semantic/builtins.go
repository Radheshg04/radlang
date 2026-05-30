package semantic

import "radlang/token"

func initBuiltins(ctx *SemanticCtx) {
	printDecl := &FuncSymbol{Params: map[string]token.TokenType{"val": token.STRING_LIT}, Returns: nil, isBuiltin: true}
	errDecl := &FuncSymbol{Params: map[string]token.TokenType{"val": token.STRING_LIT}, Returns: []token.TokenType{token.ERR}, isBuiltin: true}

	builtins := map[string]*FuncSymbol{
		"print": printDecl,
		"error": errDecl,
	}

	for funcName, funcSym := range builtins {
		ctx.Scope.SymbolTable[funcName] = funcSym
	}
}

func isBuiltinFunc(function Symbol) bool {
	funcSym, ok := function.(*FuncSymbol)
	if !ok {
		return false
	}
	return funcSym.isBuiltin
}
