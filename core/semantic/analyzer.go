package semantic

import (
	"radlang/parser"
	"radlang/token"
)

func Analyze(p *parser.Program) (*Scope, []Diagnostic) {
	globalScope := Scope{SymbolTable: make(map[string]Symbol)}
	ctx := &SemanticCtx{Scope: &globalScope, Diagnostics: []Diagnostic{}}
	initBuiltins(ctx)
	RegisterProgram(ctx, p)
	AnalyzeProgram(ctx, p)
	return &globalScope, ctx.Diagnostics
}

func AnalyzeProgram(ctx *SemanticCtx, p *parser.Program) {
	AnalyzeFunctions(ctx, p.Functions)
}

func AnalyzeFunctions(ctx *SemanticCtx, functions []*parser.Func_Decl) {
	for _, function := range functions {
		slotCounter := 0
		ctx.CurrentFunc = ctx.Scope.SymbolTable[function.Signature.Name].(*FuncSymbol)

		childCtx := newChildCtx(ctx)
		childCtx.slotCounter = &slotCounter

		for name, typ := range ctx.CurrentFunc.Params {
			childCtx.Scope.SymbolTable[name] = &VarSymbol{Slot: getNextSlot(childCtx), Type: typ, Declared: true}
		}
		for _, stmt := range function.Body.Statement_Group {
			AnalyzeStatement(childCtx, stmt)
		}
		function.Body.Scope = childCtx.Scope
		ctx.CurrentFunc.Slots = *childCtx.slotCounter
		ctx.Diagnostics = append(ctx.Diagnostics, childCtx.Diagnostics...)
	}
}

func AnalyzeBlock(ctx *SemanticCtx, block *parser.Block) {
	childCtx := newChildCtx(ctx)
	for _, stmt := range block.Statement_Group {
		AnalyzeStatement(childCtx, stmt)
	}
	block.Scope = childCtx.Scope
	ctx.Diagnostics = append(ctx.Diagnostics, childCtx.Diagnostics...)

}

func AnalyzeStatement(ctx *SemanticCtx, stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.Assign_stmt:

		var exprTypes []ValueType
		// return inference
		for _, value := range s.Values {
			exprTypes = append(exprTypes, AnalyzeExpression(ctx, value)...)
		}

		if len(exprTypes) != len(s.Targets) {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrMismatchTypesInExpr))
			return
		}

		if s.Op == token.WALRUS {
			// track if new variables on left
			newVar := false
			for i, target := range s.Targets {
				// new variable on left
				if !symbolExistAs[(*VarSymbol)](ctx, target) {
					newVar = true
					ctx.Scope.SymbolTable[target] =
						&VarSymbol{
							Slot:     getNextSlot(ctx),
							Type:     exprTypes[i],
							Declared: true,
						}
					continue
				}

				if ctx.Scope.SymbolTable[target].(*VarSymbol).Type != exprTypes[i] {
					ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrMismatchTypesInExpr))
					continue
				}
			}
			if !newVar {
				ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrNoNewVariablesOnWalrus))
			}
		} else if s.Op == token.ASSIGNMENT {
			for i, target := range s.Targets {
				if !symbolExistAs[*VarSymbol](ctx, target) {
					ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrIdentNotDeclared))
					continue
				}
				if ctx.Scope.SymbolTable[target].(*VarSymbol).Type != exprTypes[i] {
					ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrMismatchTypesInExpr))
					continue
				}
			}
		}

	case *parser.Expr_stmt:
		switch expr := s.Expression.(type) {
		case *parser.Postfix_expr, *parser.Call_expr:
			AnalyzeExpression(ctx, expr)
		default:
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrInvalidExprStmt))
		}

	case *parser.Decl_stmt:
		for _, name := range s.Name {
			if _, ok := ctx.Scope.SymbolTable[name]; ok {
				ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrRedeclaredVariables))
			}
			ctx.Scope.SymbolTable[name] =
				&VarSymbol{
					Slot:     getNextSlot(ctx),
					Type:     resolveType(s.Type),
					Declared: true,
				}
		}

	case *parser.Control_stmt:
		exprType := AnalyzeExpression(ctx, s.Expression)
		if len(exprType) != 1 {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrExpectedOneExpr))
			return
		}
		if exprType[0] != BoolType {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrIfExpressionNotBool))
			return
		}
		AnalyzeBlock(ctx, s.IfBlock)
		if s.ElseStmt != nil {
			AnalyzeStatement(ctx, s.ElseStmt)
		}
		if s.ElseBlock != nil {
			AnalyzeBlock(ctx, s.ElseBlock)
		}

	case *parser.Jump_stmt:
		if ctx.LoopDepth <= 0 {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrJumpOutsideFor))
			return
		}
		ctx.LoopDepth -= 1

	case *parser.Return_stmt:
		var exprTypes []ValueType
		for _, expr := range s.Returns {
			exprTypes = append(exprTypes, AnalyzeExpression(ctx, expr)...)
		}

		if len(exprTypes) < len(ctx.CurrentFunc.Returns) {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrNotEnoughReturnValues))
			return
		}
		if len(exprTypes) > len(ctx.CurrentFunc.Returns) {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrTooManyReturnValues))
			return
		}

		// match exprTypes with ctx.CurrentFunc.Returns
		for i, returnVal := range exprTypes {
			if returnVal != ctx.CurrentFunc.Returns[i] {
				ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrBadReturnType))
				return
			}
		}

	case *parser.Loop_stmt:
		exprType := AnalyzeExpression(ctx, s.Expression)
		if len(exprType) != 1 {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrExpectedOneExpr))
			return
		}
		// if exprType[0] != token.BOOL {
		// 	ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrIfExpressionNotBool))
		// }
		ctx.LoopDepth += 1
		AnalyzeBlock(ctx, s.Loop_block)
	}
}

func AnalyzeExpression(ctx *SemanticCtx, expr parser.Expression) []ValueType {
	var Type []ValueType
	switch e := expr.(type) {
	case *parser.Lit_val:
		var val interface{}
		var ok bool
		switch e.Type {
		case token.INT:
			val, ok = e.Value.(int64)
		case token.FLOAT:
			val, ok = e.Value.(float64)
		case token.BOOL:
			val, ok = e.Value.(bool)
		case token.STRING:
			val, ok = e.Value.(string)
		default:
			panic("unknown type found in literal value")
		}
		if !ok {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrMismatchTypesInExpr))
		}

		e.Value = val
		Type = append(Type, resolveType(e.Type))

	case *parser.Binary_expr:
		leftTok := AnalyzeExpression(ctx, e.Left)
		rightTok := AnalyzeExpression(ctx, e.Right)
		if len(leftTok) != 1 || len(rightTok) != 1 {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrExpectedOneExpr))
			return nil
		}
		if rightTok[0] != leftTok[0] {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrMismatchTypesInExpr))
		}

		switch e.Op {
		case token.EQ, token.NEQ, token.GT, token.GTE, token.LT, token.LTE:
			Type = append(Type, BoolType)
		default:
			Type = append(Type, leftTok...)
		}

	case *parser.Call_expr:
		sym, err := Resolve(ctx.Scope, e.Name)
		funcSym, ok := sym.(*FuncSymbol)
		if err != nil || !ok {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrUndefined))
			return nil
		}
		Type = append(Type, funcSym.Returns...)

	case *parser.Identifier_expr:
		sym, err := Resolve(ctx.Scope, e.Name)
		if err != nil {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrUndefined))
			return nil
		}
		varSym, ok := sym.(*VarSymbol)
		if !ok {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrUndefined))
			return nil
		}
		Type = append(Type, varSym.Type)

	case *parser.Postfix_expr:
		if !symbolExistAs[*VarSymbol](ctx, e.Target.Name) {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrIdentNotDeclared))
			return nil
		}
		varSym, _ := Resolve(ctx.Scope, e.Target.Name)
		if varSym.(*VarSymbol).Type != IntType && varSym.(*VarSymbol).Type != FloatType {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrPostfixOnNonNumeric))
			return nil
		}

	}
	return Type
}
