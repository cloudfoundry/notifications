package gobble

import (
	"fmt"
	"os"
)

type heartbeater interface {
	Beat(*Job)
	Halt()
}

type Worker struct {
	ID       string
	queue    QueueInterface
	callback func(*Job)
	beater   heartbeater
	halt     chan bool
}

func NewWorker(id int, queue QueueInterface, callback func(*Job), beater heartbeater) Worker {
	return Worker{
		ID:       fmt.Sprintf("worker-%d-%d", id, os.Getpid()),
		queue:    queue,
		callback: callback,
		beater:   beater,
		halt:     make(chan bool),
	}
}

func (worker *Worker) Perform() int {
	select {
	case job := <-worker.queue.Reserve(worker.ID):
		go worker.beater.Beat(job)
		defer worker.beater.Halt()
		worker.callback(job)

		if job.ShouldRetry {
			worker.queue.Requeue(job)
		} else {
			worker.queue.Dequeue(job)
		}
		return 0
	case <-worker.halt:
		return 1
	}
}

func (worker *Worker) Work() {
	go func() {
		for {
			if worker.Perform() != 0 {
				return
			}
		}
	}()
}

func (worker *Worker) Halt() {
	worker.halt <- true
}
