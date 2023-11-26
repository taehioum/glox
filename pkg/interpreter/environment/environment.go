package environment

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewGlobalEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (env *Environment) Assign(name string, value any) error {
	_, ok := env.values[name]
	if ok {
		env.values[name] = value
		return nil
	}
	if env.enclosing == nil {
		return fmt.Errorf("undefined variable '%s'", name)
	}
	return env.enclosing.Assign(name, value)
}

func (env *Environment) Define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) Get(name string) (any, error) {
	v, ok := env.values[name]
	if ok {
		return v, nil
	}
	if env.enclosing == nil {
		return nil, fmt.Errorf("undefined variable '%s'", name)
	}
	return env.enclosing.Get(name)
}
