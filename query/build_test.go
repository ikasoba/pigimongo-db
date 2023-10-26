package query

import (
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	tree, err := ParseQuery(`hoge == "fuga"`)
	if err != nil {
		t.Fatal(err)
	}

	ctx := NewBuildContext()
	err = ctx.BuildQueryToWhere(tree)
	if err != nil {
		t.Fatal(err)
	}

	if ctx.Query != "data -> ? = ?" {
		t.Fatal(ctx.Query)
	}

	if !reflect.DeepEqual(ctx.Values, []any{"$.hoge", `"fuga"`}) {
		t.Fatal(ctx.Values...)
	}

	ctx = NewBuildContext("fuga")
	err = ctx.BuildQueryToWhere(tree)
	if err != nil {
		t.Fatal(err)
	}

	if ctx.Query != "data -> ? = ?" {
		t.Fatal(ctx.Query)
	}

	if !reflect.DeepEqual(ctx.Values, []any{"$.hoge", `"fuga"`}) {
		t.Fatal(ctx.Values...)
	}
}
