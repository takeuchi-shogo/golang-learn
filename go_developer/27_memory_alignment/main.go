package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

func main() {
	fmt.Println("=== Memory Alignment ===")
	fmt.Println()

	paddingDemo()
	falseSharingDemo()
}

// ============================================================
// フィールド順序でサイズが変わる
// ============================================================

// Bad: パディングが多い（24 bytes）
type BadOrder struct {
	a bool   // 1 byte + 7 padding
	b int64  // 8 bytes
	c bool   // 1 byte + 7 padding
}

// Good: パディングが少ない（16 bytes）
type GoodOrder struct {
	b int64  // 8 bytes
	a bool   // 1 byte
	c bool   // 1 byte + 6 padding
}

func paddingDemo() {
	fmt.Println("--- フィールド順序とパディング ---")
	fmt.Printf("  BadOrder:  %d bytes (bool, int64, bool)\n", unsafe.Sizeof(BadOrder{}))
	fmt.Printf("  GoodOrder: %d bytes (int64, bool, bool)\n", unsafe.Sizeof(GoodOrder{}))
	fmt.Println()

	// Alignof でアライメント要件を確認
	var x int64
	var y bool
	fmt.Printf("  Alignof(int64) = %d\n", unsafe.Alignof(x))
	fmt.Printf("  Alignof(bool)  = %d\n", unsafe.Alignof(y))
	fmt.Println()
}

// ============================================================
// False Sharing
// ============================================================

// BadCounters: false sharing が発生（同じキャッシュラインに隣接）
type BadCounters struct {
	a atomic.Int64
	b atomic.Int64
}

// GoodCounters: パディングで別キャッシュラインに分離
type GoodCounters struct {
	a atomic.Int64
	_ [56]byte // 64 - 8 = 56 bytes のパディング
	b atomic.Int64
}

func falseSharingDemo() {
	fmt.Println("--- False Sharing ---")

	fmt.Printf("  BadCounters:  %d bytes (隣接 = false sharing)\n", unsafe.Sizeof(BadCounters{}))
	fmt.Printf("  GoodCounters: %d bytes (64B パディング = 分離)\n", unsafe.Sizeof(GoodCounters{}))
	fmt.Println()

	// false sharing の影響を簡易ベンチマーク
	const n = 10_000_000

	bad := &BadCounters{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); for range n { bad.a.Add(1) } }()
	go func() { defer wg.Done(); for range n { bad.b.Add(1) } }()
	wg.Wait()
	fmt.Printf("  Bad:  a=%d, b=%d\n", bad.a.Load(), bad.b.Load())

	good := &GoodCounters{}
	wg.Add(2)
	go func() { defer wg.Done(); for range n { good.a.Add(1) } }()
	go func() { defer wg.Done(); for range n { good.b.Add(1) } }()
	wg.Wait()
	fmt.Printf("  Good: a=%d, b=%d\n", good.a.Load(), good.b.Load())
	fmt.Println("  ※ GoodCounters は CPU キャッシュの競合が減り、高スループット")
}
