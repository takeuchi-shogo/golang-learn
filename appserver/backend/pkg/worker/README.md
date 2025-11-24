# Worker Package

並行処理を管理するためのWorkerパッケージです。

## 特徴

- **並行実行制御**: セマフォを使用して同時実行数を制限
- **ジョブキュー**: バッファ付きチャネルでジョブを管理
- **柔軟なタスク管理**: 関数ベースの簡単なタスク作成
- **グレースフルシャットダウン**: 実行中のジョブを待機して終了

## 基本的な使い方

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker"
    "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker/tasks"
)

func main() {
    ctx := context.Background()

    // Workerを作成
    w := worker.NewWorker(
        worker.WithName("my-worker"),
        worker.WithRunningWorkers(5), // 最大5つのジョブを同時実行
        worker.WithMaxWorkerJobs(100), // キューに最大100ジョブ
    )

    // Workerを起動
    w.Run(ctx)

    // タスクを作成
    task := tasks.NewTask(func(ctx context.Context) error {
        fmt.Println("Task executed")
        return nil
    })

    // ジョブを追加
    w.AddJob(task)

    time.Sleep(100 * time.Millisecond)

    // 終了
    w.Shutdown(ctx)
}
```

## API

### Worker Interface

```go
type Worker interface {
    Run(ctx context.Context) error
    AddJob(tasks ...Task)
    AddJobAsync(tasks ...Task) error
    Shutdown(ctx context.Context) error
}
```

#### `Run(ctx context.Context) error`
Workerを起動します。バックグラウンドでジョブキューを監視し、ジョブを処理します。

#### `AddJob(tasks ...Task)`
ジョブをキューに追加します（ブロッキング）。
キューが満杯の場合、空きができるまで待機します。

```go
w.AddJob(task1, task2, task3)
```

#### `AddJobAsync(tasks ...Task) error`
ジョブをキューに追加します（非ブロッキング）。
キューが満杯の場合、`ErrJobQueueFull` エラーを返します。

```go
if err := w.AddJobAsync(task); err != nil {
    if errors.Is(err, worker.ErrJobQueueFull) {
        // キューが満杯
    }
}
```

#### `Shutdown(ctx context.Context) error`
Workerをグレースフルシャットダウンします。
実行中のすべてのジョブが完了するのを待機します。

### Worker Options

#### `WithName(name string)`
Workerに名前を設定します。

#### `WithRunningWorkers(n int)`
同時実行するジョブの最大数を設定します（デフォルト: 5）。

#### `WithMaxWorkerJobs(n int)`
ジョブキューの最大サイズを設定します（デフォルト: 100）。

#### `WithMinWorkerJobs(n int)`
ジョブキューの最小サイズを設定します（デフォルト: 1）。

### Task Interface

```go
type Task interface {
    Execute(ctx context.Context) error
}
```

### タスクの作成

`tasks.NewTask()` を使用して、関数からタスクを作成できます。

```go
task := tasks.NewTask(func(ctx context.Context) error {
    // ビジネスロジック
    return nil
})
```

## 実践例

### データベース処理の並行化

```go
w := worker.NewWorker(
    worker.WithName("db-worker"),
    worker.WithRunningWorkers(10),
)
w.Run(ctx)

for _, user := range users {
    user := user
    task := tasks.NewTask(func(ctx context.Context) error {
        return db.SaveUser(ctx, user)
    })
    w.AddJob(task)
}
```

### API呼び出しの並行化

```go
w := worker.NewWorker(
    worker.WithName("api-worker"),
    worker.WithRunningWorkers(5),
)
w.Run(ctx)

for _, endpoint := range endpoints {
    endpoint := endpoint
    task := tasks.NewTask(func(ctx context.Context) error {
        resp, err := http.Get(endpoint)
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        // 処理
        return nil
    })
    w.AddJob(task)
}
```

### 複数タスクの順次実行

1つのジョブに複数のタスクを含めると、それらは順次実行されます。

```go
task1 := tasks.NewTask(func(ctx context.Context) error {
    // 前処理
    return nil
})

task2 := tasks.NewTask(func(ctx context.Context) error {
    // メイン処理
    return nil
})

task3 := tasks.NewTask(func(ctx context.Context) error {
    // 後処理
    return nil
})

// これらのタスクは順次実行される
w.AddJob(task1, task2, task3)
```

## アーキテクチャ

```
┌─────────────────────────────────────────────┐
│              Worker                         │
│                                             │
│  ┌─────────────┐         ┌──────────────┐  │
│  │  jobQueue   │────────>│  Goroutine   │  │
│  │  (channel)  │         │  (監視用)    │  │
│  └─────────────┘         └──────────────┘  │
│                                 │           │
│                                 v           │
│                    ┌─────────────────────┐  │
│                    │   Semaphore         │  │
│                    │   (同時実行制限)    │  │
│                    └─────────────────────┘  │
│                                 │           │
│                  ┌──────────────┴────────┐  │
│                  v              v         v  │
│          ┌───────────┐  ┌───────────┐  ...  │
│          │ Job       │  │ Job       │        │
│          │ Goroutine │  │ Goroutine │        │
│          └───────────┘  └───────────┘        │
└─────────────────────────────────────────────┘
```

### 並行実行の仕組み

1. `jobQueue` からジョブを受信
2. 各ジョブを新しいゴルーチンで起動
3. セマフォで同時実行数を制限（`runningWorkers`）
4. ジョブ内のタスクは順次実行
5. ジョブ完了後、セマフォを解放

### 設定例

```go
// 同時実行数: 5
// キューサイズ: 100
w := worker.NewWorker(
    worker.WithRunningWorkers(5),
    worker.WithMaxWorkerJobs(100),
)
```

- 最大5つのジョブが同時に処理される
- 最大100個のジョブをキューに保持できる
- 101個目のジョブは、キューに空きができるまでブロック（`AddJob`）またはエラー（`AddJobAsync`）

## エラーハンドリング

タスクがエラーを返しても、同じジョブ内の後続タスクは実行されます（`continue`）。
エラーはログに記録されます。

```go
task1 := tasks.NewTask(func(ctx context.Context) error {
    return fmt.Errorf("error") // エラー発生
})

task2 := tasks.NewTask(func(ctx context.Context) error {
    fmt.Println("Still executed") // これは実行される
    return nil
})

w.AddJob(task1, task2)
```

## ベストプラクティス

1. **適切な同時実行数**: リソースに応じて `runningWorkers` を調整
2. **タイムアウト**: 長時間実行されるタスクには context のタイムアウトを設定
3. **エラーハンドリング**: タスク内で適切にエラーを処理
4. **グレースフルシャットダウン**: アプリケーション終了時に `Shutdown()` を呼び出す

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

w := worker.NewWorker(
    worker.WithRunningWorkers(runtime.NumCPU()),
)
w.Run(ctx)

// ... タスクを追加 ...

// 終了時
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
defer shutdownCancel()
if err := w.Shutdown(shutdownCtx); err != nil {
    log.Printf("shutdown error: %v", err)
}
```
