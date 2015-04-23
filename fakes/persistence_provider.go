package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type PersistenceProvider struct {
	database                *Database
	gobbleDatabase          *GobbleDatabase
	DatabaseWasCalled       bool
	GobbleDatabaseWasCalled bool
}

func NewPersistenceProvider(database *Database, gobbleDatabase *GobbleDatabase) *PersistenceProvider {
	return &PersistenceProvider{
		database:       database,
		gobbleDatabase: gobbleDatabase,
	}
}

func (pp *PersistenceProvider) Database() models.DatabaseInterface {
	pp.DatabaseWasCalled = true
	return pp.database
}

func (pp *PersistenceProvider) GobbleDatabase() gobble.DatabaseInterface {
	pp.GobbleDatabaseWasCalled = true
	return pp.gobbleDatabase
}

type GobbleDatabase struct {
	MigrateWasCalled bool
}

func (gd *GobbleDatabase) Migrate() {
	gd.MigrateWasCalled = true
}
