package core

import (
	"testing"
)

type Hoge struct {
	A string
	B int
}

func TestDatabase(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Add(Hoge{"ぽぽやま", 100}); err != nil {
		t.Fatal(err)
	}

	hoge := &Hoge{}

	err = db.Find(
		hoge,
		`A == "ぽぽやま"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)

	hoge = &Hoge{}

	err = db.Find(
		hoge,
		`A == ?`,
		"ぽぽやま",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)

	err = db.Update(
		struct {
			A string
		}{
			A: "にょむ",
		},
		`A == ?`,
		"ぽぽやま",
	)
	if err != nil {
		t.Fatal(err)
	}

	hoge = &Hoge{}

	err = db.Find(
		hoge,
		`A == ?`,
		"にょむ",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)

	err = db.Remove(
		`A == ?`,
		"にょむ",
	)
	if err != nil {
		t.Fatal(err)
	}

	hoge = &Hoge{}

	err = db.Find(
		hoge,
		`A == ?`,
		"にょむ",
	)
	if err == nil {
		t.Fail()
	}
}
