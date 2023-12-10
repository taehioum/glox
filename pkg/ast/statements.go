package ast

import (
	"fmt"

	"github.com/taehioum/glox/pkg/token"
)

type StatementVistior interface {
	// VisitPrint(Print) error
	VisitDeclaration(Declaration) error
	VisitBlock(Block) error
	VisitIf(If) error
	VisitWhile(While) error
	VisitBreak(Break) error
	VisitContinue(Continue) error
	VisitReturn(Return) error
	VisitExpression(Expression) error
}

type Stmt interface {
	Accept(StatementVistior) error
}

type Print struct {
	Expr Expr
}

func (stmt Print) Accept(v StatementVistior) error {
	return nil
}

func (stmt Print) String() string {
	return fmt.Sprintf("Print{Expr: %v}", stmt.Expr)
}

type Expression struct {
	Expr Expr
}

func (stmt Expression) Accept(v StatementVistior) error {
	return v.VisitExpression(stmt)
}

type Declaration struct {
	Name        token.Token
	Initializer Expr
}

func (stmt Declaration) Accept(v StatementVistior) error {
	return v.VisitDeclaration(stmt)
}

func (stmt Declaration) String() string {
	return fmt.Sprintf("Declaration{Name: %s, Intializer: %s}", stmt.Name, stmt.Initializer)
}

type Block struct {
	Stmts []Stmt
}

func (stmt Block) Accept(v StatementVistior) error {
	return v.VisitBlock(stmt)
}

func (stmt Block) String() string {
	return fmt.Sprintf("Block{Stmts: %+v}", stmt.Stmts)
}

type If struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

func (stmt If) Accept(v StatementVistior) error {
	return v.VisitIf(stmt)
}

type While struct {
	Cond Expr
	Body Stmt
}

func (stmt While) Accept(v StatementVistior) error {
	return v.VisitWhile(stmt)
}

type Break struct{}

func (stmt Break) Accept(v StatementVistior) error {
	return v.VisitBreak(stmt)
}

func (stmt Break) String() string {
	return "Break{}"
}

type Continue struct{}

func (stmt Continue) Accept(v StatementVistior) error {
	return v.VisitContinue(stmt)
}

func (stmt Continue) String() string {
	return "Continue{}"
}

type Return struct {
	Keyword token.Token
	Value   Expr
}

func (stmt Return) Accept(v StatementVistior) error {
	return v.VisitReturn(stmt)
}
