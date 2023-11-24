package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taehioum/glox/pkg/expressions"
	"github.com/taehioum/glox/pkg/token"
)

func TestPrint(t *testing.T) {
	testCases := []struct {
		input    expressions.Expr
		expected string
		desc     string
	}{
		{
			input: expressions.Binary{
				Left: expressions.Literal{
					Value: 3,
				},
				Operator: token.Token{
					Type:   token.PLUS,
					Lexeme: "+",
				},
				Right: expressions.Literal{
					Value: 6,
				},
			},
			expected: "(+ 3 6)",
			desc:     "addition",
		},
		{
			input: expressions.Unary{
				Operator: token.Token{
					Type:   token.MINUS,
					Lexeme: "-",
				},
				Right: expressions.Literal{
					Value: 6,
				},
			},
			expected: "(- 6)",
			desc:     "number",
		},
		{
			expected: "(* (- 123) (group 45.67))",
			input: expressions.Binary{
				Left: expressions.Unary{
					Operator: token.Token{
						Type:   token.MINUS,
						Lexeme: "-",
					},
					Right: expressions.Literal{
						Value: 123,
					},
				},
				Operator: token.Token{
					Type:   token.STAR,
					Lexeme: "*",
				},
				Right: expressions.Grouping{
					Expr: expressions.Literal{
						Value: 45.67,
					},
				},
			},
			desc: "number",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, Print(tc.input))
		})
	}
}
