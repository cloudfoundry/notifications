package mocks

import "github.com/cloudfoundry-incubator/notifications/gobble"

type Queue struct {
	EnqueueCall struct {
		Receives struct {
			Jobs       []*gobble.Job
			Connection gobble.ConnectionInterface
		}
		Returns struct {
			Job   *gobble.Job
			Error error
		}

		Hook func()
	}

	RequeueCall struct {
		Receives struct {
			Job *gobble.Job
		}
	}

	DequeueCall struct {
		Receives struct {
			Job *gobble.Job
		}
	}

	LenCall struct {
		Returns struct {
			Length int
			Error  error
		}
	}

	ReserveCall struct {
		Receives struct {
			ID string
		}
		Returns struct {
			Chan <-chan *gobble.Job
		}
	}

	RetryQueueLengthsCall struct {
		Returns struct {
			Lengths map[int]int
			Error   error
		}
	}
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(job *gobble.Job, connection gobble.ConnectionInterface) (*gobble.Job, error) {
	q.EnqueueCall.Receives.Jobs = append(q.EnqueueCall.Receives.Jobs, job)
	q.EnqueueCall.Receives.Connection = connection

	if q.EnqueueCall.Hook != nil {
		q.EnqueueCall.Hook()
	}

	return q.EnqueueCall.Returns.Job, q.EnqueueCall.Returns.Error
}

func (q *Queue) Dequeue(job *gobble.Job) {
	q.DequeueCall.Receives.Job = job
}

func (q *Queue) Requeue(job *gobble.Job) {
	q.RequeueCall.Receives.Job = job
}

func (q *Queue) Len() (int, error) {
	return q.LenCall.Returns.Length, q.LenCall.Returns.Error
}

func (q *Queue) Reserve(id string) <-chan *gobble.Job {
	q.ReserveCall.Receives.ID = id

	return q.ReserveCall.Returns.Chan
}

func (q *Queue) RetryQueueLengths() (map[int]int, error) {
	return q.RetryQueueLengthsCall.Returns.Lengths, q.RetryQueueLengthsCall.Returns.Error
}
