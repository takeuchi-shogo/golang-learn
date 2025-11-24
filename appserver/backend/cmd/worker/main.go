package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker/tasks/healthtask"
)

/** go run cmd/worker/main.go
 * これはworkerを起動の確認するテストコマンドです。
 * 削除 or 変更予定 */
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewWorker()
	w.AddJob(healthtask.HealthTask)
	if err := w.Run(ctx); err != nil {
		fmt.Printf("failed to run worker: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("worker started")
	fmt.Printf("Queue length: %d\n", w.QueueLength())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := w.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown error", "error", err)
	}
	fmt.Printf("After shutdown - Running: %v, Status: %s\n", w.IsRunning(), w.Status())
}
