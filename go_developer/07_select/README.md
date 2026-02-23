# 07. select

## 核心ポイント
**複数の case が ready なら pseudo-random で選択される。優先順位はない。**

## 基本動作

```go
select {
case v := <-ch1:   // ch1 から受信できたら
case ch2 <- val:   // ch2 に送信できたら
case <-time.After(1*time.Second): // タイムアウト
default:           // どの case も ready でなければ（ノンブロッキング）
}
```

## 重要な挙動

### Pseudo-random 選択
複数の case が同時に ready なとき、Go ランタイムは**ランダムに1つ選ぶ**。

```go
// ch1 と ch2 の両方にデータがあるとき
select {
case <-ch1: // 50% の確率
case <-ch2: // 50% の確率
}
```

→ starvation を防ぐための設計。「上に書いた case が優先」ではない。

### default case
```go
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("データなし") // ch が空ならすぐにここに来る
}
```

### 空の select
```go
select {} // 永遠にブロック（main goroutine を止めたいときなど）
```

## よく使うパターン

1. **for-select + done**: goroutine のキャンセル
2. **select + time.After**: タイムアウト
3. **select + default**: ノンブロッキング操作
4. **select + context.Done()**: context キャンセル

## 実行方法

```bash
go run ./07_select/
```
