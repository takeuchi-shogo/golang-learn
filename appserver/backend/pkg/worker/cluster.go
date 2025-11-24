package worker

type WorkerCluster struct {
	workers []*worker
}

func NewWorkerCluster(workers []*worker) *WorkerCluster {
	return &WorkerCluster{workers: workers}
}
