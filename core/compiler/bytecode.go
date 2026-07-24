package compiler

type Opcode byte

const (
	// program
	HALT Opcode = iota
	CALL
	RETURN

	// CONSTANTS
	CONST

	// Variables
	STORE
	LOAD

	// Statements
	JMP
	JMP_IF_FALSE

	// arithmetic
	ADD
	SUB
	MUL
	DIV

	// special
	PRINT

	// stack
	POP
)

type Bytecode struct {
	Code         []byte
	Constants    []Value
	EntryPointID int
	FunctionInfo []Function
}

func (bc *Bytecode) emit(opcode Opcode, Args ...byte) {
	bc.Code = append(bc.Code, byte(opcode))
	bc.Code = append(bc.Code, Args...)
}

// takes in jmp or jmp_if_false
func (c *Compiler) emitJump(Opcode Opcode, label *Label) {
	c.bc.emit(Opcode)
	label.patchWhere = len(c.bc.Code)
	c.bc.emit(0, 0)

	if label.address != -1 {
		c.backpatch(label)
	}
}

// Adds constant to constant pool and returns idx of newly added const
func (bc *Bytecode) addConst(val Value) int {
	idx := len(bc.Constants)
	bc.Constants = append(bc.Constants, val)
	return idx
}
