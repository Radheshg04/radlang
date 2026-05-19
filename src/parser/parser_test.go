package parser

import (
	"radlang/token"
	"testing"
)

func mustParse(t *testing.T, tokens []token.Token) *Program {
	t.Helper()
	prog, err := Parse(tokens)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return prog
}

func makeTokens(pairs ...interface{}) []token.Token {
	var tokens []token.Token
	for i := 0; i < len(pairs); i += 2 {
		tokens = append(tokens, token.Token{
			Token:  pairs[i].(token.TokenType),
			Lexeme: pairs[i+1].(string),
			Line:   1,
		})
	}
	return tokens
}

func TestParseEmptyFunction(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	if len(prog.Functions) != 1 {
		t.Fatalf("want 1 function, got %d", len(prog.Functions))
	}
	if prog.Functions[0].Name != "main" {
		t.Errorf("function name: got %q, want %q", prog.Functions[0].Name, "main")
	}
	if len(prog.Functions[0].Body.Statement_Group) != 0 {
		t.Errorf("want empty body, got %d statements", len(prog.Functions[0].Body.Statement_Group))
	}
}

func TestParseDeclStatement(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.VAR, "var",
		token.IDENTIFIER, "x",
		token.INT, "int",
		token.EOL, "",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	stmts := prog.Functions[0].Body.Statement_Group
	if len(stmts) != 1 {
		t.Fatalf("want 1 statement, got %d", len(stmts))
	}
	decl, ok := stmts[0].(*Decl_stmt)
	if !ok {
		t.Fatalf("want *Decl_stmt, got %T", stmts[0])
	}
	if decl.Name != "x" {
		t.Errorf("decl name: got %q, want %q", decl.Name, "x")
	}
	if decl.Type != token.INT {
		t.Errorf("decl type: got %v, want INT", decl.Type)
	}
}

func TestParseAssignStatement(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.IDENTIFIER, "x",
		token.ASSIGNMENT, "=",
		token.INT, "5",
		token.EOL, "",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	stmts := prog.Functions[0].Body.Statement_Group
	assign, ok := stmts[0].(*Assign_stmt)
	if !ok {
		t.Fatalf("want *Assign_stmt, got %T", stmts[0])
	}
	if assign.Target != "x" {
		t.Errorf("assign target: got %q, want %q", assign.Target, "x")
	}
	lit, ok := assign.Value.(*Number_lit)
	if !ok {
		t.Fatalf("want *Number_lit, got %T", assign.Value)
	}
	if lit.Value != "5" {
		t.Errorf("number value: got %q, want %q", lit.Value, "5")
	}
}

func TestParseUpdateStatement(t *testing.T) {
	cases := []struct {
		opTok token.TokenType
		opLex string
	}{
		{token.PLUSPLUS, "++"},
		{token.MINUSMINUS, "--"},
	}
	for _, c := range cases {
		tokens := makeTokens(
			token.FUNC, "func",
			token.IDENTIFIER, "main",
			token.L_PAREN, "(",
			token.R_PAREN, ")",
			token.L_BRACE, "{",
			token.IDENTIFIER, "x",
			c.opTok, c.opLex,
			token.EOL, "",
			token.R_BRACE, "}",
			token.EOF, "",
		)
		prog := mustParse(t, tokens)
		upd, ok := prog.Functions[0].Body.Statement_Group[0].(*Update_stmt)
		if !ok {
			t.Fatalf("want *Update_stmt for %v", c.opLex)
		}
		if upd.Op != c.opTok {
			t.Errorf("update op: got %v, want %v", upd.Op, c.opTok)
		}
	}
}

func TestParseBinaryExpression(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.IDENTIFIER, "z",
		token.ASSIGNMENT, "=",
		token.INT, "1",
		token.PLUS, "+",
		token.INT, "2",
		token.EOL, "",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	assign := prog.Functions[0].Body.Statement_Group[0].(*Assign_stmt)
	bin, ok := assign.Value.(*Binary_expr)
	if !ok {
		t.Fatalf("want *Binary_expr, got %T", assign.Value)
	}
	if bin.Op != token.PLUS {
		t.Errorf("binary op: got %v, want PLUS", bin.Op)
	}
}

func TestParseCallExpression(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.IDENTIFIER, "print",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.EOL, "",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	expr := prog.Functions[0].Body.Statement_Group[0].(*Expr_stmt)
	call, ok := expr.Expr.(*Call_expr)
	if !ok {
		t.Fatalf("want *Call_expr, got %T", expr.Expr)
	}
	if call.Name != "print" {
		t.Errorf("call name: got %q, want %q", call.Name, "print")
	}
}

func TestParseMultipleFunctions(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "foo",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.FUNC, "func",
		token.IDENTIFIER, "bar",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	if len(prog.Functions) != 2 {
		t.Fatalf("want 2 functions, got %d", len(prog.Functions))
	}
}

func TestParseErrorUnexpectedToken(t *testing.T) {
	tokens := makeTokens(
		token.INT, "int",
		token.EOF, "",
	)
	_, err := Parse(tokens)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestParseErrorUnclosedBlock(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.EOF, "",
	)
	_, err := Parse(tokens)
	if err == nil {
		t.Fatal("expected parse error for unclosed block, got nil")
	}
}
