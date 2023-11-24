package token

import "fmt"

type TokenType string

const (
	// Single-character tokens.
	LEFTPAREN  TokenType = "LEFTPAREN"
	RIGHTPAREN TokenType = "RIGHTPAREN"
	LEFTBRACE  TokenType = "LEFTBRACE"
	RIGHTBRACE TokenType = "RIGHTBRACE"
	COMMA      TokenType = "COMMA"
	DOT        TokenType = "DOT"
	MINUS      TokenType = "MINUS"
	PLUS       TokenType = "PLUS"
	SEMICOLON  TokenType = "SEMICOLON"
	SLASH      TokenType = "SLASH"
	STAR       TokenType = "STAR"

	// One or two character tokens.
	BANG         TokenType = "BANG"
	BANGEQUAL    TokenType = "BANGEQUAL"
	EQUAL        TokenType = "EQUAL"
	EQUALEQUAL   TokenType = "EQUALEQUAL"
	GREATER      TokenType = "GREATER"
	GREATEREQUAL TokenType = "GREATEREQUAL"
	LESS         TokenType = "LESS"
	LESSEQUAL    TokenType = "LESSEQUAL"

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	EOF TokenType = "EOF"

	// IGNORE is assigned to tokens that are not needed for the interpreter
	// e.g. whitespace, comments...
	IGNORE TokenType = "IGNORE"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	// Line Number
	Ln int
}

func (t Token) String() string {
	if t.Type == IDENTIFIER {
		return fmt.Sprintf("%s %s", t.Type, t.Lexeme)
	}
	if t.Literal != nil {
		return fmt.Sprintf("%s %v", t.Type, t.Literal)
	}
	return string(t.Type)
}
