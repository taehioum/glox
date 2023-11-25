package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	pprint "github.com/taehioum/glox/pkg/printer"
	"github.com/taehioum/glox/pkg/scanner"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		in   string
		out  string
		desc string
	}{
		{
			in:   "-1",
			out:  "(- 1)",
			desc: "unary minus",
		},
		{
			in:   "!true",
			out:  "(! true)",
			desc: "unary bool",
		},
		{
			in:   "1 + 3",
			out:  "(+ 1 3)",
			desc: "infix plus",
		},
		{
			in:   "1 + 2 * 3",
			out:  "(+ 1 (* 2 3))",
			desc: "plus and multiply: test precedence.",
		},
		{
			in:   "-1 + 3",
			out:  "(+ (- 1) 3)",
			desc: "prefix minus, infix plus: test precedence.",
		},
		{
			in:   "(1 + 3) * 10",
			out:  "(* (group (+ 1 3)) 10)",
			desc: "grouping",
		},
		{
			in:   "!true == false",
			out:  "(== (! true) false)",
			desc: "bool",
		},
		{
			in:   "3 * 10 > 20",
			out:  "(> (* 3 10) 20)",
			desc: "comparison",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens, err := scanner.ScanTokens(tc.in)
			assert.NoError(t, err)

			parser := Parser{
				tokens: tokens,
				curr:   0,
			}
			expr, err := parser.Parse()
			assert.NoError(t, err)

			assert.Equal(t, tc.out, pprint.Print(expr))
		})
	}
}

func TestParserWithErrors(t *testing.T) {
	testCases := []struct {
		in   string
		desc string
	}{
		{
			in:   "-",
			desc: "unary minus; no operand",
		},
		{
			in:   "(1 + 2",
			desc: "no right paren",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tokens, err := scanner.ScanTokens(tc.in)
			assert.NoError(t, err)

			parser := Parser{
				tokens: tokens,
				curr:   0,
			}
			_, err = parser.Parse()
			fmt.Println(err)
			assert.Error(t, err)
		})
	}
}
