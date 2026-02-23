# 34. init() Order

## 核心ポイント
**init() 内で panic すると起動不能。重い処理を入れるな。**

## 実行順序

1. **インポートされたパッケージ**の init() が先に実行される（再帰的）
2. 同一パッケージ内: **ファイル名のアルファベット順**
3. 同一ファイル内: **定義順**
4. パッケージ変数の初期化 → init() の順序

```
main imports pkg_a → pkg_a imports pkg_b
実行順序: pkg_b.init() → pkg_a.init() → main.init() → main()
```

## init() の特徴

- 引数なし、戻り値なし
- 同一ファイルに**複数定義可能**
- 明示的に呼び出せない
- テストでも実行される

## アンチパターン

- DB接続を init() でやる → テストで困る
- 環境変数を init() で読む → テストで上書きできない
- 重い計算を init() でやる → 起動が遅くなる
- init() で panic → プロセスが起動しない

## 代替案

```go
// init() の代わりに明示的な初期化関数を使う
func Setup(cfg Config) (*DB, error) {
    return db.Connect(cfg.DSN)
}
```

## 実行方法

```bash
go run ./34_init_order/
```
