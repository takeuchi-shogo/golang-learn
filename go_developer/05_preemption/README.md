# 05. Preemption（プリエンプション）

## 核心ポイント
**Go 1.14+ は非同期プリエンプション対応。tight CPU loop も他の goroutine を飢餓させない。**

## Go のプリエンプション方式の変遷

### Go 1.13 以前: Cooperative Preemption
- goroutine は**関数呼び出し時**にのみスケジューラチェックポイントを通過
- tight loop（関数呼び出しなし）は他の goroutine を永久にブロック

```go
// Go 1.13 以前で問題になるコード
go func() {
    for { /* 関数呼び出しなし = 永久にCPUを占有 */ }
}()
// ↑ 他の goroutine が動けなくなる！
```

### Go 1.14+: Asynchronous Preemption
- OS シグナル（SIGURG）を使って goroutine を中断
- 約 10ms ごとにシグナルを送信
- tight loop でもプリエンプトされる

## safepoint とは

プリエンプションが実行できる安全な箇所：
- 関数呼び出し（Go 1.13 以前からの方式）
- チャネル操作
- GC 関連の操作
- **シグナル受信時**（Go 1.14+ asynchronous preemption）

## 実務での影響

1. CPU-bound な goroutine があっても、他の goroutine は飢餓しない（1.14+）
2. ただし `runtime.LockOSThread()` すると、その M は占有される
3. CGo 呼び出し中はプリエンプト不可（C コードは Go スケジューラの外）

## 実行方法

```bash
go run ./05_preemption/
```
