package common

import "context"

type JobEnQueuer interface {
	Enqueue(ctx context.Context, args any) error
}

type QueueJobOption func(JobEnQueuer) error

// WithQueuedJob runs fn with access to the transaction-scoped enqueuer,
// so jobs are inserted atomically with the save and roll back if it fails.
func WithQueuedJob(fn func(JobEnQueuer) error) QueueJobOption {
	return fn
}
