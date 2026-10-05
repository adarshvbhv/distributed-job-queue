package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/adarshvbhv/distributed-job-queue/internal/job"
)

type createJobRequest struct {
	Type       string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
	MaxRetries *int            `json:"max_retries,omitempty"`
}

type createJobResponse struct {
	JobID  uuid.UUID  `json:"job_id"`
	Status job.Status `json:"status"`
}

func (s *Server) createJobHandler(w http.ResponseWriter, r *http.Request) {
	// Use Error when an engineer needs to wake up or fix a bug in the code/infrastructure.

	// Use Warn/Info when a user or API client made a mistake.

	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warn("invalid JSON request payload", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "invalid json", http.StatusBadRequest)

		return
	}

	if req.Type == "" {
		http.Error(w, "job type is missing", http.StatusBadRequest)
		return
	}

	var idempotencyKey *string

	if key := r.Header.Get("Idempotency-Key"); key != "" {
		idempotencyKey = &key
	}

	maxRetries := 3

	if req.MaxRetries != nil {
		if *req.MaxRetries < 0 {
			http.Error(w, "max_retries cannot be negative", http.StatusBadRequest)
			return
		}

		maxRetries = *req.MaxRetries
	}

	newJob := job.NewJob(
		req.Type,
		req.Payload,
		idempotencyKey,
		maxRetries,
	)

	if err := s.jobRepo.Create(r.Context(), newJob); err != nil {
		if errors.Is(err, job.ErrDuplicateIdempotency) {
			http.Error(
				w,
				"duplicate idempotency key",
				http.StatusConflict,
			)
			return
		}

		s.logger.Error("failed to create job", "error", err)

		http.Error(
			w,
			"unable to create job",
			http.StatusInternalServerError,
		)
		return
	}

	response := createJobResponse{
		JobID:  newJob.ID,
		Status: newJob.Status,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode job response", "error", err)
	}

}


func (s *Server) getJobHandler(w http.ResponseWriter, r *http.Request){


	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid job ID", http.StatusBadRequest)
		return
	}

	storedJob, err := s.jobRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, job.ErrJobNotFound) {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		s.logger.Error("failed to get job", "job_id", id, "error", err)
		http.Error(w, "unable to get job", http.StatusInternalServerError)
		return
	}

	response := createJobResponse{
		JobID:  storedJob.ID,
		Status: storedJob.Status,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode job response", "error", err)
	}

}