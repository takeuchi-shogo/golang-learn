package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("=== G, M, P スケジューラモデル ===")
	fmt.Println()

	showRuntimeInfo()
	demonstrateGoroutines()
	demonstrateGosched()
}

// showRuntimeInfo: ランタイム情報を表示
func showRuntimeInfo() {
	fmt.Println("--- ランタイム情報 ---")
	fmt.Printf("  NumCPU      = %d  (CPU コア数)\n", runtime.NumCPU())
	fmt.Printf("  GOMAXPROCS  = %d  (P の数 = 同時実行可能な M の数)\n", runtime.GOMAXPROCS(0))
	fmt.Printf("  NumGoroutine= %d  (現在の G の数)\n", runtime.NumGoroutine())
	fmt.Println()
}

// demonstrateGoroutines: 大量の goroutine を起動して G の数を観察
func demonstrateGoroutines() {
	fmt.Println("--- 大量の goroutine を起動 ---")

	const n = 10000
	var wg sync.WaitGroup
	started := make(chan struct{})

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-started // 全 goroutine が起動するまで待機
		}()
	}

	fmt.Printf("  起動後の Goroutine 数 = %d\n", runtime.NumGoroutine())
	// G は軽量なので 10000 個でもメモリは約 20MB (2KB × 10000)

	close(started) // 一斉に解放
	wg.Wait()

	fmt.Printf("  完了後の Goroutine 数 = %d\n", runtime.NumGoroutine())
	fmt.Println()
}

// demonstrateGosched: runtime.Gosched() で明示的に実行権を譲る
func demonstrateGosched() {
	fmt.Println("--- runtime.Gosched() ---")

	// GOMAXPROCS(1) で直列実行を強制
	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 3 {
			fmt.Printf("  goroutine A: %d\n", i)
			runtime.Gosched() // 実行権を譲る → B が動く機会を得る
		}
	}()

	go func() {
		defer wg.Done()
		for i := range 3 {
			fmt.Printf("  goroutine B: %d\n", i)
			runtime.Gosched()
		}
	}()

	wg.Wait()
}
