# 18. panic / recover

## 核心ポイント
**recover は deferred 関数内でのみ有効。同一 goroutine のみ。**

## panic の仕組み

1. `panic(v)` が呼ばれる
2. 現在の関数の残りのコードはスキップ
3. **deferred 関数**は LIFO で実行される
4. 呼び出し元に遡って同じ処理を繰り返す（スタック巻き戻し）
5. goroutine のトップに達したらプログラムがクラッシュ

## recover のルール

1. **deferred 関数内でのみ有効**: 通常のコードで呼んでも常に `nil`
2. **同一 goroutine のみ**: 別 goroutine の panic は捕捉できない
3. 戻り値: panic に渡された値

```go
func safeCall(f func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    f()
    return nil
}
```

## runtime error

以下も panic として発生する:
- nil pointer dereference
- index out of range
- slice bounds out of range
- division by zero
- send on closed channel

## 実行方法

```bash
go run ./18_panic_recover/
```
