package main

import (
	"fmt"
	"log"

	"github.com/armadi1809/gobytecode/env"
	"github.com/armadi1809/gobytecode/expression"
	"github.com/armadi1809/gobytecode/instruction"
	"github.com/armadi1809/gobytecode/opcode"
)

func compile(exp expression.Expression) ([]*instruction.Instruction, error) {
	switch e := exp.(type) {

	case expression.IntExpr:
		inst := instruction.New(opcode.LOAD_CONST, "", int(e))
		return []*instruction.Instruction{inst}, nil
	case expression.ValExpr:
		code, err := compile(e.Expr)
		if err != nil {
			return nil, err
		}
		inst := instruction.New(opcode.STORE_NAME, e.Name, -1)
		return append(code, inst), nil
	case expression.NameExpr:
		inst := instruction.New(opcode.LOAD_NAME, string(e), -1)
		return []*instruction.Instruction{inst}, nil
	default:
		return nil, fmt.Errorf("unsupported expression type %T", exp)
	}

}

func eval(code []*instruction.Instruction, env *env.Env) int {
	pc := 0
	stack := []int{}
	for pc < len(code) {
		ins := code[pc]
		op := ins.Op
		pc += 1
		switch op {
		case opcode.LOAD_CONST:
			stack = append(stack, ins.Value)
		case opcode.STORE_NAME:
			val := stack[len(stack)-1]
			env.Define(ins.Name, val)
		case opcode.LOAD_NAME:
			val, err := env.Lookup(ins.Name)
			if err != nil {
				log.Fatalln(err)
			}
			stack = append(stack, val)
		}
	}

	if len(stack) > 0 {
		return stack[len(stack)-1]
	}

	return 0
}

func main() {
	env := env.NewEnv(nil)
	env.Define("x", 5)
	code, _ := compile(expression.NameExpr("x"))
	fmt.Println(eval(code, env))

}
