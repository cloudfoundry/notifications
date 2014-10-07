package fakes

import "github.com/cloudfoundry-incubator/notifications/gobble"

type FakeQueue struct {
    jobs         chan gobble.Job
    pk           int
    EnqueueError error
}

func NewFakeQueue() *FakeQueue {
    return &FakeQueue{
        jobs: make(chan gobble.Job),
    }
}

func (fake *FakeQueue) Enqueue(job gobble.Job) (gobble.Job, error) {
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

func (fake *FakeQueue) Reserve(string) <-chan gobble.Job {
    return fake.jobs
}

func (fake *FakeQueue) Dequeue(job gobble.Job) {
}

func (fake *FakeQueue) Requeue(job gobble.Job) {
    go func(job gobble.Job) {
        fake.jobs <- job
    }(job)
}
