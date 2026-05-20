package main

import (
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

type Queue interface {
	Enqueue(job Job) error
	Dequeue() (Job, error)
}

type Worker interface {
	Work(job Job) error
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
