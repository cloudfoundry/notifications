package mocks

import "github.com/cloudfoundry-incubator/notifications/gobble"

type Queue struct {
	Jobs          map[int]gobble.Job
	availableJobs chan gobble.Job
	pk            int
	EnqueueError  error
}

func NewQueue() *Queue {
	return &Queue{
		Jobs:          make(map[int]gobble.Job),
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

	q.Jobs[job.ID] = job

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
	delete(q.Jobs, job.ID)
}

func (q *Queue) Requeue(job gobble.Job) {
	q.Enqueue(job)
}

func (q *Queue) Unlock() {}

func (q *Queue) Len() (int, error) {
	return len(q.Jobs), nil
}

func (q *Queue) RetryQueueLengths() (map[int]int, error) {
	lengths := map[int]int{}

	for _, job := range q.Jobs {
		lengths[job.RetryCount] = lengths[job.RetryCount] + 1
	}

	return lengths, nil
}
