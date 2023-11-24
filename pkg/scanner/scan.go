package scanner

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Scanner struct {
	source    string
	currToken Token
	errors    []error

	start int
	curr  int
	line  int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source: source,
		start:  0,
		curr:   0,
		line:   0,
	}
}

func ScanTokens(source string) ([]Token, error) {
	sc := NewScanner(source)

	var tokens []Token
	for sc.Next() {
		tok := sc.Get()
		tokens = append(tokens, tok)
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scanning tokens: %w", err)
	}
	return tokens, nil
}

func (sc *Scanner) Next() bool {
	sc.start = sc.curr
	if sc.curr >= len(sc.source) {
		sc.currToken = Token{TokenType: EOF, Lexeme: sc.lexeme(), Ln: sc.line}
		return false
	}

	c := sc.advance()

	switch c {
	case '(':
		sc.currToken = Token{TokenType: LEFTPAREN, Lexeme: sc.lexeme(), Ln: sc.line}
	case ')':
		sc.currToken = Token{TokenType: RIGHTPAREN, Lexeme: sc.lexeme(), Ln: sc.line}
	case '{':
		sc.currToken = Token{TokenType: LEFTBRACE, Lexeme: sc.lexeme(), Ln: sc.line}
	case '}':
		sc.currToken = Token{TokenType: RIGHTBRACE, Lexeme: sc.lexeme(), Ln: sc.line}
	case ',':
		sc.currToken = Token{TokenType: COMMA, Lexeme: sc.lexeme(), Ln: sc.line}
	case '.':
		sc.currToken = Token{TokenType: DOT, Lexeme: sc.lexeme(), Ln: sc.line}
	case '-':
		sc.currToken = Token{TokenType: MINUS, Lexeme: sc.lexeme(), Ln: sc.line}
	case '+':
		sc.currToken = Token{TokenType: PLUS, Lexeme: sc.lexeme(), Ln: sc.line}
	case '*':
		sc.currToken = Token{TokenType: STAR, Lexeme: sc.lexeme(), Ln: sc.line}
	case ';':
		sc.currToken = Token{TokenType: SEMICOLON, Lexeme: sc.lexeme(), Ln: sc.line}
	case '!':
		if sc.match('=') {
			sc.currToken = Token{TokenType: BANGEQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			sc.currToken = Token{TokenType: BANG, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '=':
		if sc.match('=') {
			sc.currToken = Token{TokenType: EQUALEQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			sc.currToken = Token{TokenType: EQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '<':
		if sc.match('=') {
			sc.currToken = Token{TokenType: LESSEQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			sc.currToken = Token{TokenType: LESS, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '>':
		if sc.match('=') {
			sc.currToken = Token{TokenType: GREATEREQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			sc.currToken = Token{TokenType: GREATER, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '/':
		if sc.match('/') { // a comment string
			for {
				// don't consume the newline, so that the line number is incremented upon newline.
				if sc.peek() == '\n' || sc.atEnd() {
					break
				}
				sc.advance()
			}
			sc.currToken = Token{TokenType: IGNORE, Lexeme: "", Ln: sc.line}
		} else {
			sc.currToken = Token{TokenType: SLASH, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case ' ', '\r', '\t':
		sc.currToken = Token{TokenType: IGNORE, Ln: sc.line}
	case '\n':
		sc.line++
		sc.currToken = Token{TokenType: IGNORE, Ln: sc.line}
	case '"':
		val, err := sc.readString()
		if err != nil {
			sc.errors = append(sc.errors, err)
			sc.currToken = Token{TokenType: IGNORE, Ln: sc.line}
		}
		sc.currToken = Token{TokenType: STRING, Lexeme: sc.lexeme(), Literal: val, Ln: sc.line}
	// numbers
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		val, err := sc.readNumber()
		if err != nil {
			sc.errors = append(sc.errors, err)
			sc.currToken = Token{TokenType: IGNORE, Ln: sc.line}
		}
		sc.currToken = Token{TokenType: NUMBER, Lexeme: sc.lexeme(), Literal: val, Ln: sc.line}
	default:
		if unicode.IsLetter(rune(c)) {
			tok := sc.readIdentifierOrKeyword()
			sc.currToken = Token{TokenType: tok, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			sc.errors = append(sc.errors, fmt.Errorf("unexpected character: %c", c))
			sc.currToken = Token{TokenType: IGNORE, Ln: sc.line}
		}
	}
	return true
}

// readIdentifierOrKeyword consumes the rest of the identifier / keyword by advancing.
func (sc *Scanner) readIdentifierOrKeyword() TokenType {
	for (unicode.IsLetter(rune(sc.peek())) || unicode.IsDigit(rune(sc.peek()))) && !sc.atEnd() {
		sc.advance()
	}

	text := sc.lexeme()
	if keyword, ok := keywords[text]; ok {
		return keyword
	}
	return IDENTIFIER
}

// readString consumes the rest of the string by advancing, and returns its literal value
func (sc *Scanner) readString() (literal string, err error) {
	for sc.peek() != '"' && !sc.atEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}

	if sc.atEnd() {
		return "", fmt.Errorf("unterminated string")
	}

	// the closing "
	sc.advance()

	// trim the surrounding quotes
	value := sc.source[sc.start+1 : sc.curr-1]
	return value, nil
}

// readNumber consumes the rest of the number by advancing, and returns its literal value
func (sc *Scanner) readNumber() (literal float64, err error) {
	for (unicode.IsDigit(rune(sc.peek()))) && !sc.atEnd() {
		sc.advance()
	}

	// look for a fractional part
	if sc.peek() == '.' && unicode.IsDigit(rune(sc.peekNext())) {
		// consume the '.'
		sc.advance()

		for (unicode.IsDigit(rune(sc.peek()))) && !sc.atEnd() {
			sc.advance()
		}
	}

	return strconv.ParseFloat(sc.source[sc.start:sc.curr], 64)
}

func (sc *Scanner) advance() byte {
	sc.curr++
	return sc.source[sc.curr-1]
}

// match is a conditional advance.
// if the next character is not expected, do not advance
func (sc *Scanner) match(expected byte) bool {
	if sc.atEnd() {
		return false
	}
	if sc.source[sc.curr] != expected {
		return false
	}

	sc.curr++
	return true
}

// Get returns a valid token.
// If the current token is IGNORE, it skips it and returns the next valid token.
func (sc *Scanner) Get() Token {
	// skip until we don't have an IGNORE token
	for sc.currToken.TokenType == IGNORE && sc.Next() {
	}
	return sc.currToken
}

func (sc *Scanner) Err() error {
	return errors.Join(sc.errors...)
}

// peek is lookahead once
func (sc *Scanner) peek() byte {
	if sc.atEnd() {
		return '\000'
	}

	return sc.source[sc.curr]
}

// peekNext is lookahead twice.
func (sc *Scanner) peekNext() byte {
	if sc.curr+1 >= len(sc.source) {
		return '\000'
	}
	return sc.source[sc.curr+1]
}

func (sc *Scanner) atEnd() bool {
	return sc.curr >= len(sc.source)
}

func (sc *Scanner) lexeme() string {
	return sc.source[sc.start:sc.curr]
}
