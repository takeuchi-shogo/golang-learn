package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Context ===")
	fmt.Println()

	withCancelDemo()
	withTimeoutDemo()
	propagationDemo()
	withValueDemo()
}

// withCancelDemo: 手動キャンセル
func withCancelDemo() {
	fmt.Println("--- WithCancel ---")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 必ず defer cancel()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Printf("  goroutine: キャンセルされた (err=%v)\n", ctx.Err())
		}
	}(ctx)

	cancel() // 手動でキャンセル
	time.Sleep(10 * time.Millisecond)
	fmt.Println()
}

// withTimeoutDemo: タイムアウト付きコンテキスト
func withTimeoutDemo() {
	fmt.Println("--- WithTimeout ---")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("  処理完了（ここには来ない）")
	case <-ctx.Done():
		fmt.Printf("  タイムアウト! err=%v\n", ctx.Err())
	}
	fmt.Println()
}

// propagationDemo: キャンセルの親→子伝播
func propagationDemo() {
	fmt.Println("--- キャンセルの伝播（親→子→孫）---")

	parent, cancelParent := context.WithCancel(context.Background())
	child, cancelChild := context.WithCancel(parent)
	defer cancelChild()

	grandchild, cancelGrandchild := context.WithCancel(child)
	defer cancelGrandchild()

	// 親をキャンセル → 子・孫が全てキャンセルされる
	cancelParent()

	fmt.Printf("  parent:     err=%v\n", parent.Err())
	fmt.Printf("  child:      err=%v\n", child.Err())
	fmt.Printf("  grandchild: err=%v\n", grandchild.Err())
	fmt.Println()
}

// withValueDemo: リクエストスコープの値を伝播
func withValueDemo() {
	fmt.Println("--- WithValue ---")

	// key は unexported type にする（衝突防止）
	type ctxKey string
	const traceIDKey ctxKey = "traceID"

	ctx := context.WithValue(context.Background(), traceIDKey, "abc-123")

	processRequest(ctx, traceIDKey)
}

func processRequest(ctx context.Context, key any) {
	if traceID, ok := ctx.Value(key).(string); ok {
		fmt.Printf("  traceID = %s\n", traceID)
	}
}
