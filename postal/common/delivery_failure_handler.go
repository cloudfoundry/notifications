package common

import (
	"math"
	"time"

	"github.com/pivotal-golang/lager"
	"github.com/rcrowley/go-metrics"
)

type Retryable interface {
	Retry(duration time.Duration)
	State() (retryCount int, activeAt time.Time)
}

type DeliveryFailureHandler struct {
	numRetries int
}

func NewDeliveryFailureHandler(numRetries int) DeliveryFailureHandler {
	return DeliveryFailureHandler{
		numRetries: numRetries,
	}
}

func (h DeliveryFailureHandler) Handle(job Retryable, logger lager.Logger) {
	retryCount, _ := job.State()
	if retryCount >= h.numRetries {
		return
	}

	duration := time.Duration(int64(math.Pow(2, float64(retryCount)))) * time.Minute
	job.Retry(duration)

	retryCount, activeAt := job.State()
	logger.Info("delivery-failed-retrying", lager.Data{
		"retry_count": retryCount,
		"active_at":   activeAt.Format(time.RFC3339),
	})

	metrics.GetOrRegisterCounter("notifications.worker.retry", nil).Inc(1)
}
