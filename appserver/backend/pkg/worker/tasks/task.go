package tasks

import (
	"context"

	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/worker"
)

type task struct {
	fn func(ctx context.Context) error
}

func NewTask(fn func(ctx context.Context) error) worker.Task {
	return &task{fn: fn}
}

func (t *task) Execute(ctx context.Context) error {
	return t.fn(ctx)
}
