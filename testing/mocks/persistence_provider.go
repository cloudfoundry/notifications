package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
)

type PersistenceProvider struct {
	DatabaseCall struct {
		Returns struct {
			Database db.DatabaseInterface
		}
	}

	GobbleDatabaseCall struct {
		Returns struct {
			Database gobble.DatabaseInterface
		}
	}
}

func NewPersistenceProvider() *PersistenceProvider {
	return &PersistenceProvider{}
}

func (pp *PersistenceProvider) Database() db.DatabaseInterface {
	return pp.DatabaseCall.Returns.Database
}

func (pp *PersistenceProvider) GobbleDatabase() gobble.DatabaseInterface {
	return pp.GobbleDatabaseCall.Returns.Database
}

type GobbleDatabase struct {
	MigrateCall struct {
		WasCalled bool
		Receives  struct {
			MigrationsDir string
		}
	}
}

func (gd *GobbleDatabase) Migrate(migrationsDir string) {
	gd.MigrateCall.Receives.MigrationsDir = migrationsDir
	gd.MigrateCall.WasCalled = true
}
