package core

import (
	"testing"
)

type Hoge struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestDatabase(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Insert(Hoge{"うんち", 100}); err != nil {
		t.Fatal(err)
	}

	hoge := &Hoge{}

	err = db.FindEquals(
		hoge,
		EqualPair{".a", "うんち"},
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)
}
