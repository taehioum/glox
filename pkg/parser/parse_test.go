package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taehioum/glox/pkg/scanner"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		in   string
		out  string
		desc string
	}{
		{
			in:   "print 1;",
			out:  "(print 1)",
			desc: "unary minus",
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
			stmts, err := parser.Parse()
			assert.NoError(t, err)

			fmt.Printf("%+v\n", stmts)
		})
	}
}
