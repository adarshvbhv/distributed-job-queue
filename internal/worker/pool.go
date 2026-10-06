package worker

import (
	"log/slog"

	"github.com/adarshvbhv/distributed-job-queue/internal/queue"
)

type WorkerPool struct {
	queue     *queue.Queue
	workerCnt int
	logger    *slog.Logger
}

func NewWorkerPool(queue *queue.Queue, workerCnt int, logger *slog.Logger) *WorkerPool {
	return &WorkerPool{
		queue:     queue,
		workerCnt: workerCnt,
		logger:    logger,
	}
}

func (w *WorkerPool) worker() {
	for {
		job := w.queue.Dqueue()

		w.logger.Info("running job", "job id", job.ID)
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.workerCnt; i++ {
		go w.worker()
	}
}
