package worker

import (
	"context"
	"log/slog"

	"github.com/adarshvbhv/distributed-job-queue/internal/job"
	"github.com/adarshvbhv/distributed-job-queue/internal/queue"
)

type WorkerPool struct {
	queue     *queue.Queue
	workerCnt int
	logger    *slog.Logger
	registry  *job.Registry
}

func NewWorkerPool(queue *queue.Queue, workerCnt int, logger *slog.Logger, registry *job.Registry) *WorkerPool {
	return &WorkerPool{
		queue:     queue,
		workerCnt: workerCnt,
		logger:    logger,
		registry:  registry,
	}
}

func (w *WorkerPool) worker() {
	for {
		job := w.queue.Dqueue()

		handler, ok := w.registry.Get(job.Type)

		if !ok {
			w.logger.Error("no handler found for this job", "job id", job.Type)
		}

		ctx := context.Background()

		err := handler.Execute(ctx, job.Payload)

		if err != nil {
			w.logger.Error("job execution failed", "job", job.ID, "error", err)

			continue
		}

		w.logger.Info("job executed successfully", "job id", job.ID)
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.workerCnt; i++ {
		go w.worker()
	}
}
