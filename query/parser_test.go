package query

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	_, token, err := ParseAny(0, `$h0G_え`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, IdentLit("$h0G_え")) {
		t.Fatal(token)
	}

	_, token, err = ParseAny(0, `hogehoge`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, IdentLit("hogehoge")) {
		t.Fatal(token)
	}

	_, token, err = ParseAny(0, `1234`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, NumLit(1234)) {
		t.Fatal(token)
	}

	_, token, err = ParseAny(0, `"unko"`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, StrLit("unko")) {
		t.Fatal(token)
	}

	_, token, err = ParseAny(0, `"うんち"`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, StrLit("うんち")) {
		t.Fatal(token)
	}

	token, err = ParseQuery(`    $h0G_え == 1234  `)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, BinaryOperator{priority: OperatorPriority[OpEQL], op_type: OpEQL, left: IdentLit("$h0G_え"), right: NumLit(1234)}) {
		t.Fatal(token)
	}

	token, err = ParseQuery(`hoge == "fuga"`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(token, BinaryOperator{priority: OperatorPriority[OpEQL], op_type: OpEQL, left: IdentLit("hoge"), right: StrLit("fuga")}) {
		t.Fatal(token)
	}

	token, err = ParseQuery(`    $h0G_え == 1234 and name == "unko" `)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(
		token,
		BinaryOperator{
			priority: OperatorPriority[OpAND],
			op_type:  OpAND,
			left: BinaryOperator{
				priority: OperatorPriority[OpEQL],
				op_type:  OpEQL,
				left:     IdentLit("$h0G_え"),
				right:    NumLit(1234),
			},
			right: BinaryOperator{
				priority: OperatorPriority[OpEQL],
				op_type:  OpEQL,
				left:     IdentLit("name"),
				right:    StrLit("unko"),
			},
		},
	) {
		t.Fatal(token)
	}
}
