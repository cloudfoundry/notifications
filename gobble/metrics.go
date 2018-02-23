package gobble

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type QueueGauge struct {
	queue queue
	timer <-chan time.Time
}

type queue interface {
	Len() (int, error)
}

func NewQueueGauge(queue queue, timer <-chan time.Time) QueueGauge {
	return QueueGauge{
		queue: queue,
		timer: timer,
	}
}

func (g QueueGauge) Run() {
	for range g.timer {
		ql, _ := g.queue.Len()

		metrics.GetOrRegisterGauge("notifications.queue.length", nil).Update(int64(ql))
	}
}
