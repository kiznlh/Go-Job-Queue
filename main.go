package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Job struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`    // e.g. "send_email", "resize_image", "generate_report"
	Payload   []byte    `json:"payload"` // the input data (JSON usually)
	Status    string    `json:"status"`  // "pending", "running", "done", "failed"
	CreatedAt time.Time `json:"time"`
	Retries   int       `json:"retries"`
}

func (job Job) String() string {
	out, err := json.Marshal(job)
	if err != nil {
		return "Error marshaling"
	}
	return string(out)
}

type Queue interface {
	Enqueue(job Job) error
	Dequeue() (Job, error)
}

type Worker interface {
	Work(job Job) error
}
type MemoryQueue struct {
	jobs []Job
}

func (q *MemoryQueue) Enqueue(job Job) error {
	q.jobs = append(q.jobs, job)
	return nil
}

type PrintWorker struct{}

func (w *PrintWorker) Work(job Job) error {
	fmt.Println(job)
	return nil
}

func (q *MemoryQueue) Dequeue() (Job, error) {
	if len(q.jobs) == 0 {
		return Job{}, errors.New("Empty Queue")
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job, nil
}

func main() {
	router := http.NewServeMux()

	// Add Job
	router.HandleFunc("POST /submitJob", func(w http.ResponseWriter, r *http.Request) {

	})

	// Get Job by ID
	router.HandleFunc("GET /getJobStatus/{id}", func(w http.ResponseWriter, r *http.Request) {

	})

	server := http.Server{
		Addr:    ":6969",
		Handler: router,
	}

	server.ListenAndServe()
}
