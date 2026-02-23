# 09. nil Channel

## 核心ポイント
**nil channel への send/recv は永遠にブロック。select で case を無効化するテクニック。**

## nil channel の挙動

| 操作 | nil channel | closed channel |
|------|------------|----------------|
| send | **永久ブロック** | panic |
| recv | **永久ブロック** | ゼロ値 + false |
| close | panic | panic |

## select での case 無効化

nil channel を select に使うと、その case は**評価されない**（永久ブロック = 無視される）。

```go
var ch1 <-chan int = realCh1
var ch2 <-chan int = realCh2

for ch1 != nil || ch2 != nil {
    select {
    case v, ok := <-ch1:
        if !ok { ch1 = nil; continue }  // ch1 を無効化
        process(v)
    case v, ok := <-ch2:
        if !ok { ch2 = nil; continue }
        process(v)
    }
}
```

## ユースケース: Channel Merge

2つの channel を1つにマージし、片方が閉じたら nil に差し替える。

## 実行方法

```bash
go run ./09_nil_channel/
```
