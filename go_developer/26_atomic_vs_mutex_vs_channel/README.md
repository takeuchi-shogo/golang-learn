# 26. atomic vs mutex vs channel

## 核心ポイント
**シンプルなカウンタ = atomic、複数変数の一貫性 = mutex、通信 = channel。**

## 使い分け

| | atomic | mutex | channel |
|--|--------|-------|---------|
| 用途 | 単一変数の CAS | 複数変数の排他制御 | goroutine 間通信 |
| コスト | ~1-5ns | ~20-50ns | ~50-300ns |
| ロック | なし（ロックフリー） | あり | あり（内部） |
| 複数変数の一貫性 | × | ○ | ○ |
| データの所有権移転 | × | × | ○ |

## 判断フローチャート

```
単一の値を更新するだけ？
├─ Yes → atomic
└─ No
    複数の値を一貫して更新？
    ├─ Yes → mutex
    └─ No
        goroutine 間でデータを渡す？
        ├─ Yes → channel
        └─ No → mutex（迷ったらこれ）
```

## atomic.Value

設定のホットリロードに便利:
```go
var config atomic.Value
config.Store(loadConfig())
// 別 goroutine から安全に読める
cfg := config.Load().(*Config)
```

## 実行方法

```bash
go run ./26_atomic_vs_mutex_vs_channel/
```
