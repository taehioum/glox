package interpreter

import (
	"errors"
	"fmt"

	statements "github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/interpreter/environment"
)

type Function struct {
	def     statements.Lambda
	closure *environment.Environment
}

func (f Function) Arity() int {
	return len(f.def.Params)
}

func (f Function) Call(i *Interpreter, args []any) (any, error) {
	prev := i.env
	defer func() {
		// restore env
		i.env = prev
	}()
	i.env = environment.NewEnclosedEnvironment(f.closure)
	for idx, param := range f.def.Params {
		i.env.Define(param.Lexeme, args[idx])
	}

	err := i.Interprete(f.def.Body...)
	var res ErrReturn
	if errors.As(err, &res) {
		return res.Value, nil
	} else if err != nil {
		return nil, fmt.Errorf("calling %s defined on line %d: %w", f.def.Name.Lexeme, f.def.Name.Ln, err)
	}

	return nil, nil
}
