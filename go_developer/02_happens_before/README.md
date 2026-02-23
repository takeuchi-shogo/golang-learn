# 02. Happens-before

## 核心ポイント
**mutex / channel / atomic / WaitGroup が happens-before 保証を作る。**

## happens-before とは

イベント A が B に対して happens-before であるとき、A の効果（メモリ書き込み）は B から必ず観測できる。

## 各プリミティブの happens-before 保証

### 1. sync.Mutex / sync.RWMutex

```
mu.Lock()
// クリティカルセクション
mu.Unlock()  ──happens-before──▶  mu.Lock()
```

n 回目の `Unlock()` は n+1 回目の `Lock()` に対して happens-before。

### 2. Channel

**Unbuffered:**
```
ch <- value  ──happens-before──▶  <-ch
```
send の完了は対応する receive の完了に先行する。

**Buffered (cap = C):**
```
// k 回目の receive は k+C 回目の send に先行する
```

**Close:**
```
close(ch)  ──happens-before──▶  <-ch がゼロ値を返す
```

### 3. sync.WaitGroup

```
wg.Done()  ──happens-before──▶  wg.Wait() の完了
```

### 4. sync/atomic

```
atomic.Store(&x, val)  ──happens-before──▶  atomic.Load(&x) == val
```

Go 1.19+ で公式に明文化された。

### 5. sync.Once

```
once.Do(f) の完了  ──happens-before──▶  他の once.Do の return
```

### 6. Goroutine の起動

```
go f()  ──happens-before──▶  f() の実行開始
```

## 実行方法

```bash
go run ./02_happens_before/
```
