package pprint

import (
	"fmt"

	"github.com/taehioum/glox/pkg/ast/expressions"
)

func Print(e expressions.Expr) string {
	// define it as a var here, so that it can be recursively called.
	var printer func(expressions.Expr) (any, error)
	printer = func(e expressions.Expr) (any, error) {
		switch e := e.(type) {
		case expressions.Binary:
			l, _ := e.Left.Accept(printer)
			r, _ := e.Right.Accept(printer)
			return fmt.Sprintf("(%s %s %s)", e.Operator.Lexeme, l, r), nil
		case expressions.Grouping:
			l, _ := e.Expr.Accept(printer)
			return fmt.Sprintf("(group %s)", l), nil
		case expressions.Literal:
			if e.Value == nil {
				return "nil", nil
			}
			return fmt.Sprintf("%v", e.Value), nil
		case expressions.Unary:
			r, _ := e.Right.Accept(printer)
			return fmt.Sprintf("(%s %s)", e.Operator.Lexeme, r), nil
		default:
			return "unknown expression", fmt.Errorf("unknown expression %T", e)
		}
	}

	s, _ := e.Accept(printer)
	return s.(string)
}
