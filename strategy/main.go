package main

import (
	"fmt"
	"log"
	"strings"
)

// strategy は関数を引数として受け取り、新しい関数を返す高階関数
func strategy1(f func(string) string) func(string) string {
	// クロージャを返す: 引数fをキャプチャして保持
	return func(s string) string {
		// 受け取った文字列sを、キャプチャした関数fに渡して実行
		return f(s)
	}
}

func strategy2(id int, f func(int) string) (string, error) {
	return f(id), nil
}

func stringer(id int) string {
	return fmt.Sprintf("ID: %d", id+10000)
}

func main() {
	// データの流れ:
	// 1. strings.ToUpper関数を strategy に渡す
	// 2. strategy は strings.ToUpper をキャプチャしたクロージャを返す
	// 3. そのクロージャを uppercase 変数に代入
	uppercase := strategy1(strings.ToUpper)

	// 4. uppercase("hello") を呼び出す
	// 5. クロージャ内で "hello" が strings.ToUpper に渡される
	// 6. strings.ToUpper("hello") が実行され "HELLO" を返す
	// 7. "HELLO" が出力される
	fmt.Println(uppercase("hello"))

	// strategy2 の実行: 整数IDを文字列に変換
	result, err := strategy2(1, stringer)
	if err != nil {
		log.Fatalf("failed to execute strategy2: %v", err)
	}
	fmt.Println(result)
}
