package parser

import (
	"fmt"

	"github.com/taehioum/glox/pkg/ast/expressions"
	"github.com/taehioum/glox/pkg/token"
)

type UnaryOperatorParselet struct {
}

func (uop UnaryOperatorParselet) parse(parser *Parser, token token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(PrecedenceUnary)
	return expressions.Unary{
		Operator: token,
		Right:    expr,
	}, err
}

type LiteralParselet struct {
}

func (lp LiteralParselet) parse(parser *Parser, token token.Token) (expressions.Expr, error) {
	return expressions.Literal{
		Value: token.Literal,
	}, nil
}

type BoolParselet struct {
}

func (bp BoolParselet) parse(parser *Parser, tok token.Token) (expressions.Expr, error) {
	switch tok.Type {
	case token.TRUE:
		return expressions.Literal{Value: true}, nil
	case token.FALSE:
		return expressions.Literal{Value: false}, nil
	default:
		return nil, fmt.Errorf("unexpected token type %s", tok.Type)
	}
}

type GroupParselet struct{}

func (p GroupParselet) parse(parser *Parser, tok token.Token) (expressions.Expr, error) {
	expr, err := parser.parseExpr(0)
	if err != nil {
		return expr, err
	}

	_, err = parser.consumeAndCheck(token.RIGHTPAREN, "expected ')' after expression")
	if err != nil {
		return expr, err
	}

	return expressions.Grouping{
		Expr: expr,
	}, nil
}

type VariableParselet struct{}

func (p VariableParselet) parse(parser *Parser, tok token.Token) (expressions.Expr, error) {
	return expressions.Variable{
		Name: tok,
	}, nil
}

type FunctionParselet struct{}
