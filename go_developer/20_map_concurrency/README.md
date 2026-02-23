# 20. Map Concurrency

## 核心ポイント
**concurrent map writes は fatal error（recover 不可）。sync.Mutex か sync.Map で保護。**

## なぜ map はスレッドセーフでないか

パフォーマンス上の理由。大部分のユースケースは単一 goroutine からのアクセスであり、
デフォルトでロックを取ると不要なオーバーヘッドが発生する。

## fatal error

```
fatal error: concurrent map writes
fatal error: concurrent map read and map write
```

これらは**recover 不可**。通常の panic と異なりプロセスが即座に終了する。

## 対策

### 1. sync.Mutex

```go
var mu sync.Mutex
mu.Lock()
m["key"] = value
mu.Unlock()
```

### 2. sync.RWMutex（読み取りが多い場合）

```go
var mu sync.RWMutex
mu.RLock()    // 読み取りロック（複数同時可）
_ = m["key"]
mu.RUnlock()
```

### 3. sync.Map（特定のユースケース向け）

最適なユースケース:
- キーが安定（書き込みが少ない）
- 各キーが一度だけ書かれて何度も読まれる
- goroutine 間でキーの重複が少ない

## 実行方法

```bash
go run ./20_map_concurrency/
```
