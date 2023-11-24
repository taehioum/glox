package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taehioum/glox/pkg/token"
)

func TestScannerValidCode(t *testing.T) {
	testCases := []struct {
		input    string
		expected []token.Token
		desc     string
	}{
		{
			input: `
			`,
			expected: []token.Token{
				{
					Type: token.EOF,
					Ln:   1,
				},
			},
			desc: "empty file",
		},
		{
			input: `
				123
			`,
			expected: []token.Token{
				{
					Type:    token.NUMBER,
					Lexeme:  "123",
					Literal: float64(123),
					Ln:      1,
				},
				{
					Type: token.EOF,
					Ln:   2,
				},
			},
			desc: "a number",
		},
		{
			input: `
				var x=3.3
			`,
			expected: []token.Token{
				{
					Type:   token.VAR,
					Lexeme: "var",
					Ln:     1,
				},
				{
					Type:   token.IDENTIFIER,
					Lexeme: "x",
					Ln:     1,
				},
				{
					Type:   token.EQUAL,
					Lexeme: "=",
					Ln:     1,
				},
				{
					Type:    token.NUMBER,
					Lexeme:  "3.3",
					Literal: float64(3.3),
					Ln:      1,
				},
				{
					Type: token.EOF,
					Ln:   2,
				},
			},
			desc: "var assignment (number)",
		},
		{
			input: `
				var x=3.3
				var y = 4
				print x + y
			`,
			expected: []token.Token{
				{
					Type:   token.VAR,
					Lexeme: "var",
					Ln:     1,
				},
				{
					Type:   token.IDENTIFIER,
					Lexeme: "x",
					Ln:     1,
				},
				{
					Type:   token.EQUAL,
					Lexeme: "=",
					Ln:     1,
				},
				{
					Type:    token.NUMBER,
					Lexeme:  "3.3",
					Literal: float64(3.3),
					Ln:      1,
				},
				{
					Type:   token.VAR,
					Lexeme: "var",
					Ln:     2,
				},
				{
					Type:   token.IDENTIFIER,
					Lexeme: "y",
					Ln:     2,
				},
				{
					Type:   token.EQUAL,
					Lexeme: "=",
					Ln:     2,
				},
				{
					Type:    token.NUMBER,
					Lexeme:  "4",
					Literal: float64(4),
					Ln:      2,
				},
				{
					Type:   token.PRINT,
					Lexeme: "print",
					Ln:     3,
				},
				{
					Type:   token.IDENTIFIER,
					Lexeme: "x",
					Ln:     3,
				},
				{
					Type:   token.PLUS,
					Lexeme: "+",
					Ln:     3,
				},
				{
					Type:   token.IDENTIFIER,
					Lexeme: "y",
					Ln:     3,
				},
				{
					Type: token.EOF,
					Ln:   4,
				},
			},
			desc: "var assignment and addition",
		},
		{
			input: `
				print "hello"
			`,
			expected: []token.Token{
				{
					Type:   token.PRINT,
					Lexeme: "print",
					Ln:     1,
				},
				{
					Type:    token.STRING,
					Lexeme:  "\"hello\"",
					Literal: "hello",
					Ln:      1,
				},
				{
					Type: token.EOF,
					Ln:   2,
				},
			},
			desc: "print string",
		},
		{
			input: `
				if true {
					print "true"
				} else {
					print "false"
				}
			`,
			expected: []token.Token{
				{
					Type:   token.IF,
					Lexeme: "if",
					Ln:     1,
				},
				{
					Type:   token.TRUE,
					Lexeme: "true",
					Ln:     1,
				},
				{
					Type:   token.LEFTBRACE,
					Lexeme: "{",
					Ln:     1,
				},
				{
					Type:   token.PRINT,
					Lexeme: "print",
					Ln:     2,
				},
				{
					Type:    token.STRING,
					Lexeme:  "\"true\"",
					Literal: "true",
					Ln:      2,
				},
				{
					Type:   token.RIGHTBRACE,
					Lexeme: "}",
					Ln:     3,
				},
				{
					Type:   token.ELSE,
					Lexeme: "else",
					Ln:     3,
				},
				{
					Type:   token.LEFTBRACE,
					Lexeme: "{",
					Ln:     3,
				},
				{
					Type:   token.PRINT,
					Lexeme: "print",
					Ln:     4,
				},
				{
					Type:    token.STRING,
					Lexeme:  "\"false\"",
					Literal: "false",
					Ln:      4,
				},
				{
					Type:   token.RIGHTBRACE,
					Lexeme: "}",
					Ln:     5,
				},
				{
					Type: token.EOF,
					Ln:   6,
				},
			},
			desc: "print string",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			out, err := ScanTokens(tc.input)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestScannerInvalidCode(t *testing.T) {
	testCases := []struct {
		input string
		desc  string
	}{
		{
			input: `
				@
			`,
			desc: "not valid char",
		},
		{
			input: `
				"123
			`,
			desc: "unterminated string",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := ScanTokens(tc.input)
			assert.Error(t, err)
		})
	}
}
