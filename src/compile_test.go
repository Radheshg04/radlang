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
	_, err = Parse(tokens)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
}
