package parser

import (
	"fmt"

	"github.com/taehioum/glox/pkg/ast/statements"
	"github.com/taehioum/glox/pkg/token"
)

type StatementParselet struct{}

func (p StatementParselet) parse(parser *Parser, tok token.Token) (statements.Stmt, error) {
	switch tok.Type {
	case token.VAR:
		name, err := parser.consumeAndCheck(token.IDENTIFIER, "Expect variable name.")
		if err != nil {
			return nil, err
		}
		stmt := statements.Declaration{
			Name: name,
		}
		fmt.Printf("%+v\n", stmt)

		if parser.check(token.EQUAL) {
			parser.consume()
			init, err := parser.parseExpr(0)
			if err != nil {
				return nil, err
			}
			stmt.Intializer = init
		}

		_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after var declaration")
		if err != nil {
			return nil, err
		}

		return stmt, nil
	case token.PRINT:
		expr, err := parser.parseExpr(0)
		if err != nil {
			return nil, err
		}

		_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after value")
		if err != nil {
			return nil, err
		}
		return statements.Print{Expr: expr}, nil
	default:
		expr, err := parser.parseExpr(0)
		if err != nil {
			return nil, err
		}

		_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after expression")
		if err != nil {
			return nil, err
		}
		return statements.Expression{Expr: expr}, nil
	}
}
