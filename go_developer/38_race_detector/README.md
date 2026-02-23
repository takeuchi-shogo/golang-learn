# 38. Race Detector

## 核心ポイント
**実行時のみ検知。false negative はあるが false positive はない。**

## -race フラグ

```bash
go test -race ./...   # テスト時（推奨）
go build -race ./...  # ビルド時
go run -race main.go  # 実行時
```

## 仕組み

ThreadSanitizer（TSan）ベース:
- 全てのメモリアクセスを記録
- happens-before 関係を追跡
- 同期なしのアクセスを検出

## 検出できるもの

- 同一変数への非同期 read/write
- map への concurrent write
- スライスへの concurrent アクセス
- 構造体フィールドへの concurrent アクセス

## 検出できないもの

- **実行されなかったパスの race**: コードが実行されなければ検出できない
- → テストカバレッジが重要

## オーバーヘッド

| | 通常ビルド | -race ビルド |
|--|-----------|-------------|
| CPU | 1x | 5-10x |
| メモリ | 1x | 5-10x |
| バイナリサイズ | 1x | 2x |

## CI での使い方

```yaml
- name: Test with race detector
  run: go test -race -count=1 -timeout=5m ./...
```

`-count=1` でキャッシュを無効化（race は非決定的なため）。

## 実行方法

```bash
# Race detector 付きで実行
go run -race ./38_race_detector/
```
