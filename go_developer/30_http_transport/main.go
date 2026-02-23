package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== net/http Transport ===")
	fmt.Println()

	transportConfigDemo()
	correctResponseHandling()
}

// transportConfigDemo: カスタム Transport の設定例
func transportConfigDemo() {
	fmt.Println("--- カスタム Transport 設定 ---")

	transport := &http.Transport{
		MaxIdleConns:        100,              // 全体のアイドル接続数
		MaxIdleConnsPerHost: 10,               // ← デフォルト2は少なすぎる!
		IdleConnTimeout:     90 * time.Second, // アイドル接続の生存時間

		// 各種タイムアウト
		// DialContext 等は net.Dialer で設定
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // リクエスト全体のタイムアウト
	}

	fmt.Println("  MaxIdleConnsPerHost: 10 (デフォルト2は高負荷で不足)")
	fmt.Println("  Timeout: 30s (必ず設定する)")
	fmt.Println()

	_ = client // 実際のリクエストは行わない
}

// correctResponseHandling: レスポンスの正しい処理
func correctResponseHandling() {
	fmt.Println("--- レスポンスの正しい処理 ---")
	fmt.Println()

	fmt.Println("  // 悪い例: Body を読み切らない")
	fmt.Println("  resp, _ := client.Get(url)")
	fmt.Println("  resp.Body.Close()  // ← 接続が再利用されない!")
	fmt.Println()
	fmt.Println("  // 良い例: Body を完全に読み切ってから Close")
	fmt.Println("  resp, _ := client.Get(url)")
	fmt.Println("  io.Copy(io.Discard, resp.Body)  // 全て読み切る")
	fmt.Println("  resp.Body.Close()                // 接続がプールに戻る")
	fmt.Println()

	// 実際のパターンを関数で示す
	fmt.Println("  正しいレスポンス処理パターン:")
	fmt.Println("    defer func() {")
	fmt.Println("        io.Copy(io.Discard, resp.Body)")
	fmt.Println("        resp.Body.Close()")
	fmt.Println("    }()")
}

// safeGet: 正しいレスポンス処理のパターン（ユーティリティ関数例）
func safeGet(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Body を完全に読み切ってから Close
		// → TCP 接続がプールに戻り再利用される
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	return io.ReadAll(resp.Body)
}
