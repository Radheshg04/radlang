package parser

import "radlang/token"

type Program struct {
	Functions  []*Func_Decl
	Structs    []*Struct_Decl
	Interfaces []*Interface_Decl
}

// Add returntype and params
type Func_Decl struct {
	Signature *Func_Signature
	Body      *Block
}

type Func_Signature struct {
	Name    *string
	Params  []*Param
	Returns []token.TokenType
}

type Param struct {
	Name *string
	Type token.TokenType
}

type Struct_Decl struct {
	Name *string
	Body *STRUCT_DEF_BLOCK
}

type Interface_Decl struct {
	Name *string
	Body *INTERFACE_DEF_BLOCK
}

type STRUCT_DEF_BLOCK struct {
	Statement_Group []*Decl_stmt
}
type INTERFACE_DEF_BLOCK struct {
	Statement_Group []*Func_Signature
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

type Assign_stmt struct {
	Targets []string
	Values  []Expression
	Op      token.TokenType
}

func (*Assign_stmt) stmtNode() {}

type Expr_stmt struct {
	Expression Expression
}

func (*Expr_stmt) stmtNode() {}

type Decl_stmt struct {
	Name string
	Type token.TokenType
}

func (*Decl_stmt) stmtNode() {}

// Handles x++ and x--
type Update_stmt struct {
	Target string
	Op     token.TokenType
}

func (*Update_stmt) stmtNode() {}

type Control_stmt struct {
	Expression Expression
	IfBlock    *Block
	ElseStmt   *Control_stmt
	ElseBlock  *Block
}

func (*Control_stmt) stmtNode() {}

type Jump_stmt struct {
	Type token.TokenType
}

func (*Jump_stmt) stmtNode() {}

type Return_stmt struct {
	Returns []Expression
}

func (*Return_stmt) stmtNode() {}

type Loop_stmt struct {
	Expression Expression
	Loop_block *Block
}

func (*Loop_stmt) stmtNode() {}

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
	Args []Expression
}

func (*Call_expr) exprNode() {}

type Lit_val struct {
	Value interface{} // store parsed values
	Type  token.TokenType
}

func (*Lit_val) exprNode() {}
