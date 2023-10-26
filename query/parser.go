package query

import (
	"errors"
	"log"
	"regexp"
	"strconv"
)

type Parser[R any] func(int, string) (int, R, error)

func mapParser[_R any, R any](parser Parser[_R], fn func(value _R) (R, error)) Parser[R] {
	return func(i int, s string) (int, R, error) {
		i, v, err := parser(i, s)
		if err != nil {
			var zero R
			return i, zero, err
		}

		r, err := fn(v)

		return i, r, err
	}
}

func createStringParser(text string) Parser[string] {
	return func(i int, s string) (int, string, error) {
		if i+len(text) >= len(s) || s[i:i+len(text)] != text {
			return i, "", errors.New("cannot match text.")
		}

		i += len(text)

		return i, text, nil
	}
}

func createRegexpParser(pattern string) Parser[string] {
	p := regexp.MustCompile(pattern)
	return func(i int, s string) (int, string, error) {
		loc := p.FindStringIndex(s[i:])
		if loc == nil {
			return i, "", errors.New("cannot match regexp pattern.")
		}

		return i + loc[1], s[i : i+loc[1]], nil
	}
}

func skipWhiteSpace(i int, src string) int {
	p := regexp.MustCompile(`^\s+`)

	loc := p.FindStringIndex(src[i:])
	if loc == nil {
		return i
	}

	i += loc[1]

	return i
}

func skipWhiteSpace1(i int, src string) (int, error) {
	if i >= len(src) {
		return i, errors.New("out of length.")
	}

	p := regexp.MustCompile(`^\s+`)

	loc := p.FindStringIndex(src[i:])
	if loc == nil {
		return i, errors.New("cannot match white spaces.")
	}

	i += loc[1]

	return i, nil
}

var parseIdent = mapParser(
	createRegexpParser(`^[\pL_$][\pL_$0-9]*`),
	func(x string) (IdentLit, error) {
		return IdentLit(x), nil
	},
)

var parsePlaceholder = mapParser(
	createRegexpParser(`^\?`),
	func(x string) (PlaceholderLit, error) {
		return PlaceholderLit{}, nil
	},
)

var parseString Parser[StrLit] = func(i int, src string) (int, StrLit, error) {
	buf := []byte{}

	if i >= len(src) || src[i] != '"' {
		return i, "", errors.New("cannot match quote.")
	}

	i++

	for ; i < len(src); i++ {
		if src[i] == '\\' {
			i++

			if src[i] == 'n' {
				buf = append(buf, []byte("\n")...)
				continue
			} else if src[i] == 't' {
				buf = append(buf, []byte("\t")...)
				continue
			}

			buf = append(buf, src[i])
			continue
		} else if src[i] == '"' {
			i++
			return i, StrLit(buf), nil
		}

		buf = append(buf, src[i])
	}

	return i, "", errors.New("cannot match quote.")
}

var parseNumber = mapParser(
	createRegexpParser(`^[-+]?(?:[0-9]+)(?:\.(?:[0-9]+(?:[eE][+-]?[0-9]+)?)?)?`),
	func(x string) (NumLit, error) {
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, err
		}

		return NumLit(f), nil
	},
)

var (
	andOperator = mapParser(createRegexpParser(`(?i)^and`), func(x string) (OperatorType, error) { return OpAND, nil })
	orOperator  = mapParser(createRegexpParser(`(?i)^or`), func(x string) (OperatorType, error) { return OpOR, nil })
	eqlOperator = mapParser(createStringParser("=="), func(x string) (OperatorType, error) { return OpEQL, nil })
	neqOperator = mapParser(createStringParser("!="), func(x string) (OperatorType, error) { return OpNEQ, nil })
)

func ParseValue(i int, src string) (int, any, error) {
	log.Println("v s", i)
	defer log.Println("v e", i)

	if _i, v, err := parseIdent(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := parseNumber(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := parseString(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := parsePlaceholder(i, src); err == nil {
		return _i, v, nil
	}

	return i, nil, errors.New("cannot match any value.")
}

func ParseOperatorSymbol(i int, src string) (int, OperatorType, error) {
	log.Println("s s", i)
	defer log.Println("s e", i)

	if _i, v, err := andOperator(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := orOperator(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := eqlOperator(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := neqOperator(i, src); err == nil {
		return _i, v, nil
	}

	return i, 0, errors.New("cannot match any symbol.")
}

func ParseAny(i int, src string) (int, any, error) {
	if _i, v, err := ParseOperator(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := parseIdent(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := parseNumber(i, src); err == nil {
		return _i, v, nil
	} else if _i, v, err := parseString(i, src); err == nil {
		return _i, v, nil
	}

	return i, nil, errors.New("cannot match any value.")
}

func ParseOperator(_i int, src string) (int, BinaryOperator, error) {
	var prevOperator *BinaryOperator = nil

	_i, x, err := ParseValue(_i, src)
	if err != nil {
		return _i, BinaryOperator{}, err
	}

	for {
		if _i >= len(src) {
			if prevOperator != nil {
				return _i, *prevOperator, nil
			}

			break
		}

		log.Println("o l", _i)
		i := skipWhiteSpace(_i, src)

		i, op, err := ParseOperatorSymbol(i, src)
		if err != nil {
			if prevOperator != nil {
				return i, *prevOperator, nil
			}

			return i, BinaryOperator{}, err
		}

		i = skipWhiteSpace(i, src)

		i, y, err := ParseValue(i, src)
		if err != nil {
			return i, BinaryOperator{}, err
		}

		if prevOperator == nil {
			prevOperator = &BinaryOperator{
				priority: OperatorPriority[op],
				op_type:  op,
				left:     x,
				right:    y,
			}
		} else {
			var priority = OperatorPriority[op]

			if prevOperator.priority < priority {
				prevOperator.right = BinaryOperator{
					priority: priority,
					op_type:  op,
					left:     prevOperator.right,
					right:    y,
				}
			} else {
				prevOperator = &BinaryOperator{
					priority: priority,
					left:     *prevOperator,
					right:    y,
				}
			}
		}

		_i = i
	}

	return _i, BinaryOperator{}, errors.New("cannot match any operator.")
}

func ParseQuery(src string) (any, error) {
	i := skipWhiteSpace(0, src)

	i, tree, err := ParseQueryBody(i, src)
	if err != nil {
		return nil, err
	}

	i = skipWhiteSpace(i, src)

	return tree, nil
}

func ParseQueryBody(i int, src string) (int, any, error) {
	if _i, x, err := ParseOperator(i, src); err == nil {
		return _i, x, err
	}

	return i, nil, errors.New("cannot match any statement.")
}
