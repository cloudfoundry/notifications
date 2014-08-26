package gobble

import (
    "fmt"
    "os"
)

type Worker struct {
    ID       string
    queue    QueueInterface
    callback func(Job)
    halt     chan bool
}

func NewWorker(id int, queue QueueInterface, callback func(Job)) Worker {
    return Worker{
        ID:       fmt.Sprintf("worker-%d-%d", id, os.Getpid()),
        queue:    queue,
        callback: callback,
        halt:     make(chan bool),
    }
}

func (worker *Worker) Perform() int {
    select {
    case job := <-worker.queue.Reserve(worker.ID):
        worker.callback(job)
        worker.queue.Dequeue(job)
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
