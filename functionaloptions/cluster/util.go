package cluster

func WithMaxWorkers(maxWorkers int) Option {
	return func(o *workerClusterOption) {
		o.maxWorkers = maxWorkers
	}
}
