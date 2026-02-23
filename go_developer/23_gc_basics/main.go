package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== GC Basics ===")
	fmt.Println("※ GODEBUG=gctrace=1 go run ./23_gc_basics/ で GC トレースを確認")
	fmt.Println()

	memStatsDemo()
	allocationPressureDemo()
	manualGCDemo()
}

// memStatsDemo: runtime.MemStats でヒープ統計を確認
func memStatsDemo() {
	fmt.Println("--- MemStats ---")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("  Alloc      = %d KB (現在のヒープ使用量)\n", m.Alloc/1024)
	fmt.Printf("  TotalAlloc = %d KB (累計割り当て量)\n", m.TotalAlloc/1024)
	fmt.Printf("  Sys        = %d KB (OS から取得したメモリ)\n", m.Sys/1024)
	fmt.Printf("  NumGC      = %d (GC 実行回数)\n", m.NumGC)
	fmt.Printf("  GCCPUFraction = %.4f%% (GC の CPU 使用率)\n", m.GCCPUFraction*100)
	fmt.Println()
}

// allocationPressureDemo: アロケーション圧力と GC 頻度の関係
func allocationPressureDemo() {
	fmt.Println("--- アロケーション圧力 ---")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	gcBefore := m.NumGC

	// 大量のアロケーション → GC が頻繁に発生
	var sink []*[1024]byte
	for range 10000 {
		p := new([1024]byte) // 1KB ずつヒープに割り当て
		sink = append(sink, p)
	}

	runtime.ReadMemStats(&m)
	gcAfter := m.NumGC

	fmt.Printf("  10000 × 1KB 割り当て後:\n")
	fmt.Printf("  GC 回数: %d → %d (+%d)\n", gcBefore, gcAfter, gcAfter-gcBefore)
	fmt.Printf("  ヒープ使用量: %d KB\n", m.Alloc/1024)

	_ = sink // GC に回収されないようにする
	fmt.Println()
}

// manualGCDemo: runtime.GC() で手動 GC
func manualGCDemo() {
	fmt.Println("--- 手動 GC ---")

	// 大量割り当て
	for range 10000 {
		_ = make([]byte, 1024)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  GC 前: Alloc=%d KB\n", m.Alloc/1024)

	runtime.GC() // 手動 GC

	runtime.ReadMemStats(&m)
	fmt.Printf("  GC 後: Alloc=%d KB\n", m.Alloc/1024)
	fmt.Println("  ※ 実務で runtime.GC() を呼ぶことは稀（ランタイムに任せる）")
}
