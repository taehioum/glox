package parser

import (
	"github.com/taehioum/glox/pkg/expressions"
	"github.com/taehioum/glox/pkg/token"
)

type TermParselet struct{}

func (p TermParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parse(PrecedenceTerm)
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
	expr, err := parser.parse(PrecedenceFactor)
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
	expr, err := parser.parse(PrecedenceComparison)
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
	expr, err := parser.parse(PrecedenceEquality)
	return expressions.Binary{
		Left:     left,
		Operator: token,
		Right:    expr,
	}, err
}

func (p EqualityParselet) precedence() Precedence {
	return PrecedenceEquality
}
