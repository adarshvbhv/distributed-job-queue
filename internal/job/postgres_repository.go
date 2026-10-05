package job

import (
	"context"
	"errors"
	"fmt"

	"github.com/adarshvbhv/distributed-job-queue/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/google/uuid"
)

type PostgresRepository struct {
	db *storage.Postgres
}

var _ Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *storage.Postgres) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, job *Job) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO jobs (
			id,
			type,
			payload,
			status,
			attempts,
			max_retries,
			created_at,
			started_at,
			completed_at,
			next_retry_at,
			error,
			worker_id,
			idempotency_key
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		)
		`,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.Attempts,
		job.MaxRetries,
		job.CreatedAt,
		job.StartedAt,
		job.CompletedAt,
		job.NextRetryAt,
		job.Error,
		job.WorkerID,
		job.IdempotencyKey,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateIdempotency
		}
		return fmt.Errorf("create job: %w", err)
	}

	return nil

}

func (r *PostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Job, error) {
	var job Job

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			type,
			payload,
			status,
			attempts,
			max_retries,
			created_at,
			started_at,
			completed_at,
			next_retry_at,
			error,
			worker_id,
			idempotency_key
		FROM jobs
		WHERE id = $1
		`,
		id,
	).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.Attempts,
		&job.MaxRetries,
		&job.CreatedAt,
		&job.StartedAt,
		&job.CompletedAt,
		&job.NextRetryAt,
		&job.Error,
		&job.WorkerID,
		&job.IdempotencyKey,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}

		return nil, fmt.Errorf("get job: %w", err)
	}

	

	return &job, nil
}
