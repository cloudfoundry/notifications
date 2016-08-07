package mocks

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
	"gopkg.in/gorp.v1"
)

type Database struct {
	ConnectionCall struct {
		Returns struct {
			Connection db.ConnectionInterface
		}
	}

	RawConnectionCall struct {
		Returns struct {
			DB *sql.DB
		}
	}

	TraceOnCall struct {
		Receives struct {
			Prefix string
			Logger gorp.GorpLogger
		}
	}
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Connection() db.ConnectionInterface {
	return d.ConnectionCall.Returns.Connection
}

func (d *Database) RawConnection() *sql.DB {
	return d.RawConnectionCall.Returns.DB
}

func (d *Database) TraceOn(prefix string, logger gorp.GorpLogger) {
	d.TraceOnCall.Receives.Prefix = prefix
	d.TraceOnCall.Receives.Logger = logger
}
