package tests

import (
	"os"
	"testing"

	"radlang/interpreter"
	"radlang/lexer"
	"radlang/parser"
	"radlang/semantic"
)

func TestCompileHello(t *testing.T) {
	src, err := os.ReadFile("arithmetic.rad")
	if err != nil {
		t.Fatalf("read hello.rad: %v", err)
	}

	tokens := lexer.Lex(string(src))
	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	err = semantic.Analyze(ast)
	if err != nil {
		t.Fatalf("semantic analysis returned: %v", err)
	}
	interpreter.Interpret(ast)
}
