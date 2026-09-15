package instruction

import "github.com/armadi1809/gobytecode/opcode"

type Value = any

type Instruction struct {
	Op      opcode.Opcode
	Operand int
	Name    string
	Value   Value
}

func New(op opcode.Opcode, operand int, name string, val int) *Instruction {
	return &Instruction{op, operand, name, val}
}
