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
	repo      job.Repository
}

func NewWorkerPool(queue *queue.Queue, workerCnt int, logger *slog.Logger, registry *job.Registry, repo job.Repository) *WorkerPool {
	return &WorkerPool{
		queue:     queue,
		workerCnt: workerCnt,
		logger:    logger,
		registry:  registry,
		repo:      repo,
	}
}

func (w *WorkerPool) worker() {
	for {
		j := w.queue.Dqueue()
		ctx := context.Background()

		if err := w.repo.UpdateJobStatus(ctx, j.ID, job.StatusRunning, nil); err != nil {
			w.logger.Error("failed to update job to running status", "job id", j.ID, "error", err)
			continue
		}

		handler, ok := w.registry.Get(j.Type)

		if !ok {
			errMsg := "no handler found for this job"
			w.logger.Error(errMsg, "job id", j.Type)
			w.repo.UpdateJobStatus(ctx, j.ID, job.StatusFailed, &errMsg)
		}

		err := handler.Execute(ctx, j.Payload)

		if err != nil {
			errMsg := err.Error()
			w.logger.Error(errMsg, "job", j.ID, "error", err)
			w.repo.UpdateJobStatus(ctx, j.ID, job.StatusFailed, &errMsg)

			continue
		}

		w.repo.UpdateJobStatus(ctx, j.ID, job.StatusCompleted, nil)

		w.logger.Info("job executed successfully", "job id", j.ID)
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.workerCnt; i++ {
		go w.worker()
	}
}
