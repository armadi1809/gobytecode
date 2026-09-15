package builtins

import (
	"fmt"

	"github.com/armadi1809/gobytecode/instruction"
)

type NativeFunc func(args []instruction.Value) (instruction.Value, error)

func BuiltInAdd(args []instruction.Value) (instruction.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("add expects 2 arguments, got %d", len(args))
	}

	a, ok := args[0].(int)
	if !ok {
		return nil, fmt.Errorf("add: first argument must be int")
	}

	b, ok := args[1].(int)
	if !ok {
		return nil, fmt.Errorf("add: second argument must be int")
	}

	return a + b, nil
}

func BuiltInPrint(args []instruction.Value) (instruction.Value, error) {
	_, err := fmt.Print(args...)
	return nil, err
}
