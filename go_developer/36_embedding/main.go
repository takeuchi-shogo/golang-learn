package main

import "fmt"

// ============================================================
// 構造体の埋め込み
// ============================================================

type Base struct {
	ID int
}

func (b Base) GetID() int { return b.ID }
func (b Base) Type() string { return "Base" }

type User struct {
	Base // 埋め込み（has-a, NOT is-a）
	Name string
}

// User が Type() をオーバーライド（promotion より優先される）
func (u User) Type() string { return "User" }

func main() {
	fmt.Println("=== Embedding ===")
	fmt.Println()

	promotionDemo()
	overrideDemo()
	ambiguityDemo()
	interfaceEmbedDemo()
	zeroValueDemo()
}

// promotionDemo: メソッドとフィールドの promotion
func promotionDemo() {
	fmt.Println("--- Promotion（昇格）---")

	u := User{
		Base: Base{ID: 42},
		Name: "Alice",
	}

	// Base のフィールドに直接アクセス
	fmt.Printf("  u.ID:      %d (Base.ID が promoted)\n", u.ID)
	fmt.Printf("  u.Base.ID: %d (明示的アクセスも可能)\n", u.Base.ID)

	// Base のメソッドを直接呼べる
	fmt.Printf("  u.GetID(): %d (Base.GetID() が promoted)\n", u.GetID())
	fmt.Println()
}

// overrideDemo: 外側の型が同名メソッドを定義すると優先される
func overrideDemo() {
	fmt.Println("--- オーバーライド ---")

	u := User{Base: Base{ID: 1}, Name: "Alice"}

	fmt.Printf("  u.Type():      %q (User.Type() が優先)\n", u.Type())
	fmt.Printf("  u.Base.Type(): %q (Base.Type() も呼べる)\n", u.Base.Type())
	fmt.Println()
}

// ambiguityDemo: 複数の埋め込みが同名メソッドを持つ
func ambiguityDemo() {
	fmt.Println("--- Ambiguity（曖昧性）---")

	type A struct{}
	type B struct{}

	type C struct {
		A
		B
	}

	// A と B の両方にメソッドがない場合は問題なし
	// 同名メソッドがあるとコンパイルエラー:
	//   c.Hello() → ambiguous selector c.Hello
	// 解決: c.A.Hello() と明示的に指定
	fmt.Println("  2つの埋め込みが同名メソッドを持つ → コンパイルエラー")
	fmt.Println("  解決: c.A.Hello() と明示的に型を指定")
	fmt.Println()
}

// interfaceEmbedDemo: interface の埋め込み
func interfaceEmbedDemo() {
	fmt.Println("--- Interface の埋め込み ---")

	type Reader interface {
		Read() string
	}

	type Writer interface {
		Write(s string)
	}

	// interface を埋め込んで合成
	type ReadWriter interface {
		Reader
		Writer
	}

	fmt.Println("  ReadWriter = Reader + Writer（interface の合成）")
	fmt.Println()
}

// zeroValueDemo: 埋め込みフィールドのゼロ値
func zeroValueDemo() {
	fmt.Println("--- ゼロ値 ---")

	var u User // ゼロ値
	fmt.Printf("  User ゼロ値: ID=%d, Name=%q\n", u.ID, u.Name)
	fmt.Println("  ※ 埋め込みフィールドもゼロ値で初期化される")
}
