package semantic

import (
	"radlang/parser"
)

func RegisterProgram(ctx *SemanticCtx, p *parser.Program) {
	// Register functions
	for _, function := range p.Functions {
		RegisterFunction(ctx, function)
	}
}

func RegisterFunction(ctx *SemanticCtx, function *parser.Func_Decl) {
	if existing, ok := ctx.Scope.SymbolTable[function.Signature.Name]; ok {
		if isBuiltinFunc(existing) {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrRedeclaredBuiltinFunction))
			return
		}
		ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrRedeclaredSameFunc))
		return
	}

	// params can be nil
	params := resolveParams(function.Signature.Params)
	returns := resolveTypes(function.Signature.Returns)

	funcSym := &FuncSymbol{Params: params, Returns: returns, Decl: function}
	ctx.Scope.SymbolTable[function.Signature.Name] = funcSym
}
