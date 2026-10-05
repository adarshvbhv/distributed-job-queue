package job

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
}

var (
	ErrJobNotFound          = errors.New("job not found")
	ErrDuplicateIdempotency = errors.New("duplicate idempotency key")
)
