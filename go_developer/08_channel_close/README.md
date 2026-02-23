# 08. Channel Close

## 核心ポイント
**送信側が close する。受信側は ok で判断する。**

## close のルール

| 操作 | 結果 |
|------|------|
| close 済み channel に send | **panic** |
| close 済み channel を再度 close | **panic** |
| close 済み channel から recv | ゼロ値 + `false` |
| nil channel を close | **panic** |

## 誰が close すべきか

**送信側のみ** が close する。受信側は close しない。

理由: close 済み channel に send すると panic するため、
送信側が close を制御しないと安全性が保証できない。

## パターン

### ok idiom
```go
v, ok := <-ch
if !ok {
    // channel は close されている
}
```

### range over channel
```go
for v := range ch {
    // ch が close されるまでループ
}
```

### Broadcast（fan-out）
close は全ての受信者に同時にシグナルを送れる。
```go
close(done) // 全 goroutine の <-done が同時に解除される
```

## 実行方法

```bash
go run ./08_channel_close/
```
