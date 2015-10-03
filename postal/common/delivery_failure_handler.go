package common

import (
	"math"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/pivotal-golang/lager"
)

type Retryable interface {
	Retry(duration time.Duration)
	State() (retryCount int, activeAt time.Time)
}

type DeliveryFailureHandler struct{}

func NewDeliveryFailureHandler() DeliveryFailureHandler {
	return DeliveryFailureHandler{}
}

func (h DeliveryFailureHandler) Handle(job Retryable, logger lager.Logger) {
	retryCount, _ := job.State()
	if retryCount > 9 {
		return
	}

	duration := time.Duration(int64(math.Pow(2, float64(retryCount)))) * time.Minute
	job.Retry(duration)

	retryCount, activeAt := job.State()
	logger.Info("delivery-failed-retrying", lager.Data{
		"retry_count": retryCount,
		"active_at":   activeAt.Format(time.RFC3339),
	})

	// TODO: (rm) find way to test this without having to mock out globals
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.worker.retry",
	}).Log()
}
