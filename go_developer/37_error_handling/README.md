# 37. Error Handling

## 核心ポイント
**`%w` で wrap、`errors.Is` / `errors.As` でチェイン遡行。**

## エラーの基本

```go
// errors.New: 単純なエラー
var ErrNotFound = errors.New("not found")

// fmt.Errorf + %w: エラーのラッピング
return fmt.Errorf("user %d: %w", id, ErrNotFound)
```

## errors.Is（値の比較）

エラーチェインを遡ってセンチネルエラーと比較:
```go
if errors.Is(err, ErrNotFound) {
    // err のチェインのどこかに ErrNotFound がある
}
```

## errors.As（型の比較）

エラーチェインを遡って特定の型を抽出:
```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // pathErr にキャストされた値が入る
}
```

## errors.Join（Go 1.20+）

複数のエラーを1つにまとめる:
```go
err := errors.Join(err1, err2, err3)
errors.Is(err, err1) // true
```

## ベストプラクティス

1. エラーは呼び出し元にコンテキストを追加して wrap
2. `%w` で wrap（`%v` だとチェインが切れる）
3. sentinel error はパッケージレベルで `var Err... = errors.New(...)`
4. `errors.Is` / `errors.As` を使う（`==` 比較は非推奨）

## 実行方法

```bash
go run ./37_error_handling/
```
