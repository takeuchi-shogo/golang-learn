# Go Developer — 38の必須トピック

> 「これを明確に説明できないなら転職を考えろ」レベルの Go 必須知識を、
> 動くコードと解説で学ぶリポジトリ。

## トピック一覧

| # | ディレクトリ | トピック | 核心ポイント |
|---|-------------|---------|-------------|
| 01 | [memory_model](./01_memory_model/) | Go Memory Model | happens-before がない = 順序保証なし |
| 02 | [happens_before](./02_happens_before/) | Happens-before | mutex/channel/atomic/WaitGroup が保証を作る |
| 03 | [data_race](./03_data_race/) | Data Race | fatal error = recover 不可、-race で必ず検査 |
| 04 | [gmp_scheduler](./04_gmp_scheduler/) | G, M, P モデル | G=goroutine, M=OS thread, P=scheduler context |
| 05 | [preemption](./05_preemption/) | Preemption | tight loop は 1.14+ asynchronous preemption |
| 06 | [channel_internals](./06_channel_internals/) | Channel internals | unbuffered=同期点、buffered=非同期バッファ |
| 07 | [select](./07_select/) | select | 複数 ready = pseudo-random 選択 |
| 08 | [channel_close](./08_channel_close/) | Channel close | 送信側が閉じる、受信側は ok で判断 |
| 09 | [nil_channel](./09_nil_channel/) | nil channel | select での case 無効化に使える |
| 10 | [context](./10_context/) | Context | 必ず defer cancel()、伝播は親→子 |
| 11 | [sync_once](./11_sync_once/) | sync.Once | 再入不可、panic でも完了扱い |
| 12 | [sync_pool](./12_sync_pool/) | sync.Pool | GC でクリアされる = キャッシュじゃない |
| 13 | [escape_analysis](./13_escape_analysis/) | Escape analysis | -gcflags="-m" でスタック/ヒープ確認 |
| 14 | [pointer_vs_value_receiver](./14_pointer_vs_value_receiver/) | Pointer/Value receiver | \*T は全メソッド持つ、interface には \*T 安全 |
| 15 | [interface_internals](./15_interface_internals/) | Interface internals | nil pointer ≠ nil interface の罠 |
| 16 | [empty_interface](./16_empty_interface/) | Empty interface | interface{} は heap エスケープの原因 |
| 17 | [defer_internals](./17_defer_internals/) | defer internals | hot loop では関数に切り出す |
| 18 | [panic_recover](./18_panic_recover/) | panic/recover | 同一 goroutine のみ recover 可能 |
| 19 | [map_internals](./19_map_internals/) | Map internals | イテレーション順序は意図的にランダム |
| 20 | [map_concurrency](./20_map_concurrency/) | Map concurrency | fatal error = recover 不可、sync.Map |
| 21 | [slice_internals](./21_slice_internals/) | Slice internals | スライシングは配列を共有する！ |
| 22 | [string_bytes](./22_string_bytes/) | String/[]byte | 変換はコピー、Builder で最適化 |
| 23 | [gc_basics](./23_gc_basics/) | GC | 三色マーキング、GOMEMLIMIT 推奨 |
| 24 | [write_barrier](./24_write_barrier/) | Write barrier | ポインタ書き込みを GC に通知する安全装置 |
| 25 | [gomaxprocs](./25_gomaxprocs/) | GOMAXPROCS | K8s では uber/automaxprocs を入れること |
| 26 | [atomic_vs_mutex_vs_channel](./26_atomic_vs_mutex_vs_channel/) | atomic vs mutex vs channel | シンプル=atomic、複数変数=mutex、通信=channel |
| 27 | [memory_alignment](./27_memory_alignment/) | Memory alignment | 64バイトパディングでキャッシュライン分離 |
| 28 | [pprof](./28_pprof/) | pprof | \_ "net/http/pprof" import だけで有効 |
| 29 | [cgo](./29_cgo/) | cgo | 境界越え ~100-200ns、Go 関数の100倍コスト |
| 30 | [http_transport](./30_http_transport/) | net/http Transport | Body を読み切らないと接続再利用されない |
| 31 | [http_client_reuse](./31_http_client_reuse/) | http.Client reuse | リクエストごとに new するな！ |
| 32 | [json_pitfalls](./32_json_pitfalls/) | JSON pitfalls | interface{} で数値 = float64、大整数が壊れる |
| 33 | [generics](./33_generics/) | Generics | コンパイル時展開 = interface より速い |
| 34 | [init_order](./34_init_order/) | init() order | アルファベット順、パニック = 起動不能 |
| 35 | [build_tags](./35_build_tags/) | Build tags | replace はローカル/CI 全部に効く |
| 36 | [embedding](./36_embedding/) | Embedding | 継承じゃない、promotion と ambiguity に注意 |
| 37 | [error_handling](./37_error_handling/) | Error handling | %w で wrap、errors.Is/As で遡る |
| 38 | [race_detector](./38_race_detector/) | Race detector | 実行時のみ検知、不検出の可能性あり |

## 使い方

```bash
# 各トピックのコードを実行
cd go_developer
go run ./01_memory_model/

# Race detector 付きで実行（03, 20, 38 などで特に重要）
go run -race ./03_data_race/

# Escape analysis を確認（13）
go build -gcflags="-m" ./13_escape_analysis/

# テスト付きトピック
go test ./37_error_handling/
```
