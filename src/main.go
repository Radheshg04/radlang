package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("../tests/arithmetic.rad")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", string(file))
	fmt.Printf("Lexer Output:\n")
	fmt.Printf("%-4s %-16s %s\n", "Line", "Type", "Lexeme")
	fmt.Printf("%-4s %-16s %s\n", "----", "----", "------")
	tokens := Lex(string(file))
	for _, val := range tokens {
		fmt.Println(val)
	}
	fmt.Printf("\nAST:\n")
	ast, err := Parse(tokens)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}
	printAST(ast)

	fmt.Printf("\nSemantic analysis returned: %v", Analyze(ast))
}
