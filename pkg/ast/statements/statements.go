package statements

import (
	"github.com/taehioum/glox/pkg/ast/expressions"
	"github.com/taehioum/glox/pkg/token"
)

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

type Declaration struct {
	Name       token.Token
	Intializer expressions.Expr
}

func (stmt Declaration) Accept(v Visitor) error {
	return v(stmt)
}
