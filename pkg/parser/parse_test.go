package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taehioum/glox/pkg/expressions"
	"github.com/taehioum/glox/pkg/token"
)

func Test(t *testing.T) {
	testCases := []struct {
		in   []token.Token
		out  expressions.Expr
		desc string
	}{
		{
			in: []token.Token{
				{
					Type:    token.NUMBER,
					Literal: 123,
					Lexeme:  "123",
					Ln:      1,
				},
				{
					Type: token.EOF,
					Ln:   2,
				},
			},
			out: expressions.Literal{
				Value: 123,
			},
			desc: "a number",
		},
		{
			in: []token.Token{
				{
					Type:    token.NUMBER,
					Literal: 123,
					Lexeme:  "123",
					Ln:      1,
				},
				{
					Type:   token.PLUS,
					Lexeme: "+",
					Ln:     1,
				},
				{
					Type:    token.NUMBER,
					Literal: 10,
					Lexeme:  "123",
					Ln:      1,
				},
				{
					Type: token.EOF,
					Ln:   2,
				},
			},
			out: expressions.Binary{
				Left: expressions.Literal{
					Value: 123,
				},
				Operator: token.Token{
					Type:   token.PLUS,
					Lexeme: "+",
					Ln:     1,
				},
				Right: expressions.Literal{
					Value: 10,
				},
			},
			desc: "123 + 10",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			expr, err := Parse(tc.in)
			assert.NoError(t, err)
			assert.Equal(t, tc.out, expr)
		})
	}
}
