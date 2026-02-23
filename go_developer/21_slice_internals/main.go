package main

import "fmt"

func main() {
	fmt.Println("=== Slice Internals ===")
	fmt.Println()

	appendGrowthDemo()
	sharedBackingArrayDemo()
	fullSliceExprDemo()
	copyDemo()
	nilVsEmptyDemo()
}

// appendGrowthDemo: append の成長戦略を追跡
func appendGrowthDemo() {
	fmt.Println("--- append の成長戦略 ---")

	var s []int
	prevCap := 0
	for i := range 20 {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("  len=%2d, cap=%2d (成長!)\n", len(s), cap(s))
			prevCap = cap(s)
		}
	}
	fmt.Println()
}

// sharedBackingArrayDemo: スライシングは配列を共有する
func sharedBackingArrayDemo() {
	fmt.Println("--- バッキング配列の共有（罠）---")

	a := []int{1, 2, 3, 4, 5}
	b := a[1:3] // b = [2, 3], a と同じバッキング配列

	fmt.Printf("  a = %v\n", a)
	fmt.Printf("  b = a[1:3] = %v\n", b)

	b[0] = 99 // a[1] も変わる!
	fmt.Printf("  b[0] = 99 → a = %v (a[1]も変わった!)\n", a)

	// append で cap を超えなければ元の配列を書き換える
	c := a[0:2]    // c = [1, 99], cap = 5
	c = append(c, 888)
	fmt.Printf("  c = a[0:2]; append(c, 888) → a = %v (a[2]が書き換わった!)\n", a)
	fmt.Println()
}

// fullSliceExprDemo: s[low:high:max] で cap を制限
func fullSliceExprDemo() {
	fmt.Println("--- Full Slice Expression ---")

	a := []int{1, 2, 3, 4, 5}
	b := a[1:3:3] // len=2, cap=2（cap を制限）
	fmt.Printf("  a[1:3:3] → len=%d, cap=%d\n", len(b), cap(b))

	// append すると cap を超えるので新しい配列が作られる
	b = append(b, 999)
	fmt.Printf("  append 後: a = %v (元の配列は変わらない)\n", a)
	fmt.Printf("  append 後: b = %v (新しい配列)\n", b)
	fmt.Println()
}

// copyDemo: copy で独立したスライスを作る
func copyDemo() {
	fmt.Println("--- copy ---")

	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	n := copy(dst, src)

	dst[0] = 99
	fmt.Printf("  copy: %d 要素コピー\n", n)
	fmt.Printf("  src = %v (変わらない)\n", src)
	fmt.Printf("  dst = %v (独立)\n", dst)
	fmt.Println()
}

// nilVsEmptyDemo: nil slice と empty slice の違い
func nilVsEmptyDemo() {
	fmt.Println("--- nil slice vs empty slice ---")

	var nilSlice []int
	emptySlice := []int{}

	fmt.Printf("  nil slice:   len=%d, cap=%d, nil=%t\n",
		len(nilSlice), cap(nilSlice), nilSlice == nil)
	fmt.Printf("  empty slice: len=%d, cap=%d, nil=%t\n",
		len(emptySlice), cap(emptySlice), emptySlice == nil)

	// 両方とも append 可能
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)
	fmt.Printf("  append 後: nil=%v, empty=%v\n", nilSlice, emptySlice)
}
