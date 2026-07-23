package vm

import (
	"fmt"
	"radlang/compiler"
)

const MaxStackSize = 1024

type callFrame struct {
	id       int
	locals   []compiler.Value
	returnIP int
}

type vm struct {
	bc        compiler.Bytecode
	ip        int
	callStack []callFrame
	stack     [MaxStackSize]compiler.Value
	sp        int // stack pointer
}

func Execute(bc compiler.Bytecode) {
	vm := &vm{
		bc: bc,
	}

	vm.Run()
}

// vm functions
func (vm *vm) Run() {
	vm.newCallFrame(vm.bc.EntryPointID)

	for {
		op := vm.readByte()
		if !vm.executeOp(compiler.Opcode(op)) {
			break
		}
	}
}

func (vm *vm) newCallFrame(fnID int) {
	fn := vm.bc.FunctionInfo[fnID]
	vm.pushCallStack(callFrame{
		id:       fnID,
		returnIP: vm.ip,
		locals:   make([]compiler.Value, fn.Slots),
	})

	// inject args
	argc := fn.Argc
	for i := argc - 1; i >= 0; i-- {
		vm.currentCallFrame().locals[i] = vm.pop()
		argc--
	}

	// move to func exec
	vm.ip = fn.EntryIP
}

func (vm *vm) currentCallFrame() *callFrame {
	return &vm.callStack[len(vm.callStack)-1]
}

func (vm *vm) executeOp(opcode compiler.Opcode) bool {
	switch opcode {

	// program
	case compiler.HALT:
		return false
	case compiler.CALL:
		id := int(vm.readByte())
		vm.newCallFrame(id)
	case compiler.RETURN:
		frane := vm.popCallStack()
		if len(vm.callStack) == 0 {
			return false
		}
		vm.ip = frane.returnIP

	// CONSTANTS
	case compiler.CONST:
		vm.push(vm.bc.Constants[vm.readByte()])

	// Variables
	case compiler.STORE:
		slot := vm.readByte()
		vm.currentCallFrame().locals[slot] = vm.pop()
	case compiler.LOAD:
		slot := vm.readByte()
		vm.push(vm.currentCallFrame().locals[slot])

	// Statements
	case compiler.JMP:
		vm.ip = int(vm.readUint16())
	case compiler.JMP_IF_FALSE:
		addr := int(vm.readUint16())
		if val, ok := vm.pop().(compiler.BoolValue); ok && !val.Val {
			vm.ip = addr
		}

	// arithmetic
	case compiler.ADD:
		vm.BinaryOp(compiler.Value.Add)
	case compiler.SUB:
		vm.BinaryOp(compiler.Value.Sub)
	case compiler.MUL:
		vm.BinaryOp(compiler.Value.Mul)
	case compiler.DIV:
		vm.BinaryOp(compiler.Value.Div)

	// special
	case compiler.PRINT:
		fmt.Println(vm.pop())

	// stack
	case compiler.POP:
		vm.pop()
	}
	return true
}

func (vm *vm) BinaryOp(performOp func(compiler.Value, compiler.Value) compiler.Value) {
	rhs := vm.pop()
	lhs := vm.pop()
	vm.push(performOp(lhs, rhs))
}
