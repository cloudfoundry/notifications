package collections

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type DatabaseInterface interface {
	models.DatabaseInterface
}

type ConnectionInterface interface {
	models.ConnectionInterface
}
