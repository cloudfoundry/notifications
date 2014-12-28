package gobble

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/coopernurse/gorp"
)

var WaitMaxDuration = 5 * time.Second

type QueueInterface interface {
	Enqueue(Job) (Job, error)
	Reserve(string) <-chan Job
	Dequeue(Job)
	Requeue(Job)
}

type Queue struct {
	config Config
	closed bool
}

func NewQueue(config Config) *Queue {
	if config.WaitMaxDuration == 0 {
		config.WaitMaxDuration = WaitMaxDuration
	}

	return &Queue{
		config: config,
	}
}

func (queue *Queue) Enqueue(job Job) (Job, error) {
	err := Database().Connection.Insert(&job)
	if err != nil {
		return job, err
	}

	return job, nil
}

func (queue *Queue) Requeue(job Job) {
	_, err := Database().Connection.Update(&job)
	if err != nil {
		panic(err)
	}
}

func (queue *Queue) Reserve(workerID string) <-chan Job {
	channel := make(chan Job)
	go queue.reserve(channel, workerID)

	return channel
}

func (queue *Queue) Close() {
	queue.closed = true
}

func (queue *Queue) reserve(channel chan Job, workerID string) {
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
	if queue.closed {
		queue.updateJob(job, "")
		return
	}

	channel <- job
}

func (queue *Queue) Dequeue(job Job) {
	_, err := Database().Connection.Delete(&job)
	if err != nil {
		panic(err)
	}
}

func (queue Queue) Unlock() {
	_, err := Database().Connection.Exec("UPDATE `jobs` set `worker_id` = \"\" WHERE `worker_id` != \"\"")
	if err != nil {
		panic(err)
	}
}

func (queue *Queue) findJob() Job {
	job := Job{}
	for job.ID == 0 {
		err := Database().Connection.SelectOne(&job, "SELECT * FROM `jobs` WHERE `worker_id` = \"\" AND `active_at` <= ? LIMIT 1", time.Now())
		if err != nil {
			if err == sql.ErrNoRows {
				job = Job{}
				queue.waitUpTo(queue.config.WaitMaxDuration)
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
