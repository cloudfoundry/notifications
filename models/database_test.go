package models_test

import (
	"path"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	var env application.Environment
	var db *models.DB
	var connection *models.Connection

	BeforeEach(func() {
		TruncateTables()
		env = application.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		models.ClearDB()
		db = models.NewDatabase(models.Config{
			DatabaseURL:         env.DatabaseURL,
			MigrationsPath:      migrationsPath,
			DefaultTemplatePath: env.RootPath + "/templates/default.json",
		})
		connection = db.Connection().(*models.Connection)
	})

	Describe("acting as a singleton", func() {
		It("returns a connection to the database", func() {
			err := connection.Db.Ping()
			Expect(err).To(BeNil())

			_, err = connection.Db.Query("SHOW TABLES")
			Expect(err).To(BeNil())
		})

		It("returns a single connection only", func() {
			migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
			db2 := models.NewDatabase(models.Config{
				DatabaseURL:    env.DatabaseURL,
				MigrationsPath: migrationsPath,
			})

			Expect(db).To(Equal(db2))
		})
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
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))

			db.Seed()
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
				db.Seed()
				db.Seed()
				db.Seed()
			}).NotTo(Panic())
		})

		Context("when it has not been overridden", func() {
			It("re-seeds the default template when the file is updated", func() {
				db.Seed()

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

				db.Seed()

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
				db.Seed()

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

				db.Seed()

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

	Describe("Connection", func() {
		It("returns a Connection", func() {
			connection := db.Connection()
			Expect(connection).To(BeAssignableToTypeOf(&models.Connection{}))
		})
	})
})
