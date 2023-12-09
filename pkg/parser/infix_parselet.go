package parser

import (
	"errors"
	"fmt"
	"log/slog"

	expressions "github.com/taehioum/glox/pkg/ast"
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

// CallParselet parses function calls like a(b, c, d)
type CallParselet struct{}

func (p CallParselet) parse(parser *Parser, left expressions.Expr, tok token.Token) (expressions.Expr, error) {
	var args []expressions.Expr
	// parse the comma-seperated arguments until we hit a ')'
	if !parser.check(token.RIGHTPAREN) {
		ok := true
		for ok {
			expr, err := parser.parseExpr(0)
			if err != nil {
				return nil, err
			}
			if len(args) >= 255 {
				return nil, errors.New("can't have more than 255 arguments")
			}
			args = append(args, expr)

			_, err = parser.consumeAndCheck(token.COMMA, "expected ',' after argument")
			ok = err == nil
		}
	}

	rightParen, err := parser.consumeAndCheck(token.RIGHTPAREN, "parsing call: expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	return expressions.Call{
		Callee: left,
		Args:   args,
		Paren:  rightParen,
	}, nil
}

func (p CallParselet) precedence() Precedence {
	return PrecedenceCall
}
