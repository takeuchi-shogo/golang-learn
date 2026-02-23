# 19. Map Internals

## 核心ポイント
**イテレーション順序は意図的にランダム。nil map は読めるが書くと panic。**

## runtime.hmap

```
hmap {
    count      int     // 要素数（len）
    B          uint8   // バケット数 = 2^B
    buckets    *bmap   // バケット配列
    oldbuckets *bmap   // 拡張中の旧バケット
}

bmap {
    tophash [8]uint8   // 各キーのハッシュ上位8ビット
    keys    [8]K       // キー配列
    values  [8]V       // 値配列
    overflow *bmap     // オーバーフローバケット
}
```

## ハッシュとバケット選択

1. キーのハッシュを計算
2. 下位 B ビットでバケットを選択
3. 上位 8 ビット（tophash）でバケット内を線形探索

## map の成長

- **ロードファクター > 6.5**: バケット数を2倍に拡張
- **incremental growing**: 一度に全コピーせず、少しずつ移行

## イテレーション順序がランダムな理由

Go はイテレーション開始位置を**意図的にランダム化**している。
順序に依存するコードを防ぐため。

## 実行方法

```bash
go run ./19_map_internals/
```
