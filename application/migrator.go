package application

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type persistenceProvider interface {
	Database() db.DatabaseInterface
	GobbleDatabase() gobble.DatabaseInterface
}

type dbMigrator interface {
	Migrate(db *sql.DB, migrationsPath string)
	Seed(db models.DatabaseInterface, defaultTemplatePath string)
}

type Migrator struct {
	provider             persistenceProvider
	dbMigrator           dbMigrator
	shouldMigrate        bool
	gobbleMigrationsPath string
	migrationsPath       string
	defaultTemplatePath  string
}

func NewMigrator(provider persistenceProvider, dbMigrator dbMigrator, shouldMigrate bool, migrationsPath, gobbleMigrationsPath, defaultTemplatePath string) Migrator {
	return Migrator{
		provider:             provider,
		dbMigrator:           dbMigrator,
		shouldMigrate:        shouldMigrate,
		gobbleMigrationsPath: gobbleMigrationsPath,
		migrationsPath:       migrationsPath,
		defaultTemplatePath:  defaultTemplatePath,
	}
}

func (m Migrator) Migrate() {
	if m.shouldMigrate {
		m.dbMigrator.Migrate(m.provider.Database().RawConnection(), m.migrationsPath)
		m.dbMigrator.Seed(m.provider.Database(), m.defaultTemplatePath)
		m.provider.GobbleDatabase().Migrate(m.gobbleMigrationsPath)
	}
}
