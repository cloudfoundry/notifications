package models_test

import (
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesRepo", func() {
	var repo models.TemplatesRepo
	var conn models.ConnectionInterface
	var template models.Template
	var createdAt time.Time

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewTemplatesRepo()
		env := config.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		db := models.NewDatabase(env.DatabaseURL, migrationsPath)
		conn = db.Connection()
		createdAt = time.Now().Add(-1 * time.Hour).Truncate(1 * time.Second).UTC()

		template = models.Template{
			ID:        "raptor_template",
			Name:      "Raptors On The Run",
			Text:      "run and hide",
			HTML:      "<h1>containment unit breached!</h1>",
			CreatedAt: createdAt,
		}

		conn.Insert(&template)
	})

	Context("#Find", func() {
		Context("the template is in the database", func() {
			It("returns the template when it is found", func() {
				raptorTemplate, err := repo.Find(conn, "raptor_template")

				Expect(err).ToNot(HaveOccurred())
				Expect(raptorTemplate.ID).To(Equal("raptor_template"))
				Expect(raptorTemplate.Name).To(Equal("Raptors On The Run"))
				Expect(raptorTemplate.Text).To(Equal("run and hide"))
				Expect(raptorTemplate.HTML).To(Equal("<h1>containment unit breached!</h1>"))
			})
		})

		Context("the template is not in the database", func() {
			It("returns a record not found error", func() {
				sillyTemplate, err := repo.Find(conn, "silly_template")

				Expect(sillyTemplate).To(Equal(models.Template{}))
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.ErrRecordNotFound{}))
			})
		})
	})

	Describe("#Create", func() {
		It("inserts a template into the database", func() {
			newTemplate := models.Template{
				Name:      "A Nice Template",
				Subject:   "Kind Words",
				Text:      "Some kind of compliment.",
				HTML:      "<h1>Genuine Smile</h1>",
				CreatedAt: createdAt,
			}

			createdTemplate, err := repo.Create(conn, newTemplate)
			Expect(err).ToNot(HaveOccurred())

			foundTemplate, err := repo.Find(conn, createdTemplate.ID)
			if err != nil {
				panic(err)
			}

			Expect(foundTemplate.Name).To(Equal(newTemplate.Name))
			Expect(foundTemplate.Subject).To(Equal(newTemplate.Subject))
			Expect(foundTemplate.Text).To(Equal(newTemplate.Text))
			Expect(foundTemplate.HTML).To(Equal(newTemplate.HTML))
			Expect(foundTemplate.CreatedAt).To(Equal(createdAt))
			Expect(foundTemplate.UpdatedAt).To(Equal(createdAt))
		})

		XIt("generates a guid for the template ID", func() {
		})
	})

	Describe("#Upsert", func() {
		It("inserts a template into the database", func() {
			newTemplate := models.Template{
				Name:      "silly_template." + models.UserBodyTemplateName,
				Subject:   "silliness",
				Text:      "omg",
				HTML:      "<h1>OMG</h1>",
				CreatedAt: createdAt,
			}

			upsertedTemplate, err := repo.Upsert(conn, newTemplate)
			Expect(err).ToNot(HaveOccurred())

			foundTemplate, err := repo.Find(conn, upsertedTemplate.ID)
			if err != nil {
				panic(err)
			}

			Expect(foundTemplate.Name).To(Equal(newTemplate.Name))
			Expect(foundTemplate.Subject).To(Equal(newTemplate.Subject))
			Expect(foundTemplate.Text).To(Equal(newTemplate.Text))
			Expect(foundTemplate.HTML).To(Equal(newTemplate.HTML))
			Expect(foundTemplate.CreatedAt).To(Equal(createdAt))
			Expect(foundTemplate.UpdatedAt).To(Equal(createdAt))
		})

		It("updates a template currently in the database", func() {
			template.Text = "new text"
			_, err := repo.Upsert(conn, template)
			Expect(err).ToNot(HaveOccurred())

			foundTemplate, err := repo.Find(conn, template.ID)
			if err != nil {
				panic(err)
			}

			Expect(foundTemplate.Name).To(Equal(template.Name))
			Expect(foundTemplate.Subject).To(Equal(template.Subject))
			Expect(foundTemplate.Text).To(Equal(template.Text))
			Expect(foundTemplate.HTML).To(Equal(template.HTML))
			Expect(foundTemplate.CreatedAt).To(Equal(createdAt))
			Expect(foundTemplate.UpdatedAt).ToNot(Equal(createdAt))
			Expect(foundTemplate.UpdatedAt).To(BeTemporally(">", createdAt))
		})
	})

	Describe("#Destroy", func() {
		Context("the template exists in the database", func() {
			It("deletes the template", func() {
				_, err := repo.Find(conn, template.ID)
				if err != nil {
					panic(err)
				}

				err = repo.Destroy(conn, template.ID)
				Expect(err).ToNot(HaveOccurred())

				_, err = repo.Find(conn, template.ID)
				Expect(err).To(Equal(models.ErrRecordNotFound{}))

			})
		})

		Context("the template does not exist in the database", func() {
			It("does not return an error", func() {
				err := repo.Destroy(conn, "knockknock")
				Expect(err).ToNot(HaveOccurred())

				_, err = repo.Find(conn, "knockknock")
				Expect(err).To(Equal(models.ErrRecordNotFound{}))

			})
		})
	})
})
