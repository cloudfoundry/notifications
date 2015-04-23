package metrics

import (
	"log"
	"time"
)

type QueueGauge struct {
	logger *log.Logger
	queue  queue
	timer  <-chan time.Time
}

type queue interface {
	Len() (int, error)
}

func NewQueueGauge(queue queue, logger *log.Logger, timer <-chan time.Time) QueueGauge {
	return QueueGauge{
		logger: logger,
		queue:  queue,
		timer:  timer,
	}
}

func (g QueueGauge) Run() {
	for _ = range g.timer {
		length, _ := g.queue.Len()

		NewMetric("gauge", map[string]interface{}{
			"length": length,
		}).LogWith(g.logger)
	}
}
