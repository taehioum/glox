package interpreter

import (
	"fmt"
	"io"
)

type Print struct{}

func (f Print) Arity() int {
	return -1
}

func (f Print) Call(e *Interpreter, args []any) (any, error) {
	s := fmt.Sprintln(args...)
	return io.WriteString(e.writer, s)
}
