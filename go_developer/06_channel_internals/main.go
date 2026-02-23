package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Channel Internals ===")
	fmt.Println()

	unbufferedDemo()
	bufferedDemo()
	directionDemo()
	bufferStateDemo()
}

// unbufferedDemo: unbuffered channel は同期点（rendezvous）
func unbufferedDemo() {
	fmt.Println("--- Unbuffered Channel（同期通信）---")

	ch := make(chan string) // cap = 0

	go func() {
		fmt.Println("  送信側: 送信前（受信側を待つ）")
		ch <- "hello" // 受信側が来るまでブロック
		fmt.Println("  送信側: 送信完了（受信側と同期した）")
	}()

	time.Sleep(100 * time.Millisecond) // 送信側がブロックしていることを示す
	fmt.Println("  受信側: 受信前")
	msg := <-ch // 送信側のブロックが解除される
	fmt.Printf("  受信側: %q を受信\n", msg)
	fmt.Println()
}

// bufferedDemo: buffered channel はリングバッファ
func bufferedDemo() {
	fmt.Println("--- Buffered Channel（非同期通信）---")

	ch := make(chan int, 3) // cap = 3

	// バッファに空きがあれば、送信はブロックしない
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("  3つ送信（バッファ満杯）: len=%d, cap=%d\n", len(ch), cap(ch))

	// バッファが満杯なので、4つ目の送信はブロックする
	go func() {
		ch <- 4 // ← バッファに空きができるまでブロック
		fmt.Println("  4つ目の送信完了（空きができた）")
	}()

	// 1つ受信してバッファに空きを作る
	v := <-ch
	fmt.Printf("  受信: %d (空きができて4つ目が入る)\n", v)

	time.Sleep(50 * time.Millisecond)

	// 残りを受信
	for range 3 {
		fmt.Printf("  受信: %d\n", <-ch)
	}
	fmt.Println()
}

// directionDemo: channel の方向指定
func directionDemo() {
	fmt.Println("--- Channel の方向指定 ---")

	ch := make(chan string, 1)

	// send-only channel を渡す → 受信するとコンパイルエラー
	producer(ch)
	// recv-only channel を渡す → 送信するとコンパイルエラー
	consumer(ch)
	fmt.Println()
}

func producer(ch chan<- string) { // send-only
	ch <- "data from producer"
}

func consumer(ch <-chan string) { // recv-only
	msg := <-ch
	fmt.Printf("  consumer: %q\n", msg)
}

// bufferStateDemo: len() と cap() でバッファ状態を確認
func bufferStateDemo() {
	fmt.Println("--- バッファ状態の確認 ---")

	ch := make(chan int, 5)
	fmt.Printf("  初期:     len=%d, cap=%d\n", len(ch), cap(ch))

	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("  3つ追加:  len=%d, cap=%d\n", len(ch), cap(ch))

	<-ch
	fmt.Printf("  1つ受信:  len=%d, cap=%d\n", len(ch), cap(ch))
}
