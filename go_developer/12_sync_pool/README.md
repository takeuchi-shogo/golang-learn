# 12. sync.Pool

## 核心ポイント
**sync.Pool はキャッシュではない。GC でクリアされる。**

## sync.Pool とは

一時的なオブジェクトを再利用するための仕組み。
**目的はアロケーション削減**であり、永続的なキャッシュではない。

## 基本操作

```go
pool := &sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

buf := pool.Get().(*bytes.Buffer)  // 取得（なければ New が呼ばれる）
buf.Reset()                        // 使う前にリセット
// ... buf を使う ...
pool.Put(buf)                      // 返却
```

## GC との関係

- GC の**各サイクル**で Pool のオブジェクトはクリアされる
- Pool にあるオブジェクトはいつ消えるか分からない
- → 永続キャッシュとして使うな

## 適切なユースケース

- `bytes.Buffer` の再利用（JSON エンコード/デコード時）
- `fmt` パッケージ内部での利用
- 一時的なスライスバッファ
- 重い構造体の allocate/free を減らす

## 不適切なユースケース

- DB コネクションプール → `sql.DB` を使う
- 永続キャッシュ → `sync.Map` や外部キャッシュを使う
- ステートを持つオブジェクト → Put 前にリセットが必須

## 実行方法

```bash
go run ./12_sync_pool/
```
