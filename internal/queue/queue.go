package queue

import "github.com/adarshvbhv/distributed-job-queue/internal/job"

type Queue struct {
	jobs chan *job.Job
}

func New(capacity int) *Queue{
	return &Queue{
		jobs: make(chan *job.Job, capacity),
	}
}

func (q *Queue) Enquue(j *job.Job) {
	q.jobs<-j
}

func (q *Queue) Dqueue() *job.Job{
	return <- q.jobs
}