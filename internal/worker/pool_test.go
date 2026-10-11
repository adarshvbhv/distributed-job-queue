package worker

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/adarshvbhv/distributed-job-queue/internal/job"
	"github.com/adarshvbhv/distributed-job-queue/internal/queue"
)

type TestHandler struct {
	Called chan struct{}
}

func (h *TestHandler) Execute(ctx context.Context, payload []byte) error {
	h.Called <- struct{}{}
	return nil
}
func TestWorkerPool_ExecutesHandler(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	q := queue.New(1)
	registry := job.NewRegistry()
	called := make(chan struct{}, 1)
	registry.Register("test.job", &TestHandler{Called: called})
	workerPool := NewWorkerPool(q, 1, logger, registry)

	newJob := job.NewJob(
		"test.job",
		[]byte(`{"message":"hello"}`),
		nil,
		3,
	)

	q.Enquue(newJob)

	workerPool.Start()

	select {
	case <-called:
		// Handler executed successfully.
	case <-time.After(time.Second):
		t.Fatal("handler was not executed")
	}

}
