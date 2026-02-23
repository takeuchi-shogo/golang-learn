# 32. JSON Encoding Pitfalls

## 核心ポイント
**`interface{}` にデコードすると数値は float64 になる。大きな整数が壊れる。**

## float64 問題

```go
var data any
json.Unmarshal([]byte(`{"id": 123456789012345}`), &data)
// data["id"] は float64 → 精度が失われる!
```

JSON の数値を `interface{}` にデコードすると、Go は全て `float64` にする。
`float64` の仮数部は 53 ビットなので、2^53 を超える整数は精度が失われる。

## json.Number

```go
dec := json.NewDecoder(r)
dec.UseNumber()  // 数値を json.Number (string) として保持
```

## omitempty の挙動

| 型 | ゼロ値（omitempty で省略される） |
|----|------|
| bool | false |
| int | 0 |
| string | "" |
| pointer | nil |
| slice/map | nil (空でも要素0個は省略されない場合あり) |

## 実行方法

```bash
go run ./32_json_pitfalls/
```
