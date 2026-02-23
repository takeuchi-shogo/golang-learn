package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

// 本番環境では以下のように HTTP pprof を有効にする:
//
//   import _ "net/http/pprof"
//   go func() {
//       http.ListenAndServe("localhost:6060", nil)
//   }()
//
// → http://localhost:6060/debug/pprof/ でアクセス
//
// ここでは runtime/pprof を使ったファイル出力のデモを行う

func main() {
	fmt.Println("=== pprof ===")
	fmt.Println()

	cpuProfileDemo()
	memProfileDemo()
	goroutineProfileDemo()
}

// cpuProfileDemo: CPU プロファイルをファイルに書き出す
func cpuProfileDemo() {
	fmt.Println("--- CPU Profile ---")

	f, err := os.CreateTemp("", "cpu-*.prof")
	if err != nil {
		fmt.Printf("  エラー: %v\n", err)
		return
	}
	defer os.Remove(f.Name())

	pprof.StartCPUProfile(f)
	// CPU-bound な処理
	sum := 0
	for i := range 10_000_000 {
		sum += i
	}
	pprof.StopCPUProfile()
	f.Close()

	fmt.Printf("  CPU profile を %s に書き出し\n", f.Name())
	fmt.Println("  分析: go tool pprof <file>")
	fmt.Println()
}

// memProfileDemo: メモリプロファイル
func memProfileDemo() {
	fmt.Println("--- Memory Profile ---")

	// アロケーションを発生させる
	var sink [][]byte
	for range 1000 {
		sink = append(sink, make([]byte, 1024))
	}
	_ = sink

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  HeapAlloc: %d KB\n", m.Alloc/1024)
	fmt.Printf("  HeapObjects: %d\n", m.HeapObjects)
	fmt.Println("  分析: go tool pprof http://localhost:6060/debug/pprof/heap")
	fmt.Println()
}

// goroutineProfileDemo: goroutine プロファイル
func goroutineProfileDemo() {
	fmt.Println("--- Goroutine Profile ---")

	done := make(chan struct{})
	for range 5 {
		go func() { <-done }()
	}

	profile := pprof.Lookup("goroutine")
	fmt.Printf("  現在の goroutine 数: %d\n", profile.Count())
	fmt.Println("  分析: go tool pprof http://localhost:6060/debug/pprof/goroutine")

	close(done)
	fmt.Println()
	fmt.Println("=== pprof コマンド早見表 ===")
	fmt.Println("  go tool pprof <profile>")
	fmt.Println("  (pprof) top 10      # 上位10関数")
	fmt.Println("  (pprof) list <func>  # ソースコード")
	fmt.Println("  (pprof) web          # ブラウザ可視化")
}
