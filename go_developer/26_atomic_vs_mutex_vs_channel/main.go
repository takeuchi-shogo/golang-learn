package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("=== atomic vs mutex vs channel ===")
	fmt.Println()

	atomicCounterDemo()
	mutexCounterDemo()
	channelCounterDemo()
	atomicValueDemo()
}

const iterations = 100000

// atomicCounterDemo: atomic でカウンタ（最速）
func atomicCounterDemo() {
	fmt.Println("--- atomic カウンタ ---")

	var counter atomic.Int64
	var wg sync.WaitGroup

	for range iterations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	fmt.Printf("  結果: %d (ロックフリー、最速)\n", counter.Load())
	fmt.Println()
}

// mutexCounterDemo: mutex でカウンタ
func mutexCounterDemo() {
	fmt.Println("--- mutex カウンタ ---")

	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	for range iterations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("  結果: %d (排他制御、複数変数に対応)\n", counter)
	fmt.Println()
}

// channelCounterDemo: channel でカウンタ
func channelCounterDemo() {
	fmt.Println("--- channel カウンタ ---")

	counter := make(chan int, 1)
	counter <- 0

	var wg sync.WaitGroup
	for range iterations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := <-counter
			c++
			counter <- c
		}()
	}

	wg.Wait()
	result := <-counter
	fmt.Printf("  結果: %d (通信ベース、所有権移転)\n", result)
	fmt.Println()
}

// atomicValueDemo: atomic.Value で設定のホットリロード
func atomicValueDemo() {
	fmt.Println("--- atomic.Value（設定のホットリロード）---")

	type Config struct {
		Host string
		Port int
	}

	var config atomic.Value

	// 初期設定
	config.Store(&Config{Host: "localhost", Port: 8080})

	// 別 goroutine から安全に読める
	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg := config.Load().(*Config)
			fmt.Printf("  読み取り: %s:%d\n", cfg.Host, cfg.Port)
		}()
	}

	// 設定を更新（アトミックに置き換え）
	config.Store(&Config{Host: "production.example.com", Port: 443})

	wg.Wait()
	cfg := config.Load().(*Config)
	fmt.Printf("  最終: %s:%d\n", cfg.Host, cfg.Port)
}
