package main

import "fmt"

// パッケージ変数の初期化は init() より先
var globalVar = initGlobal()

func initGlobal() string {
	fmt.Println("[1] パッケージ変数の初期化")
	return "initialized"
}

// init() は複数定義できる（定義順に実行）
func init() {
	fmt.Println("[2] init() #1")
}

func init() {
	fmt.Println("[3] init() #2")
}

func init() {
	fmt.Println("[4] init() #3")
}

func main() {
	fmt.Println("[5] main()")
	fmt.Println()
	fmt.Println("=== init() の実行順序 ===")
	fmt.Println()
	fmt.Println("  実行順序:")
	fmt.Println("    1. インポート先パッケージの init()（再帰的）")
	fmt.Println("    2. パッケージ変数の初期化")
	fmt.Println("    3. init() 関数（ファイルアルファベット順 → 定義順）")
	fmt.Println("    4. main()")
	fmt.Println()
	fmt.Printf("  globalVar = %q\n", globalVar)
	fmt.Println()

	fmt.Println("=== init() のアンチパターン ===")
	fmt.Println()
	fmt.Println("  避けるべき:")
	fmt.Println("    - DB接続: テストで困る")
	fmt.Println("    - 環境変数読み取り: テストで上書きできない")
	fmt.Println("    - HTTP クライアント作成: DI できない")
	fmt.Println("    - panic する可能性がある処理: 起動不能")
	fmt.Println()
	fmt.Println("  推奨: 明示的な初期化関数を使う")
	fmt.Println("    func NewServer(cfg Config) (*Server, error) { ... }")
}
