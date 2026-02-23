# 06. Channel Internals

## 核心ポイント
**unbuffered = 同期点（rendezvous）、buffered = 非同期リングバッファ。**

## runtime.hchan 構造体

Channel の実体は `runtime.hchan`：

```
hchan {
    qcount   uint   // バッファ内の要素数
    dataqsiz uint   // バッファの容量（cap）
    buf      *array // リングバッファ
    elemsize uint16 // 要素のサイズ
    closed   uint32 // close フラグ
    sendx    uint   // 送信位置（リングバッファのインデックス）
    recvx    uint   // 受信位置
    recvq    waitq  // 受信待ちの goroutine リスト
    sendq    waitq  // 送信待ちの goroutine リスト
    lock     mutex  // 排他制御
}
```

## Unbuffered Channel（cap = 0）

- 送信側は受信側が来るまでブロック
- 受信側は送信側が来るまでブロック
- **同期点**: 送信と受信が同時に行われる（rendezvous）
- Direct send: 待機中の receiver がいれば、バッファを経由せず直接コピー

## Buffered Channel（cap > 0）

- リングバッファとして動作
- バッファに空きがあれば、送信側はブロックしない
- バッファにデータがあれば、受信側はブロックしない
- バッファが満杯なら送信ブロック、空なら受信ブロック

## パフォーマンス特性

| 操作 | コスト |
|------|-------|
| send/recv（非ブロック）| ~50ns |
| send/recv（ブロック→復帰）| ~300ns |
| select（2 case）| ~200ns |

## 実行方法

```bash
go run ./06_channel_internals/
```
