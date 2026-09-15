package instruction

import "github.com/armadi1809/gobytecode/opcode"

type Instruction struct {
	Op  opcode.Opcode
	Arg string
}

func New(op opcode.Opcode, arg string) *Instruction {
	return &Instruction{op, arg}
}
