package collections

import "github.com/cloudfoundry-incubator/notifications/models"

type DatabaseInterface interface {
	models.DatabaseInterface
}

type ConnectionInterface interface {
	models.ConnectionInterface
}
