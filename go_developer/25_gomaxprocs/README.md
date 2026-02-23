# 25. GOMAXPROCS

## 核心ポイント
**GOMAXPROCS は P の数を制御する。K8s では uber/automaxprocs を入れること。**

## GOMAXPROCS とは

同時に Go コードを実行できる OS スレッド（M）の最大数。
正確には **P（Processor）の数** を制御する。

- デフォルト: `runtime.NumCPU()`
- 設定方法: `GOMAXPROCS=N` 環境変数 or `runtime.GOMAXPROCS(n)`

## コンテナ環境の問題

Docker/Kubernetes では `runtime.NumCPU()` が**ホストの CPU 数**を返す。

```
ホスト: 64 CPU
コンテナ CPU limit: 2

→ GOMAXPROCS = 64（ホストの数！）
→ 64 個の P が 2 CPU を奪い合う
→ コンテキストスイッチ爆増 → パフォーマンス低下
```

## 解決策: uber-go/automaxprocs

```go
import _ "go.uber.org/automaxprocs"
```

コンテナの CPU limit を自動検出して GOMAXPROCS を設定する。

## 実行方法

```bash
go run ./25_gomaxprocs/
```
