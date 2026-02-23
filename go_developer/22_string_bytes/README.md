# 22. String / []byte

## 核心ポイント
**string ↔ []byte 変換はコピーが発生する。strings.Builder で最適化。**

## string の内部

```go
type StringHeader struct {
    Data uintptr  // バイト列へのポインタ
    Len  int      // バイト長
}
```

- string は **不変（immutable）**
- UTF-8 エンコーディングされたバイト列

## 変換のコスト

```go
s := "hello"
b := []byte(s)   // コピー発生! O(n)
s2 := string(b)  // コピー発生! O(n)
```

## strings.Builder

```go
var b strings.Builder
b.WriteString("hello")
b.WriteString(" world")
result := b.String()  // 内部バッファから直接 string を作る（コピー最小）
```

## rune vs byte

| | byte | rune |
|--|------|------|
| 型 | uint8 | int32 |
| 範囲 | 0-255 | Unicode コードポイント |
| `for i, b := range s` | - | rune 単位 |
| `s[i]` | byte 単位 | - |

## 実行方法

```bash
go run ./22_string_bytes/
```
