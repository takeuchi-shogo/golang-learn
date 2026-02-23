package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("=== Race Detector ===")
	fmt.Println("※ go run -race ./38_race_detector/ で実行すると race が検出される")
	fmt.Println()

	racyCodeDemo()
	fixedCodeDemo()
	commonPatternsDemo()
}

// racyCodeDemo: race のあるコード
// go run -race で実行すると WARNING: DATA RACE が出る
func racyCodeDemo() {
	fmt.Println("--- Race のあるコード ---")

	counter := 0
	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ← DATA RACE!
		}()
	}
	wg.Wait()
	fmt.Printf("  counter = %d (race あり: 結果は不定)\n", counter)
	fmt.Println()
}

// fixedCodeDemo: race を修正したコード
func fixedCodeDemo() {
	fmt.Println("--- 修正済みコード ---")

	var counter atomic.Int64
	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1) // atomic: race-free
		}()
	}
	wg.Wait()
	fmt.Printf("  counter = %d (race なし: 常に 100)\n", counter.Load())
	fmt.Println()
}

// commonPatternsDemo: よくある race パターン
func commonPatternsDemo() {
	fmt.Println("--- よくある Race パターン ---")
	fmt.Println()

	// パターン 1: ループ変数のキャプチャ（Go 1.22+ で修正済み）
	fmt.Println("  1. ループ変数のキャプチャ (Go 1.22 以前)")
	fmt.Println("     for i := 0; i < n; i++ {")
	fmt.Println("         go func() { use(i) }()  // ← i は共有変数!")
	fmt.Println("     }")
	fmt.Println("     修正: go func(i int) { use(i) }(i)")
	fmt.Println("     Go 1.22+: ループ変数は各イテレーションでコピーされる")
	fmt.Println()

	// パターン 2: 共有 map
	fmt.Println("  2. 共有 map")
	fmt.Println("     go func() { m[\"key\"] = 1 }()  // ← concurrent write!")
	fmt.Println("     修正: sync.Mutex または sync.Map を使う")
	fmt.Println()

	// パターン 3: 構造体フィールド
	fmt.Println("  3. 構造体の異なるフィールド")
	fmt.Println("     異なるフィールドでも、同じ構造体への concurrent アクセスは race")
	fmt.Println()

	fmt.Println("=== Race Detector のまとめ ===")
	fmt.Println("  - CI で go test -race ./... を必ず実行")
	fmt.Println("  - -count=1 でキャッシュ無効化（race は非決定的）")
	fmt.Println("  - false positive はない（検出されたら必ず race）")
	fmt.Println("  - false negative はある（実行されなかったパスは検出不可）")
}
