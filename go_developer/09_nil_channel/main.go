package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== nil Channel ===")
	fmt.Println()

	selectDisableDemo()
	mergeDemo()
}

// selectDisableDemo: nil channel で select case を無効化
func selectDisableDemo() {
	fmt.Println("--- select case の動的無効化 ---")

	ch1 := make(chan string, 2)
	ch2 := make(chan string, 2)

	ch1 <- "a1"
	ch1 <- "a2"
	close(ch1)

	ch2 <- "b1"
	close(ch2)

	var c1, c2 <-chan string = ch1, ch2

	for c1 != nil || c2 != nil {
		select {
		case v, ok := <-c1:
			if !ok {
				fmt.Println("  ch1 closed → nil に差し替え")
				c1 = nil // この case はもう評価されない
				continue
			}
			fmt.Printf("  ch1: %s\n", v)
		case v, ok := <-c2:
			if !ok {
				fmt.Println("  ch2 closed → nil に差し替え")
				c2 = nil
				continue
			}
			fmt.Printf("  ch2: %s\n", v)
		}
	}
	fmt.Println("  両方 close → ループ終了")
	fmt.Println()
}

// mergeDemo: 2つの channel を1つにマージ
func mergeDemo() {
	fmt.Println("--- Channel Merge パターン ---")

	ch1 := gen("fast", 30*time.Millisecond, 3)
	ch2 := gen("slow", 50*time.Millisecond, 3)

	// merge: 両方が close されるまで読む
	for v := range merge(ch1, ch2) {
		fmt.Printf("  merged: %s\n", v)
	}
}

// gen: name を count 回送信して close する channel を返す
func gen(name string, interval time.Duration, count int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := range count {
			time.Sleep(interval)
			ch <- fmt.Sprintf("%s-%d", name, i)
		}
	}()
	return ch
}

// merge: 複数 channel を1つにマージ（nil channel テクニック使用）
func merge(a, b <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
