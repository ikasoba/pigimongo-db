# ピギモンゴDB

仲間内でネタにするために作ったドキュメント指向のDB

# 使い方

## データベース初期化
```go
import (
  pigimongo "github.com/ikasoba/pigimongo-db/core"
)

db, err := pigimongo.NewDatabase(":memory:")
if err != nil {
	panic(err)
}
```

`NewDatabase`の第一引数は`":memory:"`または、任意のファイル名が使用できます。

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
err := db.Find(hoge, `Content = "ふが"`)
if err != nil {
  panic(err)
}

// プレースホルダーも使えます
hoge = &Hoge{}
err = db.Find(hoge, `Content == ?`, "ふが")
if err != nil {
  panic(err)
}
```

## 更新
```go
type Hoge struct {
  Content string
}

err := db.Update(Hoge{ "ぴよ" }, `Content == "ふが"`)
if err != nil {
  panic(err)
}

// プレースホルダーも使えます
err = db.Update(Hoge{ "ぴよ" }, `Content == ?`, "ふが")
if err != nil {
  panic(err)
}
```

## 削除
```go
err := db.Remove(`Content == "ぴよ"`)
if err != nil {
  panic(err)
}

// プレースホルダーも使えます
err = db.Remove(`Content == ?`, "ふが")
if err != nil {
  panic(err)
}
```

# クエリ

## クエリで書ける値
- 文字列 `"hoge"` `"hoge\nfuga"`
- 数値 `1234` `12.34`

## クエリで使える演算子
- `and`
- `or`
- `==`
- `!=`
- `<`
- `>`
- `<=`
- `>=`

# テクニック

## ドキュメントには一意のIDが付けられる
```go
type Hoge struct {
  // 内部で自動的に `Id_` というキー名でxidで生成されたidが付与される
  Id_ string
  Content string
}

err := db.Add(Hoge{ Content: "ぴよ" })
if err != nil {
  panic(err)
}

hoge := &Hoge{}
err = db.Find(hoge, `Content = "ぴよ"`)
if err != nil {
  panic(err)
}

log.Println(hoge.Id_) // ckt3oe822smmhr7c40eg
```
