# 35. Build Tags + Module Replace

## 核心ポイント
**replace はローカル/CI 全てのビルドに効く。ビルドタグで環境別コードを分ける。**

## Build Tags

### 新構文 (Go 1.17+)
```go
//go:build linux && amd64
```

### 旧構文
```go
// +build linux,amd64
```

### カスタムタグ
```bash
go build -tags "integration" ./...
go test -tags "e2e" ./...
```

### ファイル名ベースの制約
```
foo_linux.go     → Linux のみ
foo_windows.go   → Windows のみ
foo_test.go      → テスト時のみ
```

## go.mod の replace

```go
// ローカル開発用
replace github.com/myorg/mylib => ../mylib

// 特定バージョンの置き換え
replace github.com/broken/lib v1.0.0 => github.com/fixed/lib v1.0.1
```

### replace の影響範囲

- `go.mod` はモジュールのルートに置かれる
- replace は**そのモジュールをビルドする全ての環境**に影響
- CI でもローカルでも同じ replace が適用される
- → コミット前に不要な replace を削除すること

## 実行方法

```bash
go run ./35_build_tags/

# カスタムタグ付きビルド
go run -tags "debug" ./35_build_tags/
```
