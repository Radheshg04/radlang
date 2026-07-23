package semantic

import (
	"radlang/parser"
	"radlang/token"
	"testing"
)

func funcWith(name string, stmts ...parser.Statement) *parser.Program {
	return &parser.Program{
		Functions: []*parser.Func_Decl{
			{
				Signature: &parser.Func_Signature{Name: name},
				Body:      &parser.Block{Statement_Group: stmts},
			},
		},
	}
}

func hasErrors(diags []Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == Error {
			return true
		}
	}
	return false
}

func TestAnalyzeValidProgram(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: []string{"x"}, Type: token.INT},
		&parser.Assign_stmt{
			Targets: []string{"x"},
			Values:  []parser.Expression{&parser.Lit_val{Value: "5", Type: token.INT}},
			Op:      token.ASSIGNMENT,
		},
	)
	if _, diags := Analyze(prog); hasErrors(diags) {
		t.Fatalf("unexpected errors: %v", diags)
	}
}

func TestAnalyzeRedeclaredFunction(t *testing.T) {
	prog := &parser.Program{
		Functions: []*parser.Func_Decl{
			{Signature: &parser.Func_Signature{Name: "foo"}, Body: &parser.Block{}},
			{Signature: &parser.Func_Signature{Name: "foo"}, Body: &parser.Block{}},
		},
	}
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for redeclared function")
	}
}

func TestAnalyzeRedeclaredVariable(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: []string{"x"}, Type: token.INT},
		&parser.Decl_stmt{Name: []string{"x"}, Type: token.INT},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for redeclared variable")
	}
}

func TestAnalyzeAssignUndeclaredVariable(t *testing.T) {
	prog := funcWith("main",
		&parser.Assign_stmt{
			Targets: []string{"x"},
			Values:  []parser.Expression{&parser.Lit_val{Value: "1", Type: token.INT}},
			Op:      token.ASSIGNMENT,
		},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for assignment to undeclared variable")
	}
}

func TestAnalyzeWalrusValid(t *testing.T) {
	prog := funcWith("main",
		&parser.Assign_stmt{
			Targets: []string{"x"},
			Values:  []parser.Expression{&parser.Lit_val{Value: "5", Type: token.INT}},
			Op:      token.WALRUS,
		},
	)
	if _, diags := Analyze(prog); hasErrors(diags) {
		t.Fatalf("unexpected errors for valid walrus: %v", diags)
	}
}

func TestAnalyzeWalrusNoNewVars(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: []string{"x"}, Type: token.INT},
		&parser.Assign_stmt{
			Targets: []string{"x"},
			Values:  []parser.Expression{&parser.Lit_val{Value: "5", Type: token.INT}},
			Op:      token.WALRUS,
		},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for walrus with no new variables")
	}
}

func TestAnalyzeCallUndeclaredFunction(t *testing.T) {
	prog := funcWith("main",
		&parser.Expr_stmt{Expression: &parser.Call_expr{Name: "notDefined"}},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for call to undeclared function")
	}
}

func TestAnalyzeCallBuiltin(t *testing.T) {
	prog := funcWith("main",
		&parser.Expr_stmt{Expression: &parser.Call_expr{Name: "print"}},
	)
	if _, diags := Analyze(prog); hasErrors(diags) {
		t.Fatalf("unexpected errors for builtin print call: %v", diags)
	}
}

func TestAnalyzeCallDeclaredFunction(t *testing.T) {
	prog := &parser.Program{
		Functions: []*parser.Func_Decl{
			{Signature: &parser.Func_Signature{Name: "helper"}, Body: &parser.Block{}},
			{
				Signature: &parser.Func_Signature{Name: "main"},
				Body: &parser.Block{
					Statement_Group: []parser.Statement{
						&parser.Expr_stmt{Expression: &parser.Call_expr{Name: "helper"}},
					},
				},
			},
		},
	}
	if _, diags := Analyze(prog); hasErrors(diags) {
		t.Fatalf("unexpected errors for call to declared function: %v", diags)
	}
}

func TestAnalyzeBinaryExprTypeMismatch(t *testing.T) {
	prog := funcWith("main",
		&parser.Assign_stmt{
			Targets: []string{"x"},
			Values: []parser.Expression{
				&parser.Binary_expr{
					Left:  &parser.Lit_val{Value: "1", Type: token.INT},
					Op:    token.PLUS,
					Right: &parser.Lit_val{Value: "2.0", Type: token.FLOAT},
				},
			},
			Op: token.WALRUS,
		},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for binary expr with mismatched types")
	}
}

func TestAnalyzeUndeclaredVarInExpr(t *testing.T) {
	prog := funcWith("main",
		&parser.Assign_stmt{
			Targets: []string{"x"},
			Values:  []parser.Expression{&parser.Identifier_expr{Name: "y"}},
			Op:      token.WALRUS,
		},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for undeclared variable in expression")
	}
}

func TestAnalyzeBreakOutsideLoop(t *testing.T) {
	prog := funcWith("main",
		&parser.Jump_stmt{Type: token.BREAK},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for break outside loop")
	}
}

func TestAnalyzeReturnTypeMismatch(t *testing.T) {
	prog := &parser.Program{
		Functions: []*parser.Func_Decl{
			{
				Signature: &parser.Func_Signature{Name: "foo", Returns: []token.TokenType{token.INT}},
				Body: &parser.Block{
					Statement_Group: []parser.Statement{
						&parser.Return_stmt{Returns: []parser.Expression{
							&parser.Lit_val{Value: "3.14", Type: token.FLOAT},
						}},
					},
				},
			},
		},
	}
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for return type mismatch")
	}
}

func TestAnalyzeIfConditionNotBool(t *testing.T) {
	prog := funcWith("main",
		&parser.Control_stmt{
			Expression: &parser.Lit_val{Value: "5", Type: token.INT},
			IfBlock:    &parser.Block{},
		},
	)
	if _, diags := Analyze(prog); !hasErrors(diags) {
		t.Fatal("expected error for non-bool if condition")
	}
}
