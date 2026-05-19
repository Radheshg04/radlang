package main

import (
	"os"
	"testing"
)

func TestCompileHello(t *testing.T) {
	src, err := os.ReadFile("../tests/hello.rad")
	if err != nil {
		t.Fatalf("read hello.rad: %v", err)
	}

	tokens := Lex(string(src))
	ast, err := Parse(tokens)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	err = Analyze(ast)
	if err != nil {
		t.Fatalf("semantic analysis returned: %v", err)
	}
}
