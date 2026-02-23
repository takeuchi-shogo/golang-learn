package main

import (
	"fmt"
	"unsafe"
)

// ============================================================
// typed nil の罠を実証する
// ============================================================

type MyError struct {
	Code    int
	Message string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

func main() {
	fmt.Println("=== Interface Internals ===")
	fmt.Println()

	typedNilDemo()
	interfaceSizeDemo()
	correctNilCheckDemo()
}

// typedNilDemo: nil pointer を interface に入れると nil にならない
func typedNilDemo() {
	fmt.Println("--- typed nil の罠 ---")

	var err *MyError = nil // nil ポインタ

	// error interface に代入
	var i error = err

	fmt.Printf("  err == nil: %t (nil ポインタ)\n", err == nil)
	fmt.Printf("  i   == nil: %t ← 危険! typed nil\n", i == nil)
	// i は (*MyError, nil) なので nil ではない

	// 正しい nil を代入
	var j error = nil // (nil, nil)
	fmt.Printf("  j   == nil: %t (真の nil interface)\n", j == nil)
	fmt.Println()
}

// interfaceSizeDemo: interface のサイズを確認
func interfaceSizeDemo() {
	fmt.Println("--- interface のサイズ ---")

	var i any
	var e error

	// interface は常に 2 word（64bit 環境で 16 bytes）
	fmt.Printf("  sizeof(any)   = %d bytes (eface: _type + data)\n", unsafe.Sizeof(i))
	fmt.Printf("  sizeof(error) = %d bytes (iface: tab + data)\n", unsafe.Sizeof(e))
	fmt.Println()
}

// correctNilCheckDemo: 正しい nil チェックの方法
func correctNilCheckDemo() {
	fmt.Println("--- 正しい nil チェック ---")

	// アンチパターン: typed nil を返す関数
	err := badGetError()
	fmt.Printf("  badGetError() == nil: %t ← バグ!\n", err == nil)

	// 正しいパターン: nil を直接返す
	err = goodGetError()
	fmt.Printf("  goodGetError() == nil: %t ← 正しい\n", err == nil)
}

// badGetError: typed nil を返してしまう（バグ）
func badGetError() error {
	var err *MyError = nil
	return err // (*MyError, nil) を返す → nil チェックが false
}

// goodGetError: 明示的に nil を返す（正しい）
func goodGetError() error {
	var err *MyError = nil
	if err == nil {
		return nil // 明示的に nil を返す → (nil, nil)
	}
	return err
}
