package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Build Tags + Module Replace ===")
	fmt.Println()

	buildTagsDemo()
	replaceDemo()
}

func buildTagsDemo() {
	fmt.Println("--- Build Tags ---")
	fmt.Printf("  GOOS:   %s\n", runtime.GOOS)
	fmt.Printf("  GOARCH: %s\n", runtime.GOARCH)
	fmt.Println()

	fmt.Println("  ファイル名ベースのビルド制約:")
	fmt.Println("    foo_linux.go      → GOOS=linux のときのみコンパイル")
	fmt.Println("    foo_windows.go    → GOOS=windows のときのみコンパイル")
	fmt.Println("    foo_darwin.go     → GOOS=darwin (macOS) のときのみコンパイル")
	fmt.Println("    foo_amd64.go      → GOARCH=amd64 のときのみコンパイル")
	fmt.Println("    foo_test.go       → go test のときのみコンパイル")
	fmt.Println()

	fmt.Println("  //go:build タグの例:")
	fmt.Println("    //go:build linux && amd64        → Linux + AMD64")
	fmt.Println("    //go:build !windows              → Windows 以外")
	fmt.Println("    //go:build integration            → -tags integration で有効")
	fmt.Println()

	fmt.Println("  カスタムタグの使い方:")
	fmt.Println("    go test -tags \"integration\" ./... → integration テストを実行")
	fmt.Println("    go build -tags \"debug\" ./...      → debug ビルド")
	fmt.Println()
}

func replaceDemo() {
	fmt.Println("--- go.mod replace ---")
	fmt.Println()

	fmt.Println("  ローカル開発用:")
	fmt.Println("    replace github.com/myorg/mylib => ../mylib")
	fmt.Println()
	fmt.Println("  バージョン置き換え:")
	fmt.Println("    replace github.com/broken/lib v1.0.0 => github.com/fixed/lib v1.0.1")
	fmt.Println()
	fmt.Println("  注意点:")
	fmt.Println("    - replace はモジュールの全ビルドに影響する")
	fmt.Println("    - CI でも同じ replace が適用される")
	fmt.Println("    - コミット前に不要な replace を削除する")
	fmt.Println("    - 依存先のモジュールの replace は無視される（トップレベルのみ有効）")
}
