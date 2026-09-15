package main

import (
	"fmt"
	"log"

	"github.com/armadi1809/gobytecode/builtins"
	"github.com/armadi1809/gobytecode/env"
	"github.com/armadi1809/gobytecode/expression"
	"github.com/armadi1809/gobytecode/instruction"
	"github.com/armadi1809/gobytecode/opcode"
)

type VM struct {
	stack []instruction.Value
	env   *env.Env
}

func compile(exp expression.Expression) ([]*instruction.Instruction, error) {
	switch e := exp.(type) {

	case expression.IntExpr:
		inst := instruction.New(opcode.LOAD_CONST, -1, "", int(e))
		return []*instruction.Instruction{inst}, nil
	case expression.ValExpr:
		code, err := compile(e.Expr)
		if err != nil {
			return nil, err
		}
		inst := instruction.New(opcode.STORE_NAME, -1, e.Name, -1)
		return append(code, inst), nil
	case expression.NameExpr:
		inst := instruction.New(opcode.LOAD_NAME, -1, string(e), -1)
		return []*instruction.Instruction{inst}, nil
	case expression.CallExpr:
		code := []*instruction.Instruction{}
		functionCode, err := compile(e.Function)
		if err != nil {
			return nil, err
		}

		code = append(code, functionCode...)

		for _, arg := range e.Args {
			argCode, err := compile(arg)
			if err != nil {
				return nil, err
			}

			code = append(code, argCode...)
		}

		code = append(code, instruction.New(opcode.CALL_FUNCTION, len(e.Args), "", -1))
		return code, nil
	default:
		return nil, fmt.Errorf("unsupported expression type %T", exp)
	}

}

func (vm *VM) pop() instruction.Value {
	n := len(vm.stack)
	value := vm.stack[n-1]
	vm.stack = vm.stack[:n-1]
	return value
}

func (vm *VM) eval(code []*instruction.Instruction) (instruction.Value, error) {
	pc := 0
	for pc < len(code) {
		ins := code[pc]
		op := ins.Op
		pc += 1
		switch op {
		case opcode.LOAD_CONST:
			vm.stack = append(vm.stack, ins.Value)
		case opcode.STORE_NAME:
			val := vm.stack[len(vm.stack)-1]
			vm.env.Define(ins.Name, val)
		case opcode.LOAD_NAME:
			val, err := vm.env.Lookup(ins.Name)
			if err != nil {
				log.Fatalln(err)
			}
			vm.stack = append(vm.stack, val)
		case opcode.CALL_FUNCTION:
			args := make([]instruction.Value, ins.Operand)
			for i := len(args) - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			function := vm.pop()

			fn, ok := function.(builtins.NativeFunc)

			if !ok {
				return nil, fmt.Errorf("%T is not callable", function)
			}

			result, err := fn(args)

			if err != nil {
				return nil, err
			}
			vm.stack = append(vm.stack, result)

		}
	}

	if len(vm.stack) == 0 {
		return nil, nil
	}

	return vm.pop(), nil
}

func main() {
	exp := expression.CallExpr{
		Function: expression.NameExpr("print"),
		Args: []expression.Expression{
			expression.IntExpr(1),
			expression.IntExpr(2),
		},
	}

	code, err := compile(exp)
	if err != nil {
		panic(err)
	}

	vm := VM{
		env: env.DefaultEnv(),
	}

	_, err = vm.eval(code)
	if err != nil {
		panic(err)
	}

}
