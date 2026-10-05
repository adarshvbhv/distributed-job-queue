package job

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/adarshvbhv/distributed-job-queue/internal/storage"
	"github.com/google/uuid"
)

// Helper to set up clean DB session and automatic table truncation
func setupTestDB(t *testing.T) (*storage.Postgres, context.Context) {
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

func TestPostgresRepository_CreateAndGetByID(t *testing.T) {
	db, ctx := setupTestDB(t)
	repo := NewPostgresRepository(db)


	job := NewJob(
		"test.job",
		[]byte(`{"message":"hello"}`),
		nil,
		3,
	)

	err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	got, err := repo.GetByID(ctx, job.ID)
	if err != nil {
		t.Fatalf("get job by id: %v", err)
	}

	if got.ID != job.ID {
		t.Fatalf("expected ID %v, got %v", job.ID, got.ID)
	}

	if got.Type != job.Type {
		t.Fatalf("expected type %q, got %q", job.Type, got.Type)
	}

	if got.Status != job.Status {
		t.Fatalf("expected status %q, got %q", job.Status, got.Status)
	}
}

func TestPostgresRepository_Idempotency(t *testing.T) {
	db, ctx := setupTestDB(t)
	repo := NewPostgresRepository(db)

	key := "test-idempotency-key-2"

	job1 := NewJob(
		"test.job",
		[]byte(`{"message":"first"}`),
		&key, 
		3,
	)

	job2 := NewJob(
		"test.job",
		[]byte(`{"message":"second"}`),
		&key,
		3,
	)

	// First create must succeed
	if err := repo.Create(ctx, job1); err != nil {
		t.Fatalf("create first job: %v", err)
	}

	// Attempt second create with identical key
	err := repo.Create(ctx, job2)

	// Verify that duplicate creation returns ErrDuplicateIdempotency
	if !errors.Is(err, ErrDuplicateIdempotency) {
		t.Fatalf("expected ErrDuplicateIdempotency, got: %v", err)
	}
}

func TestPostgresRepository_GetByID_NotFound(t *testing.T) {
	// initialize DB/repository as in your existing test
	db, ctx := setupTestDB(t)
	repo := NewPostgresRepository(db)

	id := uuid.New()

	_, err := repo.GetByID(ctx, id)

	if !errors.Is(err, ErrJobNotFound) {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}
}
