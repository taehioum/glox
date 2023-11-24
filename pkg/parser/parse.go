package parser

import (
	"fmt"

	"github.com/taehioum/glox/pkg/expressions"
	"github.com/taehioum/glox/pkg/token"
)

// each method in the parser matches the grammar rules in the language.
type Parser struct {
	tokens []token.Token
	curr   int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		curr:   0,
		tokens: tokens,
	}
}

func Parse(tokens []token.Token) (expressions.Expr, error) {
	parser := NewParser(tokens)
	return parser.Parse()
}

func (p *Parser) Parse() (expressions.Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("parsing expression: %w", err)
	}
	return expr, nil
}

func (p *Parser) expression() (expressions.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (expressions.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANGEQUAL, token.EQUALEQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, nil
		}
		// note that Left is also expr.
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (expressions.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATEREQUAL, token.LESS, token.LESSEQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) term() (expressions.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) factor() (expressions.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (expressions.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expressions.Unary{Operator: operator, Right: right}, nil
	} else {
		return p.primary()
	}
}

func (p *Parser) primary() (expressions.Expr, error) {
	switch {
	case p.match(token.FALSE):
		return expressions.Literal{Value: false}, nil
	case p.match(token.TRUE):
		return expressions.Literal{Value: true}, nil
	case p.match(token.NIL):
		return expressions.Literal{Value: nil}, nil
	case p.match(token.NUMBER), p.match(token.STRING):
		t := p.previous()
		return expressions.Literal{Value: t.Literal}, nil
	case p.match(token.LEFTPAREN):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHTPAREN, "Expect ')' after expression.")
		return expressions.Grouping{Expr: expr}, err
	default:
		return nil, fmt.Errorf("line %d's %s: expect expression", p.peek().Ln, p.peek().Lexeme)
	}
}

func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.curr++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.curr]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.curr-1]
}

func (p *Parser) consume(t token.Type, msg string) (token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return token.Token{}, fmt.Errorf("line %d's %s: %s", p.peek().Ln, p.peek().Lexeme, msg)
}

// unused, yet.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
