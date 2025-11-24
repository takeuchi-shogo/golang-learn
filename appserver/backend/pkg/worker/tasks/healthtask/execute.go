package healthtask

import (
	"context"
	"log/slog"
	"time"

	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker/tasks"
)

var HealthTask = tasks.NewTask(func(ctx context.Context) error {
	slog.Info("health task executed", "time", time.Now().Format(time.RFC3339))
	return nil
})
