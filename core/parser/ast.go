package parser

import "radlang/token"

type Program struct {
	Functions []*Func_Decl
}

// Add returntype and params
type Func_Decl struct {
	Name string
	Body *Block
}

type Block struct {
	Statement_Group []Statement
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

type Decl_stmt struct {
	Name string
	Type token.TokenType
}

func (*Decl_stmt) stmtNode() {}

type Assign_stmt struct {
	Target string
	Value  Expression
}

func (*Assign_stmt) stmtNode() {}

// Handles x++ and x--
type Update_stmt struct {
	Target string
	Op     token.TokenType
}

func (*Update_stmt) stmtNode() {}

type Expr_stmt struct {
	Expr Expression
}

func (*Expr_stmt) stmtNode() {}

type Identifier_expr struct {
	Name string
}

func (*Identifier_expr) exprNode() {}

type Binary_expr struct {
	Left  Expression
	Op    token.TokenType
	Right Expression
}

func (*Binary_expr) exprNode() {}

// Function Call expr
type Call_expr struct {
	Name string
}

func (*Call_expr) exprNode() {}

type Number_lit struct {
	Value string
	Type  token.TokenType
}

func (*Number_lit) exprNode() {}
