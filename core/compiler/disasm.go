package compiler

import "fmt"

type disassembler struct {
	bc  *Bytecode
	pos int
	ip  int
}

func Disasm(bc *Bytecode) {
	d := &disassembler{
		bc: bc,
		ip: 1,
	}
	fmt.Printf("Execution starts at function ID: %d\n", bc.EntryPointID)
	for d.pos < len(bc.Code) {
		fmt.Printf("#%d %s\n", d.pos, d.opcodeToString(Opcode(d.currByte())))
		d.pos++
		d.ip++
	}
}

func (d *disassembler) currByte() byte {
	return d.bc.Code[d.pos]
}

func (d *disassembler) nextByte() byte {
	if d.pos+1 >= len(d.bc.Code) {
		panic("unexpected EOF while decoding instruction")
	}
	d.pos++
	return d.currByte()
}

func (d *disassembler) nextUint16() uint16 {
	hi := uint16(d.nextByte())
	lo := uint16(d.nextByte())
	return (hi << 8) | lo
}

func (d *disassembler) opcodeToString(opcode Opcode) string {
	switch opcode {

	// program
	case HALT:
		return "HALT"
	case CALL:
		return fmt.Sprintf("CALL %d", d.nextByte())
	case RETURN:
		return "RET"

	// CONSTANTS
	case CONST:
		constNum := d.nextByte()
		return fmt.Sprintf("CONST #%d %v", constNum, d.bc.Constants[constNum])

	// Variables
	case STORE:
		return fmt.Sprintf("STORE %d", d.nextByte())
	case LOAD:
		return fmt.Sprintf("LOAD %d", d.nextByte())

	// Statements
	case JMP:
		return fmt.Sprintf("JMP %d", d.nextUint16())
	case JMP_IF_FALSE:
		return fmt.Sprintf("JMP_IF_FALSE %d", d.nextUint16())

	// arithmetic
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"

	// special
	case PRINT:
		return "PRINT"

	// stack
	case POP:
		return "POP"
	}
	return ""
}
