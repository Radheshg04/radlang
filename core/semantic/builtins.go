package semantic

func initBuiltins(ctx *SemanticCtx) {
	printDecl := &FuncSymbol{Params: map[string]ValueType{"val": StringType}, Returns: nil, isBuiltin: true}
	errDecl := &FuncSymbol{Params: map[string]ValueType{"val": StringType}, Returns: []ValueType{ErrorType}, isBuiltin: true}

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
