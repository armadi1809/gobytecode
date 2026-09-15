package env

import (
	"fmt"
)

type Env struct {
	table  map[string]int
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{table: make(map[string]int), parent: parent}
}

func (e *Env) Define(name string, val int) {
	e.table[name] = val
}

func (e *Env) Assign(name string, val int) error {
	env, err := e.resolve(name)
	if err != nil {
		return err
	}

	env.Define(name, val)
	return nil
}

func (e *Env) Lookup(name string) (int, error) {
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
