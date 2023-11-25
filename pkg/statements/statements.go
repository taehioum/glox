package statements

import "github.com/taehioum/glox/pkg/expressions"

type Visitor func(Stmt) error

type Stmt interface {
	Accept(Visitor) error
}

type Print struct {
	Expr expressions.Expr
}

func (stmt Print) Accept(v Visitor) error {
	return v(stmt)
}

type Expression struct {
	Expr expressions.Expr
}

func (stmt Expression) Accept(v Visitor) error {
	return v(stmt)
}
