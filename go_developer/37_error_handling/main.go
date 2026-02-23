package main

import (
	"errors"
	"fmt"
)

// ============================================================
// Sentinel Errors（パッケージレベルで定義）
// ============================================================

var (
	ErrNotFound    = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// ============================================================
// Custom Error Type
// ============================================================

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

func main() {
	fmt.Println("=== Error Handling ===")
	fmt.Println()

	wrapDemo()
	errorsIsDemo()
	errorsAsDemo()
	errorsJoinDemo()
}

// wrapDemo: %w でエラーをラップ
func wrapDemo() {
	fmt.Println("--- %w でエラーラッピング ---")

	err := findUser(999)
	fmt.Printf("  エラー: %v\n", err)
	// 出力: user 999: query failed: not found
	// → 各層でコンテキストが追加されている
	fmt.Println()
}

func findUser(id int) error {
	err := queryDB(id)
	if err != nil {
		return fmt.Errorf("user %d: %w", id, err) // コンテキストを追加して wrap
	}
	return nil
}

func queryDB(id int) error {
	return fmt.Errorf("query failed: %w", ErrNotFound) // さらに wrap
}

// errorsIsDemo: errors.Is でチェインを遡る
func errorsIsDemo() {
	fmt.Println("--- errors.Is（値の比較）---")

	err := findUser(999)

	// チェインのどこかに ErrNotFound があるか
	fmt.Printf("  Is ErrNotFound:    %t\n", errors.Is(err, ErrNotFound))
	fmt.Printf("  Is ErrUnauthorized: %t\n", errors.Is(err, ErrUnauthorized))

	// == で比較するとチェインを遡れない
	fmt.Printf("  err == ErrNotFound: %t (直接比較は false)\n", err == ErrNotFound)
	fmt.Println()
}

// errorsAsDemo: errors.As で型を抽出
func errorsAsDemo() {
	fmt.Println("--- errors.As（型の比較）---")

	err := validateInput()

	var valErr *ValidationError
	if errors.As(err, &valErr) {
		fmt.Printf("  Field:   %s\n", valErr.Field)
		fmt.Printf("  Message: %s\n", valErr.Message)
	}
	fmt.Println()
}

func validateInput() error {
	err := &ValidationError{Field: "email", Message: "invalid format"}
	return fmt.Errorf("input validation: %w", err)
}

// errorsJoinDemo: errors.Join で複数エラーを結合 (Go 1.20+)
func errorsJoinDemo() {
	fmt.Println("--- errors.Join (Go 1.20+) ---")

	err1 := errors.New("file not found")
	err2 := errors.New("permission denied")
	err3 := ErrNotFound

	joined := errors.Join(err1, err2, err3)
	fmt.Printf("  joined: %v\n", joined)

	// 各エラーに対して errors.Is が使える
	fmt.Printf("  Is err1:        %t\n", errors.Is(joined, err1))
	fmt.Printf("  Is err2:        %t\n", errors.Is(joined, err2))
	fmt.Printf("  Is ErrNotFound: %t\n", errors.Is(joined, ErrNotFound))
}
