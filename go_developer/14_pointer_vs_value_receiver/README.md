# 14. Pointer vs Value Receiver

## 核心ポイント
**`*T` は T と `*T` の全メソッドを持つ。interface には `*T` を渡すのが安全。**

## 挙動の違い

| | Value Receiver `(t T)` | Pointer Receiver `(t *T)` |
|--|----------------------|------------------------|
| 操作対象 | コピー | 元の値 |
| 元を変更 | できない | できる |
| nil で呼べるか | No | Yes |

## メソッドセットの違い

```
型 T のメソッドセット:   T のメソッドのみ
型 *T のメソッドセット:  T + *T のメソッド
```

### interface 満足の罠

```go
type Stringer interface { String() string }

type MyType struct{ name string }
func (m *MyType) String() string { return m.name }  // pointer receiver

var s Stringer = MyType{}   // コンパイルエラー! T は *T のメソッドを持たない
var s Stringer = &MyType{}  // OK: *T は *T のメソッドを持つ
```

## 使い分けガイドライン

| 条件 | 推奨 |
|------|------|
| 構造体を変更する | Pointer |
| 構造体が大きい（> 64 bytes 程度） | Pointer |
| sync.Mutex 等を含む | Pointer（コピー禁止） |
| 小さくて変更不要 | Value |
| map のキーにしたい | Value |
| 一貫性（混在させない） | Pointer |

## 実行方法

```bash
go run ./14_pointer_vs_value_receiver/
```
