package cluster

import "context"

const (
	defaultMaxWorkers = 5
)

// WorkerCluster is a cluster of workers that can be used to process tasks.
type WorkerCluster struct {
	maxWorkers int
}

type workerClusterOption struct {
	maxWorkers int
}

type Option func(*workerClusterOption)

var defaultWorkerClusterOption = workerClusterOption{
	maxWorkers: defaultMaxWorkers,
}

// New returns a new WorkerCluster.
func New(ctx context.Context, opts ...Option) (*WorkerCluster, error) {
	options := defaultWorkerClusterOption
	for _, opt := range opts {
		opt(&options)
	}
	return &WorkerCluster{
		maxWorkers: options.maxWorkers,
	}, nil
}

func (c *WorkerCluster) GetMaxWorkers() int {
	return c.maxWorkers
}
