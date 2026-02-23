package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("=== sync.Pool ===")
	fmt.Println()

	basicPoolDemo()
	gcClearsPoolDemo()
	bufferPoolDemo()
}

// basicPoolDemo: Pool の基本操作
func basicPoolDemo() {
	fmt.Println("--- 基本操作 ---")

	pool := &sync.Pool{
		New: func() any {
			fmt.Println("  New() 呼び出し（プールが空のとき）")
			return "新しいオブジェクト"
		},
	}

	// Get: プールが空なので New が呼ばれる
	obj1 := pool.Get().(string)
	fmt.Printf("  Get: %s\n", obj1)

	// Put: プールに返却
	pool.Put("再利用オブジェクト")

	// Get: プールにあるので New は呼ばれない
	obj2 := pool.Get().(string)
	fmt.Printf("  Get: %s (再利用された)\n", obj2)

	// もう一度 Get: プールが空なので New が呼ばれる
	obj3 := pool.Get().(string)
	fmt.Printf("  Get: %s\n", obj3)
	fmt.Println()
}

// gcClearsPoolDemo: GC で Pool がクリアされることを確認
func gcClearsPoolDemo() {
	fmt.Println("--- GC で Pool がクリアされる ---")

	pool := &sync.Pool{
		New: func() any { return "new" },
	}

	pool.Put("cached-value")
	obj := pool.Get().(string)
	fmt.Printf("  GC 前: %s\n", obj)

	pool.Put("cached-value-2")

	// GC を強制実行 → Pool がクリアされる
	runtime.GC()

	obj = pool.Get().(string)
	fmt.Printf("  GC 後: %s (New が呼ばれた = クリア済み)\n", obj)
	fmt.Println()
}

// bufferPoolDemo: bytes.Buffer の再利用パターン（実務でよく使う）
func bufferPoolDemo() {
	fmt.Println("--- bytes.Buffer Pool（実務パターン）---")

	bufPool := &sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}

	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Pool から取得
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset() // ← 必ずリセット！前の内容が残っている可能性

			// バッファを使う
			fmt.Fprintf(buf, "worker-%d: hello", id)
			result := buf.String()

			// Pool に返却
			bufPool.Put(buf)

			fmt.Printf("  %s\n", result)
		}(i)
	}

	wg.Wait()
}
