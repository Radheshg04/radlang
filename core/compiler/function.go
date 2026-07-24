package compiler

import "radlang/semantic"

type Function struct {
	EntryIP int
	Argc    int
	Retc    int
	Slots   int
}

// Generates compiler relevant function object.
// Returns index of new fn object in bc.FunctionInfo
func (bc *Bytecode) enrichFuncInfo(function *semantic.FuncSymbol) int {
	funcInfo := Function{
		Argc:    len(function.Params),
		Retc:    len(function.Returns),
		EntryIP: len(bc.Code),
		Slots:   function.Slots,
	}
	funcidx := len(bc.FunctionInfo)
	bc.FunctionInfo = append(bc.FunctionInfo, funcInfo)
	return funcidx
}
