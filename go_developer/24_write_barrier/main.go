package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Write Barrier ===")
	fmt.Println()

	pointerVsValueDemo()
	gcPressureDemo()
}

// ============================================================
// ポインタなし vs ポインタあり構造体
// ============================================================

// NoPointer: ポインタフィールドなし → GC スキャン不要
type NoPointer struct {
	X, Y, Z int
	Data     [64]byte
}

// WithPointer: ポインタフィールドあり → write barrier + GC スキャン必要
type WithPointer struct {
	X, Y, Z int
	Data     *[64]byte // ← ポインタ: write barrier の対象
	Name     *string   // ← ポインタ
}

// pointerVsValueDemo: ポインタの有無による GC 負荷の違い
func pointerVsValueDemo() {
	fmt.Println("--- ポインタの有無と GC ---")

	fmt.Println("  NoPointer 構造体: ポインタなし → GC スキャン対象外")
	fmt.Println("  WithPointer 構造体: ポインタあり → GC スキャン対象")
	fmt.Println()
	fmt.Println("  設計指針:")
	fmt.Println("    - 可能なら値型を使う（*T より T）")
	fmt.Println("    - []byte > []*byte（ポインタの配列は GC 負荷大）")
	fmt.Println("    - map のキー/値がポインタだと GC 負荷が上がる")
	fmt.Println()
}

// gcPressureDemo: ポインタの多い構造体と少ない構造体の GC 比較
func gcPressureDemo() {
	fmt.Println("--- GC 圧力の比較 ---")

	const n = 100000

	// ポインタなし: GC スキャン不要
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	noPtr := make([]NoPointer, n)
	_ = noPtr

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("  NoPointer  ×%d: GC pause=%v\n", n, m2.PauseNs[(m2.NumGC+255)%256])

	// ポインタあり: write barrier + GC スキャン必要
	runtime.GC()
	runtime.ReadMemStats(&m1)

	withPtr := make([]WithPointer, n)
	for i := range withPtr {
		data := [64]byte{}
		name := "test"
		withPtr[i].Data = &data
		withPtr[i].Name = &name
	}
	_ = withPtr

	runtime.GC()
	runtime.ReadMemStats(&m2)
	fmt.Printf("  WithPointer×%d: GC pause=%v\n", n, m2.PauseNs[(m2.NumGC+255)%256])
	fmt.Println()
	fmt.Println("  ※ ポインタが多いほど GC のスキャン対象が増え、pause が長くなりうる")
}
