package vm

import (
	"radlang/compiler"
)

// byte readers
func (vm *vm) readByte() byte {
	if vm.ip >= len(vm.bc.Code) {
		panic("unexpected end of bytecode")
	}
	out := vm.bc.Code[vm.ip]
	vm.ip++
	return out
}
func (vm *vm) readUint16() uint16 {
	left := uint16(vm.readByte())
	right := uint16(vm.readByte())
	return (left << 8) | right
}

// stack functions
func (vm *vm) push(v compiler.Value) {
	if vm.sp >= MaxStackSize {
		panic("stack overflow")
	}
	vm.stack[vm.sp] = v
	vm.sp++
}
func (vm *vm) pop() compiler.Value {
	if vm.sp < 0 {
		panic("stack underflow")
	}
	vm.sp--
	return vm.peek()
}
func (vm *vm) peek() compiler.Value {
	return vm.stack[vm.sp]
}
func (vm *vm) stackSize() int {
	return vm.sp
}

func (vm *vm) pushCallStack(v callFrame) {
	vm.callStack = append(vm.callStack, v)
}

func (vm *vm) popCallStack() callFrame {
	frame := vm.callStack[len(vm.callStack)-1]
	vm.callStack = vm.callStack[:len(vm.callStack)-1]
	return frame
}
