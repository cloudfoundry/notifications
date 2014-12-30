package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type PersistenceProvider struct {
	database          *Database
	DatabaseWasCalled bool
	QueueWasCalled    bool
}

func NewPersistenceProvider(database *Database) *PersistenceProvider {
	return &PersistenceProvider{
		database: database,
	}
}

func (pp *PersistenceProvider) Database() models.DatabaseInterface {
	pp.DatabaseWasCalled = true
	return pp.database
}

func (pp *PersistenceProvider) Queue() gobble.QueueInterface {
	pp.QueueWasCalled = true
	return NewQueue()
}
