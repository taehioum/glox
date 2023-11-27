package interpreter

import (
	"errors"
	"fmt"

	"github.com/taehioum/glox/pkg/ast/expressions"
	"github.com/taehioum/glox/pkg/ast/statements"
	"github.com/taehioum/glox/pkg/interpreter/environment"
	"github.com/taehioum/glox/pkg/token"
)

// ErrBreak is a sentinel error to break out of a loop.
// looping constructs should catch this error and return nil.
var ErrBreak = fmt.Errorf("break")

// ErrContinue is a sentinel error to continue a loop.
// looping constructs should catch this error, and keep running.
// we get the continue behavior for free by not returning from the loop, without running any code after the continue statement.
var ErrContinue = fmt.Errorf("continue")

type Interpreter struct {
	env *environment.Environment
}

func Interprete(stmts ...statements.Stmt) error {
	i := Interpreter{
		env: environment.NewGlobalEnvironment(),
	}
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
	case statements.Declaration:
		if s.Intializer == nil {
			i.env.Define(s.Name.Lexeme, nil)
			return nil
		}

		v, err := i.evaluate(s.Intializer)
		if err != nil {
			return err
		}
		i.env.Define(s.Name.Lexeme, v)
		return nil
	case statements.Block:
		prev := i.env
		i.env = environment.NewEnclosedEnvironment(prev)
		for _, stmt := range s.Stmts {
			err := stmt.Accept(i.Interprete)
			if err != nil {
				return err
			}
		}
		// restore env
		i.env = prev
		return nil
	case statements.If:
		v, err := i.evaluate(s.Cond)
		if err != nil {
			return err
		}
		if truthy(v) {
			return s.Then.Accept(i.Interprete)
		}
		if s.Else != nil {
			return s.Else.Accept(i.Interprete)
		} else {
			return nil
		}
	case statements.While:
		for {
			v, err := i.evaluate(s.Cond)
			if err != nil {
				return err
			}
			if !truthy(v) {
				break
			}
			err = s.Body.Accept(i.Interprete)
			if errors.Is(err, ErrBreak) {
				return nil
			}
			// by continuing the for loop here, we get the continue behavior for free.
			if errors.Is(err, ErrContinue) {
				continue
			}
			if err != nil {
				return err
			}
		}
		return nil
	// if we had go-to's, break and continue would be implemented as labels and go-to's during the parse pass.
	// but we don't have go-to's, so we use sentinel errors.
	case statements.Break:
		return ErrBreak
	case statements.Continue:
		return ErrContinue
	case statements.Expression:
		_, err := i.evaluate(s.Expr)
		return err
	default:
		fmt.Println("error!")
		return fmt.Errorf("interpreting: unknown statement %T", s)
	}
}

func (i *Interpreter) evaluate(e expressions.Expr) (any, error) {
	switch e := e.(type) {
	case expressions.Literal:
		return e.Value, nil
	case expressions.Assignment:
		val, err := i.eval(e.Value)
		if err != nil {
			return nil, err
		}

		if _, err := i.env.Get(e.Name.Lexeme); err != nil {
			return nil, fmt.Errorf("assignment %w", err)
		}
		i.env.Assign(e.Name.Lexeme, val)
		return val, nil
	case expressions.Grouping:
		return i.eval(e.Expr)
	case expressions.Variable:
		v, err := i.env.Get(e.Name.Lexeme)
		if v == nil {
			return nil, fmt.Errorf("uninitialized variable '%s'", e.Name.Lexeme)
		}
		return v, err
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
		case token.EQUAL:
			// evaluate r-value
			return right, nil
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
			return nil, fmt.Errorf("unknown binary expression %T", e)
		}
	case expressions.Logical:
		lv, err := i.eval(e.Left)
		if err != nil {
			return nil, err
		}
		if e.Operator.Type == token.OR {
			if truthy(lv) {
				return lv, nil
			}
		} else {
			if !truthy(lv) {
				return lv, nil
			}
		}
		return i.eval(e.Right)
	default:
		return nil, fmt.Errorf("unknown expression %v", e)
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

func truthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}
