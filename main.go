package main

import (
	"fmt"
	"strconv"

	"github.com/armadi1809/gobytecode/expression"
	"github.com/armadi1809/gobytecode/instruction"
	"github.com/armadi1809/gobytecode/opcode"
)

func compile(exp expression.Expression) []*instruction.Instruction {
	res := []*instruction.Instruction{}

	for _, tok := range exp.Tokens {
		if _, err := strconv.Atoi(tok); err == nil {
			inst := instruction.New(opcode.LOAD_CONST, tok)
			res = append(res, inst)
		}
	}
	return res
}

func eval(code []*instruction.Instruction) int {
	pc := 0
	stack := []int{}
	for pc < len(code) {
		ins := code[pc]
		op := ins.Op
		pc += 1
		if op == opcode.LOAD_CONST {
			i, err := strconv.Atoi(ins.Arg)
			if err != nil {
				panic("Invalid argument load constant instruction")
			}
			stack = append(stack, i)
		}
	}

	if len(stack) > 0 {
		return stack[len(stack)-1]
	}

	return 0
}

func main() {
	fmt.Println(eval(compile(expression.Expression{Tokens: []string{"5"}})))
	fmt.Println(eval(compile(expression.Expression{Tokens: []string{"7"}})))

}
