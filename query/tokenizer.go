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
	p := regexp.MustCompile(`^\s*`)

	loc := p.FindStringIndex(src[i:])
	if loc == nil {
		return i
	}

	i += len(src[i : i+loc[1]])

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

	i += len(src[i : i+loc[1]])

	return i, nil
}

var parseIdent = mapParser(
	createRegexpParser(`^[\pL_$][\pL_$0-9]*`),
	func(x string) (IdentLit, error) {
		return IdentLit(x), nil
	},
)

var parseString Parser[StrLit] = func(i int, src string) (int, StrLit, error) {
	buf := ""

	if i >= len(src) || src[i] != '"' {
		return i, "", errors.New("cannot match quote.")
	}

	i++

	for ; i < len(src); i++ {
		if src[i] == '\\' {
			i++

			if src[i] == 'n' {
				buf += "\n"
				continue
			} else if src[i] == 't' {
				buf += "\t"
				continue
			}

			buf += string(src[i])
			continue
		}

		buf += string(src[i])
	}

	if i >= len(src) || src[i] != '"' {
		return i, "", errors.New("cannot match quote.")
	}

	i++

	return i, StrLit(buf), nil
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

func parseToken(i int, src string) (int, any, error) {
	log.Println("s", i)
	defer log.Println("e", i)

	if i, v, err := parseIdent(i, src); err == nil {
		return i, v, nil
	} else if i, v, err := parseNumber(i, src); err == nil {
		return i, v, nil
	} else if i, v, err := parseString(i, src); err == nil {
		return i, v, nil
	} else if i, v, err := andOperator(i, src); err == nil {
		return i, v, nil
	} else if i, v, err := orOperator(i, src); err == nil {
		return i, v, nil
	} else if i, v, err := eqlOperator(i, src); err == nil {
		return i, v, nil
	} else if i, v, err := neqOperator(i, src); err == nil {
		return i, v, nil
	}
	return i, nil, errors.New("cannot match token.")
}

func Tokenize(src string) ([]any, error) {
	res := []any{}
	i := 0

	i = skipWhiteSpace(i, src)

	for i < len(src) {
		_i, v, err := parseToken(i, src)
		if err != nil {
			return nil, err
		}

		i = _i

		res = append(res, v)

		if i < len(src) {
			log.Println("--c-", string(src[i]))
		}

		i, err = skipWhiteSpace1(i, src)
		if i >= len(src) {
			break
		} else if err != nil {
			break
		}
	}

	return res, nil
}
