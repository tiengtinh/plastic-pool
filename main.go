package main

import (
	"flag"
	"net/http"
	"os"
	"time"
)

func requestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the delay.
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate delay is in range 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Set name and validate value.
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	// Create Job and push the work onto the jobQueue.
	job := Job{Name: name, Delay: delay}
	jobQueue <- job

	// Render success.
	w.WriteHeader(http.StatusCreated)
}

func main() {
	// Parse the command-line flags.
	maxWorkers := flag.Int("max_workers", 5, "The number of workers to start")
	maxQueueSize := flag.Int("max_queue_size", 100, "The size of job queue")
	flag.Parse()

	// Create the job queue.
	jobQueue := make(chan Job, *maxQueueSize)

	// Start the dispatcher.
	dispatcher := NewDispatcher(jobQueue, *maxWorkers)
	dispatcher.run()

	// Start the HTTP handler.
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, jobQueue)
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
