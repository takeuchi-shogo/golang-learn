package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== Channel Close ===")
	fmt.Println()

	okIdiomDemo()
	rangeDemo()
	broadcastDemo()
}

// okIdiomDemo: close されたかを ok で判断
func okIdiomDemo() {
	fmt.Println("--- ok idiom ---")

	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	close(ch)

	// close 後もバッファのデータは受信できる
	for {
		v, ok := <-ch
		if !ok {
			fmt.Println("  channel closed!")
			break
		}
		fmt.Printf("  v=%d, ok=%t\n", v, ok)
	}

	// close 後の受信はゼロ値 + false
	v, ok := <-ch
	fmt.Printf("  close 後: v=%d, ok=%t (ゼロ値)\n", v, ok)
	fmt.Println()
}

// rangeDemo: range で channel を読み切る
func rangeDemo() {
	fmt.Println("--- range over channel ---")

	ch := make(chan string, 3)
	ch <- "a"
	ch <- "b"
	ch <- "c"
	close(ch) // close しないと range が永遠にブロック

	// close されるまで自動的にループ
	for v := range ch {
		fmt.Printf("  %s\n", v)
	}
	fmt.Println()
}

// broadcastDemo: close で全 goroutine に一斉通知
func broadcastDemo() {
	fmt.Println("--- Broadcast with close ---")

	done := make(chan struct{})
	var wg sync.WaitGroup

	// 5つの goroutine が done を待つ
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-done // close されるまでブロック
			fmt.Printf("  worker %d: 通知を受信!\n", id)
		}(i)
	}

	fmt.Println("  close(done) で一斉通知...")
	close(done) // 全員が同時に解除される

	wg.Wait()
}
