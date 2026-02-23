package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Preemption（プリエンプション） ===")
	fmt.Println()

	asyncPreemption()
	explicitYield()
}

// asyncPreemption: Go 1.14+ の非同期プリエンプションを確認
// tight loop を回しても他の goroutine が動けることを証明する
func asyncPreemption() {
	fmt.Println("--- 非同期プリエンプション (Go 1.14+) ---")

	// P を1つに制限して、プリエンプションの効果を明確にする
	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: tight CPU loop（関数呼び出しなし）
	go func() {
		defer wg.Done()
		start := time.Now()
		count := 0
		for time.Since(start) < 100*time.Millisecond {
			count++ // tight loop: 関数呼び出しなしの CPU-bound 処理
		}
		fmt.Printf("  tight loop: %d 回のイテレーション\n", count)
	}()

	// Goroutine 2: 定期的にメッセージを出す
	// Go 1.14+ なら、tight loop があっても実行される
	go func() {
		defer wg.Done()
		for range 3 {
			time.Sleep(30 * time.Millisecond)
			fmt.Println("  他の goroutine: 実行された！（プリエンプションが機能している）")
		}
	}()

	wg.Wait()
	fmt.Println()
}

// explicitYield: runtime.Gosched() で明示的に譲る
func explicitYield() {
	fmt.Println("--- runtime.Gosched()（明示的な譲渡）---")

	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)

	done := make(chan struct{})

	go func() {
		for range 5 {
			fmt.Print("A")
			// Gosched() なしだと、A が連続して出力されやすい
			// Gosched() ありだと、A と B が交互に出やすい
			runtime.Gosched()
		}
		close(done)
	}()

	for range 5 {
		fmt.Print("B")
		runtime.Gosched()
	}

	<-done
	fmt.Println()
	fmt.Println("  ※ A と B が交互に出やすい（Gosched で譲渡）")
}
