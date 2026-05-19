package main

import (
	"fmt"
	"strings"
)

func printAST(program *Program) {
	fmt.Println("Program")
	for _, fn := range program.Functions {
		fmt.Printf("  Func_Decl: %s\n", fn.Name)
		printBlock(fn.Body, 2)
	}
}

func printBlock(block *Block, depth int) {
	pad := strings.Repeat("  ", depth)
	fmt.Printf("%sBlock\n", pad)
	for _, stmt := range block.Statement_Group {
		printStatement(stmt, depth+1)
	}
}

func printStatement(stmt Statement, depth int) {
	pad := strings.Repeat("  ", depth)
	switch s := stmt.(type) {
	case *Decl_stmt:
		fmt.Printf("%sDecl_stmt: %s %v\n", pad, s.Name, s.Type)
	case *Assign_stmt:
		fmt.Printf("%sAssign_stmt: %s =\n", pad, s.Target)
		printExpr(s.Value, depth+1)
	case *Update_stmt:
		fmt.Printf("%sUpdate_stmt: %s %v\n", pad, s.Target, s.Op)
	case *Expr_stmt:
		fmt.Printf("%sExpr_stmt:\n", pad)
		printExpr(s.Expr, depth+1)
	}
}

func printExpr(expr Expression, depth int) {
	pad := strings.Repeat("  ", depth)
	switch e := expr.(type) {
	case *Binary_expr:
		fmt.Printf("%sBinary_expr: %v\n", pad, e.Op)
		printExpr(e.Left, depth+1)
		printExpr(e.Right, depth+1)
	case *Call_expr:
		fmt.Printf("%sCall_expr: %s()\n", pad, e.Name)
	case *Identifier_expr:
		fmt.Printf("%sIdentifier_expr: %s\n", pad, e.Name)
	case *Number_lit:
		fmt.Printf("%sNumber_lit: %s %s\n", pad, e.Type, e.Value)
	}
}
