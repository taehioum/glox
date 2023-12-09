package expressions

import (
	"fmt"

	"github.com/taehioum/glox/pkg/token"
)

type Visitor interface {
	VisitAssignment(Assignment) (any, error)
	VisitBinary(Binary) (any, error)
	VisitGrouping(Grouping) (any, error)
	VisitLiteral(Literal) (any, error)
	VisitUnary(Unary) (any, error)
	VisitVariable(Variable) (any, error)
	VisitLogical(Logical) (any, error)
	VisitPostUnary(PostUnary) (any, error)
	VisitCall(Call) (any, error)
}

type Expr interface {
	Accept(Visitor) (any, error)
}

type Assignment struct {
	Name  token.Token
	Value Expr
}

func (e Assignment) Accept(v Visitor) (any, error) {
	return v.VisitAssignment(e)
}

func (e Assignment) String() string {
	return fmt.Sprintf("Assignment{Name: %s, Value: %s}", e.Name, e.Value)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e Binary) Accept(v Visitor) (any, error) {
	return v.VisitBinary(e)
}

type Grouping struct {
	Expr Expr
}

func (e Grouping) Accept(v Visitor) (any, error) {
	return v.VisitGrouping(e)
}

type Literal struct {
	Value any
}

func (e Literal) Accept(v Visitor) (any, error) {
	return v.VisitLiteral(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (e Unary) Accept(v Visitor) (any, error) {
	return v.VisitUnary(e)
}

type Variable struct {
	Name token.Token
}

func (e Variable) Accept(v Visitor) (any, error) {
	return v.VisitVariable(e)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e Logical) Accept(v Visitor) (any, error) {
	return v.VisitLogical(e)
}

type PostUnary struct {
	Left     Expr
	Operator token.Token
}

func (e PostUnary) Accept(v Visitor) (any, error) {
	return v.VisitPostUnary(e)
}

type Call struct {
	Callee Expr
	Args   []Expr

	// used to report error on the location of the closing paren
	Paren token.Token
}

func (e Call) Accept(v Visitor) (any, error) {
	return v.VisitCall(e)
}
