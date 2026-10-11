package worker

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/adarshvbhv/distributed-job-queue/internal/job"
	"github.com/adarshvbhv/distributed-job-queue/internal/queue"
	"github.com/adarshvbhv/distributed-job-queue/internal/storage"
)

type TestHandler struct {
	Called chan struct{}
}

func (h *TestHandler) Execute(ctx context.Context, payload []byte) error {
	h.Called <- struct{}{}
	return nil
}

// Helper to set up clean DB session and automatic table truncation
func SetupTestDB(t *testing.T) (*storage.Postgres, context.Context) {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := storage.NewPostgres(ctx, logger, databaseURL)
	if err != nil {
		t.Fatalf("connect to postgres: %v", err)
	}

	// Close connection pool when test finishes
	t.Cleanup(func() {
		db.Pool.Close()
	})

	// Clean table before starting test
	_, err = db.Pool.Exec(ctx, "TRUNCATE TABLE jobs;")
	if err != nil {
		t.Fatalf("truncate table: %v", err)
	}

	return db, ctx
}
func TestWorkerPool_ExecutesHandler(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	q := queue.New(1)
	registry := job.NewRegistry()
	called := make(chan struct{}, 1)
	registry.Register("test.job", &TestHandler{Called: called})
	db, _ := SetupTestDB(t)
	repo := job.NewPostgresRepository(db)

	workerPool := NewWorkerPool(q, 1, logger, registry, repo)

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
