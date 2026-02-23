# 16. Empty Interface（any）

## 核心ポイント
**`interface{}` に入れるとヒープエスケープの原因になる。**

## 内部表現

`interface{}` (= `any`) は `eface` 構造体:
```
eface {
    _type *_type  // 型情報
    data  *void   // 値へのポインタ
}
```

値を `any` に代入すると:
1. 型情報が `_type` に格納される
2. 値が**ヒープにコピー**されて `data` がそのポインタを持つ

→ **スタックにあった値がヒープにエスケープする**

## type assertion のコスト

```go
var i any = 42
v := i.(int)        // 型チェック + 値の取り出し（~2ns）
v, ok := i.(int)    // パニックしない安全版
```

## type switch

```go
switch v := i.(type) {
case int:     // v は int
case string:  // v は string
default:      // 不明な型
}
```

## any vs generics

| | `any` | Generics `[T any]` |
|--|-------|-------------------|
| 型安全 | 実行時チェック | コンパイル時チェック |
| パフォーマンス | ヒープエスケープ | スタックに留まりうる |
| 使い所 | 本当に任意の型 | 型パラメータで制約 |

## 実行方法

```bash
go run ./16_empty_interface/
```
