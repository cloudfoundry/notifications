package gobble

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID          int       `db:"id"`
	WorkerID    string    `db:"worker_id"`
	Payload     string    `db:"payload"`
	Version     int64     `db:"version"`
	RetryCount  int       `db:"retry_count"`
	ActiveAt    time.Time `db:"active_at"`
	ShouldRetry bool      `db:"-"`
}

func NewJob(data interface{}) *Job {
	payload, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return &Job{
		Payload: string(payload),
	}
}

func (job Job) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(job.Payload), v)
}

func (job *Job) Retry(duration time.Duration) {
	job.WorkerID = ""
	job.RetryCount++
	job.ActiveAt = time.Now().Add(duration)
	job.ShouldRetry = true
}

func (job *Job) State() (int, time.Time) {
	return job.RetryCount, job.ActiveAt
}
