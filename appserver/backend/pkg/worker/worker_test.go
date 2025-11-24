package worker_test

import (
	"context"
	"fmt"
	"time"

	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker/tasks"
)

func ExampleWorker_basic() {
	ctx := context.Background()

	// Workerを作成
	w := worker.NewWorker(
		worker.WithName("example-worker"),
		worker.WithRunningWorkers(3), // 最大3つのジョブを同時実行
	)

	// Workerを起動
	w.Run(ctx)

	// タスクを作成してジョブを追加
	task1 := tasks.NewTask(func(ctx context.Context) error {
		fmt.Println("Task 1 executed")
		return nil
	})

	task2 := tasks.NewTask(func(ctx context.Context) error {
		fmt.Println("Task 2 executed")
		return nil
	})

	// 複数のタスクを1つのジョブとして追加
	w.AddJob(task1, task2)

	time.Sleep(100 * time.Millisecond)

	// 終了
	w.Shutdown(ctx)
}

func ExampleWorker_concurrent() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5つのジョブを同時実行可能なWorkerを作成
	w := worker.NewWorker(
		worker.WithName("concurrent-worker"),
		worker.WithRunningWorkers(5),
		worker.WithMaxWorkerJobs(100),
	)

	w.Run(ctx)

	// 10個のタスクを追加（5つずつ並行実行される）
	for i := 0; i < 10; i++ {
		i := i // キャプチャ
		task := tasks.NewTask(func(ctx context.Context) error {
			fmt.Printf("Processing task %d\n", i)
			time.Sleep(100 * time.Millisecond)
			return nil
		})
		w.AddJob(task)
	}

	time.Sleep(500 * time.Millisecond)

	w.Shutdown(ctx)
}

func ExampleWorker_withError() {
	ctx := context.Background()

	w := worker.NewWorker(
		worker.WithName("error-worker"),
		worker.WithRunningWorkers(2),
	)

	w.Run(ctx)

	// 成功するタスク
	successTask := tasks.NewTask(func(ctx context.Context) error {
		fmt.Println("Success task executed")
		return nil
	})

	// 失敗するタスク
	errorTask := tasks.NewTask(func(ctx context.Context) error {
		return fmt.Errorf("task failed")
	})

	// エラーが発生しても後続のタスクは実行される
	w.AddJob(errorTask, successTask)

	time.Sleep(100 * time.Millisecond)

	w.Shutdown(ctx)
}
