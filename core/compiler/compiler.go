package compiler

import (
	"radlang/parser"
	"radlang/semantic"
	"radlang/token"
)

type Compiler struct {
	bc           Bytecode
	currentScope *semantic.Scope
	loopCtx      *LoopContext // stores loop ctx (used for jump stmt)
	funcBound    *Label       // marks the current function end (used for ret stmt)
}

func Compile(ast *parser.Program, globalScope *semantic.Scope) (*Bytecode, error) {
	c := Compiler{
		bc:           Bytecode{},
		currentScope: globalScope,
		loopCtx:      &LoopContext{},
		funcBound:    &Label{},
	}
	c.compileProgram(ast)
	return &c.bc, nil
}

func (c *Compiler) compileProgram(program *parser.Program) {
	// register pass
	for _, function := range program.Functions {
		funcSym, _ := function.Symbol.(*semantic.FuncSymbol)
		funcSym.ID = c.bc.enrichFuncInfo(funcSym)

		if function.Signature.Name == "main" {
			c.bc.EntryPointID = funcSym.ID
		}
	}

	// compilation pass
	for _, function := range program.Functions {
		c.compileFunction(function)
	}
}

func (c *Compiler) compileFunction(function *parser.Func_Decl) {
	end := c.newLabel()
	c.funcBound = end

	c.compileBlock(function.Body)
	c.markLabel(end)

	if function.Signature.Name == "main" {
		c.bc.emit(HALT)
	}
}

func (c *Compiler) compileBlock(block *parser.Block) {
	parentScope := c.currentScope
	c.currentScope = block.Scope.(*semantic.Scope)
	for _, stmt := range block.Statement_Group {
		c.compileStatement(stmt)
	}
	c.currentScope = parentScope
}

func (c *Compiler) compileStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.Decl_stmt:
		// no bytecode emitted

	case *parser.Assign_stmt:
		for _, value := range s.Values {
			c.compileExpression(value)
		}

		for i := len(s.Targets) - 1; i >= 0; i-- {
			slot, err := c.lookupVar(s.Targets[i])
			if err != nil {
				panic("could not find var in scope")
			}
			c.bc.emit(STORE, byte(slot))
		}

	case *parser.Expr_stmt:
		c.compileExpression(s.Expression)
		c.bc.emit(POP)

	case *parser.Jump_stmt:
		switch s.Type {
		case token.CONTINUE:
			c.emitJump(JMP, c.loopCtx.continueLabel)
		case token.BREAK:
			c.emitJump(JMP, c.loopCtx.breakLabel)
		}

	case *parser.Return_stmt:
		for _, value := range s.Returns {
			c.compileExpression(value)
		}
		c.bc.emit(RETURN)
		c.emitJump(JMP, c.funcBound)

	case *parser.Control_stmt:
		falseLabel := c.newLabel()
		end := c.newLabel()

		c.compileExpression(s.Expression)
		c.emitJump(JMP_IF_FALSE, falseLabel)
		c.compileBlock(s.IfBlock)
		c.emitJump(JMP, end)

		c.markLabel(falseLabel)
		if s.ElseStmt != nil {
			c.compileStatement(s.ElseStmt)
		} else if s.ElseBlock != nil {
			c.compileBlock(s.ElseBlock)
		}

		c.markLabel(end)

	case *parser.Loop_stmt:
		parent := c.loopCtx
		defer func() {
			c.loopCtx = parent
		}()

		start := c.newLabel()
		end := c.newLabel()

		c.loopCtx = &LoopContext{
			breakLabel:    end,
			continueLabel: start,
		}

		c.markLabel(start)
		c.compileExpression(s.Expression)
		c.emitJump(JMP_IF_FALSE, end)
		c.compileBlock(s.Loop_block)
		c.emitJump(JMP, start)
		c.markLabel(end)
	}
}

func (c *Compiler) compileExpression(expr parser.Expression) {

	switch e := expr.(type) {

	case *parser.Lit_val:
		var val Value

		switch e.Type {
		case token.INT:
			val = IntValue{Val: e.Value.(int64)}
		case token.FLOAT:
			val = FloatValue{Val: e.Value.(float64)}
		case token.BOOL:
			val = BoolValue{Val: e.Value.(bool)}
		case token.STRING:
			val = StringValue{Val: e.Value.(string)}
		default:
			panic("unknown type found in literal value during compilation")
		}
		idx := c.bc.addConst(val)
		c.bc.emit(CONST, byte(idx))

	case *parser.Call_expr:
		for _, arg := range e.Args {
			c.compileExpression(arg)
		}
		if e.Name == "print" {
			c.bc.emit(PRINT)
			return
		}
		id, _ := c.lookupFunc(e.Name)
		c.bc.emit(CALL, byte(id))

	case *parser.Binary_expr:
		c.compileExpression(e.Left)
		c.compileExpression(e.Right)

		switch e.Op {
		case token.PLUS:
			c.bc.emit(ADD)
		case token.MINUS:
			c.bc.emit(SUB)
		case token.ASTERISK:
			c.bc.emit(MUL)
		case token.SLASH:
			c.bc.emit(DIV)
		}

	case *parser.Identifier_expr:
		sym, _ := semantic.Resolve(c.currentScope, e.Name)
		varSym, _ := sym.(*semantic.VarSymbol)
		c.bc.emit(LOAD, byte(varSym.Slot))

	case *parser.Postfix_expr:
		sym, _ := semantic.Resolve(c.currentScope, e.Target.Name)
		varSym, _ := sym.(*semantic.VarSymbol)
		c.bc.emit(LOAD, byte(varSym.Slot))

		// TODO: this is temporary, update before PR
		idx := c.bc.addConst(IntValue{Val: 1})

		c.bc.emit(CONST, byte(idx))
		switch e.Op {
		case token.PLUSPLUS:
			c.bc.emit(ADD)
		case token.MINUSMINUS:
			c.bc.emit(SUB)
		}
		c.bc.emit(STORE, byte(varSym.Slot))
	}
}
