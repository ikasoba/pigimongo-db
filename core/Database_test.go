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

	if err := db.Add(Hoge{"ぺっぽこ", 100}); err != nil {
		t.Fatal(err)
	}

	hoge := &Hoge{}

	err = db.Find(
		hoge,
		`A == "ぺっぽこ"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)

	hoge = &Hoge{}

	err = db.Find(
		hoge,
		`A == ?`,
		"ぺっぽこ",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)

	err = db.Update(
		struct {
			A string
		}{
			A: "にょろん",
		},
		`A == ?`,
		"ぺっぽこ",
	)
	if err != nil {
		t.Fatal(err)
	}

	hoge = &Hoge{}

	err = db.Find(
		hoge,
		`A == ?`,
		"にょろん",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hoge)

	err = db.Remove(
		`A == ?`,
		"にょろん",
	)
	if err != nil {
		t.Fatal(err)
	}

	hoge = &Hoge{}

	err = db.Find(
		hoge,
		`A == ?`,
		"にょろん",
	)
	if err == nil {
		t.Fail()
	}
}
