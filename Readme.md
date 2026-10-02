# Distributed Job Queue — Go

## 1. Project Objective

Build a production-oriented **Distributed Job Queue and Worker System in Go** that demonstrates:

* Go concurrency
* Goroutines and channels
* Worker pools
* Context cancellation
* Job scheduling and execution
* Retry and exponential backoff
* Idempotent job execution
* Distributed workers
* Worker heartbeats
* Failure detection
* Redis-based coordination
* PostgreSQL-based durable job state
* Observability
* Performance benchmarking
* Docker-based deployment

The project should be a **focused infrastructure/distributed-systems project**, not another generic microservices application.

The primary purpose is to demonstrate why Go is particularly suitable for concurrent backend and infrastructure systems.

---

# 2. Relationship With My Other Project

I already have another project called **Bangalore Rentals**.

Bangalore Rentals is my application/backend project and already covers:

* REST APIs
* gRPC
* JWT authentication
* Authorization
* WebSockets
* Redis Pub/Sub
* RabbitMQ
* OpenSearch
* PostgreSQL
* PostGIS
* Database-per-service
* Docker
* MinIO/S3
* Prometheus
* Grafana
* OpenTelemetry
* Distributed tracing

Therefore, DO NOT duplicate those responsibilities in this project.

This project must focus on:

> **Go concurrency + job execution + distributed workers + reliability + fault tolerance + infrastructure engineering.**

---

# 3. Hard Scope Boundary

## MUST NOT BUILD

Do NOT add the following unless explicitly requested:

* Frontend
* React/Next.js
* Mobile application
* Chat
* WebSockets
* Property/search functionality
* OpenSearch/Elasticsearch
* PostGIS
* RabbitMQ
* Kafka
* API Gateway
* Service Discovery platform
* Notification platform
* Payment system
* Custom database
* Full SaaS dashboard
* Complex authentication system
* Terraform
* Large AWS architecture

Do not introduce technologies merely to make the project look impressive.

Every technology must have a clear architectural reason.

---

# 4. Core Architecture

Initial architecture:

Client
|
v
Go API Server
|
v
Job Queue
|
v
Worker Pool
|
v
Job Handler
|
v
PostgreSQL

Final architecture:

```
                Client
                  |
                  v
           +--------------+
           |   Go API     |
           | REST / gRPC  |
           +------+-------+
                  |
                  v
           +--------------+
           |    Redis     |
           | Job Queue /  |
           | Coordination |
           +------+-------+
                  |
      +-----------+-----------+
      |           |           |
      v           v           v
  Worker 1    Worker 2    Worker N
      |           |           |
      +-----------+-----------+
                  |
                  v
           +--------------+
           | PostgreSQL   |
           | Job State    |
           +--------------+
```

Observability:

Prometheus
OpenTelemetry
Structured Logging
pprof

---

# 5. Core Concepts

The system revolves around a Job.

A Job represents a unit of asynchronous work.

Example jobs:

* Send an email
* Generate a report
* Resize an image
* Process a CSV
* Generate a PDF
* Run a data-processing task

The system should support multiple job types through a handler/registry mechanism.

Example:

```
email.send
report.generate
image.resize
csv.process
```

The system must NOT implement complex business logic for these jobs.

Use simple demo handlers.

---

# 6. Job Lifecycle

A job should have the following lifecycle:

```
QUEUED
   |
   v
RUNNING
   |
   +--------> COMPLETED
   |
   +--------> FAILED
                 |
                 v
              RETRYING
                 |
                 v
              QUEUED
```

If maximum retries are exceeded:

```
FAILED
   |
   v
DEAD_LETTER
```

Supported states:

* QUEUED
* RUNNING
* COMPLETED
* FAILED
* RETRYING
* CANCELLED
* DEAD_LETTER

Do not create unnecessary states.

---

# 7. Job Model

A job should contain approximately:

```
ID
Type
Payload
Status
Attempts
MaxRetries
CreatedAt
StartedAt
CompletedAt
NextRetryAt
Error
WorkerID
IdempotencyKey
```

Use appropriate PostgreSQL types.

Payload can initially be stored as JSON/JSONB.

---

# 8. API Requirements

Implement a REST API.

## Create Job

POST /jobs

Example request:

{
"type": "email.send",
"payload": {
"to": "[user@example.com](mailto:user@example.com)",
"subject": "Welcome"
},
"max_retries": 3
}

Response:

{
"job_id": "uuid",
"status": "QUEUED"
}

---

## Get Job

GET /jobs/{id}

Response:

{
"id": "uuid",
"type": "email.send",
"status": "COMPLETED",
"attempts": 1,
"created_at": "...",
"completed_at": "..."
}

---

## Cancel Job

POST /jobs/{id}/cancel

Cancellation should work for queued jobs and, where possible, running jobs using Go context cancellation.

---

## Health

GET /health

Returns service health.

---

## Readiness

GET /ready

Returns whether the service is ready to accept work.

---

# 9. Go Concurrency Requirements

Go concurrency is a central part of this project.

The implementation must demonstrate:

* Goroutines
* Channels
* Buffered channels
* select
* sync.WaitGroup
* Mutex/RWMutex where appropriate
* context.Context
* Atomic operations where appropriate
* Graceful shutdown
* Worker pools

Do NOT use concurrency primitives unnecessarily.

Every synchronization mechanism should have a clear reason.

---

# 10. Worker Pool

Implement a configurable worker pool.

Example:

```
WorkerPool
    |
    +-- Worker 1
    +-- Worker 2
    +-- Worker 3
    +-- Worker 4
    +-- Worker N
```

Configuration:

```
WORKER_COUNT=10
```

Workers should:

1. Receive a job.
2. Mark it RUNNING.
3. Execute the registered handler.
4. Record success/failure.
5. Retry when appropriate.
6. Update persistent state.
7. Continue processing.

Workers should terminate gracefully.

---

# 11. Job Handler Architecture

Create a handler interface similar to:

```
type Handler interface {
    Execute(ctx context.Context, payload []byte) error
}
```

Create a registry:

```
HandlerRegistry
```

Example:

```
email.send
report.generate
image.resize
```

Handlers should be independently testable.

Do not couple handlers to the API layer.

---

# 12. Context and Cancellation

Use context.Context for:

* HTTP request cancellation
* Job execution timeout
* Worker shutdown
* Job cancellation

Example concept:

```
Job Context
    |
    +-- timeout
    +-- cancellation
    +-- worker shutdown
```

A long-running handler should stop when its context is cancelled.

---

# 13. Retry System

Implement configurable retries.

Default:

```
max_retries = 3
```

Retry strategy:

```
Attempt 1
   |
   | failure
   v
Backoff
   |
Attempt 2
   |
   | failure
   v
Backoff
   |
Attempt 3
```

Use exponential backoff.

Example:

```
delay = baseDelay * 2^attempt
```

Add a maximum backoff limit.

Avoid busy waiting.

---

# 14. Idempotency

The system should support idempotency keys.

Example:

POST /jobs

Header:

```
Idempotency-Key: abc123
```

If the same request is submitted again with the same idempotency key, the system should not create duplicate work.

This demonstrates an important real-world distributed-systems concept.

---

# 15. PostgreSQL

PostgreSQL is the source of truth for durable job state.

Use PostgreSQL for:

* Job metadata
* Job state
* Attempt count
* Errors
* timestamps
* worker assignment where necessary

Implement:

* migrations
* indexes
* transactions
* connection pooling

Important indexes should be based on actual query patterns.

Do not add random indexes.

---

# 16. Redis

Redis should NOT simply duplicate PostgreSQL.

Use Redis for distributed coordination and fast queue/worker operations.

Possible responsibilities:

* Queue coordination
* Worker heartbeat
* Worker registration
* Distributed locks where necessary
* Temporary state

Clearly document what belongs in Redis vs PostgreSQL.

PostgreSQL = durable source of truth.

Redis = fast coordination/runtime state.

---

# 17. Distributed Worker Architecture

Eventually the system must support:

```
API Server
     |
   Redis
     |
+----+----+----+
|    |    |    |
W1   W2   W3   W4
```

Workers should be separate processes.

The system must work even when workers are started independently.

Example:

```
./worker --id worker-1

./worker --id worker-2

./worker --id worker-3
```

---

# 18. Worker Registration

Workers should register themselves.

Worker metadata should include:

```
worker_id
hostname/container_id
started_at
last_heartbeat
status
```

Possible worker states:

```
ACTIVE
DEAD
DRAINING
```

Keep the model simple.

---

# 19. Worker Heartbeat

Workers periodically send heartbeats.

Example:

```
heartbeat interval = 5 seconds
```

If a worker stops sending heartbeats beyond a configured threshold, it should be considered dead.

Example:

```
heartbeat interval = 5s
failure threshold = 15s
```

These values must be configurable.

---

# 20. Worker Failure Detection

Example:

Worker 1 receives:

```
Job A
```

Worker 1 crashes.

The system must eventually detect:

```
Worker 1 DEAD
```

and make the job eligible for recovery/reprocessing.

Avoid silently losing jobs.

This is one of the most important distributed-systems features of the project.

---

# 21. Delivery Semantics

The project should explicitly document its delivery guarantee.

Target:

> At-least-once job execution.

This means a job may execute more than once under certain failure scenarios.

Therefore:

> Job handlers must be designed to be idempotent where necessary.

Do NOT claim exactly-once execution.

---

# 22. Dead Letter Queue

If a job exceeds maximum retry attempts:

```
FAILED
   |
   v
DEAD_LETTER
```

Expose an API to inspect dead-letter jobs.

Optional:

```
POST /jobs/{id}/retry
```

to manually retry a dead-letter job.

---

# 23. Graceful Shutdown

When receiving SIGTERM/SIGINT:

1. Stop accepting new work.
2. Stop fetching new jobs.
3. Allow active jobs to finish within a configurable timeout.
4. Cancel remaining jobs if necessary.
5. Persist state.
6. Close Redis.
7. Close PostgreSQL.
8. Exit cleanly.

This must be tested.

---

# 24. Observability

Implement:

## Metrics

Prometheus metrics such as:

```
jobs_submitted_total
jobs_completed_total
jobs_failed_total
jobs_retried_total
jobs_dead_letter_total
jobs_running
job_execution_duration_seconds
queue_depth
worker_count
worker_failures_total
```

Avoid creating dozens of useless metrics.

---

# 25. Logging

Use structured logging.

Every important event should include relevant fields:

```
job_id
worker_id
job_type
attempt
status
```

Example:

```
job_id=123
worker_id=worker-2
type=email.send
attempt=2
status=failed
```

Do not log sensitive payloads.

---

# 26. OpenTelemetry

Add tracing after the core system works.

Trace:

```
HTTP Request
     |
     v
Job Creation
     |
     v
Queue
     |
     v
Worker
     |
     v
Job Handler
```

Do not spend excessive time building a sophisticated tracing platform.

---

# 27. Profiling

Use Go pprof.

Profile:

* CPU
* memory
* goroutines

Use profiling to identify actual bottlenecks.

Do not optimize without measurements.

---

# 28. Testing Requirements

Minimum testing:

## Unit tests

Test:

* job state transitions
* retry calculation
* backoff
* idempotency
* handlers
* worker behavior
* cancellation

## Integration tests

Test:

* API + PostgreSQL
* API + Redis
* worker + Redis
* worker + PostgreSQL

## Failure tests

Simulate:

* worker crash
* Redis failure
* PostgreSQL failure
* handler failure
* timeout
* cancellation
* duplicate submission

---

# 29. Race Detection

The project must be tested with:

```
go test -race ./...
```

Any race condition discovered must be fixed rather than ignored.

---

# 30. Benchmarking

Benchmark:

1. Single worker
2. 5 workers
3. 10 workers
4. 25 workers
5. 50 workers

Measure:

* throughput
* average latency
* P95 latency
* P99 latency
* CPU
* memory
* queue depth

Use actual measured values in the README.

Never invent benchmark numbers.

---

# 31. Docker

Provide Docker support.

At minimum:

```
API
Worker
PostgreSQL
Redis
```

Use Docker Compose for local development.

Example:

```
docker compose up
```

should start the complete local system.

---

# 32. Kubernetes

Kubernetes is OPTIONAL and should only be implemented after the application is stable.

If implemented, deploy:

```
API
Worker
Redis
PostgreSQL
```

The main purpose is to demonstrate:

* multiple worker replicas
* horizontal scaling
* graceful termination
* health checks

Do not build a complex Kubernetes platform.

---

# 33. Project Structure

Preferred structure:

```
distributed-job-queue/

├── cmd/
│   ├── api/
│   └── worker/
│
├── internal/
│   ├── job/
│   │   ├── model.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── handler.go
│   │
│   ├── queue/
│   │   ├── queue.go
│   │   └── worker.go
│   │
│   ├── retry/
│   │   └── retry.go
│   │
│   ├── storage/
│   │   ├── postgres.go
│   │   └── redis.go
│   │
│   └── http/
│       └── handlers.go
│
├── migrations/
│
├── tests/
│
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

The structure may evolve if there is a strong reason.

Do not create unnecessary abstractions.

---

# 34. Development Order

The AI agent MUST follow this order.

## Phase 1 — Foundation

* Initialize Go project
* Basic HTTP server
* Configuration
* Logging
* PostgreSQL connection
* Database migration
* Job model

## Phase 2 — Local Queue

* In-memory queue
* Worker pool
* Job handler registry
* Job submission
* Job status
* Graceful shutdown

## Phase 3 — Persistence

* PostgreSQL repository
* Durable job state
* Transactions
* Job state transitions
* Integration tests

## Phase 4 — Reliability

* Retry
* Exponential backoff
* Timeout
* Cancellation
* Idempotency
* Dead-letter state

## Phase 5 — Distributed Workers

* Redis
* Separate worker process
* Worker registration
* Heartbeat
* Worker failure detection
* Job recovery

## Phase 6 — Production Engineering

* Prometheus
* Structured logs
* OpenTelemetry
* pprof
* Docker
* Load testing
* Benchmarks

## Phase 7 — Optional

* gRPC
* Scheduled jobs
* Priority queues
* Kubernetes

Do not jump ahead.

---

# 35. Definition of MVP

The MVP is COMPLETE when:

* [ ] Client can submit a job.
* [ ] Job is persisted in PostgreSQL.
* [ ] Worker receives the job.
* [ ] Worker executes the job concurrently.
* [ ] Job status can be queried.
* [ ] Successful jobs become COMPLETED.
* [ ] Failed jobs become FAILED.
* [ ] Jobs can retry.
* [ ] Retry uses exponential backoff.
* [ ] Jobs can timeout.
* [ ] Jobs can be cancelled.
* [ ] Server shuts down gracefully.
* [ ] Unit tests exist.
* [ ] Integration tests exist.
* [ ] `go test -race ./...` passes.

The MVP should NOT require Kubernetes, OpenTelemetry, or distributed workers.

---

# 36. Definition of Final Version

The final project is COMPLETE when:

* [ ] Multiple worker processes can execute jobs.
* [ ] Redis is used for distributed coordination.
* [ ] Workers have heartbeats.
* [ ] Worker failures are detected.
* [ ] Jobs can be recovered after worker failure.
* [ ] At-least-once semantics are documented.
* [ ] Idempotency is implemented.
* [ ] Dead-letter handling exists.
* [ ] Prometheus metrics exist.
* [ ] Structured logging exists.
* [ ] OpenTelemetry tracing exists.
* [ ] pprof profiling is available.
* [ ] Docker Compose works.
* [ ] Load tests exist.
* [ ] Benchmarks exist.
* [ ] Race detector passes.
* [ ] README contains architecture and design decisions.
* [ ] Actual benchmark results are documented.

---

# 37. Time Constraint

I can spend approximately:

```
2 hours/day
```

on this project.

I already know Go basics.

Therefore the AI agent must optimize for:

```
Learning + implementation + interview value
```

rather than maximum feature count.

Target completion:

```
~40–45 days
```

Do NOT expand the project beyond this target without explicit approval.

---

# 38. Daily Development Rule

For every development session:

1. Explain the concept I need.
2. Explain why the project needs it.
3. Show the design.
4. Implement the smallest useful version.
5. Test it.
6. Explain the important code.
7. Give me a small exercise if useful.
8. Update the project progress.
9. Tell me what remains.
10. Stop when the planned work for the day is complete.

Do not dump huge amounts of code without explanation.

---

# 39. AI Agent Tracking

Maintain a project progress file:

```
PROJECT_PROGRESS.md
```

It should contain:

## Current Phase

Example:

```
Phase 3 — Persistence
```

## Completed

* [x] Go project initialized
* [x] HTTP server
* [x] Worker pool
* [x] Job model

## In Progress

* [ ] PostgreSQL repository

## Next

* [ ] Database migrations
* [ ] Job persistence tests

## Architecture Decisions

Record important decisions such as:

```
PostgreSQL = durable source of truth
Redis = runtime coordination
At-least-once delivery
Go channels = local worker queue
```

## Known Issues

Track bugs and technical debt.

## Time Tracking

Record approximately:

```
Day 1: 2h
Day 2: 1.5h
Day 3: 2h
```

The agent must NOT assume work was completed unless I confirm it.

---

# 40. AI Agent Rules

The AI agent must:

### Rule 1

Never expand scope automatically.

### Rule 2

Never introduce a technology just because it is popular.

### Rule 3

Prefer standard Go libraries when appropriate.

### Rule 4

Explain concurrency carefully because concurrency is a core learning objective.

### Rule 5

Never hide complexity behind libraries if implementing the concept manually would provide meaningful learning value.

### Rule 6

Do not implement Kafka/RabbitMQ because the purpose is to understand job queues and distributed workers rather than reuse an existing message broker.

### Rule 7

Do not build a frontend.

### Rule 8

Do not rewrite working code without a measurable reason.

### Rule 9

Prefer simple architecture first, then evolve it.

### Rule 10

Every major feature must have tests.

### Rule 11

Measure performance before optimizing.

### Rule 12

Do not claim scalability without benchmarks.

### Rule 13

Do not claim exactly-once delivery.

### Rule 14

Keep the project achievable within approximately 40–45 days at 2 hours/day.

### Rule 15

If a requested feature conflicts with the scope, explain the tradeoff and ask before expanding the project.

---

# 41. Learning Objectives

By completing this project, I should be able to confidently explain:

### Go

* Goroutines
* Channels
* Worker pools
* Mutexes
* WaitGroups
* Context
* Cancellation
* Race conditions
* Graceful shutdown
* Profiling

### Backend

* REST APIs
* PostgreSQL
* Transactions
* Connection pooling
* Redis
* Structured logging
* Metrics

### Distributed Systems

* Producer/consumer
* At-least-once delivery
* Idempotency
* Retries
* Exponential backoff
* Dead-letter queues
* Heartbeats
* Failure detection
* Worker recovery
* Distributed coordination
* Failure scenarios

### Production Engineering

* Docker
* Observability
* Load testing
* Benchmarking
* Kubernetes fundamentals

---

# 42. Final Resume Positioning

The project should eventually be described as:

**Distributed Job Queue & Worker Platform — Go**

The final resume bullets should emphasize:

* Go concurrency
* distributed workers
* fault tolerance
* retries
* idempotency
* Redis
* PostgreSQL
* observability
* performance benchmarks

Do NOT fill the resume bullets with every technology used.

The project should communicate:

> Built infrastructure that reliably executes asynchronous workloads across concurrent and distributed Go workers.

---

# 43. Success Criteria

The project is successful if, during an interview, I can explain:

1. Why Go was chosen.
2. How the worker pool works.
3. How goroutines and channels are used.
4. How jobs are persisted.
5. What happens when a worker crashes.
6. How duplicate execution is handled.
7. Why the system uses at-least-once delivery.
8. How retries work.
9. How idempotency works.
10. Why Redis and PostgreSQL have different responsibilities.
11. How the system scales horizontally.
12. What happens during graceful shutdown.
13. Where bottlenecks occur.
14. How the system was benchmarked.
15. What tradeoffs were made.

The goal is **deep understanding and a strong engineering project**, not maximum feature count.

# END OF PROJECT SPECIFICATION
