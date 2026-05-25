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

// funcWrap wraps pairs in `func main() { <pairs> }` token stream.
func funcWrap(inner ...interface{}) []token.Token {
	pairs := []interface{}{
		token.FUNC, "func",
		token.IDENTIFIER, "main",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
	}
	pairs = append(pairs, inner...)
	pairs = append(pairs, token.R_BRACE, "}", token.EOF, "")
	return makeTokens(pairs...)
}

func stmts(t *testing.T, tokens []token.Token) []Statement {
	t.Helper()
	prog := mustParse(t, tokens)
	return prog.Functions[0].Body.Statement_Group
}

// --- Function declarations ---

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
	if prog.Functions[0].Signature.Name != "main" {
		t.Errorf("function name: got %q, want %q", prog.Functions[0].Signature.Name, "main")
	}
	if len(prog.Functions[0].Body.Statement_Group) != 0 {
		t.Errorf("want empty body, got %d statements", len(prog.Functions[0].Body.Statement_Group))
	}
}

func TestParseFuncWithParams(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "add",
		token.L_PAREN, "(",
		token.IDENTIFIER, "a",
		token.INT, "int",
		token.COMMA, ",",
		token.IDENTIFIER, "b",
		token.INT, "int",
		token.R_PAREN, ")",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	sig := prog.Functions[0].Signature
	if sig.Name != "add" {
		t.Errorf("name: got %q, want add", sig.Name)
	}
	if len(sig.Params) != 2 {
		t.Fatalf("want 2 params, got %d", len(sig.Params))
	}
	if sig.Params[0].Name != "a" || sig.Params[0].Type != token.INT {
		t.Errorf("param[0]: got {%q, %v}", sig.Params[0].Name, sig.Params[0].Type)
	}
	if sig.Params[1].Name != "b" || sig.Params[1].Type != token.INT {
		t.Errorf("param[1]: got {%q, %v}", sig.Params[1].Name, sig.Params[1].Type)
	}
}

func TestParseFuncWithReturnType(t *testing.T) {
	tokens := makeTokens(
		token.FUNC, "func",
		token.IDENTIFIER, "getX",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.INT, "int",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.EOF, "",
	)
	prog := mustParse(t, tokens)
	sig := prog.Functions[0].Signature
	if len(sig.Returns) != 1 || sig.Returns[0] != token.INT {
		t.Errorf("returns: got %v, want [INT]", sig.Returns)
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
	if prog.Functions[0].Signature.Name != "foo" {
		t.Errorf("func[0]: got %q, want foo", prog.Functions[0].Signature.Name)
	}
	if prog.Functions[1].Signature.Name != "bar" {
		t.Errorf("func[1]: got %q, want bar", prog.Functions[1].Signature.Name)
	}
}

// --- Decl statements ---

func TestParseDeclStatement(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.VAR, "var",
		token.IDENTIFIER, "x",
		token.INT, "int",
		token.EOL, "",
	))
	if len(ss) != 1 {
		t.Fatalf("want 1 statement, got %d", len(ss))
	}
	decl, ok := ss[0].(*Decl_stmt)
	if !ok {
		t.Fatalf("want *Decl_stmt, got %T", ss[0])
	}
	if len(decl.Name) != 1 || decl.Name[0] != "x" {
		t.Errorf("decl name: got %v, want [x]", decl.Name)
	}
	if decl.Type != token.INT {
		t.Errorf("decl type: got %v, want INT", decl.Type)
	}
}

func TestParseDeclMultiName(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.VAR, "var",
		token.IDENTIFIER, "x",
		token.COMMA, ",",
		token.IDENTIFIER, "y",
		token.INT, "int",
		token.EOL, "",
	))
	decl, ok := ss[0].(*Decl_stmt)
	if !ok {
		t.Fatalf("want *Decl_stmt, got %T", ss[0])
	}
	if len(decl.Name) != 2 || decl.Name[0] != "x" || decl.Name[1] != "y" {
		t.Errorf("decl names: got %v, want [x y]", decl.Name)
	}
	if decl.Type != token.INT {
		t.Errorf("decl type: got %v, want INT", decl.Type)
	}
}

// --- Assign statements ---

func TestParseAssignStatement(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "x",
		token.ASSIGNMENT, "=",
		token.INT_LIT, "5",
		token.EOL, "",
	))
	assign, ok := ss[0].(*Assign_stmt)
	if !ok {
		t.Fatalf("want *Assign_stmt, got %T", ss[0])
	}
	if len(assign.Targets) != 1 || assign.Targets[0] != "x" {
		t.Errorf("targets: got %v, want [x]", assign.Targets)
	}
	if assign.Op != token.ASSIGNMENT {
		t.Errorf("op: got %v, want ASSIGNMENT", assign.Op)
	}
	if len(assign.Values) != 1 {
		t.Fatalf("want 1 value, got %d", len(assign.Values))
	}
	lit, ok := assign.Values[0].(*Lit_val)
	if !ok {
		t.Fatalf("want *Lit_val, got %T", assign.Values[0])
	}
	if lit.Value != "5" || lit.Type != token.INT {
		t.Errorf("lit: got {%v, %v}, want {5, INT}", lit.Value, lit.Type)
	}
}

func TestParseWalrusAssign(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "x",
		token.WALRUS, ":=",
		token.INT_LIT, "42",
		token.EOL, "",
	))
	assign, ok := ss[0].(*Assign_stmt)
	if !ok {
		t.Fatalf("want *Assign_stmt, got %T", ss[0])
	}
	if assign.Op != token.WALRUS {
		t.Errorf("op: got %v, want WALRUS", assign.Op)
	}
}

func TestParseMultiAssign(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "x",
		token.COMMA, ",",
		token.IDENTIFIER, "y",
		token.ASSIGNMENT, "=",
		token.INT_LIT, "1",
		token.COMMA, ",",
		token.INT_LIT, "2",
		token.EOL, "",
	))
	assign, ok := ss[0].(*Assign_stmt)
	if !ok {
		t.Fatalf("want *Assign_stmt, got %T", ss[0])
	}
	if len(assign.Targets) != 2 || assign.Targets[0] != "x" || assign.Targets[1] != "y" {
		t.Errorf("targets: got %v, want [x y]", assign.Targets)
	}
	if len(assign.Values) != 2 {
		t.Fatalf("want 2 values, got %d", len(assign.Values))
	}
}

// --- Expressions ---

func TestParseBinaryExpression(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "z",
		token.ASSIGNMENT, "=",
		token.INT_LIT, "1",
		token.PLUS, "+",
		token.INT_LIT, "2",
		token.EOL, "",
	))
	assign := ss[0].(*Assign_stmt)
	bin, ok := assign.Values[0].(*Binary_expr)
	if !ok {
		t.Fatalf("want *Binary_expr, got %T", assign.Values[0])
	}
	if bin.Op != token.PLUS {
		t.Errorf("op: got %v, want PLUS", bin.Op)
	}
	left, ok := bin.Left.(*Lit_val)
	if !ok || left.Value != "1" {
		t.Errorf("left: got %T %v, want Lit_val 1", bin.Left, bin.Left)
	}
	right, ok := bin.Right.(*Lit_val)
	if !ok || right.Value != "2" {
		t.Errorf("right: got %T %v, want Lit_val 2", bin.Right, bin.Right)
	}
}

func TestParseArithmeticPrecedence(t *testing.T) {
	// z = 1 + 2 * 3 → (1 + (2 * 3))
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "z",
		token.ASSIGNMENT, "=",
		token.INT_LIT, "1",
		token.PLUS, "+",
		token.INT_LIT, "2",
		token.ASTERISK, "*",
		token.INT_LIT, "3",
		token.EOL, "",
	))
	assign := ss[0].(*Assign_stmt)
	add, ok := assign.Values[0].(*Binary_expr)
	if !ok || add.Op != token.PLUS {
		t.Fatalf("want PLUS at root, got %T", assign.Values[0])
	}
	mul, ok := add.Right.(*Binary_expr)
	if !ok || mul.Op != token.ASTERISK {
		t.Fatalf("want ASTERISK on right, got %T", add.Right)
	}
}

func TestParseComparisonExpr(t *testing.T) {
	cases := []struct {
		tok token.TokenType
		lex string
	}{
		{token.EQ, "=="},
		{token.NEQ, "!="},
		{token.GT, ">"},
		{token.GTE, ">="},
		{token.LT, "<"},
		{token.LTE, "<="},
	}
	for _, c := range cases {
		ss := stmts(t, funcWrap(
			token.IDENTIFIER, "z",
			token.ASSIGNMENT, "=",
			token.INT_LIT, "1",
			c.tok, c.lex,
			token.INT_LIT, "2",
			token.EOL, "",
		))
		assign := ss[0].(*Assign_stmt)
		bin, ok := assign.Values[0].(*Binary_expr)
		if !ok {
			t.Fatalf("%v: want *Binary_expr, got %T", c.lex, assign.Values[0])
		}
		if bin.Op != c.tok {
			t.Errorf("%v: op: got %v", c.lex, bin.Op)
		}
	}
}

func TestParseGroupedExpr(t *testing.T) {
	// z = (1 + 2) * 3 → ASTERISK at root
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "z",
		token.ASSIGNMENT, "=",
		token.L_PAREN, "(",
		token.INT_LIT, "1",
		token.PLUS, "+",
		token.INT_LIT, "2",
		token.R_PAREN, ")",
		token.ASTERISK, "*",
		token.INT_LIT, "3",
		token.EOL, "",
	))
	assign := ss[0].(*Assign_stmt)
	mul, ok := assign.Values[0].(*Binary_expr)
	if !ok || mul.Op != token.ASTERISK {
		t.Fatalf("want ASTERISK at root, got %T", assign.Values[0])
	}
	add, ok := mul.Left.(*Binary_expr)
	if !ok || add.Op != token.PLUS {
		t.Fatalf("want PLUS on left, got %T", mul.Left)
	}
}

func TestParseLiterals(t *testing.T) {
	cases := []struct {
		tok   token.TokenType
		lex   string
		wtype token.TokenType
	}{
		{token.INT_LIT, "7", token.INT},
		{token.FLOAT_LIT, "3.14", token.FLOAT},
		{token.STRING_LIT, "hello", token.STRING},
		{token.BOOL_LIT, "true", token.BOOL},
	}
	for _, c := range cases {
		ss := stmts(t, funcWrap(
			token.IDENTIFIER, "z",
			token.ASSIGNMENT, "=",
			c.tok, c.lex,
			token.EOL, "",
		))
		assign := ss[0].(*Assign_stmt)
		lit, ok := assign.Values[0].(*Lit_val)
		if !ok {
			t.Fatalf("%v: want *Lit_val, got %T", c.tok, assign.Values[0])
		}
		if lit.Value != c.lex || lit.Type != c.wtype {
			t.Errorf("%v: got {%v, %v}, want {%v, %v}", c.tok, lit.Value, lit.Type, c.lex, c.wtype)
		}
	}
}

func TestParseErrLiteral(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "z",
		token.ASSIGNMENT, "=",
		token.ERR, "err",
		token.L_PAREN, "(",
		token.STRING_LIT, "oops",
		token.R_PAREN, ")",
		token.EOL, "",
	))
	assign := ss[0].(*Assign_stmt)
	lit, ok := assign.Values[0].(*Lit_val)
	if !ok {
		t.Fatalf("want *Lit_val, got %T", assign.Values[0])
	}
	if lit.Value != "oops" || lit.Type != token.ERR {
		t.Errorf("err lit: got {%v, %v}, want {oops, ERR}", lit.Value, lit.Type)
	}
}

func TestParseIdentifierExpr(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "z",
		token.ASSIGNMENT, "=",
		token.IDENTIFIER, "x",
		token.EOL, "",
	))
	assign := ss[0].(*Assign_stmt)
	ident, ok := assign.Values[0].(*Identifier_expr)
	if !ok {
		t.Fatalf("want *Identifier_expr, got %T", assign.Values[0])
	}
	if ident.Name != "x" {
		t.Errorf("ident name: got %q, want x", ident.Name)
	}
}

func TestParsePostfixStmt(t *testing.T) {
	cases := []struct {
		opTok token.TokenType
		opLex string
	}{
		{token.PLUSPLUS, "++"},
		{token.MINUSMINUS, "--"},
	}
	for _, c := range cases {
		ss := stmts(t, funcWrap(
			token.IDENTIFIER, "x",
			c.opTok, c.opLex,
			token.EOL, "",
		))
		expr, ok := ss[0].(*Expr_stmt)
		if !ok {
			t.Fatalf("%v: want *Expr_stmt, got %T", c.opLex, ss[0])
		}
		post, ok := expr.Expression.(*Postfix_expr)
		if !ok {
			t.Fatalf("%v: want *Postfix_expr, got %T", c.opLex, expr.Expression)
		}
		if post.Op != c.opTok {
			t.Errorf("%v: op: got %v, want %v", c.opLex, post.Op, c.opTok)
		}
		ident, ok := post.Target.(*Identifier_expr)
		if !ok || ident.Name != "x" {
			t.Errorf("%v: target: got %T %v", c.opLex, post.Target, post.Target)
		}
	}
}

func TestParseCallExpression(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "print",
		token.L_PAREN, "(",
		token.R_PAREN, ")",
		token.EOL, "",
	))
	expr, ok := ss[0].(*Expr_stmt)
	if !ok {
		t.Fatalf("want *Expr_stmt, got %T", ss[0])
	}
	call, ok := expr.Expression.(*Call_expr)
	if !ok {
		t.Fatalf("want *Call_expr, got %T", expr.Expression)
	}
	if call.Name != "print" {
		t.Errorf("call name: got %q, want print", call.Name)
	}
	if len(call.Args) != 0 {
		t.Errorf("want 0 args, got %d", len(call.Args))
	}
}

func TestParseCallWithArgs(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IDENTIFIER, "add",
		token.L_PAREN, "(",
		token.INT_LIT, "1",
		token.COMMA, ",",
		token.INT_LIT, "2",
		token.R_PAREN, ")",
		token.EOL, "",
	))
	expr, ok := ss[0].(*Expr_stmt)
	if !ok {
		t.Fatalf("want *Expr_stmt, got %T", ss[0])
	}
	call, ok := expr.Expression.(*Call_expr)
	if !ok {
		t.Fatalf("want *Call_expr, got %T", expr.Expression)
	}
	if call.Name != "add" {
		t.Errorf("call name: got %q, want add", call.Name)
	}
	if len(call.Args) != 2 {
		t.Fatalf("want 2 args, got %d", len(call.Args))
	}
}

// --- Control flow ---

func TestParseReturnStmt(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.RETURN, "return",
		token.INT_LIT, "5",
		token.EOL, "",
	))
	ret, ok := ss[0].(*Return_stmt)
	if !ok {
		t.Fatalf("want *Return_stmt, got %T", ss[0])
	}
	if len(ret.Returns) != 1 {
		t.Fatalf("want 1 return value, got %d", len(ret.Returns))
	}
	lit, ok := ret.Returns[0].(*Lit_val)
	if !ok || lit.Value != "5" {
		t.Errorf("return value: got %T %v, want Lit_val 5", ret.Returns[0], ret.Returns[0])
	}
}

func TestParseReturnNoValue(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.RETURN, "return",
		token.EOL, "",
	))
	ret, ok := ss[0].(*Return_stmt)
	if !ok {
		t.Fatalf("want *Return_stmt, got %T", ss[0])
	}
	if len(ret.Returns) != 0 {
		t.Errorf("want 0 return values, got %d", len(ret.Returns))
	}
}

func TestParseJumpStmt(t *testing.T) {
	cases := []struct {
		tok token.TokenType
		lex string
	}{
		{token.BREAK, "break"},
		{token.CONTINUE, "continue"},
	}
	for _, c := range cases {
		ss := stmts(t, funcWrap(
			c.tok, c.lex,
			token.EOL, "",
		))
		jmp, ok := ss[0].(*Jump_stmt)
		if !ok {
			t.Fatalf("%v: want *Jump_stmt, got %T", c.lex, ss[0])
		}
		if jmp.Type != c.tok {
			t.Errorf("%v: type: got %v", c.lex, jmp.Type)
		}
	}
}

func TestParseLoopStmt(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.FOR, "for",
		token.IDENTIFIER, "x",
		token.LT, "<",
		token.INT_LIT, "10",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
	))
	loop, ok := ss[0].(*Loop_stmt)
	if !ok {
		t.Fatalf("want *Loop_stmt, got %T", ss[0])
	}
	cond, ok := loop.Expression.(*Binary_expr)
	if !ok || cond.Op != token.LT {
		t.Errorf("loop cond: want Binary_expr LT, got %T", loop.Expression)
	}
	if loop.Loop_block == nil {
		t.Error("loop block is nil")
	}
}

func TestParseControlIf(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IF, "if",
		token.IDENTIFIER, "x",
		token.GT, ">",
		token.INT_LIT, "0",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.EOL, "",
	))
	ctrl, ok := ss[0].(*Control_stmt)
	if !ok {
		t.Fatalf("want *Control_stmt, got %T", ss[0])
	}
	if ctrl.IfBlock == nil {
		t.Error("IfBlock is nil")
	}
	if ctrl.ElseBlock != nil || ctrl.ElseStmt != nil {
		t.Error("want no else branch")
	}
}

func TestParseControlIfElse(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IF, "if",
		token.IDENTIFIER, "x",
		token.GT, ">",
		token.INT_LIT, "0",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.ELSE, "else",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
	))
	ctrl, ok := ss[0].(*Control_stmt)
	if !ok {
		t.Fatalf("want *Control_stmt, got %T", ss[0])
	}
	if ctrl.ElseBlock == nil {
		t.Error("want ElseBlock, got nil")
	}
	if ctrl.ElseStmt != nil {
		t.Error("want no ElseStmt")
	}
}

func TestParseControlElseIf(t *testing.T) {
	ss := stmts(t, funcWrap(
		token.IF, "if",
		token.IDENTIFIER, "x",
		token.GT, ">",
		token.INT_LIT, "0",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.ELSE, "else",
		token.IF, "if",
		token.IDENTIFIER, "x",
		token.LT, "<",
		token.INT_LIT, "0",
		token.L_BRACE, "{",
		token.R_BRACE, "}",
		token.EOL, "",
	))
	ctrl, ok := ss[0].(*Control_stmt)
	if !ok {
		t.Fatalf("want *Control_stmt, got %T", ss[0])
	}
	if ctrl.ElseStmt == nil {
		t.Error("want ElseStmt, got nil")
	}
	if ctrl.ElseBlock != nil {
		t.Error("want no ElseBlock")
	}
}

// --- Error cases ---

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
