package main

import (
	"fmt"
	"net/http"
	"time"
)

// ============================================================
// 正しいパターン: Client をグローバルまたは DI で再利用
// ============================================================

// 推奨: アプリケーション起動時に一度だけ作成
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

func main() {
	fmt.Println("=== http.Client Reuse ===")
	fmt.Println()

	showAntiPatterns()
	showCorrectPatterns()
}

func showAntiPatterns() {
	fmt.Println("--- アンチパターン ---")
	fmt.Println()

	fmt.Println("  // 1. リクエストごとに新しい Client（最悪）")
	fmt.Println("  func fetch(url string) {")
	fmt.Println("      client := &http.Client{}  // ← 毎回新しい接続プール!")
	fmt.Println("      resp, _ := client.Get(url)")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  問題:")
	fmt.Println("    - 毎リクエストで TCP + TLS ハンドシェイク")
	fmt.Println("    - 接続プールが共有されない")
	fmt.Println("    - 高負荷時にソケット枯渇")
	fmt.Println()

	fmt.Println("  // 2. http.DefaultClient を使う")
	fmt.Println("  http.Get(url)")
	fmt.Println()
	fmt.Println("  問題:")
	fmt.Println("    - Timeout が 0（= タイムアウトなし!）")
	fmt.Println("    - レスポンスが返らないと永遠にブロック")
	fmt.Println()
}

func showCorrectPatterns() {
	fmt.Println("--- 正しいパターン ---")
	fmt.Println()

	fmt.Println("  // 1. グローバル変数（シンプル）")
	fmt.Println("  var httpClient = &http.Client{")
	fmt.Println("      Timeout: 30 * time.Second,")
	fmt.Println("  }")
	fmt.Println()

	fmt.Println("  // 2. DI で注入（テスタブル）")
	fmt.Println("  type Service struct {")
	fmt.Println("      client *http.Client")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  func NewService(client *http.Client) *Service {")
	fmt.Println("      return &Service{client: client}")
	fmt.Println("  }")
	fmt.Println()

	// 設定値の確認
	fmt.Printf("  現在の httpClient 設定:\n")
	fmt.Printf("    Timeout:             %v\n", httpClient.Timeout)
	t := httpClient.Transport.(*http.Transport)
	fmt.Printf("    MaxIdleConns:        %d\n", t.MaxIdleConns)
	fmt.Printf("    MaxIdleConnsPerHost: %d\n", t.MaxIdleConnsPerHost)
	fmt.Printf("    IdleConnTimeout:     %v\n", t.IdleConnTimeout)
}
