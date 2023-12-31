package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/taehioum/glox/pkg/ast"
	expressions "github.com/taehioum/glox/pkg/ast"
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

	Locals map[any]int

	writer io.Writer
	reader bufio.Reader
}

func New(writer io.Writer) *Interpreter {
	global := environment.NewGlobalEnvironment()
	i := &Interpreter{
		env:    global,
		global: global,
		writer: writer,
		reader: *bufio.NewReader(os.Stdin),
		Locals: make(map[any]int),
	}

	i.global.Define("clock", Clock{})
	i.global.Define("print", Print{})
	i.global.Define("input", Input{})

	return i
}

func (i *Interpreter) Interprete(stmts ...ast.Stmt) error {
	for _, stmt := range stmts {
		err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Eval(e ast.Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) Resolve(e ast.Expr, depth int) {
	slog.Debug("resolving", slog.Attr{Key: "locals", Value: slog.AnyValue(i.Locals)})
	i.Locals[e] = depth
}

func (i *Interpreter) lookup(e expressions.Variable) (any, error) {
	if distance, ok := i.Locals[e]; ok {
		return i.env.GetAt(distance, e.Name.Lexeme)
	}
	return i.global.Get(e.Name.Lexeme)
}

type Callable interface {
	Call(i *Interpreter, args []any) (any, error)
	Arity() int
}

// TODO: we might return type checked value, so we don't have to type check twice.
func checkNumberOperands(l, r any) bool {
	_, ok := l.(float64)
	if !ok {
		return false
	}
	_, ok = r.(float64)
	return ok
}
func checkStringOperands(l, r any) bool {
	_, ok := l.(string)
	if !ok {
		return false
	}
	_, ok = r.(string)
	return ok
}

func truthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}
