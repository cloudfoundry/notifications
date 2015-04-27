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
	provider             PersistenceProvider
	shouldMigrate        bool
	gobbleMigrationsPath string
	migrationsPath       string
}

func NewMigrator(provider PersistenceProvider, shouldMigrate bool, migrationsPath, gobbleMigrationsPath string) Migrator {
	return Migrator{
		provider:             provider,
		shouldMigrate:        shouldMigrate,
		gobbleMigrationsPath: gobbleMigrationsPath,
		migrationsPath:       migrationsPath,
	}
}

func (m Migrator) Migrate() {
	if m.shouldMigrate {
		m.provider.Database().Migrate(m.migrationsPath)
		m.provider.Database().Seed()
		m.provider.GobbleDatabase().Migrate(m.gobbleMigrationsPath)
	}
}
