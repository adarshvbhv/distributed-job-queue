package job

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusQueued     Status = "QUEUED"
	StatusRunning    Status = "RUNNING"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusRetrying   Status = "RETRYING"
	StatusCancelled  Status = "CANCELLED"
	StatusDeadLetter Status = "DEAD_LETTER"
)

type Job struct {
	ID             uuid.UUID
	Type           string
	Payload        []byte
	Status         Status
	Attempts       int
	MaxRetries     int
	CreatedAt      time.Time
	StartedAt      *time.Time
	CompletedAt    *time.Time
	NextRetryAt    *time.Time
	Error          *string
	WorkerID       *string
	IdempotencyKey *string
}

func NewJob(jobType string, payload []byte, idempotencyKey *string, maxRetries int) *Job {
	job := &Job{
		ID:             uuid.New(),
		Type:           jobType,
		Payload:        payload,
		Status:         StatusQueued,
		Attempts:       0,
		MaxRetries:     maxRetries,
		CreatedAt:      time.Now(),
		StartedAt:      nil,
		CompletedAt:    nil,
		NextRetryAt:    nil,
		Error:          nil,
		WorkerID:       nil,
		IdempotencyKey: idempotencyKey,
	}

	return job
}
