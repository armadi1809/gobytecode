package instruction

import "github.com/armadi1809/gobytecode/opcode"

type Instruction struct {
	Op    opcode.Opcode
	Name  string
	Value int
}

func New(op opcode.Opcode, name string, val int) *Instruction {
	return &Instruction{op, name, val}
}
