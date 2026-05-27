package semantic

import "radlang/token"

func initBuiltins(ctx *SemanticCtx) {
	printDecl := &FuncSymbol{Params: map[string]token.TokenType{"val": token.STRING_LIT}, Returns: nil}
	errDecl := &FuncSymbol{Params: map[string]token.TokenType{"val": token.STRING_LIT}, Returns: nil}

	builtins := map[string]*FuncSymbol{
		"print": printDecl,
		"err":   errDecl,
	}

	for funcName, funcSym := range builtins {
		ctx.Scope.SymbolTable[funcName] = funcSym
	}
}
