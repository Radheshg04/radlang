package main

type Parser struct {
	tokens []Token
	pos    int
}

type Node interface{}

type Program struct {
	Functions []*Func_Decl
}

// Add returntype and params
type Func_Decl struct {
	Name string
	Body *Block
}

type Block struct {
	Statement_Group []*Statement
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

type Assign_stmt struct {
	Target string
	Value  Expression
}

func (*Assign_stmt) stmtNode() {}

type Expr_stmt struct {
	Expr Expression
}

func (*Expr_stmt) stmtNode() {}
