package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== panic / recover ===")
	fmt.Println()

	basicRecoverDemo()
	onlyInDeferDemo()
	goroutineIsolationDemo()
	httpHandlerPattern()
}

// basicRecoverDemo: 基本的な panic/recover
func basicRecoverDemo() {
	fmt.Println("--- 基本的な recover ---")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  recovered: %v\n", r)
		}
	}()

	fmt.Println("  panic 前")
	panic("something went wrong")
	// ここには到達しない
}

// onlyInDeferDemo: defer 内でのみ recover が有効
func onlyInDeferDemo() {
	fmt.Println()
	fmt.Println("--- defer 内でのみ有効 ---")

	// 通常のコードで recover() を呼んでも常に nil
	r := recover()
	fmt.Printf("  通常コードでの recover(): %v (常に nil)\n", r)

	// defer 内で呼ぶと有効
	result := safeCall(func() {
		panic("error in function")
	})
	fmt.Printf("  safeCall の結果: %v\n", result)
	fmt.Println()
}

// safeCall: panic を error に変換するラッパー（実務パターン）
func safeCall(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught: %v", r)
		}
	}()
	f()
	return nil
}

// goroutineIsolationDemo: 別 goroutine の panic は recover できない
func goroutineIsolationDemo() {
	fmt.Println("--- goroutine 間の分離 ---")
	fmt.Println("  別 goroutine の panic は recover できない")
	fmt.Println("  → 各 goroutine に個別の recover が必要")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  goroutine 内で recover: %v\n", r)
			}
		}()
		panic("goroutine panic")
	}()

	wg.Wait()
	fmt.Println()
}

// httpHandlerPattern: HTTP handler を保護するパターン
func httpHandlerPattern() {
	fmt.Println("--- HTTP handler 保護パターン ---")

	handler := func() {
		panic("unexpected error in handler")
	}

	// ミドルウェアで handler をラップ
	recoveryMiddleware(handler)
}

func recoveryMiddleware(handler func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  [middleware] 500 Internal Server Error: %v\n", r)
			// 実務では: レスポンス 500 を返す + エラーログ出力
		}
	}()
	handler()
}
