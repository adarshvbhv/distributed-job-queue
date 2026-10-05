package queue

import "github.com/adarshvbhv/distributed-job-queue/internal/job"

type Queue struct {
	jobs chan *job.Job
}