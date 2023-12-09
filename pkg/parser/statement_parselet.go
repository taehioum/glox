package parser

import (
	"fmt"
	"log/slog"

	"github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/token"
)

type PrintStatmentParselet struct{}

func (p PrintStatmentParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume PRINT
	expr, err := parser.parseExpr(0)
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after value")
	if err != nil {
		return nil, err
	}
	return ast.Print{Expr: expr}, nil
}

type DeclarationStatementParselet struct{}

func (p DeclarationStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume VAR
	name, err := parser.consumeAndCheck(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	stmt := ast.Declaration{
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

func (p BlockStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume LEFTBRACE
	var stmts []ast.Stmt
	for !parser.isAtEnd() && !parser.check(token.RIGHTBRACE) {
		stmt, err := parser.parseSingleStatement()
		if err != nil {
			return ast.Block{Stmts: stmts}, err
		}
		stmts = append(stmts, stmt)
	}
	_, err := parser.consumeAndCheck(token.RIGHTBRACE, "expected ';' after value")
	if err != nil {
		return ast.Block{Stmts: stmts}, err
	}
	return ast.Block{Stmts: stmts}, nil
}

type IfStatementParselet struct{}

func (p IfStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume IF
	parser.consumeAndCheck(token.LEFTPAREN, "expected '(' after if")
	cond, err := parser.parseExpr(0)
	if err != nil {
		return ast.If{}, fmt.Errorf("if condition: %w", err)
	}
	parser.consumeAndCheck(token.RIGHTPAREN, "Expect ')' after if condition.")
	then, err := parser.parseSingleStatement()
	if err != nil {
		return ast.If{}, fmt.Errorf("if then: %w", err)
	}

	var elseBranch ast.Stmt
	if parser.check(token.ELSE) {
		parser.consume() // consume ELSE
		elseBranch, err = parser.parseSingleStatement()
		if err != nil {
			return ast.If{}, fmt.Errorf("if else: %w", err)
		}
	}

	return ast.If{
		Cond: cond,
		Then: then,
		Else: elseBranch,
	}, nil
}

type WhileStatementParselet struct{}

func (p WhileStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume WHILE
	_, err := parser.consumeAndCheck(token.LEFTPAREN, "expected '(' after while's condition expression")
	if err != nil {
		return nil, err
	}
	cond, err := parser.parseExpr(0)
	if err != nil {
		return ast.If{}, fmt.Errorf("if condition: %w", err)
	}
	_, err = parser.consumeAndCheck(token.RIGHTPAREN, "Expect ')' after while's condition expression")
	if err != nil {
		return nil, err
	}
	body, err := parser.parseSingleStatement()
	if err != nil {
		return ast.While{}, fmt.Errorf("while body: %w", err)
	}

	return ast.While{
		Cond: cond,
		Body: body,
	}, nil
}

type ForStatementParselet struct{}

func (p ForStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume FOR
	_, err := parser.consumeAndCheck(token.LEFTPAREN, "expected '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var init ast.Stmt
	if parser.check(token.SEMICOLON) {
		init = nil
		parser.consume()
	} else if parser.check(token.VAR) {
		init, err = DeclarationStatementParselet{}.parse(parser)
		if err != nil {
			return nil, err
		}
	} else {
		init, err = ExpressionStatementParselet{}.parse(parser)
		if err != nil {
			return nil, err
		}
	}
	// _, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after init condition")
	// if err != nil {
	// 	return nil, err
	// }

	var cond ast.Expr
	if !parser.check(token.SEMICOLON) {
		cond, err = parser.parseExpr(0)
		if err != nil {
			return nil, err
		}
	}
	_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var incr ast.Expr
	if !parser.check(token.RIGHTPAREN) {
		incr, err = parser.parseExpr(0)
		if err != nil {
			return nil, err
		}
	}
	_, err = parser.consumeAndCheck(token.RIGHTPAREN, "expected ')' after incr expression")
	if err != nil {
		return nil, err
	}

	body, err := parser.parseSingleStatement()
	if err != nil {
		return nil, err
	}

	var res ast.Stmt
	if incr != nil {
		res = ast.Block{
			Stmts: []ast.Stmt{
				body,
				ast.Expression{Expr: incr},
			},
		}
	} else {
		res = body
	}

	if cond == nil {
		cond = ast.Literal{Value: true}
	}
	res = ast.While{Cond: cond, Body: res}

	if init != nil {
		res = ast.Block{
			Stmts: []ast.Stmt{
				init,
				res,
			},
		}
	}

	return res, nil
}

type ExpressionStatementParselet struct{}

func (p ExpressionStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	slog.Debug("default parse stmt")
	expr, err := parser.parseExpr(0)
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeAndCheck(token.SEMICOLON, "expected ';' after expression")
	if err != nil {
		return nil, err
	}
	return ast.Expression{Expr: expr}, nil
}

// TODO
type BreakStatementParselet struct{}

func (p BreakStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume BREAK
	_, err := parser.consumeAndCheck(token.SEMICOLON, "expected ';' after break")
	if err != nil {
		return nil, err
	}
	// check that we are in a loop block
	return ast.Break{}, nil
}

// TODO
type ContinueStatementParslet struct{}

func (p ContinueStatementParslet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume CONTINUE
	_, err := parser.consumeAndCheck(token.SEMICOLON, "expected ';' after continue")
	if err != nil {
		return nil, err
	}
	// check that we are in a loop block
	return ast.Continue{}, nil
}

type FunctionDeclarationStatementParselet struct{}

func (p FunctionDeclarationStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	parser.consume() // consume FUN
	name, err := parser.consumeAndCheck(token.IDENTIFIER, "Expect function name.")
	if err != nil {
		return nil, err
	}

	lambda, err := LambdaParselet{}.parse(parser, name)
	if err != nil {
		return nil, err
	}

	return ast.Declaration{
		Name:       name,
		Intializer: lambda,
	}, nil
}

type ReturnStatementParselet struct{}

func (p ReturnStatementParselet) parse(parser *Parser) (ast.Stmt, error) {
	t := parser.consume() // consume RETURN
	var expr ast.Expr
	if !parser.check(token.SEMICOLON) {
		var err error
		expr, err = parser.parseExpr(0)
		if err != nil {
			return nil, err
		}
	}
	_, err := parser.consumeAndCheck(token.SEMICOLON, "expected ';' after return")
	if err != nil {
		return nil, err
	}
	return ast.Return{Keyword: t, Value: expr}, nil
}
