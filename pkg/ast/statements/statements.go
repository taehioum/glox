package statements

import (
	"fmt"

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

func (stmt Print) String() string {
	return fmt.Sprintf("Print{Expr: %v}", stmt.Expr)
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

func (stmt Declaration) String() string {
	return fmt.Sprintf("Declaration{Name: %s, Intializer: %s}", stmt.Name, stmt.Intializer)
}

type Block struct {
	Stmts []Stmt
}

func (stmt Block) Accept(v Visitor) error {
	return v(stmt)
}

func (stmt Block) String() string {
	return fmt.Sprintf("Block{Stmts: %+v}", stmt.Stmts)
}
