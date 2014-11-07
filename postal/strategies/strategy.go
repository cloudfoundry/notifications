package strategies

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type StrategyInterface interface {
    Dispatch(clientID string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]Response, error)
    Trim([]byte) []byte
}
