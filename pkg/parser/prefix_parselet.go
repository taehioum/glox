package parser

import (
	"errors"
	"fmt"

	"github.com/taehioum/glox/pkg/ast"
	expressions "github.com/taehioum/glox/pkg/ast"
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

type LambdaParselet struct{}

func (p LambdaParselet) parse(parser *Parser, tok token.Token) (expressions.Expr, error) {
	_, err := parser.consumeAndCheck(token.LEFTPAREN, "Expect '(' after fun.")
	if err != nil {
		return nil, err
	}

	var params []token.Token
	// parse the comma-seperated arguments until we hit a ')'
	if !parser.check(token.RIGHTPAREN) {
		ok := true
		for ok {
			id, err := parser.consumeAndCheck(token.IDENTIFIER, "expected identifier")
			if err != nil {
				return nil, err
			}
			if len(params) >= 255 {
				return nil, errors.New("can't have more than 255 parameters")
			}
			params = append(params, id)

			_, err = parser.consumeAndCheck(token.COMMA, "expected ',' after argument")
			ok = err == nil
		}
	}

	_, err = parser.consumeAndCheck(token.RIGHTPAREN, "expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	ok := parser.check(token.LEFTBRACE)
	if !ok {
		return nil, errors.New("expected '{' after function declaration")
	}

	b, err := BlockStatementParselet{}.parse(parser)
	if err != nil {
		return nil, err
	}

	body, ok := b.(ast.Block)
	if !ok {
		return nil, errors.New("expected block statement")
	}

	return expressions.Lambda{
		Name:   tok,
		Params: params,
		Body:   body.Stmts,
	}, nil
}
