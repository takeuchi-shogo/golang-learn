# 17. defer Internals

## 核心ポイント
**hot loop 内の defer は遅い。関数に切り出して対処する。**

## defer の動作

1. **LIFO（後入れ先出し）**: 後に defer した関数が先に実行される
2. **引数は即時評価**: defer 文が実行された時点で引数が評価される
3. **名前付き戻り値の変更可能**: defer 内で名前付き戻り値を変更できる

## 内部実装の変遷

| バージョン | 方式 | コスト |
|-----------|------|--------|
| Go 1.12 以前 | ヒープ割当（_defer 構造体） | ~35ns/defer |
| Go 1.13 | スタック割当 | ~10ns/defer |
| Go 1.14+ | Open-coded defer | ~0ns（インライン展開） |

### Open-coded defer（Go 1.14+）
- defer をコンパイル時にインライン展開
- 8個以下の defer かつループ外なら適用
- ほぼゼロコスト

### ループ内 defer が遅い理由
Open-coded defer が適用されず、毎イテレーションで `_defer` 構造体を割り当てるため。

## 実行方法

```bash
go run ./17_defer_internals/
```
