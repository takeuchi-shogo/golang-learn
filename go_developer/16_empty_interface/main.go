package main

import "fmt"

func main() {
	fmt.Println("=== Empty Interface (any) ===")
	fmt.Println()

	typeAssertionDemo()
	typeSwitchDemo()
	genericVsAnyDemo()
}

// typeAssertionDemo: type assertion の使い方
func typeAssertionDemo() {
	fmt.Println("--- Type Assertion ---")

	var i any = "hello"

	// 安全な type assertion（ok パターン）
	s, ok := i.(string)
	fmt.Printf("  string: %q, ok=%t\n", s, ok)

	n, ok := i.(int)
	fmt.Printf("  int:    %d, ok=%t (型が違う)\n", n, ok)

	// 危険な type assertion（panic する）
	// _ = i.(int) // ← panic: interface conversion

	fmt.Println()
}

// typeSwitchDemo: type switch で型ごとに処理
func typeSwitchDemo() {
	fmt.Println("--- Type Switch ---")

	values := []any{42, "hello", 3.14, true, nil}

	for _, v := range values {
		switch val := v.(type) {
		case int:
			fmt.Printf("  int:     %d\n", val)
		case string:
			fmt.Printf("  string:  %q\n", val)
		case float64:
			fmt.Printf("  float64: %f\n", val)
		case bool:
			fmt.Printf("  bool:    %t\n", val)
		case nil:
			fmt.Println("  nil")
		default:
			fmt.Printf("  unknown: %T\n", val)
		}
	}
	fmt.Println()
}

// ============================================================
// any vs generics の比較
// ============================================================

// anySum: any を使った汎用関数（型安全でない、ヒープエスケープする）
func anySum(values []any) float64 {
	var sum float64
	for _, v := range values {
		switch n := v.(type) {
		case int:
			sum += float64(n)
		case float64:
			sum += n
		}
	}
	return sum
}

// genericSum: generics を使った汎用関数（型安全、スタックに留まりうる）
type Number interface {
	~int | ~float64
}

func genericSum[T Number](values []T) T {
	var sum T
	for _, v := range values {
		sum += v
	}
	return sum
}

func genericVsAnyDemo() {
	fmt.Println("--- any vs generics ---")

	// any: 型安全でない
	result1 := anySum([]any{1, 2.5, 3})
	fmt.Printf("  anySum:     %.1f (any → 型チェック必要、ヒープエスケープ)\n", result1)

	// generics: 型安全
	result2 := genericSum([]int{1, 2, 3})
	fmt.Printf("  genericSum: %d   (generics → コンパイル時チェック、高速)\n", result2)

	result3 := genericSum([]float64{1.0, 2.5, 3.0})
	fmt.Printf("  genericSum: %.1f (float64 版も自動生成)\n", result3)
}
