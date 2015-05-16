package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

type StrategyInterface interface {
	Dispatch(clientID, guid, vcapRequestID string, options postal.Options, conn models.ConnectionInterface) ([]Response, error)
}
