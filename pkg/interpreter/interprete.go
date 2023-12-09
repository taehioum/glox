package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/taehioum/glox/pkg/ast/expressions"
	"github.com/taehioum/glox/pkg/ast/statements"
	"github.com/taehioum/glox/pkg/interpreter/environment"
)

// ErrBreak is a sentinel error to break out of a loop.
// looping constructs should catch this error and return nil.
var ErrBreak = fmt.Errorf("break")

// ErrContinue is a sentinel error to continue a loop.
// looping constructs should catch this error, and keep running.
// we get the continue behavior for free by not returning from the loop, without running any code after the continue statement.
var ErrContinue = fmt.Errorf("continue")

// ErrReturn is a sentinel error to return from a function.
type ErrReturn struct {
	Value any
}

func (e ErrReturn) Error() string {
	return fmt.Sprintf("return %v", e.Value)
}

type Interpreter struct {
	env    *environment.Environment
	global *environment.Environment

	writer io.Writer
	reader bufio.Reader
}

func NewInterpreter(writer io.Writer) *Interpreter {
	global := environment.NewGlobalEnvironment()
	i := &Interpreter{
		env:    global,
		global: global,
		writer: writer,
		reader: *bufio.NewReader(os.Stdin),
	}

	i.env.Define("clock", Clock{})
	i.env.Define("print", Print{})
	i.env.Define("input", Input{})

	return i
}

func (i *Interpreter) Interprete(stmts ...statements.Stmt) error {
	for _, stmt := range stmts {
		err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Eval(e expressions.Expr) (any, error) {
	return e.Accept(i)
}
