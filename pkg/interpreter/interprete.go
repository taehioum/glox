package interpreter

import (
	"fmt"

	"github.com/taehioum/glox/pkg/expressions"
	"github.com/taehioum/glox/pkg/statements"
	"github.com/taehioum/glox/pkg/token"
)

type Interpreter struct {
}

func Interprete(stmts ...statements.Stmt) error {
	i := Interpreter{}
	for _, stmt := range stmts {
		err := stmt.Accept(i.Interprete)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Interprete(s statements.Stmt) error {
	switch s := s.(type) {
	case statements.Print:
		v, err := i.evaluate(s.Expr)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", v)
		return nil
	default:
		panic(fmt.Sprintf("unknown statement %T", s))
	}
}

func (i *Interpreter) evaluate(e expressions.Expr) (any, error) {
	switch e := e.(type) {
	case expressions.Literal:
		return e.Value, nil
	case expressions.Grouping:
		return i.eval(e.Expr)
	case expressions.Unary:
		right, err := i.eval(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Type {
		case token.MINUS:
			if n, ok := right.(float64); ok {
				return -n, nil
			}
			return nil, fmt.Errorf("expected numbers, got %T", right)
		case token.BANG:
			if right == nil { // nil is falsy
				return true, nil
			}
			if b, ok := right.(bool); ok { // negate the bool
				return !b, nil
			}
			// anything non-bool is truthy
			return false, nil
		default:
			return nil, fmt.Errorf("unknown expression %T", e)
		}
	case expressions.Binary:
		l, err := i.eval(e.Left)
		if err != nil {
			return nil, err
		}
		r, err := i.eval(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Type {
		case token.MINUS:
			if !checkNumberOperands(l, r) {
				return nil, fmt.Errorf("operands must be numbers")
			}
			return l.(float64) - r.(float64), nil
		case token.SLASH:
			if !checkNumberOperands(l, r) {
				return nil, fmt.Errorf("operands must be numbers")
			}
			return l.(float64) / r.(float64), nil
		case token.STAR:
			if !checkNumberOperands(l, r) {
				return nil, fmt.Errorf("operands must be numbers")
			}
			return l.(float64) * r.(float64), nil
		case token.PLUS: // todo: tidy
			if checkNumberOperands(l, r) {
				return l.(float64) + r.(float64), nil
			} else if checkStringOperands(l, r) {
				return l.(string) + r.(string), nil
			} else {
				return nil, fmt.Errorf("operands must be numbers or strings")
			}
		case token.GREATER:
			return l.(float64) > r.(float64), nil
		case token.GREATEREQUAL:
			return l.(float64) >= r.(float64), nil
		case token.LESS:
			return l.(float64) < r.(float64), nil
		case token.LESSEQUAL:
			return l.(float64) <= r.(float64), nil
		case token.BANGEQUAL: // deep equality for numbers, bools, strings
			return l != r, nil
		case token.EQUALEQUAL: // deep equality for numbers, bools, strings
			return l == r, nil
		default:
			return nil, fmt.Errorf("unknown expression %T", e)
		}
	default:
		return nil, fmt.Errorf("unknown expression %T", e)
	}
}

func (i *Interpreter) eval(e expressions.Expr) (any, error) {
	return e.Accept(i.evaluate)
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
