package mocks

import "time"

type GobbleJob struct {
	RetryCall struct {
		WasCalled bool
		Receives  struct {
			Duration time.Duration
		}
	}

	StateCall struct {
		Returns struct {
			Count int
			Time  time.Time
		}
	}
}

func NewGobbleJob() *GobbleJob {
	return &GobbleJob{}
}

func (j *GobbleJob) Retry(duration time.Duration) {
	j.RetryCall.WasCalled = true
	j.RetryCall.Receives.Duration = duration
}

func (j *GobbleJob) State() (int, time.Time) {
	return j.StateCall.Returns.Count, j.StateCall.Returns.Time
}
