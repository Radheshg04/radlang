package semantic

import (
	"fmt"
	"radlang/parser"
	"radlang/token"
)

func Analyze(p *parser.Program) []Diagnostic {
	globalScope := Scope{SymbolTable: make(map[string]Symbol)}
	ctx := &SemanticCtx{Scope: &globalScope, Diagnostics: []Diagnostic{}}
	initBuiltins(ctx)
	RegisterProgram(ctx, p)
	AnalyzeProgram(ctx, p)
	return ctx.Diagnostics
}

func AnalyzeProgram(ctx *SemanticCtx, p *parser.Program) {
	AnalyzeFunctions(ctx, p.Functions)
}

func AnalyzeFunctions(ctx *SemanticCtx, functions []*parser.Func_Decl) {
	for _, function := range functions {
		ctx.CurrentFunc = ctx.Scope.SymbolTable[function.Signature.Name].(*FuncSymbol)

		childCtx := newChildCtx(ctx)
		for name, typ := range ctx.CurrentFunc.Params {
			childCtx.Scope.SymbolTable[name] = &VarSymbol{Type: typ, Declared: true}
		}
		for _, stmt := range function.Body.Statement_Group {
			AnalyzeStatement(childCtx, stmt)
		}
		ctx.Diagnostics = append(ctx.Diagnostics, childCtx.Diagnostics...)
	}
}

func AnalyzeBlock(ctx *SemanticCtx, block *parser.Block) {
	childCtx := newChildCtx(ctx)
	for _, stmt := range block.Statement_Group {
		AnalyzeStatement(childCtx, stmt)
	}
	ctx.Diagnostics = append(ctx.Diagnostics, childCtx.Diagnostics...)

}

func AnalyzeStatement(ctx *SemanticCtx, stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.Assign_stmt:
		if s.Op == token.WALRUS {
			// track if new variables on left
			newVar := false
			var exprTypes []token.TokenType
			// return inference
			for _, value := range s.Values {
				exprTypes = append(exprTypes, AnalyzeExpression(ctx, value)...)
			}

			if len(exprTypes) != len(s.Targets) {
				ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrMismatchTypesInExpr))
				return
			}

			for i, target := range s.Targets {
				// new variable on left
				if !symbolExistAs[(*VarSymbol)](ctx, target) {
					newVar = true
				}

				ctx.Scope.SymbolTable[target] =
					&VarSymbol{
						Type:     exprTypes[i],
						Value:    s.Values[i],
						Declared: true,
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
				ctx.Scope.SymbolTable[target].(*VarSymbol).Value = s.Values[i]

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
					Type:     s.Type,
					Declared: true,
				}
		}

	case *parser.Control_stmt:
		exprType := AnalyzeExpression(ctx, s.Expression)
		if len(exprType) != 1 {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrExpectedOneExpr))
			return
		}
		if exprType[0] != token.BOOL_LIT {
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
		var exprTypes []token.TokenType
		for _, expr := range s.Returns {
			exprTypes = append(exprTypes, AnalyzeExpression(ctx, expr)...)
		}

		if len(exprTypes) < len(ctx.CurrentFunc.Returns) {
			fmt.Print(exprTypes, ctx.CurrentFunc.Returns)
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

func AnalyzeExpression(ctx *SemanticCtx, expr parser.Expression) []token.TokenType {
	var tokens []token.TokenType
	switch e := expr.(type) {
	case *parser.Lit_val:
		tokens = append(tokens, e.Type)

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
			tokens = append(tokens, token.BOOL_LIT)
		default:
			tokens = append(tokens, leftTok...)
		}

	case *parser.Call_expr:
		sym, err := resolve(ctx.Scope, e.Name)
		funcSym, ok := sym.(*FuncSymbol)
		if err != nil || !ok {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrUndefined))
			return nil
		}
		tokens = append(tokens, funcSym.Returns...)

	case *parser.Identifier_expr:
		sym, err := resolve(ctx.Scope, e.Name)
		if err != nil {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrUndefined))
			return nil
		}
		varSym, ok := sym.(*VarSymbol)
		if !ok {
			ctx.Diagnostics = append(ctx.Diagnostics, *NewRLDiagnostic(ErrUndefined))
			return nil
		}
		tokens = append(tokens, varSym.Type)

	case *parser.Postfix_expr:

	}
	return tokens
}
