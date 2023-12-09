package parser

import (
	expressions "github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/token"
)

type PostfixParselet struct{}

func (p PostfixParselet) parse(parser *Parser, left expressions.Expr, token token.Token) (expressions.Expr, error) {
	return expressions.PostUnary{Left: left, Operator: token}, nil
}

func (p PostfixParselet) precedence() Precedence {
	return PrecedencePostfix
}
