package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type DatabaseInterface interface {
	models.DatabaseInterface
}

type ConnectionInterface interface {
	models.ConnectionInterface
}
