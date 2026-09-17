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

type Function struct {
	Params []string
	Body   []*instruction.Instruction
	Env    *env.Env
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
	case expression.BeginExpr:
		return compileProgram(e.Exps)
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
	case expression.IfExpr:
		ifTrueCode, err := compile(e.IfTrue)
		if err != nil {
			return nil, err
		}
		ifFalseCode, err := compile(e.IfFalse)
		if err != nil {
			return nil, err
		}
		relJumpFalseInst := &instruction.Instruction{
			Op:      opcode.RELATIVE_JUMP,
			Operand: len(ifTrueCode),
		}
		ifFalseCode = append(ifFalseCode, relJumpFalseInst)
		condCode, err := compile(e.Cond)
		if err != nil {
			return nil, err
		}
		relJumpTrueInst := &instruction.Instruction{
			Op:      opcode.RELATIVE_JUMP_IF_TRUE,
			Operand: len(ifFalseCode),
		}
		code := []*instruction.Instruction{}
		code = append(code, condCode...)
		code = append(code, relJumpTrueInst)
		code = append(code, ifFalseCode...)
		code = append(code, ifTrueCode...)
		return code, nil
	case expression.LambdaExpr:
		body, err := compile(e.Body)
		if err != nil {
			return nil, err
		}

		return []*instruction.Instruction{
			{Op: opcode.LOAD_CONST, Value: e.Params},
			{Op: opcode.LOAD_CONST, Value: body},
			{Op: opcode.MAKE_FUNCTION, Operand: len(e.Params)},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported expression type %T", exp)
	}

}

func compileProgram(prog []expression.Expression) ([]*instruction.Instruction, error) {
	res := []*instruction.Instruction{}
	for _, exp := range prog {
		code, err := compile(exp)
		if err != nil {
			return nil, err
		}
		res = append(res, code...)

	}
	return res, nil
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
			switch fn := function.(type) {
			case builtins.NativeFunc:
				result, err := fn(args)

				if err != nil {
					return nil, err
				}
				vm.stack = append(vm.stack, result)
			case Function:
				if len(args) != len(fn.Params) {
					return nil, fmt.Errorf("function expects %d arguments, got %d",
						len(fn.Params), len(args))
				}

				callEnv := env.NewEnv(fn.Env)
				for index, param := range fn.Params {
					callEnv.Define(param, args[index])
				}

				callVM := VM{env: callEnv}
				result, err := callVM.eval(fn.Body)
				if err != nil {
					return nil, err
				}
				vm.stack = append(vm.stack, result)

			default:
				return nil, fmt.Errorf("%T is not callable", function)
			}

		case opcode.RELATIVE_JUMP_IF_TRUE:
			condVal := vm.pop()
			if cond := condVal.(bool); cond {
				pc += ins.Operand
			}
		case opcode.RELATIVE_JUMP:
			pc += ins.Operand
		case opcode.MAKE_FUNCTION:
			body := vm.pop().([]*instruction.Instruction)
			params := vm.pop().([]string)

			if len(params) != ins.Operand {
				return nil, fmt.Errorf("function parameter count mismatch")
			}

			vm.stack = append(vm.stack, Function{
				Params: params,
				Body:   body,
				Env:    vm.env,
			})
		}
	}

	if len(vm.stack) == 0 {
		return nil, nil
	}

	return vm.pop(), nil
}

func main() {
	vm := VM{
		env: env.DefaultEnv(),
	}

	vm.env.Define("*", builtins.NativeFunc(func(args []instruction.Value) (instruction.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("* expects 2 arguments, got %d", len(args))
		}

		left, leftOK := args[0].(int)
		right, rightOK := args[1].(int)
		if !leftOK || !rightOK {
			return nil, fmt.Errorf("* expects integer arguments")
		}

		return left * right, nil
	}))

	vm.env.Define("-", builtins.NativeFunc(func(args []instruction.Value) (instruction.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("- expects 2 arguments, got %d", len(args))
		}

		left, leftOK := args[0].(int)
		right, rightOK := args[1].(int)
		if !leftOK || !rightOK {
			return nil, fmt.Errorf("- expects integer arguments")
		}

		return left - right, nil
	}))

	vm.env.Define("eq", builtins.NativeFunc(func(args []instruction.Value) (instruction.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("eq expects 2 arguments, got %d", len(args))
		}

		return args[0] == args[1], nil
	}))

	/*
		Would be nice to run something like

		Begin(
			Define("factorial", ["x"], IF(Call("eq", x, 0), Int(1), Call("factorial", [Call("-", x, 1)]))),
			Call("factorial", Int(5))
		)


	*/

	program := expression.BeginExpr{
		Exps: []expression.Expression{
			expression.ValExpr{
				Name: "factorial",
				Expr: expression.LambdaExpr{
					Params: []string{"x"},
					Body: expression.IfExpr{
						Cond: expression.CallExpr{
							Function: expression.NameExpr("eq"),
							Args: []expression.Expression{
								expression.NameExpr("x"),
								expression.IntExpr(0),
							},
						},
						IfTrue: expression.IntExpr(1),
						IfFalse: expression.CallExpr{
							Function: expression.NameExpr("*"),
							Args: []expression.Expression{
								expression.NameExpr("x"),
								expression.CallExpr{
									Function: expression.NameExpr("factorial"),
									Args: []expression.Expression{
										expression.CallExpr{
											Function: expression.NameExpr("-"),
											Args: []expression.Expression{
												expression.NameExpr("x"),
												expression.IntExpr(1),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expression.CallExpr{
				Function: expression.NameExpr("factorial"),
				Args: []expression.Expression{
					expression.IntExpr(5),
				},
			},
		},
	}

	code, err := compile(program)
	if err != nil {
		panic(err)
	}

	result, err := vm.eval(code)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
