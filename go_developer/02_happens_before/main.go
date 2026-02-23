package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("=== Happens-before 保証 ===")
	fmt.Println()

	mutexHB()
	channelHB()
	waitGroupHB()
	atomicHB()
}

// mutexHB: Unlock → Lock が happens-before を作る
func mutexHB() {
	fmt.Println("--- Mutex happens-before ---")
	var mu sync.Mutex
	var data string

	go func() {
		mu.Lock()
		data = "hello" // (1) 書き込み
		mu.Unlock()    // (2) Unlock ──HB──▶ 次の Lock
	}()

	mu.Lock()          // (3) この Lock は (2) の後
	fmt.Printf("  data = %q\n", data) // (1) の結果が見える保証
	mu.Unlock()
	fmt.Println()
}

// channelHB: send → receive が happens-before を作る
func channelHB() {
	fmt.Println("--- Channel happens-before ---")

	// Unbuffered: send の完了は receive の完了に先行
	done := make(chan struct{})
	var data string

	go func() {
		data = "world"  // (1) 書き込み
		done <- struct{}{} // (2) send ──HB──▶ receive
	}()

	<-done              // (3) receive: (1) の結果が見える保証
	fmt.Printf("  data = %q (unbuffered channel)\n", data)

	// Close: close → receive(ゼロ値) も happens-before
	ready := make(chan struct{})
	var value int

	go func() {
		value = 42
		close(ready) // close ──HB──▶ <-ready
	}()

	<-ready
	fmt.Printf("  value = %d (channel close)\n", value)
	fmt.Println()
}

// waitGroupHB: Done → Wait の完了が happens-before
func waitGroupHB() {
	fmt.Println("--- WaitGroup happens-before ---")

	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done() // Done ──HB──▶ Wait の完了
			results[idx] = idx * idx
		}(i)
	}

	wg.Wait() // 全ての Done の後に完了する
	// 全 goroutine の書き込みが見える保証
	fmt.Printf("  results = %v\n", results)
	fmt.Println()
}

// atomicHB: Store → Load が happens-before (Go 1.19+)
func atomicHB() {
	fmt.Println("--- Atomic happens-before ---")

	var flag atomic.Bool
	var payload string
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		payload = "secret" // (1)
		flag.Store(true)   // (2) Store ──HB──▶ Load
	}()

	go func() {
		defer wg.Done()
		for !flag.Load() { // (3) Load が true を返したら
		}
		// (1) の payload = "secret" が見える保証
		fmt.Printf("  payload = %q (atomic flag)\n", payload)
	}()

	wg.Wait()
}
