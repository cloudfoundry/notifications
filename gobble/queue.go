package gobble

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/gorp.v1"
)

var WaitMaxDuration = 5 * time.Second

type QueueInterface interface {
	Enqueue(*Job, ConnectionInterface) (*Job, error)
	Reserve(string) <-chan *Job
	Dequeue(*Job)
	Requeue(*Job)
	Len() (int, error)
	RetryQueueLengths() (map[int]int, error)
}

type clock interface {
	Now() time.Time
}

type Queue struct {
	config   Config
	database *DB
	clock    clock
	closed   bool
}

func NewQueue(database DatabaseInterface, clock clock, config Config) *Queue {
	if config.WaitMaxDuration == 0 {
		config.WaitMaxDuration = WaitMaxDuration
	}

	return &Queue{
		database: database.(*DB),
		clock:    clock,
		config:   config,
	}
}

func (queue *Queue) Enqueue(job *Job, connection ConnectionInterface) (*Job, error) {
	if (job.ActiveAt == time.Time{}) {
		job.ActiveAt = queue.clock.Now()
	}

	err := connection.Insert(job)
	if err != nil {
		return job, err
	}

	return job, nil
}

func (queue *Queue) Requeue(job *Job) {
	_, err := queue.database.Connection.Update(job)
	if err != nil {
		panic(err)
	}
}

func (queue *Queue) Len() (int, error) {
	length, err := queue.database.Connection.SelectInt("SELECT COUNT(*) FROM `jobs`")
	return int(length), err
}

func (queue *Queue) RetryQueueLengths() (map[int]int, error) {
	lengths := map[int]int{}

	type QueueLength struct {
		RetryCount int `db:"retry_count"`
		Count      int `db:"count"`
	}

	records, err := queue.database.Connection.Select(QueueLength{}, "SELECT retry_count, COUNT(*) AS count FROM `jobs` GROUP BY retry_count")
	if err != nil {
		return lengths, err
	}

	for _, value := range records {
		length := value.(*QueueLength)
		lengths[length.RetryCount] = length.Count
	}

	return lengths, nil
}

func (queue *Queue) Close() {
	queue.closed = true
}

func (queue *Queue) Reserve(workerID string) <-chan *Job {
	channel := make(chan *Job)
	go queue.reserve(channel, workerID)

	return channel
}

func (queue *Queue) reserve(channel chan *Job, workerID string) {
	var job *Job
	for job == nil {
		var err error

		job = queue.findJob()
		if queue.closed {
			return
		}

		job, err = queue.updateJob(job, workerID)
		if err != nil {
			if _, ok := err.(gorp.OptimisticLockError); ok {
				job = nil
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

func (queue *Queue) Dequeue(job *Job) {
	_, err := queue.database.Connection.Delete(job)
	if err != nil {
		if _, ok := err.(gorp.OptimisticLockError); ok && strings.Contains(err.Error(), "no row found") {
			return
		}
		panic(err)
	}
}

func (queue *Queue) findJob() *Job {
	var job *Job
	for job == nil {
		job = &Job{}
		now := time.Now()
		expired := now.Add(-2 * time.Minute)
		err := queue.database.Connection.SelectOne(job, "SELECT * FROM `jobs` WHERE ( `worker_id` = \"\" AND `active_at` <= ? ) OR `active_at` <= ? LIMIT 1", now, expired)
		if err != nil {
			if err == sql.ErrNoRows {
				job = nil
				queue.waitUpTo(queue.config.WaitMaxDuration)
				continue
			}
			panic(err)
		}
	}
	return job
}

func (queue *Queue) updateJob(job *Job, workerID string) (*Job, error) {
	if job == nil {
		return job, nil
	}

	job.WorkerID = workerID
	job.ActiveAt = time.Now()
	_, err := queue.database.Connection.Update(job)
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
