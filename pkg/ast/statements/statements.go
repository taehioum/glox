package statements

import (
	"fmt"

	"github.com/taehioum/glox/pkg/ast/expressions"
	"github.com/taehioum/glox/pkg/token"
)

type Visitor interface {
	VisitPrint(Print) error
	VisitDeclaration(Declaration) error
	VisitBlock(Block) error
	VisitIf(If) error
	VisitWhile(While) error
	VisitBreak(Break) error
	VisitContinue(Continue) error
	VisitExpression(Expression) error
}

// type Visitor func(Stmt) error

type Stmt interface {
	Accept(Visitor) error
}

type Print struct {
	Expr expressions.Expr
}

func (stmt Print) Accept(v Visitor) error {
	return v.VisitPrint(stmt)
}

func (stmt Print) String() string {
	return fmt.Sprintf("Print{Expr: %v}", stmt.Expr)
}

type Expression struct {
	Expr expressions.Expr
}

func (stmt Expression) Accept(v Visitor) error {
	return v.VisitExpression(stmt)
}

type Declaration struct {
	Name       token.Token
	Intializer expressions.Expr
}

func (stmt Declaration) Accept(v Visitor) error {
	return v.VisitDeclaration(stmt)
}

func (stmt Declaration) String() string {
	return fmt.Sprintf("Declaration{Name: %s, Intializer: %s}", stmt.Name, stmt.Intializer)
}

type Block struct {
	Stmts []Stmt
}

func (stmt Block) Accept(v Visitor) error {
	return v.VisitBlock(stmt)
}

func (stmt Block) String() string {
	return fmt.Sprintf("Block{Stmts: %+v}", stmt.Stmts)
}

type If struct {
	Cond expressions.Expr
	Then Stmt
	Else Stmt
}

func (stmt If) Accept(v Visitor) error {
	return v.VisitIf(stmt)
}

type While struct {
	Cond expressions.Expr
	Body Stmt
}

func (stmt While) Accept(v Visitor) error {
	return v.VisitWhile(stmt)
}

type Break struct{}

func (stmt Break) Accept(v Visitor) error {
	return v.VisitBreak(stmt)
}

func (stmt Break) String() string {
	return "Break{}"
}

type Continue struct{}

func (stmt Continue) Accept(v Visitor) error {
	return v.VisitContinue(stmt)
}

func (stmt Continue) String() string {
	return "Continue{}"
}
