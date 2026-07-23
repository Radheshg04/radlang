package compiler

import (
	"fmt"
	"radlang/semantic"
)

func (c *Compiler) lookupFunc(name string) (int, error) {
	sym, err := semantic.Resolve(c.currentScope, name)
	if err != nil {
		return 0, err
	}
	funcSym, exists := sym.(*semantic.FuncSymbol)
	if !exists {
		return 0, fmt.Errorf("sym is not funcsym")
	}
	return funcSym.ID, nil
}

func (c *Compiler) lookupVar(name string) (int, error) {
	sym, err := semantic.Resolve(c.currentScope, name)
	if err != nil {
		return 0, err
	}
	varSym, exists := sym.(*semantic.VarSymbol)
	if !exists {
		return 0, fmt.Errorf("sym is not varsym")
	}
	return varSym.Slot, nil
}
