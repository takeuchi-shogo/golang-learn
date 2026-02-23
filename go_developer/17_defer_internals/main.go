package main

import "fmt"

func main() {
	fmt.Println("=== defer Internals ===")
	fmt.Println()

	lifoDemo()
	immediateEvalDemo()
	namedReturnDemo()
	loopDeferDemo()
}

// lifoDemo: LIFO（後入れ先出し）
func lifoDemo() {
	fmt.Println("--- LIFO 順序 ---")
	fmt.Print("  ")
	for i := range 5 {
		defer fmt.Printf("%d ", i) // main 終了時に 4,3,2,1,0 の順で実行
	}
	// ※ この defer は main 関数の終了時に実行される
	// 以降のデモの後に出力される
	fmt.Println()
}

// immediateEvalDemo: 引数は defer 文の時点で評価される
func immediateEvalDemo() {
	fmt.Println("--- 引数の即時評価 ---")

	x := 10
	defer fmt.Printf("  defer 時点の x: %d\n", x) // x=10 が記録される
	x = 20
	fmt.Printf("  現在の x: %d\n", x)
	// 出力: 現在の x: 20 → defer 時点の x: 10
}

// namedReturnDemo: 名前付き戻り値を defer で変更
func namedReturnDemo() {
	fmt.Println("--- 名前付き戻り値の変更 ---")

	result := doubleReturn()
	fmt.Printf("  doubleReturn() = %d (defer が戻り値を変更)\n", result)
	fmt.Println()
}

func doubleReturn() (result int) {
	defer func() {
		result *= 2 // 戻り値を変更できる
	}()
	return 10 // result = 10 → defer で result = 20
}

// loopDeferDemo: ループ内 defer の問題と対策
func loopDeferDemo() {
	fmt.Println("--- ループ内 defer ---")

	// 悪い例: ループ内で defer（関数終了まで解放されない）
	fmt.Println("  悪い例: ループ内 defer（全イテレーション分がスタックに溜まる）")
	// for i := range 10000 {
	//     f, _ := os.Open(files[i])
	//     defer f.Close() // ← 関数終了まで全ファイルが開いたまま!
	// }

	// 良い例: 関数に切り出す（各イテレーションで defer が実行される）
	fmt.Println("  良い例: 関数に切り出す")
	for i := range 3 {
		processItem(i)
	}
}

func processItem(i int) {
	// defer はこの関数の終了時に実行される
	defer fmt.Printf("  cleanup item %d\n", i)
	fmt.Printf("  process item %d\n", i)
}
