package main

import "fmt"

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,

		cmd: make(chan int),
	}
}

type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job

	cmd chan int
}

func (d *Dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		fmt.Printf("START select dispatch \n")
		select {
		case job := <-d.jobQueue: //receive job request
			fmt.Printf("START dispatch got job %s \n", job.Name)
			go func() {
				fmt.Printf("fetching workerJobQueue for: %s\n", job.Name)
				workerJobQueue := <-d.workerPool
				fmt.Printf("adding %s to workerJobQueue\n", job.Name)
				workerJobQueue <- job
			}()
		case cmd := <-d.cmd: //add/remove worker
			switch cmd {
			case addWorker:
				fmt.Println("add worker")
				//TODO worker id
				worker := NewWorker(0, d.workerPool)
				worker.start()
			case removeWorker:
				fmt.Println("remove worker")
				//TODO
			}
		}

	}
}

const (
	addWorker    = iota
	removeWorker = iota
)

func (d *Dispatcher) Add(workers int) {
	if workers == 0 {
		return
	}

	cmd := addWorker
	if workers < 0 {
		workers = workers * -1
		cmd = removeWorker
	}

	for i := 0; i < workers; i++ {
		d.cmd <- cmd
	}
}
