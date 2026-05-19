package interpreter

import (
	"radlang/parser"
	"radlang/token"
	"testing"
)

func resetEnv() {
	env = make(map[string]Value)
}

func progMain(stmts ...parser.Statement) *parser.Program {
	return &parser.Program{
		Functions: []*parser.Func_Decl{
			{Name: "main", Body: &parser.Block{Statement_Group: stmts}},
		},
	}
}

func TestInterpretDeclInt(t *testing.T) {
	resetEnv()
	prog := progMain(&parser.Decl_stmt{Name: "x", Type: token.INT})
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["x"].(IntValue); !ok {
		t.Errorf("expected IntValue for x, got %T", env["x"])
	}
}

func TestInterpretDeclFloat(t *testing.T) {
	resetEnv()
	prog := progMain(&parser.Decl_stmt{Name: "f", Type: token.FLOAT})
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["f"].(FloatValue); !ok {
		t.Errorf("expected FloatValue for f, got %T", env["f"])
	}
}

func TestInterpretDeclBool(t *testing.T) {
	resetEnv()
	prog := progMain(&parser.Decl_stmt{Name: "b", Type: token.BOOL})
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["b"].(BoolValue); !ok {
		t.Errorf("expected BoolValue for b, got %T", env["b"])
	}
}

func TestInterpretDeclString(t *testing.T) {
	resetEnv()
	prog := progMain(&parser.Decl_stmt{Name: "s", Type: token.STRING})
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["s"].(StringValue); !ok {
		t.Errorf("expected StringValue for s, got %T", env["s"])
	}
}

func TestInterpretAssignInt(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{
			Target: "x",
			Value:  &parser.Number_lit{Value: "10", Type: token.INT},
		},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["x"].(IntValue).Val != 10 {
		t.Errorf("x = %d, want 10", env["x"].(IntValue).Val)
	}
}

func TestInterpretAssignFloat(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "f", Type: token.FLOAT},
		&parser.Assign_stmt{
			Target: "f",
			Value:  &parser.Number_lit{Value: "2.5", Type: token.FLOAT},
		},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["f"].(FloatValue).Val != 2.5 {
		t.Errorf("f = %v, want 2.5", env["f"].(FloatValue).Val)
	}
}

func TestInterpretBinaryAdd(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "z", Type: token.INT},
		&parser.Assign_stmt{
			Target: "z",
			Value: &parser.Binary_expr{
				Left:  &parser.Number_lit{Value: "3", Type: token.INT},
				Op:    token.PLUS,
				Right: &parser.Number_lit{Value: "4", Type: token.INT},
			},
		},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["z"].(IntValue).Val != 7 {
		t.Errorf("z = %d, want 7", env["z"].(IntValue).Val)
	}
}

func TestInterpretBinaryMul(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "z", Type: token.INT},
		&parser.Assign_stmt{
			Target: "z",
			Value: &parser.Binary_expr{
				Left:  &parser.Number_lit{Value: "6", Type: token.INT},
				Op:    token.ASTERISK,
				Right: &parser.Number_lit{Value: "7", Type: token.INT},
			},
		},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["z"].(IntValue).Val != 42 {
		t.Errorf("z = %d, want 42", env["z"].(IntValue).Val)
	}
}

func TestInterpretIncrement(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{
			Target: "x",
			Value:  &parser.Number_lit{Value: "5", Type: token.INT},
		},
		&parser.Update_stmt{Target: "x", Op: token.PLUSPLUS},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["x"].(IntValue).Val != 6 {
		t.Errorf("x = %d, want 6", env["x"].(IntValue).Val)
	}
}

func TestInterpretDecrement(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "x", Type: token.INT},
		&parser.Assign_stmt{
			Target: "x",
			Value:  &parser.Number_lit{Value: "5", Type: token.INT},
		},
		&parser.Update_stmt{Target: "x", Op: token.MINUSMINUS},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["x"].(IntValue).Val != 4 {
		t.Errorf("x = %d, want 4", env["x"].(IntValue).Val)
	}
}

func TestInterpretFloatIncrement(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "f", Type: token.FLOAT},
		&parser.Assign_stmt{
			Target: "f",
			Value:  &parser.Number_lit{Value: "1.5", Type: token.FLOAT},
		},
		&parser.Update_stmt{Target: "f", Op: token.PLUSPLUS},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["f"].(FloatValue).Val != 2.5 {
		t.Errorf("f = %v, want 2.5", env["f"].(FloatValue).Val)
	}
}

func TestInterpretIdentifierExpr(t *testing.T) {
	resetEnv()
	prog := progMain(
		&parser.Decl_stmt{Name: "a", Type: token.INT},
		&parser.Assign_stmt{
			Target: "a",
			Value:  &parser.Number_lit{Value: "9", Type: token.INT},
		},
		&parser.Decl_stmt{Name: "b", Type: token.INT},
		&parser.Assign_stmt{
			Target: "b",
			Value:  &parser.Identifier_expr{Name: "a"},
		},
	)
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["b"].(IntValue).Val != 9 {
		t.Errorf("b = %d, want 9", env["b"].(IntValue).Val)
	}
}

func TestInterpretSkipsNonMainFunctions(t *testing.T) {
	resetEnv()
	prog := &parser.Program{
		Functions: []*parser.Func_Decl{
			{
				Name: "other",
				Body: &parser.Block{
					Statement_Group: []parser.Statement{
						&parser.Decl_stmt{Name: "shouldNotExist", Type: token.INT},
					},
				},
			},
		},
	}
	if err := Interpret(prog); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := env["shouldNotExist"]; exists {
		t.Error("non-main function should not be executed")
	}
}
