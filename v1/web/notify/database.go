package notify

import "github.com/cloudfoundry-incubator/notifications/v1/services"

type DatabaseInterface interface {
	services.DatabaseInterface
}

type ConnectionInterface interface {
	services.ConnectionInterface
}
