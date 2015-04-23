package application

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type PersistenceProvider interface {
	Database() models.DatabaseInterface
	GobbleDatabase() gobble.DatabaseInterface
}

type Migrator struct {
	provider      PersistenceProvider
	shouldMigrate bool
}

func NewMigrator(provider PersistenceProvider, shouldMigrate bool) Migrator {
	return Migrator{
		provider:      provider,
		shouldMigrate: shouldMigrate,
	}
}

func (m Migrator) Migrate() {
	if m.shouldMigrate {
		m.provider.Database().Seed()
		m.provider.GobbleDatabase().Migrate()
	}
}
