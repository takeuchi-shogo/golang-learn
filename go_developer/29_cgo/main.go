package main

/*
#include <stdlib.h>
#include <math.h>
*/
import "C"

import (
	"fmt"
	"math"
	"time"
)

// ============================================================
// cgo: Go から C を呼ぶ
// ============================================================
//
// 注意: import "C" は C コメントブロックの直後に書く。
//       間に空行を入れるとコンパイルエラーになる。
//
// ビルド: CGO_ENABLED=1 go run ./29_cgo/

func main() {
	fmt.Println("=== cgo ===")
	fmt.Println()

	basicCgoDemo()
	performanceDemo()
}

// basicCgoDemo: C の関数を呼ぶ基本例
func basicCgoDemo() {
	fmt.Println("--- 基本的な cgo 呼び出し ---")

	// C.sqrt を呼ぶ（math.h の sqrt）
	result := C.sqrt(C.double(144.0))
	fmt.Printf("  C.sqrt(144) = %f\n", float64(result))

	// C.abs を呼ぶ
	absResult := C.abs(C.int(-42))
	fmt.Printf("  C.abs(-42)  = %d\n", int(absResult))
	fmt.Println()
}

// performanceDemo: Go 関数 vs cgo のコスト比較
func performanceDemo() {
	fmt.Println("--- パフォーマンス比較 ---")

	const n = 1_000_000

	// Go の math.Sqrt
	start := time.Now()
	for i := range n {
		_ = math.Sqrt(float64(i))
	}
	goTime := time.Since(start)

	// C の sqrt（cgo 経由）
	start = time.Now()
	for i := range n {
		_ = C.sqrt(C.double(i))
	}
	cgoTime := time.Since(start)

	fmt.Printf("  Go  math.Sqrt ×%d: %v\n", n, goTime)
	fmt.Printf("  cgo C.sqrt   ×%d: %v\n", n, cgoTime)
	fmt.Printf("  cgo / Go = %.1fx 遅い\n", float64(cgoTime)/float64(goTime))
	fmt.Println()
	fmt.Println("  教訓:")
	fmt.Println("    - cgo は境界越えコストが大きい (~100-200ns/call)")
	fmt.Println("    - 高頻度呼び出しは避ける（バッチ化を検討）")
	fmt.Println("    - Pure Go の代替があればそちらを使う")
}
