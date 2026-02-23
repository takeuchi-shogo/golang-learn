package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== sync.Once ===")
	fmt.Println()

	basicOnceDemo()
	concurrentOnceDemo()
	onceFuncDemo()
	panicOnceDemo()
}

// basicOnceDemo: 一度だけ実行される
func basicOnceDemo() {
	fmt.Println("--- 基本: 一度だけ実行 ---")

	var once sync.Once

	for range 5 {
		once.Do(func() {
			fmt.Println("  初期化処理（1回だけ実行される）")
		})
	}
	fmt.Println()
}

// concurrentOnceDemo: 複数 goroutine から呼んでも1回だけ
func concurrentOnceDemo() {
	fmt.Println("--- 並行呼び出し ---")

	var once sync.Once
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(func() {
				fmt.Printf("  goroutine %d が初期化を実行\n", id)
			})
		}(i)
	}

	wg.Wait()
	fmt.Println("  ※ 10個の goroutine のうち1つだけが実行")
	fmt.Println()
}

// onceFuncDemo: Go 1.21+ の sync.OnceFunc / sync.OnceValue
func onceFuncDemo() {
	fmt.Println("--- sync.OnceFunc / sync.OnceValue (Go 1.21+) ---")

	// OnceFunc: 引数なし・戻り値なし
	init := sync.OnceFunc(func() {
		fmt.Println("  OnceFunc: 初期化!")
	})
	init() // 実行される
	init() // 実行されない

	// OnceValue: 戻り値あり
	getConfig := sync.OnceValue(func() string {
		fmt.Println("  OnceValue: 設定を読み込み")
		return "production"
	})
	fmt.Printf("  config = %s\n", getConfig())
	fmt.Printf("  config = %s (キャッシュ)\n", getConfig())
	fmt.Println()
}

// panicOnceDemo: panic しても完了扱い
func panicOnceDemo() {
	fmt.Println("--- panic しても完了扱い ---")

	var once sync.Once

	// 1回目: panic する
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  1回目: panic recovered: %v\n", r)
			}
		}()
		once.Do(func() {
			panic("初期化失敗!")
		})
	}()

	// 2回目: 実行されない（panic しても「完了」扱い）
	executed := false
	once.Do(func() {
		executed = true
	})
	fmt.Printf("  2回目: executed=%v (panic でも完了扱い)\n", executed)

	// デッドロックの例（実行するとハングするのでコメントのみ）:
	// var once2 sync.Once
	// once2.Do(func() {
	//     once2.Do(func() {}) // ← デッドロック!
	// })
}
