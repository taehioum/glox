package parser

import (
	"fmt"
	"log/slog"

	"github.com/taehioum/glox/pkg/ast/expressions"
	"github.com/taehioum/glox/pkg/token"
)

type TermParselet struct{}

func (p TermParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceTerm)
	return expressions.Binary{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p TermParselet) precedence() Precedence {
	return PrecedenceTerm
}

type FactorParselet struct{}

func (p FactorParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceFactor)
	return expressions.Binary{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p FactorParselet) precedence() Precedence {
	return PrecedenceFactor
}

type ComparsionParselet struct{}

func (p ComparsionParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceComparison)
	return expressions.Binary{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p ComparsionParselet) precedence() Precedence {
	return PrecedenceComparison
}

type EqualityParselet struct{}

func (p EqualityParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceEquality)
	return expressions.Binary{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p EqualityParselet) precedence() Precedence {
	return PrecedenceEquality
}

type AssignmentParselet struct{}

func (p AssignmentParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	slog.Debug("assignment parselet: ", "left", left)
	expr, err := parser.parseExpr(PrecedenceAssignment - 1)

	variable, ok := left.(expressions.Variable)
	if !ok {
		return nil, fmt.Errorf("line %d's %s: left hand side of assignment must be a variable", token.Ln, token.Lexeme)
	}
	return expressions.Assignment{
		Name:  variable.Name,
		Value: expr,
	}, err
}

func (p AssignmentParselet) precedence() Precedence {
	return PrecedenceAssignment
}

type OrParselet struct{}

func (p OrParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceOr)
	return expressions.Logical{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p OrParselet) precedence() Precedence {
	return PrecedenceOr
}

type AndParselet struct{}

func (p AndParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceAnd)
	return expressions.Logical{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p AndParselet) precedence() Precedence {
	return PrecedenceAnd
}
