package object

type Environment struct {
	vars  map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.vars[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.vars[name] = val
	return val
}

func NewEnvironment() *Environment {
	env := &Environment{}
	env.vars = make(map[string]Object)
	env.outer = nil
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
