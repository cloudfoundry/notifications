package application_test

import (
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrator", func() {
	Describe("Migrate", func() {
		var (
			migrator       application.Migrator
			provider       *mocks.PersistenceProvider
			database       *mocks.Database
			gobbleDatabase *mocks.GobbleDatabase
			dbMigrator     *mocks.DatabaseMigrator
		)

		BeforeEach(func() {
			database = mocks.NewDatabase()
			gobbleDatabase = &mocks.GobbleDatabase{}
			provider = mocks.NewPersistenceProvider()
			provider.DatabaseCall.Returns.Database = database
			provider.GobbleDatabaseCall.Returns.Database = gobbleDatabase

			dbMigrator = mocks.NewDatabaseMigrator()
		})

		Context("when configured to run migrations", func() {
			BeforeEach(func() {
				migrator = application.NewMigrator(provider, dbMigrator, true, "/my-migrations/dir", "/my-gobble/dir", "/my-templates/dir")
				migrator.Migrate()
			})

			It("migrates the gobble database", func() {
				Expect(gobbleDatabase.MigrateCall.Receives.MigrationsDir).To(Equal("/my-gobble/dir"))
			})

			It("migrates the notifications database", func() {
				Expect(dbMigrator.MigrateCall.Called).To(BeTrue())
				Expect(dbMigrator.MigrateCall.Receives.DB).To(Equal(database.RawConnection()))
				Expect(dbMigrator.MigrateCall.Receives.MigrationsPath).To(Equal("/my-migrations/dir"))
			})

			It("seeds the database", func() {
				Expect(dbMigrator.SeedCall.Called).To(BeTrue())
				Expect(dbMigrator.SeedCall.Receives.Database).To(Equal(database))
				Expect(dbMigrator.SeedCall.Receives.DefaultTemplatePath).To(Equal("/my-templates/dir"))
			})
		})

		Context("when configured to skip migrations", func() {
			BeforeEach(func() {
				migrator = application.NewMigrator(provider, dbMigrator, false, "these-dont-matter", "these-dont-matter", "these-dont-matter")
				migrator.Migrate()
			})

			It("does not migrate the gobble database", func() {
				Expect(gobbleDatabase.MigrateCall.WasCalled).To(BeFalse())
			})

			It("does not migrate the notifications database", func() {
				Expect(dbMigrator.MigrateCall.Called).To(BeFalse())
			})

			It("does not seed the database", func() {
				Expect(dbMigrator.SeedCall.Called).To(BeFalse())
			})
		})
	})
})
