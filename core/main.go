package main

import (
	"fmt"
	"os"
	"radlang/interpreter"
	"radlang/lexer"
	"radlang/parser"
	"radlang/semantic"
	"strings"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Fprintln(os.Stderr, "usage: radlang <file>")
	// 	os.Exit(1)
	// }
	// fileName := os.Args[1]
	// file, err := os.ReadFile(fileName)
	file, err := os.ReadFile("../tests/arithmetic.rad")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", string(file))
	fmt.Printf("Lexer Output:\n")
	fmt.Printf("%-4s %-16s %s\n", "Line", "Type", "Lexeme")
	fmt.Printf("%-4s %-16s %s\n", "----", "----", "------")
	tokens := lexer.Lex(string(file))
	for _, val := range tokens {
		fmt.Println(val)
	}
	fmt.Printf("\nAST:\n")
	ast, err := parser.Parse(tokens)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}
	printAST(ast)

	err = semantic.Analyze(ast)
	if err != nil {
		fmt.Printf("\nSemantic analysis returned: %v\n", err)
		return
	}

	interpreter.Interpret(ast)
}

func printAST(program *parser.Program) {
	fmt.Println("Program")
	for _, fn := range program.Functions {
		fmt.Printf("  Func_Decl: %s\n", fn.Name)
		printBlock(fn.Body, 2)
	}
}

func printBlock(block *parser.Block, depth int) {
	pad := strings.Repeat("  ", depth)
	fmt.Printf("%sBlock\n", pad)
	for _, stmt := range block.Statement_Group {
		printStatement(stmt, depth+1)
	}
}

func printStatement(stmt parser.Statement, depth int) {
	pad := strings.Repeat("  ", depth)
	switch s := stmt.(type) {
	case *parser.Decl_stmt:
		fmt.Printf("%sDecl_stmt: %s %v\n", pad, s.Name, s.Type)
	case *parser.Assign_stmt:
		fmt.Printf("%sAssign_stmt: %s =\n", pad, s.Target)
		printExpr(s.Value, depth+1)
	case *parser.Update_stmt:
		fmt.Printf("%sUpdate_stmt: %s %v\n", pad, s.Target, s.Op)
	case *parser.Expr_stmt:
		fmt.Printf("%sExpr_stmt:\n", pad)
		printExpr(s.Expr, depth+1)
	}
}

func printExpr(expr parser.Expression, depth int) {
	pad := strings.Repeat("  ", depth)
	switch e := expr.(type) {
	case *parser.Binary_expr:
		fmt.Printf("%sBinary_expr: %v\n", pad, e.Op)
		printExpr(e.Left, depth+1)
		printExpr(e.Right, depth+1)
	case *parser.Call_expr:
		fmt.Printf("%sCall_expr: %s()\n", pad, e.Name)
	case *parser.Identifier_expr:
		fmt.Printf("%sIdentifier_expr: %s\n", pad, e.Name)
	case *parser.Number_lit:
		fmt.Printf("%sNumber_lit: %s %s\n", pad, e.Type, e.Value)
	}
}
