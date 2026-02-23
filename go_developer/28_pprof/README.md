# 28. pprof

## 核心ポイント
**`_ "net/http/pprof"` を import するだけで HTTP プロファイリングが有効になる。**

## プロファイルの種類

| プロファイル | 用途 |
|------------|------|
| CPU | CPU を消費している関数 |
| heap | メモリ割り当て |
| goroutine | goroutine のスタックトレース |
| mutex | mutex の競合 |
| block | ブロッキング操作 |
| threadcreate | OS スレッド作成 |

## net/http/pprof（HTTP 経由）

```go
import _ "net/http/pprof"

go func() {
    http.ListenAndServe("localhost:6060", nil)
}()
```

→ `http://localhost:6060/debug/pprof/` でアクセス可能

## go tool pprof の使い方

```bash
# CPU プロファイル（30秒間）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# ヒーププロファイル
go tool pprof http://localhost:6060/debug/pprof/heap

# goroutine ダンプ
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 対話モードで分析
(pprof) top 10     # 上位10関数
(pprof) list func  # 関数のソースコード
(pprof) web        # ブラウザで可視化
```

## 実行方法

```bash
go run ./28_pprof/
# 別ターミナルで:
# go tool pprof http://localhost:6060/debug/pprof/heap
```
