package worker

import (
	"testing"

	"github.com/adarshvbhv/distributed-job-queue/internal/queue"
)


func TestWorkerPool_ExecutesHandler(t *testing.T){

	q:= queue.New(10)

}