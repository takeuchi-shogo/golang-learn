# 33. Generics

## 核心ポイント
**コンパイル時展開 = interface より速い。Go 1.18+ で利用可能。**

## 基本構文

```go
func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}
```

## 型制約（Constraints）

```go
type Number interface {
    ~int | ~int64 | ~float64  // ~ は underlying type を含む
}

type Comparable interface {
    comparable  // == と != が使える型
}
```

## GC Shape Stenciling + Dictionaries

Go のジェネリクスは **GC shape** が同じ型をまとめて1つの関数にする:
- ポインタ型は全て同じ shape → 1つの関数
- 値型（int, float64 等）は shape が異なる → 別々の関数

→ interface の動的ディスパッチより速い（間接参照がない）

## いつ使うべきか

- コレクション操作（Map, Filter, Reduce）
- データ構造（Stack, Queue, Tree）
- 型安全なユーティリティ

## いつ避けるべきか

- 型パラメータが1つの具象型しか取らない
- interface で十分（メソッドの多態性）
- 過度な抽象化になる場合

## 実行方法

```bash
go run ./33_generics/
```
