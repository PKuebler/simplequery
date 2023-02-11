package simplequery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenToString(t *testing.T) {
	t.Parallel()

	tokens := 15
	for i := 0; i < tokens; i++ {
		token := Token(i)
		assert.Greater(t, len(token.String()), 0)
	}
}

func TestNewLexer(t *testing.T) {
	t.Parallel()

	l := NewLexer("abc")
	assert.NotNil(t, l)
}

func TestLexer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		query  string
		tokens []Token
		texts  []string
	}{
		{
			query:  "q:variableName",
			tokens: []Token{IDENT, EOF},
			texts:  []string{"q:variableName", ""},
		},
		{
			query:  "variableName",
			tokens: []Token{IDENT, EOF},
			texts:  []string{"variableName", ""},
		},
		{
			query:  "variableName AND test",
			tokens: []Token{IDENT, AND, IDENT, EOF},
			texts:  []string{"variableName", "AND", "test", ""},
		},
		{
			query:  "variableName OR test",
			tokens: []Token{IDENT, OR, IDENT, EOF},
			texts:  []string{"variableName", "OR", "test", ""},
		},
		{
			query:  "variableName=b",
			tokens: []Token{IDENT, EQ, IDENT, EOF},
			texts:  []string{"variableName", "=", "b", ""},
		},
		{
			query:  "variableName>b",
			tokens: []Token{IDENT, GT, IDENT, EOF},
			texts:  []string{"variableName", ">", "b", ""},
		},
		{
			query:  "variableName<b",
			tokens: []Token{IDENT, LT, IDENT, EOF},
			texts:  []string{"variableName", "<", "b", ""},
		},
		{
			query:  "variableName>=b",
			tokens: []Token{IDENT, GTE, IDENT, EOF},
			texts:  []string{"variableName", ">=", "b", ""},
		},
		{
			query:  "variableName<=b",
			tokens: []Token{IDENT, LTE, IDENT, EOF},
			texts:  []string{"variableName", "<=", "b", ""},
		},
		{
			query:  "variableName<1234",
			tokens: []Token{IDENT, LT, NUMBER, EOF},
			texts:  []string{"variableName", "<", "1234", ""},
		},
		{
			query:  "variableName<12.34",
			tokens: []Token{IDENT, LT, NUMBER, EOF},
			texts:  []string{"variableName", "<", "12.34", ""},
		},
		{
			query:  "variableName<12,34",
			tokens: []Token{IDENT, LT, NUMBER, EOF},
			texts:  []string{"variableName", "<", "12,34", ""},
		},
		{
			query:  "a=b g>123",
			tokens: []Token{IDENT, EQ, IDENT, IDENT, GT, NUMBER, EOF},
			texts:  []string{"a", "=", "b", "g", ">", "123", ""},
		},
		{
			query:  "a=b g>123 AND (bla OR f != b) OR g",
			tokens: []Token{IDENT, EQ, IDENT, IDENT, GT, NUMBER, AND, BRACKET_LEFT, IDENT, OR, IDENT, NE, IDENT, BRACKET_RIGHT, OR, IDENT, EOF},
			texts:  []string{"a", "=", "b", "g", ">", "123", "AND", "(", "bla", "OR", "f", "!=", "b", ")", "OR", "g", ""},
		},
		{
			query:  "vari.able,Na(m)e<.1234",
			tokens: []Token{IDENT, ILLEGAL, IDENT, ILLEGAL, IDENT, BRACKET_LEFT, IDENT, BRACKET_RIGHT, IDENT, LT, ILLEGAL, NUMBER, EOF},
			texts:  []string{"vari", ".", "able", ",", "Na", "(", "m", ")", "e", "<", ".", "1234", ""},
		},
	}

	for _, testCase := range testCases {
		tokens := []Token{}
		texts := []string{}

		lexer := NewLexer(testCase.query)
		for {
			_, tok, text := lexer.Lex()

			tokens = append(tokens, tok)
			texts = append(texts, text)

			if tok == EOF {
				break
			}
		}

		assert.Len(t, tokens, len(testCase.tokens), testCase.query)
		assert.EqualValues(t, testCase.tokens, tokens, testCase.query)
		assert.Len(t, texts, len(testCase.texts), testCase.query)
		assert.EqualValues(t, testCase.texts, texts, testCase.query)
	}
}
