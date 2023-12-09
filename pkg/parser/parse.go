package parser

import (
	"fmt"

	"github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/token"
)

type Parser struct {
	tokens []token.Token
	curr   int
}

/**
 * One of the two parselet interfaces used by the Pratt parser. An
 * InfixParselet is associated with a token that appears in the middle of the
 * expression it parses. Its parse() method will be called after the left-hand
 * side has been parsed, and it in turn is responsible for parsing everything
 * that comes after the token. This is also used for postfix ast. in
 * which case it simply doesn't consume any more tokens in its parse() call.
 */
type InfixParselet interface {
	parse(parser *Parser, left ast.Expr, token token.Token) (ast.Expr, error)
	precedence() Precedence
}

/**
 * One of the two interfaces used by the Pratt parser. A PrefixParselet is
 * associated with a token that appears at the beginning of an expression. Its
 * parse() method will be called with the consumed leading token, and the
 * parselet is responsible for parsing anything that comes after that token.
 * This interface is also used for single-token ast.like variables, in
 * which case parse() simply doesn't consume any more tokens.
 * @author rnystrom
 *
 */
type PrefixParselet interface {
	parse(parser *Parser, token token.Token) (ast.Expr, error)
}

type StatementParselet interface {
	parse(parser *Parser) (ast.Stmt, error)
}

var statementParselets = map[token.Type]StatementParselet{
	token.LEFTBRACE: BlockStatementParselet{},
	// token.PRINT:     PrintStatmentParselet{},
	token.FUN:      FunctionDeclarationStatementParselet{},
	token.VAR:      DeclarationStatementParselet{},
	token.IF:       IfStatementParselet{},
	token.WHILE:    WhileStatementParselet{},
	token.FOR:      ForStatementParselet{},
	token.BREAK:    BreakStatementParselet{},
	token.CONTINUE: ContinueStatementParslet{},
	token.RETURN:   ReturnStatementParselet{},
}

var prefixPraseletsbyTokenType = map[token.Type]PrefixParselet{
	token.PLUS:       UnaryOperatorParselet{},
	token.MINUS:      UnaryOperatorParselet{},
	token.BANG:       UnaryOperatorParselet{},
	token.NUMBER:     LiteralParselet{},
	token.STRING:     LiteralParselet{},
	token.NIL:        LiteralParselet{},
	token.IDENTIFIER: VariableParselet{},
	token.TRUE:       BoolParselet{},
	token.FALSE:      BoolParselet{},
	token.LEFTPAREN:  GroupParselet{},
	token.FUN:        LambdaParselet{},
}

var infixPraseletsbyTokenType = map[token.Type]InfixParselet{
	token.EQUAL:        AssignmentParselet{},
	token.PLUS:         TermParselet{},
	token.MINUS:        TermParselet{},
	token.STAR:         FactorParselet{},
	token.SLASH:        FactorParselet{},
	token.OR:           OrParselet{},
	token.AND:          AndParselet{},
	token.EQUALEQUAL:   EqualityParselet{},
	token.BANGEQUAL:    EqualityParselet{},
	token.LESS:         ComparsionParselet{},
	token.LESSEQUAL:    ComparsionParselet{},
	token.GREATER:      ComparsionParselet{},
	token.GREATEREQUAL: ComparsionParselet{},
	token.PLUSPLUS:     PostfixParselet{},
	token.MINUSMINUS:   PostfixParselet{},
	token.LEFTPAREN:    CallParselet{},
}

func Parse(tokens []token.Token) ([]ast.Stmt, error) {
	parser := Parser{
		tokens: tokens,
		curr:   0,
	}
	stmts, err := parser.Parse()
	return stmts, err
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var stmts []ast.Stmt
	for !p.isAtEnd() {
		stmt, err := p.parseSingleStatement()
		if err != nil {
			return stmts, err
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *Parser) parseSingleStatement() (ast.Stmt, error) {
	tok := p.peek()
	parselet, ok := statementParselets[tok.Type]
	if !ok { // the default parselet for statments is expression statement
		return ExpressionStatementParselet{}.parse(p)
	}
	return parselet.parse(p)
}

func (p *Parser) parseExpr(precendence Precedence) (ast.Expr, error) {
	tok := p.consume()
	prefix, ok := prefixPraseletsbyTokenType[tok.Type]
	if !ok {
		return nil, fmt.Errorf("line %d's %s: no prefix parselet for token type %s", tok.Ln, tok.Lexeme, tok.Type)
	}

	left, err := prefix.parse(p, tok)
	if err != nil {
		return left, err
	}

	for precendence < p.precendence() {
		tok := p.consume()

		infix, ok := infixPraseletsbyTokenType[tok.Type]
		if !ok {
			return nil, fmt.Errorf("line %d's %s: no infix parselet for token type %s", tok.Ln, tok.Lexeme, tok.Type)
		}

		left, err = infix.parse(p, left, tok)
		if err != nil {
			return left, err
		}
	}

	return left, nil
}

func (p *Parser) precendence() Precedence {
	infix, ok := infixPraseletsbyTokenType[p.peek().Type]
	if !ok {
		return 0
	}

	return infix.precedence()
}

// lookahead of distance zero.
func (p *Parser) peek() token.Token {
	return p.tokens[p.curr]
}

func (p *Parser) consume() token.Token {
	tok := p.tokens[p.curr]
	p.curr++
	return tok
}

func (p *Parser) consumeAndCheck(t token.Type, msg string) (token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return token.Token{}, fmt.Errorf("line %d's %s: %s", p.peek().Ln, p.peek().Lexeme, msg)
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.curr++
	}
	return p.previous()
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.curr-1]
}
