# 21. Slice Internals

## 核心ポイント
**スライシングはバッキング配列を共有する！意図しない副作用に注意。**

## SliceHeader

```go
type SliceHeader struct {
    Data uintptr  // バッキング配列へのポインタ
    Len  int      // 長さ
    Cap  int      // 容量
}
```

## append の成長戦略

- `cap < 256`: 容量を **2倍**
- `cap >= 256`: 容量を **1.25倍 + 定数**
- append 後に新しいスライスが返る（cap を超えると別配列になる）

## バッキング配列の共有

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]  // b = [2, 3], バッキング配列は a と同じ!
b[0] = 99    // a[1] も 99 になる!
```

## 安全なコピー

- `copy(dst, src)`: 独立した配列にコピー
- `s[low:high:max]`: full slice expression で cap を制限

## nil slice vs empty slice

| | nil slice | empty slice |
|--|----------|-------------|
| 宣言 | `var s []int` | `s := []int{}` |
| len | 0 | 0 |
| cap | 0 | 0 |
| == nil | true | false |
| JSON | `null` | `[]` |

## 実行方法

```bash
go run ./21_slice_internals/
```
