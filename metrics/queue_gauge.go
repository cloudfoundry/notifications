package metrics

import (
	"log"
	"strconv"
	"time"
)

type QueueGauge struct {
	logger *log.Logger
	queue  queue
	timer  <-chan time.Time
}

type queue interface {
	Len() (int, error)
	RetryQueueLengths() (map[int]int, error)
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
		retryCounts, _ := g.queue.RetryQueueLengths()

		NewMetric("gauge", map[string]interface{}{
			"name":  "notifications.queue.length",
			"value": length,
		}).LogWith(g.logger)

		for index := range make([]int, 11) {
			NewMetric("gauge", map[string]interface{}{
				"name": "notifications.queue.retry",
				"tags": map[string]interface{}{
					"count": strconv.Itoa(index),
				},
				"value": retryCounts[index],
			}).LogWith(g.logger)
		}
	}
}
