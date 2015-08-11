package collections

import "github.com/cloudfoundry-incubator/notifications/db"

type DatabaseInterface interface {
	db.DatabaseInterface
}

type ConnectionInterface interface {
	db.ConnectionInterface
}
