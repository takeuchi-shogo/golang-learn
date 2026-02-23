# 10. Context

## 核心ポイント
**必ず `defer cancel()`。キャンセルは親→子に伝播する。**

## context の種類

| 関数 | 用途 |
|------|------|
| `context.Background()` | ルートコンテキスト。main, init, テストで使う |
| `context.TODO()` | まだ決まっていないとき（リファクタリング中） |
| `WithCancel(parent)` | 手動キャンセル |
| `WithTimeout(parent, d)` | 期限付き（相対時間） |
| `WithDeadline(parent, t)` | 期限付き（絶対時間） |
| `WithValue(parent, k, v)` | リクエストスコープの値を伝播 |

## 必ず defer cancel()

```go
ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
defer cancel() // ← 忘れると goroutine リーク
```

`cancel()` を呼ばないと、タイマー用の goroutine がリークする。

## キャンセルの伝播

```
parent (cancel) ──▶ child1 (キャンセルされる)
                ──▶ child2 (キャンセルされる)
                        ──▶ grandchild (キャンセルされる)
```

親をキャンセルすると、全ての子孫が自動的にキャンセルされる。
子をキャンセルしても、親には影響しない。

## WithValue のルール

- **リクエストスコープ**の値にのみ使う（traceID, userID, etc.）
- 関数のパラメータの代わりに使ってはいけない
- key は unexported type にする（衝突防止）

## 実行方法

```bash
go run ./10_context/
```
