package fakes

import "github.com/cloudfoundry-incubator/notifications/gobble"

type Queue struct {
	jobs         chan gobble.Job
	pk           int
	EnqueueError error
}

func NewQueue() *Queue {
	return &Queue{
		jobs: make(chan gobble.Job),
	}
}

func (fake *Queue) Enqueue(job gobble.Job) (gobble.Job, error) {
	if fake.EnqueueError != nil {
		return job, fake.EnqueueError
	}
	fake.pk++
	job.ID = fake.pk
	go func(job gobble.Job) {
		fake.jobs <- job
	}(job)

	return job, nil
}

func (fake *Queue) Reserve(string) <-chan gobble.Job {
	return fake.jobs
}

func (fake *Queue) Dequeue(job gobble.Job) {
}

func (fake *Queue) Requeue(job gobble.Job) {
	go func(job gobble.Job) {
		fake.jobs <- job
	}(job)
}

func (fake *Queue) Unlock() {}
