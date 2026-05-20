package semantic

import (
	"radlang/parser"
	"radlang/token"
	"testing"
)

func funcWith(name string, stmts ...parser.Statement) *parser.Program {
	return &parser.Program{
		Functions: []*parser.Func_Decl{
			{Name: name, Body: &parser.Block{Statement_Group: stmts}},
		},
	}
}

func TestAnalyzeValidProgram(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{Target: "x", Value: &parser.Number_lit{Value: "5", Type: token.INT}},
	)
	if err := Analyze(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnalyzeRedeclaredFunction(t *testing.T) {
	prog := &parser.Program{
		Functions: []*parser.Func_Decl{
			{Name: "foo", Body: &parser.Block{}},
			{Name: "foo", Body: &parser.Block{}},
		},
	}
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for redeclared function")
	}
}

func TestAnalyzeRedeclaredVariable(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Decl_stmt{Name: "x", Type: token.INT},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for redeclared variable")
	}
}

func TestAnalyzeAssignBeforeDeclaration(t *testing.T) {
	prog := funcWith("main",
		&parser.Assign_stmt{Target: "x", Value: &parser.Number_lit{Value: "1", Type: token.INT}},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for assignment before declaration")
	}
}

func TestAnalyzeTypeMismatch(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{Target: "x", Value: &parser.Number_lit{Value: "3.14", Type: token.FLOAT}},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for type mismatch")
	}
}

func TestAnalyzeUpdateUndeclaredVariable(t *testing.T) {
	prog := funcWith("main",
		&parser.Update_stmt{Target: "x", Op: token.PLUSPLUS},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for update on undeclared variable")
	}
}

func TestAnalyzeUpdateNonNumeric(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "s", Type: token.STRING},
		&parser.Update_stmt{Target: "s", Op: token.PLUSPLUS},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for update on string variable")
	}
}

func TestAnalyzeUpdateInt(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "i", Type: token.INT},
		&parser.Update_stmt{Target: "i", Op: token.PLUSPLUS},
	)
	if err := Analyze(prog); err != nil {
		t.Fatalf("unexpected error for int update: %v", err)
	}
}

func TestAnalyzeUpdateFloat(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "f", Type: token.FLOAT},
		&parser.Update_stmt{Target: "f", Op: token.MINUSMINUS},
	)
	if err := Analyze(prog); err != nil {
		t.Fatalf("unexpected error for float update: %v", err)
	}
}

func TestAnalyzeCallUndeclaredFunction(t *testing.T) {
	prog := funcWith("main",
		&parser.Expr_stmt{Expr: &parser.Call_expr{Name: "notDefined"}},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for call to undeclared function")
	}
}

func TestAnalyzeCallBuiltin(t *testing.T) {
	prog := funcWith("main",
		&parser.Expr_stmt{Expr: &parser.Call_expr{Name: "print"}},
	)
	if err := Analyze(prog); err != nil {
		t.Fatalf("unexpected error for builtin print call: %v", err)
	}
}

func TestAnalyzeCallDeclaredFunction(t *testing.T) {
	prog := &parser.Program{
		Functions: []*parser.Func_Decl{
			{Name: "helper", Body: &parser.Block{}},
			{Name: "main", Body: &parser.Block{
				Statement_Group: []parser.Statement{
					&parser.Expr_stmt{Expr: &parser.Call_expr{Name: "helper"}},
				},
			}},
		},
	}
	if err := Analyze(prog); err != nil {
		t.Fatalf("unexpected error for call to declared function: %v", err)
	}
}

func TestAnalyzeBinaryExprTypeMismatch(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{
			Target: "x",
			Value: &parser.Binary_expr{
				Left:  &parser.Number_lit{Value: "1", Type: token.INT},
				Op:    token.PLUS,
				Right: &parser.Number_lit{Value: "2.0", Type: token.FLOAT},
			},
		},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for binary expr with mismatched types")
	}
}

func TestAnalyzeUndeclaredVarInExpr(t *testing.T) {
	prog := funcWith("main",
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{
			Target: "x",
			Value:  &parser.Identifier_expr{Name: "y"},
		},
	)
	if err := Analyze(prog); err == nil {
		t.Fatal("expected error for undeclared variable in expression")
	}
}
