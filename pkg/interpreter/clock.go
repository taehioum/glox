package interpreter

import (
	"time"
)

type Clock struct{}

func (f Clock) Arity() int {
	return 0
}

func (f Clock) Call(e *Interpreter, args []any) (any, error) {
	return time.Now().Second(), nil
}
