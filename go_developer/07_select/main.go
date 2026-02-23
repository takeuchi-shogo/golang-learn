package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== select ===")
	fmt.Println()

	pseudoRandomDemo()
	nonBlockingDemo()
	timeoutDemo()
	forSelectDoneDemo()
}

// pseudoRandomDemo: 複数 case が ready なとき、ランダムに選ばれることを証明
func pseudoRandomDemo() {
	fmt.Println("--- Pseudo-random 選択 ---")

	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)
	count1, count2 := 0, 0

	for range 1000 {
		// 両方の channel にデータを入れる（両方 ready）
		ch1 <- struct{}{}
		ch2 <- struct{}{}

		select {
		case <-ch1:
			count1++
		case <-ch2:
			count2++
		}

		// 残りを消費
		select {
		case <-ch1:
		case <-ch2:
		}
	}

	fmt.Printf("  ch1 選択: %d回, ch2 選択: %d回 (≒50:50)\n", count1, count2)
	fmt.Println()
}

// nonBlockingDemo: default でノンブロッキング操作
func nonBlockingDemo() {
	fmt.Println("--- ノンブロッキング操作（default）---")

	ch := make(chan string, 1)

	// ノンブロッキング受信
	select {
	case msg := <-ch:
		fmt.Printf("  受信: %s\n", msg)
	default:
		fmt.Println("  データなし（ブロックしない）")
	}

	// ノンブロッキング送信
	ch <- "hello"
	select {
	case ch <- "world": // バッファ満杯なので送信できない
		fmt.Println("  送信成功")
	default:
		fmt.Println("  バッファ満杯（ブロックしない）")
	}
	fmt.Println()
}

// timeoutDemo: time.After でタイムアウト
func timeoutDemo() {
	fmt.Println("--- タイムアウトパターン ---")

	ch := make(chan string)

	// 遅い処理をシミュレート
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "完了"
	}()

	select {
	case result := <-ch:
		fmt.Printf("  結果: %s\n", result)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  タイムアウト！（100ms 以内に完了しなかった）")
	}
	fmt.Println()
}

// forSelectDoneDemo: for-select + done channel パターン
func forSelectDoneDemo() {
	fmt.Println("--- for-select + done channel ---")

	done := make(chan struct{})
	data := make(chan int)

	// Producer
	go func() {
		defer close(data)
		for i := range 5 {
			select {
			case data <- i:
			case <-done:
				fmt.Println("  producer: キャンセルされた")
				return
			}
		}
	}()

	// Consumer: 3つ受信したら done を閉じて停止
	received := 0
	for v := range data {
		fmt.Printf("  received: %d\n", v)
		received++
		if received >= 3 {
			close(done) // producer にキャンセルを通知
			break
		}
	}
}
