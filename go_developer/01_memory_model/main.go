package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================
// 01. Go Memory Model
// happens-before がない読み書きは順序保証がない
// ============================================================

func main() {
	fmt.Println("=== Go Memory Model ===")
	fmt.Println()

	brokenExample()
	correctMutexExample()
	correctAtomicExample()
}

// brokenExample は同期なしで変数を共有する危険な例
// 実行結果は不定（0, 1, 2 のいずれも有り得る）
func brokenExample() {
	fmt.Println("--- 壊れた例（同期なし）---")

	var x, y int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		x = 1 // (A)
		y = 1 // (B)
	}()

	go func() {
		defer wg.Done()
		// (B) が見えるのに (A) が見えない可能性がある
		// コンパイラ/CPU のリオーダリングにより:
		//   - r1=1, r2=0 が起こり得る
		r1 := y
		r2 := x
		fmt.Printf("  y=%d, x=%d (同期なし: 結果は不定)\n", r1, r2)
	}()

	wg.Wait()
	fmt.Println()
}

// correctMutexExample は Mutex で happens-before を保証する
func correctMutexExample() {
	fmt.Println("--- 正しい例（Mutex）---")

	var mu sync.Mutex
	var x, y int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		mu.Lock()
		x = 1
		y = 1
		mu.Unlock() // Unlock happens-before 次の Lock
	}()

	go func() {
		defer wg.Done()
		mu.Lock() // この Lock は前の Unlock の後
		r1 := y
		r2 := x
		mu.Unlock()
		// Mutex により x=1, y=1 が保証される
		// （ただし最初の goroutine が先にロックを取った場合）
		fmt.Printf("  y=%d, x=%d (Mutex で保護)\n", r1, r2)
	}()

	wg.Wait()
	fmt.Println()
}

// correctAtomicExample は atomic で happens-before を保証する
func correctAtomicExample() {
	fmt.Println("--- 正しい例（atomic）---")

	var ready atomic.Bool
	var data int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		data = 42              // データを準備
		ready.Store(true)      // Store happens-before 後続の Load
	}()

	go func() {
		defer wg.Done()
		// ready が true になるまでスピン
		for !ready.Load() {
			// busy wait（実際のコードでは channel を使う）
		}
		// ready.Store(true) happens-before ready.Load() == true
		// よって data = 42 が見える保証がある
		fmt.Printf("  data=%d (atomic で順序保証)\n", data)
	}()

	wg.Wait()
}
