# Go Job Queue — Checkpoint Roadmap

> Assuming **1.5–2 hours/day** on work days, **3–4 hours** on weekends.

---

## Checkpoint 1 — Core Loop
**Week 1 · ~8 hours**

**Goal:** Job goes in, job comes out, worker processes it. Nothing fancy.

- Define `Job` struct and status constants (`pending`, `running`, `done`, `failed`)
- `Queue` interface with in-memory channel implementation
- `Enqueue` / `Dequeue` working
- Single worker that picks up job and processes it (just print for now)
- `SubmitJob` HTTP endpoint — returns job ID
- `GetJobStatus` HTTP endpoint — returns current status

✅ **Done when:** You can `curl` to submit a job and check its status.

---

## Checkpoint 2 — Worker Pool
**Week 1–2 · ~6 hours**

**Goal:** Multiple workers running concurrently, controlled concurrency.

- Replace single worker with a goroutine pool (N workers from env var `WORKER_COUNT`)
- Use `errgroup` to manage worker lifecycle
- Each worker loops: dequeue → process → update status
- Basic structured logging with `slog` on every state change

✅ **Done when:** 10 jobs submitted, processed concurrently across N workers, logs show which worker handled which job.

---

## Checkpoint 3 — Graceful Shutdown
**Week 2 · ~4 hours**

**Goal:** App shuts down cleanly without abandoning in-flight jobs.

- Listen for `SIGTERM` / `SIGINT` using `signal.NotifyContext`
- Cancel root context on signal
- Workers finish current job before exiting
- HTTP server drains connections with timeout

✅ **Done when:** Kill the app mid-processing, all started jobs complete before exit, no jobs lost.

---

## Checkpoint 4 — Retry + Dead Letter Queue
**Week 2–3 · ~6 hours**

**Goal:** Failed jobs retry automatically, hopeless jobs are parked.

- Add `Retries` and `MaxRetries` to `Job` struct
- Implement exponential backoff with jitter
- After `MaxRetries` exceeded → move to dead letter queue (in-memory list for now)
- `GetDeadLetterJobs` endpoint — list all dead jobs for inspection

✅ **Done when:** A job that always fails retries 3 times with increasing delay then appears in dead letter queue.

---

## Checkpoint 5 — Swap to Redis
**Week 3 · ~8 hours**

**Goal:** Replace in-memory queue with Redis. App survives restarts.

- Add Redis client (`go-redis`)
- Implement `RedisQueue` satisfying your existing `Queue` interface
- Store job data as Redis Hash (`job:{id}`)
- Queue as Redis List (`queue:pending`, `queue:dead`)
- Use `BRPOP` for blocking dequeue
- Swap implementation via config flag — in-memory still works

✅ **Done when:** Restart the app, previously submitted jobs still exist and get processed.

---

## Checkpoint 6 — Crash Recovery (Visibility Timeout)
**Week 3–4 · ~6 hours**

**Goal:** Jobs don't disappear if a worker crashes mid-processing.

- When worker dequeues a job, move it to `queue:processing` with a timestamp
- Reaper goroutine runs every 10s — finds jobs stuck in `processing` longer than 30s
- Reaper requeues stalled jobs back to `queue:pending`
- This gives you **at-least-once delivery**

✅ **Done when:** Kill a worker mid-job, reaper detects it within 30s and requeues, job eventually completes.

---

## Checkpoint 7 — Rate Limiting
**Week 4 · ~5 hours**

**Goal:** Clients can't spam the queue. Implement token bucket from scratch.

- `TokenBucket` struct — tokens, refill rate, last refill time
- Middleware on `SubmitJob` — check bucket before accepting job
- Per client IP or API key
- Return `429 Too Many Requests` when bucket empty
- No libraries — implement the math yourself

✅ **Done when:** Spam 20 requests quickly, first N succeed, rest get 429.

---

## Checkpoint 8 — Observability
**Week 4–5 · ~5 hours**

**Goal:** You can see what's happening without reading logs line by line.

- Prometheus metrics endpoint `/metrics`
- Track: jobs submitted, jobs completed, jobs failed, queue depth, processing duration
- Structured log fields on every job event: `job_id`, `status`, `worker_id`, `duration_ms`
- `docker-compose.yml` with app + Redis + Prometheus

✅ **Done when:** Open Prometheus, query `jobs_completed_total`, see real numbers.

---

## Checkpoint 9 — Docker + Kubernetes
**Week 5 · ~5 hours**

**Goal:** Ship it like a real service.

- `Dockerfile` — multi-stage build, small final image
- `docker-compose.yml` — app + Redis together
- K8s manifests: `Deployment`, `Service`, `ConfigMap` for env vars
- `WORKER_COUNT` and Redis URL via env — no hardcoding

✅ **Done when:** `kubectl apply -f k8s/` and the whole system runs in a cluster.

---

## Checkpoint 10 — Polish + README
**Week 5–6 · ~4 hours**

**Goal:** Make it something you're proud to show in an interview.

- `README.md` with architecture diagram, design decisions, tradeoffs
- Document every env var
- Write table-driven tests for `TokenBucket`, retry logic, job state transitions
- Prepare your **3 talking points** for interviews

✅ **Done when:** A stranger can read the README and understand what you built and why.

---

## Summary

| # | Checkpoint | Hours | Week |
|---|---|---|---|
| 1 | Core loop | 8h | 1 |
| 2 | Worker pool | 6h | 1–2 |
| 3 | Graceful shutdown | 4h | 2 |
| 4 | Retry + Dead Letter Queue | 6h | 2–3 |
| 5 | Redis swap | 8h | 3 |
| 6 | Crash recovery | 6h | 3–4 |
| 7 | Rate limiting | 5h | 4 |
| 8 | Observability | 5h | 4–5 |
| 9 | Docker + Kubernetes | 5h | 5 |
| 10 | Polish + README | 4h | 5–6 |
| **Total** | | **~57h** | **~6 weeks** |

---

> **Rule:** Don't move to the next checkpoint until the current one fully works.
> A half-built system with 10 features is worse than a solid system with 5.
