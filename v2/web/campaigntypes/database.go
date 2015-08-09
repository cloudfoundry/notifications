package campaigntypes

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type DatabaseInterface interface {
	collections.DatabaseInterface
}

type ConnectionInterface interface {
	collections.ConnectionInterface
}
