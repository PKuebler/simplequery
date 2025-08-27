package simplequery

import (
	"unicode"
	"unicode/utf8"
)

// Token type
type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	NUMBER

	// Infix ops
	EQ  // =
	GT  // >
	GTE // >=
	LT  // <
	LTE // <=
	N   // !
	NE  // !=

	BRACKET_LEFT  // (
	BRACKET_RIGHT // )

	OR  // or
	AND // and
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	NUMBER:  "NUMBER",

	// Infix ops
	EQ:  "=",
	GT:  ">",
	GTE: ">=",
	LT:  "<",
	LTE: "<=",
	N:   "!",
	NE:  "!=",

	BRACKET_LEFT:  "(",
	BRACKET_RIGHT: ")",

	AND: "AND",
	OR:  "OR",
}

// String name of a token
func (t Token) String() string {
	return tokens[t]
}

// Lexer breaks down the input as tokens
type Lexer struct {
	input string
	pos   int
}

// NewLexer create a lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
	}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return EOF
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	return r
}

func (l *Lexer) backup() {
	if l.pos > 0 {
		_, w := utf8.DecodeRuneInString(l.input[:l.pos])
		l.pos -= w
	}
}

// Lex returns the next token, the position and the content.
func (l *Lexer) Lex() (position int, token Token, text string) {
	for {
		switch r := l.next(); {
		case r == EOF:
			return l.pos, EOF, ""
		case unicode.IsSpace(r):
			continue
		case r == '(':
			return l.pos, BRACKET_LEFT, "("
		case r == ')':
			return l.pos, BRACKET_RIGHT, ")"
		case r == '=':
			return l.pos, EQ, "="
		case r == '>':
			startPos := l.pos
			if l.next() == '=' {
				return startPos, GTE, ">="
			} else {
				l.backup()
			}
			return l.pos, GT, ">"
		case r == '<':
			startPos := l.pos
			if l.next() == '=' {
				return startPos, LTE, "<="
			} else {
				l.backup()
			}
			return l.pos, LT, "<"
		case r == '!':
			startPos := l.pos
			if l.next() == '=' {
				return startPos, NE, "!="
			} else {
				l.backup()
			}
			return startPos, N, "!"
		case unicode.IsDigit(r):
			startPos := l.pos
			l.backup()
			lit := l.lexNumber()
			return startPos, NUMBER, lit
		case unicode.IsLetter(r):
			startPos := l.pos
			switch r {
			case 'A', 'a':
				nextN := l.next()
				if nextN == 'N' || nextN == 'n' {
					nextD := l.next()
					if nextD == 'D' || nextD == 'd' {
						return startPos, AND, "AND"
					} else {
						l.backup()
						l.backup()
					}
				} else {
					l.backup()
				}
			case 'O', 'o':
				nextR := l.next()
				if nextR == 'R' || nextR == 'r' {
					return startPos, OR, "OR"
				} else {
					l.backup()
				}
			}

			l.backup()
			lit := l.lexIdent()
			return startPos, IDENT, lit
		default:
			return l.pos, ILLEGAL, string(r)
		}
	}
}

func (l *Lexer) lexNumber() string {
	var lit string
	for {
		switch r := l.next(); {
		case r == EOF:
			return lit
		case unicode.IsDigit(r):
			lit = lit + string(r)
		case r == '.':
			lit = lit + string(r)
		case r == ',':
			lit = lit + string(r)
		default:
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		switch r := l.next(); {
		case r == EOF:
			return lit
		case unicode.IsLetter(r) || r == ':':
			lit = lit + string(r)
		default:
			l.backup()
			return lit
		}
	}
}
