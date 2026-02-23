# 01. Go Memory Model

## 核心ポイント
**happens-before 関係がなければ、変数の読み書き順序は保証されない。**

## Go Memory Model とは

Go Memory Model は「ある goroutine での変数の書き込みが、別の goroutine での読み取りで観測される条件」を定義する。

### コンパイラとCPUのリオーダリング

```
// あなたが書いたコード:
x = 1
y = 2

// コンパイラ/CPU が実行するかもしれない順序:
y = 2
x = 1
```

単一 goroutine 内ではリオーダリングは観測されない（as-if-serial）。
しかし **別の goroutine から見ると、書き込み順序が入れ替わる可能性がある。**

## happens-before 関係

2つのイベント A, B について：
- **A happens-before B**: A の結果は B に必ず見える
- **関係なし**: どちらが先かは不定。コンパイラとCPUが自由にリオーダリングできる

### happens-before を作る同期プリミティブ

| プリミティブ | 保証 |
|------------|------|
| `sync.Mutex` | Unlock → 次の Lock |
| channel 送信 | send → 対応する receive |
| `sync.WaitGroup` | Done → Wait の完了 |
| `sync/atomic` | Store → 後続の Load |
| `sync.Once` | Do の完了 → 他の Do の開始 |
| goroutine 起動 | go 文 → goroutine の開始 |

## なぜ「動いてるように見える」が嘘なのか

1. **テスト環境**: GOMAXPROCS=1 で goroutine が直列に動く
2. **x86 の TSO**: x86 は比較的強いメモリモデルを持つ（ARM は弱い）
3. **タイミング**: race が発生する確率が低いだけ
4. **最適化レベル**: ビルドオプションで最適化が変わる

→ **プロダクションの ARM サーバーで突然壊れる**

## 実践的な対策

1. 共有変数へのアクセスには必ず同期プリミティブを使う
2. `go test -race ./...` を CI に組み込む
3. 「ロックなしで読むだけだから大丈夫」は嘘
4. `sync/atomic` は単一変数のみ。複数変数の一貫性には mutex を使う

## 実行方法

```bash
go run ./01_memory_model/
```

## 参考

- [The Go Memory Model (公式)](https://go.dev/ref/mem)
