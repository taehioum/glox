package pprint

import (
	"fmt"

	"github.com/taehioum/glox/pkg/expressions"
)

func Print(e expressions.Expr) string {
	// define it as a var here, so that it can be recursively called.
	var printer func(expressions.Expr) any
	printer = func(e expressions.Expr) any {
		switch e := e.(type) {
		case expressions.Binary:
			return fmt.Sprintf("(%s %s %s)", e.Operator.Lexeme, e.Left.Accept(printer), e.Right.Accept(printer))
		case expressions.Grouping:
			return fmt.Sprintf("(group %s)", e.Expr.Accept(printer))
		case expressions.Literal:
			if e.Value == nil {
				return "nil"
			}
			return fmt.Sprintf("%v", e.Value)
		case expressions.Unary:
			return fmt.Sprintf("(%s %s)", e.Operator.Lexeme, e.Right.Accept(printer))
		default:
			return "unknown expression"
		}
	}
	return e.Accept(printer).(string)
}
