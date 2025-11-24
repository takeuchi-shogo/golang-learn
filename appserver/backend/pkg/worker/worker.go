package worker

import (
	"context"
	"errors"
	"log/slog"

	"golang.org/x/sync/semaphore"
)

var (
	ErrJobQueueFull = errors.New("job queue is full")
)

type WorkerStatus int

func (w WorkerStatus) String() string {
	switch w {
	case WorkerStatusRunning:
		return "running"
	case WorkerStatusProcessing:
		return "processing"
	case WorkerStatusStopped:
		return "stopped"
	}
	return "unknown"
}

func (w WorkerStatus) IsRunning() bool {
	return w == WorkerStatusRunning
}

func (w WorkerStatus) IsProcessing() bool {
	return w == WorkerStatusProcessing
}

func (w WorkerStatus) IsStopped() bool {
	return w == WorkerStatusStopped
}

const (
	defaultWorkerName     = "default-worker"
	defaultMinWorkerJobs  = 1
	defaultMaxWorkerJobs  = 100
	defaultRunningWorkers = 5
	defaultWorkerStatus   = WorkerStatusRunning

	WorkerStatusRunning WorkerStatus = iota
	WorkerStatusProcessing
	WorkerStatusStopped
)

type WorkerOption func(*worker)

func WithName(name string) WorkerOption {
	return func(w *worker) {
		w.name = name
	}
}

func WithMaxWorkerJobs(maxWorkerJobs int) WorkerOption {
	if maxWorkerJobs <= 0 {
		return func(w *worker) {
			w.maxWorkerJobs = defaultMaxWorkerJobs
		}
	}
	return func(w *worker) {
		w.maxWorkerJobs = maxWorkerJobs
	}
}

func WithMinWorkerJobs(minWorkerJobs int) WorkerOption {
	if minWorkerJobs <= 0 {
		return func(w *worker) {
			w.minWorkerJobs = defaultMinWorkerJobs
		}
	}
	return func(w *worker) {
		w.minWorkerJobs = minWorkerJobs
	}
}

func WithRunningWorkers(runningWorkers int) WorkerOption {
	if runningWorkers <= 0 {
		return func(w *worker) {
			w.runningWorkers = defaultRunningWorkers
		}
	}
	return func(w *worker) {
		w.runningWorkers = runningWorkers
	}
}

type Worker interface {
	Run(ctx context.Context) error
	AddJob(tasks ...Task)
	AddJobAsync(tasks ...Task) error
	Shutdown(ctx context.Context) error
	Status() WorkerStatus
	QueueLength() int
	IsRunning() bool
}

type job struct {
	task []Task
}

type Task interface {
	Execute(ctx context.Context) error
}

type worker struct {
	name           string
	minWorkerJobs  int
	maxWorkerJobs  int
	status         WorkerStatus
	runningWorkers int
	sem            *semaphore.Weighted
	jobQueue       chan job
	running        bool
}

var defaultWorker = worker{
	name:           defaultWorkerName,
	minWorkerJobs:  defaultMinWorkerJobs,
	maxWorkerJobs:  defaultMaxWorkerJobs,
	status:         defaultWorkerStatus,
	runningWorkers: defaultRunningWorkers,
	sem:            semaphore.NewWeighted(int64(defaultRunningWorkers)),
	jobQueue:       make(chan job, defaultMaxWorkerJobs),
}

func NewWorker(opts ...WorkerOption) Worker {
	options := defaultWorker
	for _, opt := range opts {
		opt(&options)
	}
	return &worker{
		name:           options.name,
		minWorkerJobs:  options.minWorkerJobs,
		maxWorkerJobs:  options.maxWorkerJobs,
		status:         options.status,
		runningWorkers: options.runningWorkers,
		sem:            semaphore.NewWeighted(int64(options.runningWorkers)),
		jobQueue:       make(chan job, options.maxWorkerJobs),
	}
}

func (w *worker) processJob(ctx context.Context, job *job) {
	if err := w.sem.Acquire(ctx, 1); err != nil {
		slog.Error("failed to acquire semaphore", "error", err)
		return
	}
	defer w.sem.Release(1)

	for _, task := range job.task {
		if err := task.Execute(ctx); err != nil {
			slog.Error("task execution failed", "error", err)
			continue
		}
		slog.Info("task executed successfully")
	}
}

func (w *worker) AddJob(tasks ...Task) {
	w.jobQueue <- job{task: tasks}
}

func (w *worker) AddJobAsync(tasks ...Task) error {
	select {
	case w.jobQueue <- job{task: tasks}:
		return nil
	default:
		return ErrJobQueueFull
	}
}

func (w *worker) Status() WorkerStatus {
	return w.status
}

func (w *worker) QueueLength() int {
	return len(w.jobQueue)
}

func (w *worker) IsRunning() bool {
	return w.running
}

func (w *worker) Shutdown(ctx context.Context) error {
	w.running = false
	close(w.jobQueue)

	// セマフォが解放されるまで待機
	if err := w.sem.Acquire(ctx, int64(w.runningWorkers)); err != nil {
		return err
	}
	defer w.sem.Release(int64(w.runningWorkers))

	w.status = WorkerStatusStopped
	slog.Info("worker shutdown completed", "name", w.name)
	return nil
}

func (w *worker) Run(ctx context.Context) error {
	w.running = true
	w.status = WorkerStatusRunning
	slog.Info("worker started", "name", w.name, "minWorkerJobs", w.minWorkerJobs, "maxWorkerJobs", w.maxWorkerJobs, "runningWorkers", w.runningWorkers)
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.running = false
				w.status = WorkerStatusStopped
				slog.Info("worker stopped", "name", w.name)
				return
			case job, ok := <-w.jobQueue:
				if !ok {
					w.running = false
					w.status = WorkerStatusStopped
					slog.Info("job queue closed", "name", w.name)
					return
				}
				go w.processJob(ctx, &job)
			}
		}
	}()
	return nil
}
