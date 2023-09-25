package query

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tokens, err := Tokenize(`$h0G_え`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tokens, []any{IdentLit("$h0G_え")}) {
		t.Fatal(tokens)
	}

	tokens, err = Tokenize(`    $h0G_え 1234  `)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tokens, []any{IdentLit("$h0G_え"), NumLit(1234)}) {
		t.Fatal(tokens)
	}
}
