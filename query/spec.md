# 仕様のメモ

```go
err := db.Find(out, "name == ? and is_verified == true", "...")

err := db.Remove("id == ?", "...")
```

リテラルは
- 文字列
- 数値
- 論理値 (まだ実装できてない)
