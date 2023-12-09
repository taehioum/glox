package interpreter

type Input struct{}

func (f Input) Arity() int {
	return 0
}

func (f Input) Call(e *Interpreter, args []any) (any, error) {
	return e.reader.ReadString('\n')
}
