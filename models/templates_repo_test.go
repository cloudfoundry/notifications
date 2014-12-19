package models_test

import (
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
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
		env := application.NewEnvironment()
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

	Context("#FindByID", func() {
		Context("the template is in the database", func() {
			It("returns the template when it is found", func() {
				raptorTemplate, err := repo.FindByID(conn, "raptor_template")

				Expect(err).ToNot(HaveOccurred())
				Expect(raptorTemplate.ID).To(Equal("raptor_template"))
				Expect(raptorTemplate.Name).To(Equal("Raptors On The Run"))
				Expect(raptorTemplate.Text).To(Equal("run and hide"))
				Expect(raptorTemplate.HTML).To(Equal("<h1>containment unit breached!</h1>"))
			})
		})

		Context("the template is not in the database", func() {
			It("returns a record not found error", func() {
				sillyTemplate, err := repo.FindByID(conn, "silly_template")

				Expect(sillyTemplate).To(Equal(models.Template{}))
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})
		})
	})

	Context("#Find", func() {
		Context("the template is in the database", func() {
			It("returns the template when it is found", func() {
				raptorTemplate, err := repo.Find(conn, "Raptors On The Run")

				Expect(err).ToNot(HaveOccurred())
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
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
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

			foundTemplate, err := repo.Find(conn, createdTemplate.Name)
			if err != nil {
				panic(err)
			}

			Expect(foundTemplate.ID).NotTo(BeNil())
			Expect(foundTemplate.Name).To(Equal(newTemplate.Name))
			Expect(foundTemplate.Subject).To(Equal(newTemplate.Subject))
			Expect(foundTemplate.Text).To(Equal(newTemplate.Text))
			Expect(foundTemplate.HTML).To(Equal(newTemplate.HTML))
			Expect(foundTemplate.CreatedAt).To(Equal(createdAt))
			Expect(foundTemplate.UpdatedAt).To(Equal(createdAt))
		})
	})

	Describe("Update", func() {
		var aNewTemplate models.Template

		BeforeEach(func() {
			aNewTemplate = models.Template{
				Name:    "a brand new name",
				Subject: "Some new subject",
				Text:    "some newer text",
				HTML:    "<p>new HTML</p>",
			}
		})

		Context("the template exists in the database", func() {
			It("updates a template currently in the database", func() {
				_, err := repo.Update(conn, template.ID, aNewTemplate)
				Expect(err).ToNot(HaveOccurred())

				foundTemplate, err := repo.FindByID(conn, template.ID)
				if err != nil {
					panic(err)
				}

				Expect(foundTemplate.Name).To(Equal(aNewTemplate.Name))
				Expect(foundTemplate.Subject).To(Equal(aNewTemplate.Subject))
				Expect(foundTemplate.Text).To(Equal(aNewTemplate.Text))
				Expect(foundTemplate.HTML).To(Equal(aNewTemplate.HTML))
				Expect(foundTemplate.CreatedAt).To(Equal(createdAt))
				Expect(foundTemplate.UpdatedAt).ToNot(Equal(createdAt))
				Expect(foundTemplate.UpdatedAt).To(BeTemporally(">", createdAt))
			})
		})

		Context("the template does not exist in the database", func() {
			It("bubbles up the error", func() {
				_, err := repo.Update(conn, "a-bad-id", aNewTemplate)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.TemplateFindError{Message: "Template a-bad-id not found"}))
			})
		})
	})

	Describe("#Upsert", func() {
		It("inserts a template into the database", func() {
			var err error

			newTemplate := models.Template{
				Name:      "silly_template." + models.UserBodyTemplateName,
				Subject:   "silliness",
				Text:      "omg",
				HTML:      "<h1>OMG</h1>",
				CreatedAt: createdAt,
			}

			_, err = repo.Upsert(conn, newTemplate)
			Expect(err).ToNot(HaveOccurred())

			foundTemplate, err := repo.Find(conn, newTemplate.Name)
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

			foundTemplate, err := repo.Find(conn, template.Name)
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

	Describe("#ListIDsAndNames", func() {
		Context("there are templates in the database", func() {
			It("returns a list of templates - ID and Name only", func() {
				secondTemplate := models.Template{
					ID:        "star_template",
					Name:      "Shooting Stars",
					Text:      "pretty",
					HTML:      "<h1>Awe</h1>",
					CreatedAt: createdAt,
				}

				conn.Insert(&secondTemplate)

				expectedMetadata := []models.Template{
					models.Template{
						ID:   "raptor_template",
						Name: "Raptors On The Run",
					},
					models.Template{
						ID:   "star_template",
						Name: "Shooting Stars",
					},
				}
				templatesMetadata, err := repo.ListIDsAndNames(conn)

				Expect(err).ToNot(HaveOccurred())
				Expect(templatesMetadata).To(Equal(expectedMetadata))
				Expect(templatesMetadata[0].Text).To(BeEmpty())
				Expect(templatesMetadata[0].HTML).To(BeEmpty())
				Expect(templatesMetadata[0].Subject).To(BeEmpty())
			})
		})
	})

	Describe("#Destroy", func() {
		Context("the template exists in the database", func() {
			It("deletes the template by templateID", func() {
				_, err := repo.FindByID(conn, template.ID)
				if err != nil {
					panic(err)
				}

				err = repo.Destroy(conn, template.ID)
				Expect(err).ToNot(HaveOccurred())

				_, err = repo.FindByID(conn, template.ID)
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})
		})

		Context("the template does not exist in the database", func() {
			It("returns an RecordNotFoundError", func() {
				err := repo.Destroy(conn, "knockknock")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))

				_, err = repo.FindByID(conn, "knockknock")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})
		})
	})

	Describe("#DeprecatedDestroy", func() {
		Context("the template exists in the database", func() {
			It("deletes the template by name", func() {
				_, err := repo.Find(conn, template.Name)
				if err != nil {
					panic(err)
				}

				err = repo.DeprecatedDestroy(conn, template.Name)
				Expect(err).ToNot(HaveOccurred())

				_, err = repo.Find(conn, template.Name)
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))

			})
		})

		Context("the template does not exist in the database", func() {
			It("does not return an error", func() {
				err := repo.DeprecatedDestroy(conn, "knockknock")
				Expect(err).ToNot(HaveOccurred())

				_, err = repo.Find(conn, "knockknock")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))

			})
		})
	})
})
