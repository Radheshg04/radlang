package lexer

import (
	"radlang/token"
	"testing"
)

type tokenExpect struct {
	tok    token.TokenType
	lexeme string
}

func assertTokens(t *testing.T, input string, want []tokenExpect) {
	t.Helper()
	got := Lex(input)
	if len(got) != len(want) {
		t.Fatalf("input %q: got %d tokens, want %d\ngot: %v", input, len(got), len(want), got)
	}
	for i, w := range want {
		if got[i].Token != w.tok || got[i].Lexeme != w.lexeme {
			t.Errorf("token[%d]: got {%v %q}, want {%v %q}", i, got[i].Token, got[i].Lexeme, w.tok, w.lexeme)
		}
	}
}

func TestLexKeywords(t *testing.T) {
	cases := []struct {
		input string
		tok   token.TokenType
	}{
		{"func", token.FUNC},
		{"var", token.VAR},
		{"int", token.INT},
		{"float", token.FLOAT},
		{"bool", token.BOOL},
		{"string", token.STRING},
	}
	for _, c := range cases {
		tokens := Lex(c.input)
		if len(tokens) != 2 {
			t.Fatalf("Lex(%q): want 2 tokens (keyword+EOF), got %d", c.input, len(tokens))
		}
		if tokens[0].Token != c.tok {
			t.Errorf("Lex(%q): got token type %v, want %v", c.input, tokens[0].Token, c.tok)
		}
		if tokens[0].Lexeme != c.input {
			t.Errorf("Lex(%q): got lexeme %q, want %q", c.input, tokens[0].Lexeme, c.input)
		}
	}
}

func TestLexIdentifier(t *testing.T) {
	assertTokens(t, "myVar", []tokenExpect{
		{token.IDENTIFIER, "myVar"},
		{token.EOF, ""},
	})
	assertTokens(t, "_x1", []tokenExpect{
		{token.IDENTIFIER, "_x1"},
		{token.EOF, ""},
	})
}

func TestLexIntLiteral(t *testing.T) {
	assertTokens(t, "42", []tokenExpect{
		{token.INT, "42"},
		{token.EOF, ""},
	})
}

func TestLexFloatLiteral(t *testing.T) {
	assertTokens(t, "3.14", []tokenExpect{
		{token.FLOAT, "3.14"},
		{token.EOF, ""},
	})
}

func TestLexOperators(t *testing.T) {
	assertTokens(t, "+ - * /", []tokenExpect{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.EOF, ""},
	})
}

func TestLexIncDecrement(t *testing.T) {
	assertTokens(t, "++ --", []tokenExpect{
		{token.PLUSPLUS, "++"},
		{token.MINUSMINUS, "--"},
		{token.EOF, ""},
	})
}

func TestLexAssignment(t *testing.T) {
	assertTokens(t, "=", []tokenExpect{
		{token.ASSIGNMENT, "="},
		{token.EOF, ""},
	})
}

func TestLexDelimiters(t *testing.T) {
	assertTokens(t, "( ) { }", []tokenExpect{
		{token.L_PAREN, "("},
		{token.R_PAREN, ")"},
		{token.L_BRACE, "{"},
		{token.R_BRACE, "}"},
		{token.EOF, ""},
	})
}

func TestLexStringLiteral(t *testing.T) {
	assertTokens(t, `"hello"`, []tokenExpect{
		{token.STRING_LITERAL, `"hello"`},
		{token.EOF, ""},
	})
}

func TestLexNewline(t *testing.T) {
	tokens := Lex("x\ny")
	if tokens[1].Token != token.EOL {
		t.Errorf("expected EOL between identifiers, got %v", tokens[1].Token)
	}
}

func TestLexLineNumbers(t *testing.T) {
	tokens := Lex("x\ny")
	if tokens[0].Line != 1 {
		t.Errorf("first token line: got %d, want 1", tokens[0].Line)
	}
	if tokens[2].Line != 2 {
		t.Errorf("second ident line: got %d, want 2", tokens[2].Line)
	}
}

func TestLexIllegalToken(t *testing.T) {
	tokens := Lex("@")
	if tokens[0].Token != token.ILLEGAL {
		t.Errorf("expected ILLEGAL for '@', got %v", tokens[0].Token)
	}
}

func TestLexEOF(t *testing.T) {
	tokens := Lex("x")
	last := tokens[len(tokens)-1]
	if last.Token != token.EOF {
		t.Errorf("last token should be EOF, got %v", last.Token)
	}
}

func TestLexSimpleFunction(t *testing.T) {
	src := "func main() {\n}\n"
	tokens := Lex(src)
	want := []token.TokenType{
		token.FUNC, token.IDENTIFIER, token.L_PAREN, token.R_PAREN,
		token.L_BRACE, token.EOL, token.R_BRACE, token.EOL, token.EOF,
	}
	if len(tokens) != len(want) {
		t.Fatalf("got %d tokens, want %d: %v", len(tokens), len(want), tokens)
	}
	for i, w := range want {
		if tokens[i].Token != w {
			t.Errorf("token[%d]: got %v, want %v", i, tokens[i].Token, w)
		}
	}
}
