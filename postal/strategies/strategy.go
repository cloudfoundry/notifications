package strategies

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

type StrategyInterface interface {
	Dispatch(clientID, guid, vcapRequestID string, requestReceived time.Time, options postal.Options, conn models.ConnectionInterface) ([]Response, error)
}
