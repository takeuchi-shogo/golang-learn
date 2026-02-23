# 11. sync.Once

## 核心ポイント
**再入するとデッドロック。panic しても「完了扱い」。**

## 内部実装

```go
type Once struct {
    done atomic.Uint32  // 高速パス用
    m    Mutex          // 低速パス用
}
```

1. `atomic.Load(&done)` が 1 なら即 return（高速パス）
2. 0 なら `Mutex.Lock()` して再チェック（ダブルチェックロッキング）
3. `f()` を実行して `atomic.Store(&done, 1)`

## 注意点

### panic しても完了扱い
```go
once.Do(func() { panic("oops") }) // panic するが「完了」とマークされる
once.Do(func() { /* 実行されない */ })
```

### 再入でデッドロック
```go
once.Do(func() {
    once.Do(func() { /* ← デッドロック! */ })
})
```

## Go 1.21+ の便利関数

- `sync.OnceFunc(f)`: func() を返す
- `sync.OnceValue(f)`: func() T を返す
- `sync.OnceValues(f)`: func() (T, error) を返す

## 実行方法

```bash
go run ./11_sync_once/
```
