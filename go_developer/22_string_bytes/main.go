package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== String / []byte ===")
	fmt.Println()

	immutabilityDemo()
	conversionDemo()
	builderDemo()
	runeDemo()
}

// immutabilityDemo: string は不変
func immutabilityDemo() {
	fmt.Println("--- string は不変 ---")

	s := "hello"
	// s[0] = 'H' // ← コンパイルエラー! string は変更できない

	// 変更するには []byte に変換（コピー発生）
	b := []byte(s)
	b[0] = 'H'
	s2 := string(b)
	fmt.Printf("  元: %q, 変更後: %q\n", s, s2)
	fmt.Println()
}

// conversionDemo: string ↔ []byte の変換コスト
func conversionDemo() {
	fmt.Println("--- 変換はコピーが発生する ---")

	s := "hello"
	b := []byte(s) // コピー
	b[0] = 'X'

	fmt.Printf("  string: %q (変わらない)\n", s)
	fmt.Printf("  []byte: %q (変更された)\n", string(b))
	// → コピーなので元の string は影響を受けない
	fmt.Println()
}

// builderDemo: strings.Builder で効率的に連結
func builderDemo() {
	fmt.Println("--- strings.Builder vs + 連結 ---")

	// 悪い例: + 連結（毎回新しい string を作る = O(n²)）
	result := ""
	for range 5 {
		result += "hello" // 毎回コピーが発生
	}
	fmt.Printf("  + 連結:         %q\n", result)

	// 良い例: strings.Builder（内部バッファに追記 = O(n)）
	var b strings.Builder
	for range 5 {
		b.WriteString("hello") // バッファに追記（コピーなし）
	}
	fmt.Printf("  strings.Builder: %q\n", b.String())

	// strings.Join も効率的
	parts := []string{"a", "b", "c", "d"}
	fmt.Printf("  strings.Join:    %q\n", strings.Join(parts, "-"))
	fmt.Println()
}

// runeDemo: rune（Unicode コードポイント）とバイトの違い
func runeDemo() {
	fmt.Println("--- rune vs byte ---")

	s := "Hello, 世界!"

	fmt.Printf("  string:  %q\n", s)
	fmt.Printf("  len():   %d (バイト数)\n", len(s))
	fmt.Printf("  runes:   %d (文字数)\n", len([]rune(s)))

	// range は rune 単位でイテレート
	fmt.Print("  range:   ")
	for i, r := range s {
		if r > 127 {
			fmt.Printf("[%d]=%c(U+%04X) ", i, r, r)
		}
	}
	fmt.Println()

	// s[i] は byte 単位
	fmt.Printf("  s[0]=%c (byte), s[7]=0x%02X (マルチバイトの1バイト目)\n", s[0], s[7])
}
