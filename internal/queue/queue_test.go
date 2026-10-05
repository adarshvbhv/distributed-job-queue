package queue

import (
	"testing"

	"github.com/adarshvbhv/distributed-job-queue/internal/job"
)

func TestEnqueDeque_Jobs(t *testing.T) {

	newJob := job.NewJob(
		"test.job",
		[]byte(`{"message":"hello"}`),
		nil,
		3,
	)

	id := newJob.ID

	q := New(10)
	q.Enquue(newJob)

	jobGot := q.Dqueue()

	if id != jobGot.ID {
		t.Error("id enqueued is not equal to job id queued")
	}

}
