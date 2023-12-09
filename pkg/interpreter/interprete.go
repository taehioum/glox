package interpreter

import (
	"errors"
	"fmt"
	"io"
	"os"

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

	writer io.Writer
}

func Interprete(stmts ...statements.Stmt) error {
	i := &Interpreter{
		env:    environment.NewGlobalEnvironment(),
		writer: os.Stdout,
	}
	for _, stmt := range stmts {
		err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) eval(e expressions.Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) VisitPrint(stmt statements.Print) error {
	v, err := i.eval(stmt.Expr)
	if err != nil {
		return err
	}
	i.writer.Write([]byte(fmt.Sprintf("%v\n", v)))
	return nil
}

func (i *Interpreter) VisitDeclaration(stmt statements.Declaration) error {
	if stmt.Intializer == nil {
		i.env.Define(stmt.Name.Lexeme, nil)
		return nil
	}

	v, err := i.eval(stmt.Intializer)
	if err != nil {
		return err
	}
	i.env.Define(stmt.Name.Lexeme, v)
	return nil
}

func (i *Interpreter) VisitBlock(stmt statements.Block) error {
	prev := i.env
	i.env = environment.NewEnclosedEnvironment(prev)
	for _, stmt := range stmt.Stmts {
		err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}
	// restore env
	i.env = prev
	return nil
}

func (i *Interpreter) VisitIf(stmt statements.If) error {
	v, err := i.eval(stmt.Cond)
	if err != nil {
		return err
	}
	if truthy(v) {
		return stmt.Then.Accept(i)
	}
	if stmt.Else != nil {
		return stmt.Else.Accept(i)
	}
	return nil
}

func (i *Interpreter) VisitWhile(stmt statements.While) error {

	for {
		v, err := i.eval(stmt.Cond)
		if err != nil {
			return err
		}
		if !truthy(v) {
			break
		}
		err = stmt.Body.Accept(i)
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
}

func (i *Interpreter) VisitBreak(stmt statements.Break) error {
	return ErrBreak
}

func (i *Interpreter) VisitContinue(stmt statements.Continue) error {
	return ErrContinue
}

func (i *Interpreter) VisitExpression(stmt statements.Expression) error {
	_, err := i.eval(stmt.Expr)
	return err
}

func (i *Interpreter) VisitAssignment(e expressions.Assignment) (any, error) {
	v, err := i.eval(e.Value)
	if err != nil {
		return nil, err
	}
	i.env.Assign(e.Name.Lexeme, v)
	return v, nil
}

func (i *Interpreter) VisitLiteral(e expressions.Literal) (any, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGrouping(e expressions.Grouping) (any, error) {
	return i.eval(e.Expr)
}

func (i *Interpreter) VisitVariable(e expressions.Variable) (any, error) {
	v, err := i.env.Get(e.Name.Lexeme)
	if v == nil {
		return nil, fmt.Errorf("uninitialized variable '%s'", e.Name.Lexeme)
	}
	return v, err
}

func (i *Interpreter) VisitUnary(e expressions.Unary) (any, error) {
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
}

func (i *Interpreter) VisitPostUnary(e expressions.PostUnary) (any, error) {
	v, err := i.eval(e.Left)
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
}

func (i *Interpreter) VisitLogical(e expressions.Logical) (any, error) {
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
