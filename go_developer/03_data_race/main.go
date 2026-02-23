package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("=== Data Race ===")
	fmt.Println("※ go run -race ./03_data_race/ で実行すると race が検出される")
	fmt.Println()

	raceExample()
	fixedMutexExample()
	fixedAtomicExample()
}

// raceExample: data race のあるコード
// go run -race で実行すると WARNING: DATA RACE が出る
func raceExample() {
	fmt.Println("--- Data Race（壊れた例）---")

	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ← DATA RACE: read-modify-write が非原子的
		}()
	}

	wg.Wait()
	// 期待値は1000だが、race により1000未満になることがある
	fmt.Printf("  counter = %d (期待値: 1000, 実際: 不定)\n", counter)
	fmt.Println()
}

// fixedMutexExample: Mutex で修正した例
func fixedMutexExample() {
	fmt.Println("--- Mutex で修正 ---")

	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // Mutex で保護
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("  counter = %d (常に 1000)\n", counter)
	fmt.Println()
}

// fixedAtomicExample: atomic で修正した例
func fixedAtomicExample() {
	fmt.Println("--- atomic で修正 ---")

	var counter atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1) // atomic 操作: ロックフリー
		}()
	}

	wg.Wait()
	fmt.Printf("  counter = %d (常に 1000)\n", counter.Load())
}
