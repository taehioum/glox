package parser

import (
	"log/slog"

	"github.com/taehioum/glox/pkg/ast/statements"
	"github.com/taehioum/glox/pkg/token"
)

type PrintStatmentParselet struct{}

func (p PrintStatmentParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume PRINT
	expr, err := parser.parseExpr(0)
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after value")
	if err != nil {
		return nil, err
	}
	return statements.Print{Expr: expr}, nil
}

type DeclarationStatementParselet struct{}

func (p DeclarationStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume VAR
	name, err := parser.consumeAndCheck(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	stmt := statements.Declaration{
		Name: name,
	}

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
}

type BlockStatementParselet struct{}

func (p BlockStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume LEFTBRACE
	var stmts []statements.Stmt
	for !parser.isAtEnd() && !parser.check(token.RIGHTBRACE) {
		stmt, err := parser.parseSingleStatement()
		if err != nil {
			return statements.Block{Stmts: stmts}, err
		}
		stmts = append(stmts, stmt)
	}
	_, err := parser.consumeAndCheck(token.RIGHTBRACE, "expected ';' after value")
	if err != nil {
		return statements.Block{Stmts: stmts}, err
	}
	return statements.Block{Stmts: stmts}, nil
}

type ExpressionStatementParselet struct{}

// TODO: split the switch statement into separate parselets
func (p ExpressionStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	slog.Debug("default parse stmt")
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
