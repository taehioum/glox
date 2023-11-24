package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScannerValidCode(t *testing.T) {
	testCases := []struct {
		input    string
		expected []Token
		desc     string
	}{
		{
			input: `
			`,
			expected: []Token{
				{
					TokenType: EOF,
					Ln:        1,
				},
			},
			desc: "empty file",
		},
		{
			input: `
				123
			`,
			expected: []Token{
				{
					TokenType: NUMBER,
					Lexeme:    "123",
					Literal:   float64(123),
					Ln:        1,
				},
				{
					TokenType: EOF,
					Ln:        2,
				},
			},
			desc: "a number",
		},
		{
			input: `
				var x=3.3
			`,
			expected: []Token{
				{
					TokenType: VAR,
					Lexeme:    "var",
					Ln:        1,
				},
				{
					TokenType: IDENTIFIER,
					Lexeme:    "x",
					Ln:        1,
				},
				{
					TokenType: EQUAL,
					Lexeme:    "=",
					Ln:        1,
				},
				{
					TokenType: NUMBER,
					Lexeme:    "3.3",
					Literal:   float64(3.3),
					Ln:        1,
				},
				{
					TokenType: EOF,
					Ln:        2,
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
			expected: []Token{
				{
					TokenType: VAR,
					Lexeme:    "var",
					Ln:        1,
				},
				{
					TokenType: IDENTIFIER,
					Lexeme:    "x",
					Ln:        1,
				},
				{
					TokenType: EQUAL,
					Lexeme:    "=",
					Ln:        1,
				},
				{
					TokenType: NUMBER,
					Lexeme:    "3.3",
					Literal:   float64(3.3),
					Ln:        1,
				},
				{
					TokenType: VAR,
					Lexeme:    "var",
					Ln:        2,
				},
				{
					TokenType: IDENTIFIER,
					Lexeme:    "y",
					Ln:        2,
				},
				{
					TokenType: EQUAL,
					Lexeme:    "=",
					Ln:        2,
				},
				{
					TokenType: NUMBER,
					Lexeme:    "4",
					Literal:   float64(4),
					Ln:        2,
				},
				{
					TokenType: PRINT,
					Lexeme:    "print",
					Ln:        3,
				},
				{
					TokenType: IDENTIFIER,
					Lexeme:    "x",
					Ln:        3,
				},
				{
					TokenType: PLUS,
					Lexeme:    "+",
					Ln:        3,
				},
				{
					TokenType: IDENTIFIER,
					Lexeme:    "y",
					Ln:        3,
				},
				{
					TokenType: EOF,
					Ln:        4,
				},
			},
			desc: "var assignment and addition",
		},
		{
			input: `
				print "hello"
			`,
			expected: []Token{
				{
					TokenType: PRINT,
					Lexeme:    "print",
					Ln:        1,
				},
				{
					TokenType: STRING,
					Lexeme:    "\"hello\"",
					Literal:   "hello",
					Ln:        1,
				},
				{
					TokenType: EOF,
					Ln:        2,
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
			expected: []Token{
				{
					TokenType: IF,
					Lexeme:    "if",
					Ln:        1,
				},
				{
					TokenType: TRUE,
					Lexeme:    "true",
					Ln:        1,
				},
				{
					TokenType: LEFTBRACE,
					Lexeme:    "{",
					Ln:        1,
				},
				{
					TokenType: PRINT,
					Lexeme:    "print",
					Ln:        2,
				},
				{
					TokenType: STRING,
					Lexeme:    "\"true\"",
					Literal:   "true",
					Ln:        2,
				},
				{
					TokenType: RIGHTBRACE,
					Lexeme:    "}",
					Ln:        3,
				},
				{
					TokenType: ELSE,
					Lexeme:    "else",
					Ln:        3,
				},
				{
					TokenType: LEFTBRACE,
					Lexeme:    "{",
					Ln:        3,
				},
				{
					TokenType: PRINT,
					Lexeme:    "print",
					Ln:        4,
				},
				{
					TokenType: STRING,
					Lexeme:    "\"false\"",
					Literal:   "false",
					Ln:        4,
				},
				{
					TokenType: RIGHTBRACE,
					Lexeme:    "}",
					Ln:        5,
				},
				{
					TokenType: EOF,
					Ln:        6,
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
