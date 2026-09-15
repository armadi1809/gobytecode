package env

import (
	"fmt"

	"github.com/armadi1809/gobytecode/builtins"
	"github.com/armadi1809/gobytecode/instruction"
)

type Env struct {
	table  map[string]instruction.Value
	parent *Env
}

func DefaultEnv() *Env {
	env := &Env{table: make(map[string]instruction.Value), parent: nil}
	env.Define("+", builtins.NativeFunc(builtins.BuiltInAdd))
	env.Define("print", builtins.NativeFunc(builtins.BuiltInPrint))

	return env
}

func NewEnv(parent *Env) *Env {
	return &Env{table: make(map[string]instruction.Value), parent: parent}
}

func (e *Env) Define(name string, val instruction.Value) {
	e.table[name] = val
}

func (e *Env) Assign(name string, val instruction.Value) error {
	env, err := e.resolve(name)
	if err != nil {
		return err
	}

	env.Define(name, val)
	return nil
}

func (e *Env) Lookup(name string) (instruction.Value, error) {
	env, err := e.resolve(name)
	if err != nil {
		return -1, err
	}
	return env.table[name], nil
}

func (e *Env) resolve(name string) (*Env, error) {
	if _, ok := e.table[name]; ok {
		return e, nil
	}
	if e.parent == nil {
		return nil, fmt.Errorf("unreferenced variable %s", name)
	}
	return e.parent.resolve(name)
}

func (e *Env) IsDefined(name string) bool {
	_, err := e.resolve(name)
	return err == nil
}
