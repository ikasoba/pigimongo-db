# ピギモンゴDB

仲間内でネタにするために作ったドキュメント指向のDB

# 使い方

## データベース初期化
```go
db, err := NewDatabase(":memory:")
if err != nil {
	panic(err)
}
```

`NewDatabase`の第一引数は`":memory:"`または、任意のパスが使用できます。

## 書き込み
```go
type Hoge struct {
  Content string
}

err := db.Add(Hoge{ "ふが" })
if err != nil {
  panic(err)
}
```

## 読み取り
```go
type Hoge struct {
  Content string
}

hoge := &Hoge{}
err := db.Find(`Content = "ふが"`)
if err != nil {
  panic(err)
}

// プレースホルダーも使えます
hoge = &Hoge{}
err = db.Find(`Content = ?`, "ふが")
if err != nil {
  panic(err)
}
```

## 更新
```go
type Hoge struct {
  Content string
}

err := db.Update(Hoge{ "ぴよ" }, `Content = "ふが"`)
if err != nil {
  panic(err)
}

// プレースホルダーも使えます
err = db.Update(Hoge{ "ぴよ" }, `Content = ?`, "ふが")
if err != nil {
  panic(err)
}
```

## 削除
```go
err := db.Remove(`Content = "ぴよ"`)
if err != nil {
  panic(err)
}

// プレースホルダーも使えます
err = db.Remove(`Content = ?`, "ふが")
if err != nil {
  panic(err)
}
```