package main

import "fmt"

func main() {
	fmt.Println("=== Escape Analysis ===")
	fmt.Println("※ go build -gcflags=\"-m\" ./13_escape_analysis/ で確認")
	fmt.Println()

	// 1. スタックに留まる例
	stackExample()

	// 2. ヒープにエスケープする例
	p := heapExample()
	fmt.Printf("ヒープのポインタ: %v\n", *p)

	// 3. interface{} でエスケープ
	interfaceEscape()

	// 4. クロージャでエスケープ
	fn := closureEscape()
	fmt.Printf("クロージャ: %d\n", fn())
}

// stackExample: 変数がスタックに留まる
// x は関数外に出ないのでスタック割り当て
func stackExample() {
	x := 42        // ← スタック（エスケープしない）
	y := x + 1     // ← スタック
	fmt.Println("スタック:", y)
}

// heapExample: ポインタを返すとヒープにエスケープ
func heapExample() *int {
	x := 100       // ← ヒープにエスケープ（ポインタが関数外に出る）
	return &x      // &x を返す → x は関数終了後も生きている必要がある
}

// interfaceEscape: interface{} に代入するとエスケープ
func interfaceEscape() {
	x := 42
	// any (interface{}) に代入 → x がヒープにエスケープ
	// interface は (type, pointer) のペアなので、値のアドレスが必要
	var i any = x
	fmt.Println("interface{}:", i)
}

// closureEscape: クロージャが変数をキャプチャするとエスケープ
func closureEscape() func() int {
	count := 0     // ← ヒープにエスケープ（クロージャが参照を保持）
	return func() int {
		count++    // count はクロージャの寿命まで生きる必要がある
		return count
	}
}

// noEscape: コンパイラが証明できればスタックに留まる
//
//go:noinline
func noEscape(x int) int {
	// ポインタを取るが、関数内で完結する
	p := &x
	return *p + 1  // p は関数外に出ないのでスタック
}
