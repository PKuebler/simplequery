package simplequery

import (
	"fmt"
	"strconv"
)

// Match the input to the data. Returns whether it is a successful match,
// an array with the individual results and an error if the input query contains errors.
func Match(input string, data map[string]string) (ok bool, details []bool, err error) {
	lexer := NewLexer(input)

	return query(lexer, data)
}

func query(lexer *Lexer, data map[string]string) (ok bool, details []bool, err error) {
	key := ""
	operator := Token(ILLEGAL)
	value := ""
	isPositive := true

	lastWaiting := false
	lastResult := true

	bracket := false
	bracketResult := true

	ok = true

	for {
		pos, tok, lit := lexer.Lex()

		if tok == ILLEGAL {
			return false, nil, fmt.Errorf("illegal query party %s on %d: %s", tok.String(), pos, lit)
		}

		// stop token or next pair without right part
		if isStopToken(tok) || (key != "" && operator == ILLEGAL && tok == IDENT) || bracket {
			var result bool
			if bracket {
				result = bracketResult
				bracket = false
				bracketResult = true
			} else {
				// process current pair
				result, err = processPair(isPositive, key, operator, value, data)
				if err != nil {
					return false, details, err
				}
			}

			details = append(details, result)

			if tok != OR && lastWaiting {
				ok = (result || lastResult)
				lastWaiting = false
				lastResult = true
			} else if tok == OR && lastWaiting {
				lastWaiting = true
				lastResult = (result || lastResult)
			} else if tok == OR {
				lastWaiting = true
				lastResult = result
			} else if ok && !result {
				ok = result
			}

			key = ""
			operator = Token(ILLEGAL)
			value = ""
			isPositive = true

			if tok == EOF {
				break
			}

			if tok == BRACKET_RIGHT {
				return ok, details, nil
			}

			continue
		}

		// left
		if key == "" {
			switch tok {
			case BRACKET_LEFT:
				subOk, subDetails, subErr := query(lexer, data)
				if subErr != nil {
					return false, nil, subErr
				}
				details = append(details, subDetails...)
				bracket = true
				bracketResult = subOk
			case N:
				isPositive = false
			case IDENT:
				key = lit
			default:
				return false, nil, fmt.Errorf("illegal query party %s on %d: %s", tok.String(), pos, lit)
			}

			continue
		}

		// operator
		if operator == ILLEGAL {
			if isOperator(tok) {
				operator = tok
			} else {
				return false, nil, fmt.Errorf("illegal query party %s on %d: %s", tok.String(), pos, lit)
			}

			continue
		}

		// right
		switch tok {
		case IDENT:
			value = lit
		case NUMBER:
			value = lit
		default:
			return false, nil, fmt.Errorf("illegal query party %s on %d: %s", tok.String(), pos, lit)
		}
	}

	return ok, details, nil
}

func processPair(isPositive bool, key string, operator Token, value string, data map[string]string) (bool, error) {
	dataValue, keyFound := data[key]

	// only key
	if operator == ILLEGAL {
		return keyFound == isPositive, nil
	}

	// key not found
	if !keyFound {
		return keyFound == isPositive, nil
	}

	// operator
	result := false
	switch operator {
	case EQ:
		result = value == dataValue
	case GT:
		expectedFloat, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return false, err
		}
		actualFloat, err := strconv.ParseFloat(dataValue, 32)
		if err != nil {
			return false, err
		}

		result = actualFloat > expectedFloat
	case GTE:
		expectedFloat, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return false, err
		}
		actualFloat, err := strconv.ParseFloat(dataValue, 32)
		if err != nil {
			return false, err
		}

		result = actualFloat >= expectedFloat
	case LT:
		expectedFloat, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return false, err
		}
		actualFloat, err := strconv.ParseFloat(dataValue, 32)
		if err != nil {
			return false, err
		}

		result = actualFloat < expectedFloat
	case LTE:
		expectedFloat, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return false, err
		}
		actualFloat, err := strconv.ParseFloat(dataValue, 32)
		if err != nil {
			return false, err
		}

		result = actualFloat <= expectedFloat
	case NE:
		result = value != dataValue
	}

	return result == isPositive, nil
}

func isStopToken(token Token) bool {
	return token == EOF || token == AND || token == OR || token == BRACKET_RIGHT
}

func isOperator(token Token) bool {
	return token == EQ || token == GT || token == GTE || token == LT || token == LTE || token == NE
}
