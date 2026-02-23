package main

import "fmt"

func main() {
	fmt.Println("=== Map Internals ===")
	fmt.Println()

	basicOpsDemo()
	iterationOrderDemo()
	nilMapDemo()
	deleteDemo()
}

// basicOpsDemo: map の基本操作
func basicOpsDemo() {
	fmt.Println("--- 基本操作 ---")

	// make vs リテラル
	m1 := make(map[string]int)    // 空 map
	m2 := map[string]int{"a": 1}  // リテラル初期化

	m1["x"] = 10
	fmt.Printf("  make: %v\n", m1)
	fmt.Printf("  literal: %v\n", m2)

	// 存在チェック（ok idiom）
	v, ok := m1["x"]
	fmt.Printf("  m1[\"x\"] = %d, ok=%t\n", v, ok)
	v, ok = m1["y"]
	fmt.Printf("  m1[\"y\"] = %d, ok=%t (ゼロ値)\n", v, ok)
	fmt.Println()
}

// iterationOrderDemo: イテレーション順序がランダム
func iterationOrderDemo() {
	fmt.Println("--- イテレーション順序（意図的にランダム）---")

	m := map[string]int{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
	}

	// 3回繰り返して順序が変わることを確認
	for trial := range 3 {
		fmt.Printf("  試行 %d: ", trial)
		for k := range m {
			fmt.Printf("%s ", k)
		}
		fmt.Println()
	}
	fmt.Println()
}

// nilMapDemo: nil map の挙動
func nilMapDemo() {
	fmt.Println("--- nil map ---")

	var m map[string]int // nil map

	// 読み取り: OK（ゼロ値を返す）
	v := m["key"]
	fmt.Printf("  nil map の読み取り: %d (ゼロ値)\n", v)

	// len: OK
	fmt.Printf("  nil map の len: %d\n", len(m))

	// 書き込み: panic!
	// m["key"] = 1 // ← panic: assignment to entry in nil map
	fmt.Println("  nil map への書き込み: panic! (必ず make が必要)")
	fmt.Println()
}

// deleteDemo: delete の挙動
func deleteDemo() {
	fmt.Println("--- delete ---")

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("  削除前: %v\n", m)

	delete(m, "b")
	fmt.Printf("  delete(m, \"b\"): %v\n", m)

	// 存在しないキーの delete: 何も起きない（panic しない）
	delete(m, "z")
	fmt.Printf("  delete(m, \"z\"): %v (何も起きない)\n", m)
}
