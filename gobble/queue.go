package gobble

import (
    "database/sql"
    "math/rand"
    "time"

    "github.com/coopernurse/gorp"
)

var WaitMaxDuration = 5 * time.Second

type QueueInterface interface {
    Enqueue(Job) Job
    Reserve(string) <-chan Job
    Dequeue(Job)
}

type Queue struct{}

func NewQueue() *Queue {
    return &Queue{}
}

func (queue *Queue) Enqueue(job Job) Job {
    err := Database().Connection.Insert(&job)
    if err != nil {
        panic(err)
    }

    return job
}

func (queue *Queue) Reserve(workerID string) <-chan Job {
    channel := make(chan Job)

    go func() {
        job := Job{}
        for job.ID == 0 {
            var err error

            job = queue.findJob()
            job, err = queue.updateJob(job, workerID)
            if err != nil {
                if _, ok := err.(gorp.OptimisticLockError); ok {
                    job = Job{}
                    continue
                } else {
                    panic(err)
                }
            }
        }

        channel <- job
    }()

    return channel
}

func (queue *Queue) Dequeue(job Job) {
    _, err := Database().Connection.Delete(&job)
    if err != nil {
        panic(err)
    }
}

func (queue *Queue) findJob() Job {
    job := Job{}
    for job.ID == 0 {
        err := Database().Connection.SelectOne(&job, "SELECT * FROM `jobs` WHERE `worker_id` = \"\" LIMIT 1")
        if err != nil {
            if err == sql.ErrNoRows {
                job = Job{}
                queue.waitUpTo(WaitMaxDuration)
                continue
            }
            panic(err)
        }
    }
    return job
}

func (queue *Queue) updateJob(job Job, workerID string) (Job, error) {
    job.WorkerID = workerID
    _, err := Database().Connection.Update(&job)
    if err != nil {
        return job, err
    }
    return job, nil
}

func (queue *Queue) waitUpTo(max time.Duration) {
    rand.Seed(time.Now().UnixNano())
    waitTime := rand.Int63n(int64(max))
    <-time.After(time.Duration(waitTime))
}
