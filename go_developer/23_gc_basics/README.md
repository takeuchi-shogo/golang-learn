# 23. GC Basics（ガベージコレクション）

## 核心ポイント
**三色マーキング + concurrent GC。GOMEMLIMIT (Go 1.19+) を推奨。**

## 三色マーキング

| 色 | 意味 |
|----|------|
| 白 | 未到達（GC 対象候補） |
| 灰 | 到達済みだが子を未探索 |
| 黒 | 到達済みかつ子も探索済み |

### アルゴリズム
1. 全オブジェクトを白にする
2. ルート（スタック、グローバル変数）を灰にする
3. 灰オブジェクトの子を灰にし、自身を黒にする
4. 灰がなくなるまで繰り返す
5. 残った白オブジェクトを回収

## GC フェーズ

1. **Mark Setup (STW)**: write barrier を有効化（短い）
2. **Concurrent Mark**: アプリと並行してマーキング
3. **Mark Termination (STW)**: マーキング完了を確認（短い）
4. **Concurrent Sweep**: メモリ回収（アプリと並行）

## GC トリガー

- **GOGC=100**（デフォルト）: ヒープが前回 GC 後の2倍になったら GC
- **GOMEMLIMIT**: メモリ上限を設定（Go 1.19+、推奨）

## GOMEMLIMIT（推奨）

```bash
GOMEMLIMIT=1GiB ./myapp
```

コンテナ環境では OOM を防ぐために必ず設定する。

## 実行方法

```bash
go run ./23_gc_basics/

# GC トレースを有効にして実行
GODEBUG=gctrace=1 go run ./23_gc_basics/
```
