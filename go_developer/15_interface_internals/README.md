# 15. Interface Internals

## 核心ポイント
**nil pointer ≠ nil interface。typed nil の罠に注意。**

## interface の内部表現

### iface（メソッドあり interface）
```
iface {
    tab  *itab  // 型情報 + メソッドテーブル
    data *void  // 実際の値へのポインタ
}
```

### eface（空 interface = `any`）
```
eface {
    _type *_type  // 型情報
    data  *void   // 実際の値へのポインタ
}
```

### itab の内部
```
itab {
    inter *interfacetype  // interface の型情報
    _type *_type          // 具体型の型情報
    hash  uint32          // 型のハッシュ（高速 type switch 用）
    fun   [1]uintptr      // メソッドテーブル（可変長）
}
```

## typed nil の罠

```go
var p *MyType = nil          // nil ポインタ
var i interface{} = p        // i は nil ではない!

fmt.Println(i == nil)        // false ← 危険！
fmt.Println(p == nil)        // true
```

interface 変数は `(型情報, 値)` のペア。
`p` を代入すると `(*MyType, nil)` になる。型情報が入っているので `nil` ではない。

**真の nil interface** = `(nil, nil)` = 型情報も値も nil。

## 実行方法

```bash
go run ./15_interface_internals/
```
