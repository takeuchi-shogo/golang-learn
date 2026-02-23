# 36. Embedding

## 核心ポイント
**埋め込みは継承ではない。promotion と ambiguity に注意。**

## 構造体の埋め込み

```go
type Base struct { ID int }
func (b Base) GetID() int { return b.ID }

type User struct {
    Base        // 埋め込み（has-a 関係）
    Name string
}

u := User{Base: Base{ID: 1}, Name: "Alice"}
u.GetID()  // Base.GetID() が promotion される
u.ID       // Base.ID も直接アクセス可能
```

## Promotion（昇格）

埋め込まれた型のメソッドとフィールドが、外側の型で直接使える。

## Ambiguity（曖昧性）

複数の埋め込みが同名メソッドを持つ場合:
```go
type A struct{}
func (A) Hello() string { return "A" }

type B struct{}
func (B) Hello() string { return "B" }

type C struct { A; B }
// c.Hello() → コンパイルエラー（曖昧）
// c.A.Hello() → 明示的に指定すれば OK
```

## Interface の埋め込み

```go
type ReadWriter interface {
    io.Reader  // 埋め込み
    io.Writer  // 埋め込み
}
```

## 実行方法

```bash
go run ./36_embedding/
```
