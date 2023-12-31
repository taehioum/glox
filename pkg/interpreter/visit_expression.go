package interpreter

import (
	"fmt"

	expressions "github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/token"
)

func (i *Interpreter) VisitAssignment(e expressions.Assignment) (any, error) {
	v, err := i.Eval(e.Value)
	if err != nil {
		return nil, err
	}

	if distance, ok := i.Locals[e]; ok {
		i.env.AssignAt(distance, e.Name.Lexeme, v)
	} else {
		i.global.Assign(e.Name.Lexeme, v)
	}
	return v, nil
}

func (i *Interpreter) VisitLiteral(e expressions.Literal) (any, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGrouping(e expressions.Grouping) (any, error) {
	return i.Eval(e.Expr)
}

func (i *Interpreter) VisitVariable(e expressions.Variable) (any, error) {
	return i.lookup(e)
}

func (i *Interpreter) VisitUnary(e expressions.Unary) (any, error) {
	right, err := i.Eval(e.Right)
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
}

func (i *Interpreter) VisitPostUnary(e expressions.PostUnary) (any, error) {
	v, err := i.Eval(e.Left)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case token.PLUSPLUS:
		if n, ok := v.(float64); ok {
			i.env.Assign(e.Left.(expressions.Variable).Name.Lexeme, n+1)
			return n + 1, nil
		}
		return nil, fmt.Errorf("expected numbers, got %T", v)
	case token.MINUSMINUS:
		if n, ok := v.(float64); ok {
			i.env.Assign(e.Left.(expressions.Variable).Name.Lexeme, n-1)
			return n - 1, nil
		}
		return nil, fmt.Errorf("expected numbers, got %T", v)
	default:
		return nil, fmt.Errorf("unknown expression %T", e)
	}
}

func (i *Interpreter) VisitBinary(e expressions.Binary) (any, error) {
	l, err := i.Eval(e.Left)
	if err != nil {
		return nil, err
	}
	r, err := i.Eval(e.Right)
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
}

func (i *Interpreter) VisitLogical(e expressions.Logical) (any, error) {
	lv, err := i.Eval(e.Left)
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
	return i.Eval(e.Right)
}

func (i *Interpreter) VisitCall(e expressions.Call) (any, error) {
	callee, err := i.Eval(e.Callee)
	if err != nil {
		return nil, err
	}

	args := make([]any, len(e.Args))
	for idx, arg := range e.Args {
		v, err := i.Eval(arg)
		if err != nil {
			return nil, err
		}
		args[idx] = v
	}

	if fn, ok := callee.(Callable); ok {
		if fn.Arity() != -1 && len(args) != fn.Arity() {
			return nil, fmt.Errorf("expected %d arguments, got %d", fn.Arity(), len(args))
		}
		v, err := fn.Call(i, args)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", e.Paren.Ln, err)
		}
		return v, nil
	} else {
		return nil, fmt.Errorf("can only call functions and classes")
	}
}

func (i *Interpreter) VisitLambda(e expressions.Lambda) (any, error) {
	return Function{def: e, closure: i.env}, nil
}
