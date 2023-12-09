package parser

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/taehioum/glox/pkg/ast/expressions"
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

type IfStatementParselet struct{}

func (p IfStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume IF
	parser.consumeAndCheck(token.LEFTPAREN, "expected '(' after if")
	cond, err := parser.parseExpr(0)
	if err != nil {
		return statements.If{}, fmt.Errorf("if condition: %w", err)
	}
	parser.consumeAndCheck(token.RIGHTPAREN, "Expect ')' after if condition.")
	then, err := parser.parseSingleStatement()
	if err != nil {
		return statements.If{}, fmt.Errorf("if then: %w", err)
	}

	var elseBranch statements.Stmt
	if parser.check(token.ELSE) {
		parser.consume() // consume ELSE
		elseBranch, err = parser.parseSingleStatement()
		if err != nil {
			return statements.If{}, fmt.Errorf("if else: %w", err)
		}
	}

	return statements.If{
		Cond: cond,
		Then: then,
		Else: elseBranch,
	}, nil
}

type WhileStatementParselet struct{}

func (p WhileStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume WHILE
	_, err := parser.consumeAndCheck(token.LEFTPAREN, "expected '(' after while's condition expression")
	if err != nil {
		return nil, err
	}
	cond, err := parser.parseExpr(0)
	if err != nil {
		return statements.If{}, fmt.Errorf("if condition: %w", err)
	}
	_, err = parser.consumeAndCheck(token.RIGHTPAREN, "Expect ')' after while's condition expression")
	if err != nil {
		return nil, err
	}
	body, err := parser.parseSingleStatement()
	if err != nil {
		return statements.While{}, fmt.Errorf("while body: %w", err)
	}

	return statements.While{
		Cond: cond,
		Body: body,
	}, nil
}

type ForStatementParselet struct{}

func (p ForStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume FOR
	_, err := parser.consumeAndCheck(token.LEFTPAREN, "expected '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var init statements.Stmt
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

	var cond expressions.Expr
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

	var incr expressions.Expr
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

	var res statements.Stmt
	if incr != nil {
		res = statements.Block{
			Stmts: []statements.Stmt{
				body,
				statements.Expression{Expr: incr},
			},
		}
	} else {
		res = body
	}

	if cond == nil {
		cond = expressions.Literal{Value: true}
	}
	res = statements.While{Cond: cond, Body: res}

	if init != nil {
		res = statements.Block{
			Stmts: []statements.Stmt{
				init,
				res,
			},
		}
	}

	return res, nil
}

type ExpressionStatementParselet struct{}

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

// TODO
type BreakStatementParselet struct{}

func (p BreakStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume BREAK
	_, err := parser.consumeAndCheck(token.SEMICOLON, "expected ';' after break")
	if err != nil {
		return nil, err
	}
	// check that we are in a loop block
	return statements.Break{}, nil
}

// TODO
type ContinueStatementParslet struct{}

func (p ContinueStatementParslet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume CONTINUE
	_, err := parser.consumeAndCheck(token.SEMICOLON, "expected ';' after continue")
	if err != nil {
		return nil, err
	}
	// check that we are in a loop block
	return statements.Continue{}, nil
}

type FunctionStatementParselet struct{}

func (p FunctionStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	parser.consume() // consume FUN
	name, err := parser.consumeAndCheck(token.IDENTIFIER, "Expect function name.")
	if err != nil {
		return nil, err
	}

	_, err = parser.consumeAndCheck(token.LEFTPAREN, "Expect '(' after function name.")
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

	body, ok := b.(statements.Block)
	if !ok {
		return nil, errors.New("expected block statement")
	}

	return statements.Function{
		Name:   name,
		Params: params,
		Body:   body.Stmts,
	}, nil
}

type ReturnStatementParselet struct{}

func (p ReturnStatementParselet) parse(parser *Parser) (statements.Stmt, error) {
	t := parser.consume() // consume RETURN
	var expr expressions.Expr
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
	return statements.Return{Keyword: t, Value: expr}, nil
}
