package fakes

import "github.com/cloudfoundry-incubator/notifications/gobble"

type Queue struct {
	jobs          map[int]gobble.Job
	availableJobs chan gobble.Job
	pk            int
	EnqueueError  error
}

func NewQueue() *Queue {
	return &Queue{
		jobs:          make(map[int]gobble.Job),
		availableJobs: make(chan gobble.Job),
	}
}

func (q *Queue) Enqueue(job gobble.Job) (gobble.Job, error) {
	if q.EnqueueError != nil {
		return job, q.EnqueueError
	}

	if job.ID == 0 {
		q.pk++
		job.ID = q.pk
	}

	go func() {
		q.availableJobs <- job
	}()

	q.jobs[job.ID] = job

	return job, nil
}

func (q *Queue) Reserve(string) <-chan gobble.Job {
	jobs := make(chan gobble.Job, 1)

	go func() {
		job := <-q.availableJobs
		jobs <- job
	}()

	return jobs
}

func (q *Queue) Dequeue(job gobble.Job) {
	delete(q.jobs, job.ID)
}

func (q *Queue) Requeue(job gobble.Job) {
	q.Enqueue(job)
}

func (q *Queue) Unlock() {}

func (q *Queue) Len() (int, error) {
	return len(q.jobs), nil
}

func (q *Queue) RetryQueueLengths() (map[int]int, error) {
	lengths := map[int]int{}

	for _, job := range q.jobs {
		lengths[job.RetryCount] = lengths[job.RetryCount] + 1
	}

	return lengths, nil
}
