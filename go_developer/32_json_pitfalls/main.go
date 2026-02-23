package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== JSON Encoding Pitfalls ===")
	fmt.Println()

	float64ProblemDemo()
	jsonNumberDemo()
	omitemptyDemo()
	decoderVsUnmarshalDemo()
}

// float64ProblemDemo: interface{} にデコードすると数値が float64 になる
func float64ProblemDemo() {
	fmt.Println("--- float64 問題 ---")

	data := `{"id": 123456789012345678, "name": "test"}`

	// interface{} にデコード → 数値は全て float64
	var result map[string]any
	json.Unmarshal([]byte(data), &result)

	id := result["id"]
	fmt.Printf("  型:   %T\n", id)
	fmt.Printf("  値:   %.0f\n", id)
	fmt.Printf("  元:   123456789012345678\n")
	fmt.Printf("  差異: %.0f ← 精度が失われた!\n", id.(float64)-123456789012345678)
	fmt.Println()
}

// jsonNumberDemo: json.Number で精度を保つ
func jsonNumberDemo() {
	fmt.Println("--- json.Number で解決 ---")

	data := `{"id": 123456789012345678}`

	dec := json.NewDecoder(strings.NewReader(data))
	dec.UseNumber() // 数値を json.Number (string) として保持

	var result map[string]any
	dec.Decode(&result)

	id := result["id"].(json.Number)
	fmt.Printf("  型:   %T\n", id)
	fmt.Printf("  値:   %s (文字列として正確に保持)\n", id.String())

	// int64 に変換
	n, _ := id.Int64()
	fmt.Printf("  int64: %d\n", n)
	fmt.Println()
}

// omitemptyDemo: omitempty の各型での挙動
func omitemptyDemo() {
	fmt.Println("--- omitempty ---")

	type Example struct {
		Name    string  `json:"name,omitempty"`
		Age     int     `json:"age,omitempty"`
		Active  bool    `json:"active,omitempty"`
		Score   float64 `json:"score,omitempty"`
		Tags    []string `json:"tags,omitempty"`
		Address *string `json:"address,omitempty"`
	}

	// ゼロ値の構造体
	empty := Example{}
	b, _ := json.Marshal(empty)
	fmt.Printf("  ゼロ値: %s\n", string(b))
	// → {} （全て省略される）

	// 値が入っている構造体
	addr := "Tokyo"
	filled := Example{
		Name:    "Alice",
		Age:     0,      // ← 0 は省略される!
		Active:  false,  // ← false は省略される!
		Tags:    []string{},
		Address: &addr,
	}
	b, _ = json.Marshal(filled)
	fmt.Printf("  注意:   %s\n", string(b))
	fmt.Println("  ※ Age=0, Active=false は omitempty で消える（意図的?）")
	fmt.Println()
}

// decoderVsUnmarshalDemo: Decoder vs Unmarshal
func decoderVsUnmarshalDemo() {
	fmt.Println("--- Decoder vs Unmarshal ---")

	fmt.Println("  json.Unmarshal:")
	fmt.Println("    - []byte を一括デコード")
	fmt.Println("    - メモリに全データが必要")
	fmt.Println()
	fmt.Println("  json.Decoder:")
	fmt.Println("    - io.Reader からストリーミングデコード")
	fmt.Println("    - 大きな JSON や HTTP body に最適")
	fmt.Println("    - UseNumber() が使える")
}
