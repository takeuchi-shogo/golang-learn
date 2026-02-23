package main

import (
	"cmp"
	"fmt"
)

func main() {
	fmt.Println("=== Generics ===")
	fmt.Println()

	genericFuncDemo()
	genericTypeDemo()
	constraintDemo()
}

// ============================================================
// ジェネリック関数
// ============================================================

// Min: Ordered な型の最小値
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Map: スライスの各要素を変換
func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Filter: 条件を満たす要素を抽出
func Filter[T any](s []T, f func(T) bool) []T {
	var result []T
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Contains: comparable な型でスライスに要素が含まれるか
func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

func genericFuncDemo() {
	fmt.Println("--- ジェネリック関数 ---")

	// Min: int でも float64 でも string でも使える
	fmt.Printf("  Min(3, 7)         = %d\n", Min(3, 7))
	fmt.Printf("  Min(3.14, 2.71)   = %.2f\n", Min(3.14, 2.71))
	fmt.Printf("  Min(\"b\", \"a\")     = %q\n", Min("b", "a"))

	// Map
	nums := []int{1, 2, 3, 4, 5}
	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Printf("  Map(doubled):     %v\n", doubled)

	// Filter
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("  Filter(evens):    %v\n", evens)

	// Contains
	fmt.Printf("  Contains(3):      %t\n", Contains(nums, 3))
	fmt.Printf("  Contains(99):     %t\n", Contains(nums, 99))
	fmt.Println()
}

// ============================================================
// ジェネリック型
// ============================================================

// Stack: ジェネリックなスタック
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func genericTypeDemo() {
	fmt.Println("--- ジェネリック型（Stack）---")

	// int スタック
	intStack := &Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	for intStack.Len() > 0 {
		v, _ := intStack.Pop()
		fmt.Printf("  pop: %d\n", v)
	}

	// string スタック（同じコードが再利用される）
	strStack := &Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")
	v, _ := strStack.Pop()
	fmt.Printf("  string pop: %q\n", v)
	fmt.Println()
}

// ============================================================
// 型制約
// ============================================================

// Number: underlying type で制約
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

// カスタム型でも ~ があれば使える
type Celsius float64
type Meter int

func constraintDemo() {
	fmt.Println("--- 型制約（~ underlying type）---")

	// 基本型
	fmt.Printf("  Sum(int):     %d\n", Sum([]int{1, 2, 3}))
	fmt.Printf("  Sum(float64): %.1f\n", Sum([]float64{1.5, 2.5}))

	// カスタム型（~float64 で Celsius も受け入れる）
	temps := []Celsius{36.5, 37.0, 38.2}
	fmt.Printf("  Sum(Celsius): %.1f\n", Sum(temps))

	meters := []Meter{100, 200, 300}
	fmt.Printf("  Sum(Meter):   %d\n", Sum(meters))
}
