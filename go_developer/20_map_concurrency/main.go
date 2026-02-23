package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== Map Concurrency ===")
	fmt.Println()

	mutexMapDemo()
	rwMutexMapDemo()
	syncMapDemo()
}

// mutexMapDemo: sync.Mutex で保護
func mutexMapDemo() {
	fmt.Println("--- sync.Mutex で保護 ---")

	var mu sync.Mutex
	m := make(map[string]int)
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id%10)
			mu.Lock()
			m[key] = id
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("  map size: %d\n", len(m))
	fmt.Println()
}

// rwMutexMapDemo: sync.RWMutex で読み取り並行性を改善
func rwMutexMapDemo() {
	fmt.Println("--- sync.RWMutex（読み取り並行性）---")

	var mu sync.RWMutex
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	var wg sync.WaitGroup

	// 読み取り: RLock（複数同時に可能）
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.RLock() // 読み取りロック
			_ = m["a"]
			mu.RUnlock()
		}(i)
	}

	// 書き込み: Lock（排他）
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock() // 書き込みロック
		m["d"] = 4
		mu.Unlock()
	}()

	wg.Wait()
	fmt.Printf("  map: %v\n", m)
	fmt.Println()
}

// syncMapDemo: sync.Map の基本操作
func syncMapDemo() {
	fmt.Println("--- sync.Map ---")

	var sm sync.Map

	// Store: 書き込み
	sm.Store("name", "Alice")
	sm.Store("age", 30)

	// Load: 読み取り
	if v, ok := sm.Load("name"); ok {
		fmt.Printf("  Load(\"name\"): %v\n", v)
	}

	// LoadOrStore: 存在しなければ保存
	actual, loaded := sm.LoadOrStore("name", "Bob")
	fmt.Printf("  LoadOrStore(\"name\", \"Bob\"): %v, loaded=%t (既存値を返す)\n", actual, loaded)

	// Delete
	sm.Delete("age")

	// Range: 全要素を走査
	fmt.Print("  Range: ")
	sm.Range(func(key, value any) bool {
		fmt.Printf("%v=%v ", key, value)
		return true // false で走査停止
	})
	fmt.Println()

	// 並行アクセスの例
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sm.Store(fmt.Sprintf("key-%d", id), id)
		}(i)
	}
	wg.Wait()
	fmt.Println("  100 goroutine からの同時書き込み: 成功（ロック不要）")
}
