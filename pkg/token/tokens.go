package token

import "fmt"

type Type string

const (
	// Single-character tokens.
	LEFTPAREN  Type = "LEFTPAREN"
	RIGHTPAREN Type = "RIGHTPAREN"
	LEFTBRACE  Type = "LEFTBRACE"
	RIGHTBRACE Type = "RIGHTBRACE"
	COMMA      Type = "COMMA"
	DOT        Type = "DOT"
	MINUS      Type = "MINUS"
	PLUS       Type = "PLUS"
	SEMICOLON  Type = "SEMICOLON"
	SLASH      Type = "SLASH"
	STAR       Type = "STAR"

	// One or two character tokens.
	BANG         Type = "BANG"
	BANGEQUAL    Type = "BANGEQUAL"
	EQUAL        Type = "EQUAL"
	EQUALEQUAL   Type = "EQUALEQUAL"
	GREATER      Type = "GREATER"
	GREATEREQUAL Type = "GREATEREQUAL"
	LESS         Type = "LESS"
	LESSEQUAL    Type = "LESSEQUAL"

	// Literals.
	IDENTIFIER Type = "IDENTIFIER"
	STRING     Type = "STRING"
	NUMBER     Type = "NUMBER"

	AND    Type = "AND"
	CLASS  Type = "CLASS"
	ELSE   Type = "ELSE"
	FALSE  Type = "FALSE"
	FUN    Type = "FUN"
	FOR    Type = "FOR"
	IF     Type = "IF"
	NIL    Type = "NIL"
	OR     Type = "OR"
	PRINT  Type = "PRINT"
	RETURN Type = "RETURN"
	SUPER  Type = "SUPER"
	THIS   Type = "THIS"
	TRUE   Type = "TRUE"
	VAR    Type = "VAR"
	WHILE  Type = "WHILE"

	EOF Type = "EOF"

	// IGNORE is assigned to tokens that are not needed for the interpreter
	// e.g. whitespace, comments...
	IGNORE Type = "IGNORE"
)

type Token struct {
	Type    Type
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
