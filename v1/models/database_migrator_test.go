package models_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DatabaseMigrator", func() {
	var (
		database            *db.DB
		connection          *db.Connection
		dbMigrator          models.DatabaseMigrator
		defaultTemplatePath string
	)

	BeforeEach(func() {
		env, err := application.NewEnvironment()
		Expect(err).NotTo(HaveOccurred())

		defaultTemplatePath = env.RootPath + "/templates/default.json"
		database = db.NewDatabase(sqlDB, db.Config{
			DefaultTemplatePath: defaultTemplatePath,
		})
		helpers.TruncateTables(database)
		connection = database.Connection().(*db.Connection)
		dbMigrator = models.DatabaseMigrator{}
	})

	Describe("migrating the database", func() {
		It("has the correct tables", func() {
			err := connection.Db.Ping()
			Expect(err).To(BeNil())

			rows, err := connection.Db.Query("SHOW TABLES")
			Expect(err).To(BeNil())

			tables := []string{}
			for rows.Next() {
				var table string
				err = rows.Scan(&table)
				if err != nil {
					panic(err)
				}
				tables = append(tables, table)
			}
			err = rows.Err()
			if err != nil {
				panic(err)
			}

			rows.Close()

			Expect(tables).To(ContainElement("clients"))
			Expect(tables).To(ContainElement("kinds"))
			Expect(tables).To(ContainElement("receipts"))
			Expect(tables).To(ContainElement("unsubscribes"))
			Expect(tables).To(ContainElement("global_unsubscribes"))
			Expect(tables).To(ContainElement("templates"))
		})
	})

	Describe("seeding the default template", func() {
		var repo models.TemplatesRepo

		BeforeEach(func() {
			repo = models.NewTemplatesRepo()
		})

		It("has the default template pre-seeded", func() {
			_, err := repo.FindByID(connection, models.DefaultTemplateID)
			Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Template with ID \"default\" could not be found")}))

			dbMigrator.Seed(database, defaultTemplatePath)
			template, err := repo.FindByID(connection, models.DefaultTemplateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(template.Name).To(Equal("Default Template"))
			Expect(template.Subject).To(Equal("CF Notification: {{.Subject}}"))
			Expect(template.HTML).To(Equal("<p>{{.Endorsement}}</p>{{.HTML}}"))
			Expect(template.Text).To(Equal("{{.Endorsement}}\n{{.Text}}"))
			Expect(template.Metadata).To(Equal("{}"))
		})

		It("can be called multiple times without panicking", func() {
			Expect(func() {
				dbMigrator.Seed(database, defaultTemplatePath)
				dbMigrator.Seed(database, defaultTemplatePath)
				dbMigrator.Seed(database, defaultTemplatePath)
			}).NotTo(Panic())
		})

		Context("when it has not been overridden", func() {
			It("re-seeds the default template when the file is updated", func() {
				dbMigrator.Seed(database, defaultTemplatePath)

				template, err := repo.FindByID(connection, models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())

				template.Name = "Updated Default"
				template.Subject = "Updated Subject"
				template.Text = "Updated Text"
				template.HTML = "Updated HTML"
				template.Metadata = `{"test": true}`
				template.Overridden = false
				_, err = connection.Update(&template)
				Expect(err).NotTo(HaveOccurred())

				dbMigrator.Seed(database, defaultTemplatePath)

				template, err = repo.FindByID(connection, models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())
				Expect(template.Name).To(Equal("Default Template"))
				Expect(template.Subject).To(Equal("CF Notification: {{.Subject}}"))
				Expect(template.HTML).To(Equal("<p>{{.Endorsement}}</p>{{.HTML}}"))
				Expect(template.Text).To(Equal("{{.Endorsement}}\n{{.Text}}"))
				Expect(template.Metadata).To(Equal("{}"))
				Expect(template.Overridden).To(BeFalse())
			})
		})

		Context("when it has been overridden", func() {
			It("does not re-seed the default template", func() {
				dbMigrator.Seed(database, defaultTemplatePath)

				template, err := repo.FindByID(connection, models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())

				template.Name = "Updated Default"
				template.Subject = "Updated Subject"
				template.Text = "Updated Text"
				template.HTML = "Updated HTML"
				template.Metadata = `{"test": true}`
				template.Overridden = true
				_, err = repo.Update(connection, models.DefaultTemplateID, template)
				Expect(err).NotTo(HaveOccurred())

				dbMigrator.Seed(database, defaultTemplatePath)

				template, err = repo.FindByID(connection, models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())
				Expect(template.Name).To(Equal("Updated Default"))
				Expect(template.Subject).To(Equal("Updated Subject"))
				Expect(template.HTML).To(Equal("Updated HTML"))
				Expect(template.Text).To(Equal("Updated Text"))
				Expect(template.Metadata).To(Equal(`{"test": true}`))
				Expect(template.Overridden).To(BeTrue())
			})
		})
	})
})
