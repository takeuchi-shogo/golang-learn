package main

import "fmt"

// ============================================================
// Point: value receiver と pointer receiver の違いを示す
// ============================================================

type Point struct {
	X, Y int
}

// Value receiver: コピーを操作（元は変わらない）
func (p Point) ScaleValue(factor int) Point {
	p.X *= factor
	p.Y *= factor
	return p // 新しい値を返す必要がある
}

// Pointer receiver: 元を直接変更
func (p *Point) ScalePointer(factor int) {
	p.X *= factor
	p.Y *= factor
}

// ============================================================
// interface 満足とメソッドセット
// ============================================================

type Stringer interface {
	String() string
}

type User struct {
	Name string
}

// Pointer receiver で定義
func (u *User) String() string {
	if u == nil {
		return "<nil user>"
	}
	return u.Name
}

func main() {
	fmt.Println("=== Pointer vs Value Receiver ===")
	fmt.Println()

	receiverDiffDemo()
	methodSetDemo()
	nilReceiverDemo()
}

func receiverDiffDemo() {
	fmt.Println("--- Value vs Pointer Receiver ---")

	p := Point{X: 1, Y: 2}

	// Value receiver: 元は変わらない
	p2 := p.ScaleValue(10)
	fmt.Printf("  Value receiver: 元=%v, 新=%v\n", p, p2)

	// Pointer receiver: 元が変わる
	p.ScalePointer(10)
	fmt.Printf("  Pointer receiver: 元=%v (変更された)\n", p)
	fmt.Println()
}

func methodSetDemo() {
	fmt.Println("--- メソッドセットと interface ---")

	// *User は *User のメソッド(String)を持つ → Stringer を満たす
	var s Stringer = &User{Name: "Alice"}
	fmt.Printf("  &User → Stringer: %s\n", s.String())

	// User は *User のメソッド(String)を持たない → Stringer を満たさない
	// var s2 Stringer = User{Name: "Bob"} // ← コンパイルエラー!
	// cannot use User{} (value of type User) as Stringer value in variable declaration:
	//   User does not implement Stringer (method String has pointer receiver)
	fmt.Println("  User → Stringer: コンパイルエラー（pointer receiver のため）")
	fmt.Println()
}

// nilReceiverDemo: nil ポインタでもメソッド呼び出し可能
func nilReceiverDemo() {
	fmt.Println("--- nil receiver ---")

	var u *User // nil
	// Pointer receiver は nil で呼べる（メソッド内で nil チェック可能）
	fmt.Printf("  nil *User: %s\n", u.String())

	// Value receiver は nil で呼ぶと panic (nil pointer dereference)
	// なぜなら値のコピーを作れないから
}
