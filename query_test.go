package simplequery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	t.Parallel()

	ok, details, err := Match("(existingKey AND foo) OR (abc AND def)", map[string]string{})
	assert.False(t, ok)
	assert.Len(t, details, 6)
	assert.NoError(t, err)
}

func TestQuery(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc    string
		query   string
		data    map[string]string
		ok      bool
		details []bool
	}{
		{
			query:   "existingKey",
			data:    map[string]string{"existingKey": "value"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "!existingKey",
			data:    map[string]string{"nonExistingKey": "value"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "existingKey",
			data:    map[string]string{"nonExistingKey": "value"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey AND foo",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, true},
		},
		{
			query:   "existingKey AND !foo",
			data:    map[string]string{"existingKey": "value"},
			ok:      true,
			details: []bool{true, true},
		},
		{
			query:   "existingKey AND foo",
			data:    map[string]string{"existingKey": "value"},
			ok:      false,
			details: []bool{true, false},
		},
		{
			query:   "existingKey OR foo",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, true},
		},
		{
			query:   "existingKey OR foo",
			data:    map[string]string{"existingKey": "value"},
			ok:      true,
			details: []bool{true, false},
		},
		{
			query:   "existingKey OR !foo",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, false},
		},
		{
			query:   "!existingKey OR !foo",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      false,
			details: []bool{false, false},
		},
		{
			query:   "existingKey=abc",
			data:    map[string]string{"existingKey": "abc"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "existingKey=abc",
			data:    map[string]string{"nonExistingKey": "value"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey=value AND foo=abc",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, true},
		},
		{
			query:   "existingKey!=value AND foo=abc",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      false,
			details: []bool{false, true},
		},
		{
			query:   "existingKey!=value AND foo=abc",
			data:    map[string]string{"existingKey": "foo", "foo": "abc"},
			ok:      true,
			details: []bool{true, true},
		},
		{
			query:   "existingKey=1234",
			data:    map[string]string{"existingKey": "1234"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "existingKey<1234",
			data:    map[string]string{"nonExistingKey": "1234"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey<12345",
			data:    map[string]string{"existingKey": "12345"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey<12346",
			data:    map[string]string{"existingKey": "12345"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "existingKey<=1234",
			data:    map[string]string{"nonExistingKey": "1234"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey<=12345",
			data:    map[string]string{"existingKey": "12346"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey<=12345",
			data:    map[string]string{"existingKey": "12345"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "existingKey>1234",
			data:    map[string]string{"nonExistingKey": "1234"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey>12345",
			data:    map[string]string{"existingKey": "12345"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey>12346",
			data:    map[string]string{"existingKey": "12347"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "existingKey>=1234",
			data:    map[string]string{"nonExistingKey": "1234"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey>=12346",
			data:    map[string]string{"existingKey": "12345"},
			ok:      false,
			details: []bool{false},
		},
		{
			query:   "existingKey>=12347",
			data:    map[string]string{"existingKey": "12347"},
			ok:      true,
			details: []bool{true},
		},
		{
			query:   "sdfs OR foo OR sdf OR sdfsdf OR dsfsdf",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{false, true, false, false, false},
		},
		{
			query:   "(existingKey AND foo)",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, true, true},
		},
		{
			query:   "(existingKey OR foo)",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, true, true},
		},
		{
			query:   "(!existingKey OR !foo)",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      false,
			details: []bool{false, false, false},
		},
		{
			query:   "(existingKey AND foo) OR (abc AND def)",
			data:    map[string]string{"existingKey": "value", "foo": "abc"},
			ok:      true,
			details: []bool{true, true, true, false, false, false},
		},
	}

	for i, testCase := range testCases {
		lex := NewLexer(testCase.query)
		ok, details, err := query(lex, testCase.data)

		testCase.desc = fmt.Sprintf("%d: %s (%s)", i, testCase.query, testCase.desc)
		assert.Equal(t, testCase.ok, ok, testCase.desc)
		assert.Equal(t, testCase.details, details, testCase.desc)
		assert.NoError(t, err, testCase.desc)
	}
}

func TestProcessPair(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		isPositive bool
		key        string
		operator   Token
		value      string
		data       map[string]string
		desc       string
		ok         bool
		isError    bool
	}{
		{
			isPositive: true,
			key:        "abc",
			operator:   ILLEGAL,
			value:      "",
			data:       map[string]string{"abc": ""},
			desc:       "only operator, positive, found",
			ok:         true,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   ILLEGAL,
			value:      "",
			data:       map[string]string{"foo": ""},
			desc:       "only operator, positive, not found",
			ok:         false,
			isError:    false,
		},
		{
			isPositive: false,
			key:        "abc",
			operator:   ILLEGAL,
			value:      "",
			data:       map[string]string{"abc": ""},
			desc:       "only operator, not positive, found",
			ok:         false,
			isError:    false,
		},
		{
			isPositive: false,
			key:        "abc",
			operator:   ILLEGAL,
			value:      "",
			data:       map[string]string{"foo": ""},
			desc:       "only operator, not positive, not found",
			ok:         true,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   EQ,
			value:      "",
			data:       map[string]string{"foo": ""},
			desc:       "positive, not found",
			ok:         false,
			isError:    false,
		},
		{
			isPositive: false,
			key:        "abc",
			operator:   EQ,
			value:      "",
			data:       map[string]string{"foo": ""},
			desc:       "not positive, not found",
			ok:         true,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   EQ,
			value:      "foo",
			data:       map[string]string{"abc": "foo"},
			desc:       "EQ, positive, found",
			ok:         true,
			isError:    false,
		},
		{
			isPositive: false,
			key:        "abc",
			operator:   EQ,
			value:      "foo",
			data:       map[string]string{"abc": "foo"},
			desc:       "EQ, not positive, found",
			ok:         false,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   EQ,
			value:      "foo",
			data:       map[string]string{"foo": "foo"},
			desc:       "EQ, positive, not found",
			ok:         false,
			isError:    false,
		},
		{
			isPositive: false,
			key:        "abc",
			operator:   EQ,
			value:      "foo",
			data:       map[string]string{"foo": "foo"},
			desc:       "EQ, not positive, not found",
			ok:         true,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   NE,
			value:      "foo",
			data:       map[string]string{"abc": "foo"},
			desc:       "NE, positive, found, match",
			ok:         false,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   NE,
			value:      "foo",
			data:       map[string]string{"abc": "abc"},
			desc:       "NE, positive, found, no match",
			ok:         true,
			isError:    false,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   GT,
			value:      "foo",
			data:       map[string]string{"abc": "abc"},
			desc:       "GT, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   GT,
			value:      "12.34",
			data:       map[string]string{"abc": "abc"},
			desc:       "GT, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   GTE,
			value:      "foo",
			data:       map[string]string{"abc": "abc"},
			desc:       "GTE, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   GTE,
			value:      "12.34",
			data:       map[string]string{"abc": "abc"},
			desc:       "GTE, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   LT,
			value:      "foo",
			data:       map[string]string{"abc": "abc"},
			desc:       "LT, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   LT,
			value:      "12.34",
			data:       map[string]string{"abc": "abc"},
			desc:       "LT, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   LTE,
			value:      "foo",
			data:       map[string]string{"abc": "abc"},
			desc:       "LTE, positive, found, no number",
			ok:         false,
			isError:    true,
		},
		{
			isPositive: true,
			key:        "abc",
			operator:   LTE,
			value:      "12.34",
			data:       map[string]string{"abc": "abc"},
			desc:       "LTE, positive, found, no number",
			ok:         false,
			isError:    true,
		},
	}

	for _, testCase := range testCases {
		ok, err := processPair(testCase.isPositive, testCase.key, testCase.operator, testCase.value, testCase.data)
		assert.Equal(t, testCase.ok, ok, testCase.desc)
		if testCase.isError {
			assert.Error(t, err, testCase.desc)
		} else {
			assert.NoError(t, err, testCase.desc)
		}
	}
}
