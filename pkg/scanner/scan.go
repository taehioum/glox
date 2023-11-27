package scanner

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/taehioum/glox/pkg/token"
)

type Scanner struct {
	source string
	errors []error

	start int
	curr  int
	line  int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source: source,
		start:  0,
		curr:   0,
		line:   1,
	}
}

func ScanTokens(source string) ([]token.Token, error) {
	sc := NewScanner(source)

	var tokens []token.Token
	for tok := sc.Scan(); tok.Type != token.EOF; tok = sc.Scan() {
		tokens = append(tokens, tok)
	}
	tokens = append(tokens, token.Token{Type: token.EOF, Lexeme: sc.lexeme(), Ln: sc.line})

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scanning tokens: %w", err)
	}
	return tokens, nil
}

func (sc *Scanner) Scan() token.Token {
	sc.start = sc.curr
	if sc.curr >= len(sc.source) {
		return token.Token{Type: token.EOF, Lexeme: sc.lexeme(), Ln: sc.line}
	}

	c := sc.advance()

	switch c {
	case '(':
		return token.Token{Type: token.LEFTPAREN, Lexeme: sc.lexeme(), Ln: sc.line}
	case ')':
		return token.Token{Type: token.RIGHTPAREN, Lexeme: sc.lexeme(), Ln: sc.line}
	case '{':
		return token.Token{Type: token.LEFTBRACE, Lexeme: sc.lexeme(), Ln: sc.line}
	case '}':
		return token.Token{Type: token.RIGHTBRACE, Lexeme: sc.lexeme(), Ln: sc.line}
	case ',':
		return token.Token{Type: token.COMMA, Lexeme: sc.lexeme(), Ln: sc.line}
	case '.':
		return token.Token{Type: token.DOT, Lexeme: sc.lexeme(), Ln: sc.line}
	case '-':
		return token.Token{Type: token.MINUS, Lexeme: sc.lexeme(), Ln: sc.line}
	case '+':
		return token.Token{Type: token.PLUS, Lexeme: sc.lexeme(), Ln: sc.line}
	case '*':
		return token.Token{Type: token.STAR, Lexeme: sc.lexeme(), Ln: sc.line}
	case ';':
		return token.Token{Type: token.SEMICOLON, Lexeme: sc.lexeme(), Ln: sc.line}
	case '!':
		if sc.match('=') {
			return token.Token{Type: token.BANGEQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			return token.Token{Type: token.BANG, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '=':
		if sc.match('=') {
			return token.Token{Type: token.EQUALEQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			return token.Token{Type: token.EQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '<':
		if sc.match('=') {
			return token.Token{Type: token.LESSEQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			return token.Token{Type: token.LESS, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case '>':
		if sc.match('=') {
			return token.Token{Type: token.GREATEREQUAL, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			return token.Token{Type: token.GREATER, Lexeme: sc.lexeme(), Ln: sc.line}
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
			return sc.Scan()
		} else {
			return token.Token{Type: token.SLASH, Lexeme: sc.lexeme(), Ln: sc.line}
		}
	case ' ', '\r', '\t':
		return sc.Scan()
	case '\n':
		sc.line++
		return sc.Scan()
	case '"':
		val, err := sc.readString()
		if err != nil {
			sc.errors = append(sc.errors, err)
			return sc.Scan()
		}
		return token.Token{Type: token.STRING, Lexeme: sc.lexeme(), Literal: val, Ln: sc.line}
	// numbers
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		val, err := sc.readNumber()
		if err != nil {
			sc.errors = append(sc.errors, err)
			return sc.Scan()
		}
		return token.Token{Type: token.NUMBER, Lexeme: sc.lexeme(), Literal: val, Ln: sc.line}
	default:
		if unicode.IsLetter(rune(c)) {
			tok := sc.readIdentifierOrKeyword()
			return token.Token{Type: tok, Lexeme: sc.lexeme(), Ln: sc.line}
		} else {
			sc.errors = append(sc.errors, fmt.Errorf("unexpected character: %c", c))
			return sc.Scan()
		}
	}
}

// readIdentifierOrKeyword consumes the rest of the identifier / keyword by advancing.
func (sc *Scanner) readIdentifierOrKeyword() token.Type {
	for (unicode.IsLetter(rune(sc.peek())) || unicode.IsDigit(rune(sc.peek()))) && !sc.atEnd() {
		sc.advance()
	}

	text := sc.lexeme()
	if keyword, ok := keywords[text]; ok {
		return keyword
	}
	return token.IDENTIFIER
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

var keywords = map[string]token.Type{
	"and":      token.AND,
	"class":    token.CLASS,
	"else":     token.ELSE,
	"false":    token.FALSE,
	"for":      token.FOR,
	"fun":      token.FUN,
	"if":       token.IF,
	"nil":      token.NIL,
	"or":       token.OR,
	"print":    token.PRINT,
	"return":   token.RETURN,
	"super":    token.SUPER,
	"this":     token.THIS,
	"true":     token.TRUE,
	"var":      token.VAR,
	"while":    token.WHILE,
	"break":    token.BREAK,
	"continue": token.CONTINUE,
}
