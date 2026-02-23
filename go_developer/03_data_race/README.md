# 03. Data Race

## 核心ポイント
**Data race は fatal error（recover 不可）。`-race` フラグで必ず検査せよ。**

## Data Race とは

以下の3条件を同時に満たすとき data race が発生する：

1. 2つ以上の goroutine が**同じ変数**にアクセス
2. 少なくとも1つが**書き込み**
3. **同期が無い**（happens-before 関係がない）

## なぜ「自分のマシンでは動く」は嘘か

- **x86 の TSO**: 比較的強いメモリモデルなので race が顕在化しにくい
- **GOMAXPROCS=1**: goroutine が直列に動くと race が起きない
- **タイミング依存**: race window が短く、確率的にしか発生しない
- **ARM/M1**: x86 より弱いメモリモデルのため、x86 で動くコードが壊れる

## fatal error は recover 不可

```go
// concurrent map writes → fatal error（recover できない）
// concurrent map read + write → fatal error
```

通常の `panic` と異なり、`recover()` で捕捉できない。プロセスが即座に終了する。

## Race Detector の使い方

```bash
# テスト時（推奨）
go test -race ./...

# ビルド時
go build -race -o myapp .

# 実行時
go run -race main.go
```

### CI への組み込み

```yaml
# GitHub Actions の例
- name: Test with race detector
  run: go test -race -count=1 ./...
```

## 典型的な data race パターン

1. **カウンタの非同期更新**: `count++` は read-modify-write（3操作）
2. **map の同時書き込み**: fatal error
3. **スライスへの同時 append**: backing array の破損
4. **構造体フィールドの非同期アクセス**: フィールド単位でも race

## 実行方法

```bash
# 通常実行（問題が見えないかもしれない）
go run ./03_data_race/

# Race detector 付き（推奨）
go run -race ./03_data_race/
```
