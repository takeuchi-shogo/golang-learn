package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== GOMAXPROCS ===")
	fmt.Println()

	infoDemo()
	gomaxprocsEffectDemo()
}

func infoDemo() {
	fmt.Println("--- 現在の設定 ---")
	fmt.Printf("  NumCPU:     %d\n", runtime.NumCPU())
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)) // 0 = 変更せず取得
	fmt.Println()
	fmt.Println("  K8s/Docker での推奨:")
	fmt.Println("    import _ \"go.uber.org/automaxprocs\"")
	fmt.Println("    → コンテナの CPU limit を自動検出")
	fmt.Println()
}

// gomaxprocsEffectDemo: GOMAXPROCS の値が並列実行に与える影響
func gomaxprocsEffectDemo() {
	fmt.Println("--- GOMAXPROCS の効果 ---")

	for _, procs := range []int{1, 2, runtime.NumCPU()} {
		elapsed := benchWithProcs(procs)
		fmt.Printf("  GOMAXPROCS=%d: %v\n", procs, elapsed)
	}
	fmt.Println()
	fmt.Println("  ※ CPU-bound タスクでは GOMAXPROCS を増やすと速くなる")
	fmt.Println("  ※ ただしコンテナの CPU limit を超えると逆効果")
}

func benchWithProcs(procs int) time.Duration {
	prev := runtime.GOMAXPROCS(procs)
	defer runtime.GOMAXPROCS(prev)

	start := time.Now()
	var wg sync.WaitGroup

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// CPU-bound work
			sum := 0
			for i := range 10_000_000 {
				sum += i
			}
			_ = sum
		}()
	}

	wg.Wait()
	return time.Since(start)
}
