package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *Counter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ミューテックスを使用して、カウンターをインクリメントする
// mutex がない場合、Goroutine が同時にカウンターをインクリメントしようとすると、競合が発生して、正しいカウントが得られない
// Goroutine は、sync.WaitGroup を使用して、すべての Goroutine が完了するのを待ち、最終的なカウントを取得するようにする必要がある
func main() {
	counter := Counter{}
	wg := sync.WaitGroup{} // WaitGroupを宣言
	for i := 0; i < 20; i++ {
		wg.Add(1) // WaitGroupに1を追加
		go func(i int) {
			counter.Increment()
			wg.Done() // WaitGroupから1を減算
		}(i)
	}
	wg.Wait() // WaitGroupが0になるまで待機
	fmt.Println("Done", counter.GetCount())
}
