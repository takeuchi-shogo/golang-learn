# 13. Escape Analysis（エスケープ解析）

## 核心ポイント
**`-gcflags="-m"` でスタック/ヒープの割り当てを確認できる。**

## スタック vs ヒープ

| | スタック | ヒープ |
|--|---------|-------|
| 速度 | 超高速（ポインタ移動のみ） | 遅い（GC が管理） |
| 寿命 | 関数スコープ | GC が回収するまで |
| コスト | ほぼゼロ | GC 負荷 |

## エスケープとは

変数が関数のスコープ外で参照される可能性がある場合、ヒープに割り当てられる（エスケープする）。

## ヒープにエスケープする典型パターン

1. **ポインタを返す**: `return &x`
2. **interface{} に代入**: `var i any = x`
3. **クロージャでキャプチャ**: `func() { use(x) }`
4. **スライスの成長**: `append` で容量超過
5. **チャネルに送信**: `ch <- &x`
6. **大きすぎる変数**: スタックに収まらない

## 確認方法

```bash
# エスケープ解析の結果を表示
go build -gcflags="-m" ./13_escape_analysis/

# より詳細
go build -gcflags="-m -m" ./13_escape_analysis/
```

出力例:
```
./main.go:10:6: moved to heap: x    ← ヒープにエスケープ
./main.go:15:6: x does not escape   ← スタックに留まる
```

## 実行方法

```bash
go run ./13_escape_analysis/

# エスケープ解析を確認
go build -gcflags="-m" ./13_escape_analysis/
```
