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
	if function, ok := ctx.Scope.SymbolTable[function.Signature.Name]; ok {
		if isBuiltinFunc(function) {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrRedeclaredBuiltinFunction))
			return
		}
		ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrRedeclaredSameFunc))
		return
	}

	// params can be nil
	params := resolveParams(function.Signature.Params)

	funcSym := &FuncSymbol{Params: params, Returns: function.Signature.Returns, Decl: function}
	ctx.CurrentFunc, ctx.Scope.SymbolTable[function.Signature.Name] = funcSym, funcSym

	childCtx := newChildCtx(ctx)
	for name, typ := range ctx.CurrentFunc.Params {
		childCtx.Scope.SymbolTable[name] = &VarSymbol{Type: typ, Declared: true}
	}
	ctx.Diagnostics = append(ctx.Diagnostics, childCtx.Diagnostics...)
}
