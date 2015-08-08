package application_test

import (
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrator", func() {
	Describe("Migrate", func() {
		var (
			migrator       application.Migrator
			provider       *fakes.PersistenceProvider
			database       *fakes.Database
			gobbleDatabase *fakes.GobbleDatabase
		)

		BeforeEach(func() {
			database = fakes.NewDatabase()
			gobbleDatabase = &fakes.GobbleDatabase{}
			provider = fakes.NewPersistenceProvider(database, gobbleDatabase)
		})

		Context("when configured to run migrations", func() {
			BeforeEach(func() {
				migrator = application.NewMigrator(provider, true, "/my-migrations/dir", "/my-gobble/dir")
				migrator.Migrate()
			})

			It("migrates the gobble database", func() {
				Expect(gobbleDatabase.MigrateWasCalled).To(BeTrue())
				Expect(gobbleDatabase.MigrationsDir).To(Equal("/my-gobble/dir"))
			})

			It("migrates the notifications database", func() {
				Expect(database.MigrateWasCalled).To(BeTrue())
				Expect(database.MigrationsPath).To(Equal("/my-migrations/dir"))
			})

			It("seeds the database", func() {
				Expect(database.SeedWasCalled).To(BeTrue())
			})
		})

		Context("when configured to skip migrations", func() {
			BeforeEach(func() {
				migrator = application.NewMigrator(provider, false, "these-dont-matter", "these-dont-matter")
				migrator.Migrate()
			})

			It("does not migrate the gobble database", func() {
				Expect(gobbleDatabase.MigrateWasCalled).To(BeFalse())
			})

			It("does not migrate the notifications database", func() {
				Expect(database.MigrateWasCalled).To(BeFalse())
			})

			It("does not seed the database", func() {
				Expect(database.SeedWasCalled).To(BeFalse())
			})
		})
	})
})
