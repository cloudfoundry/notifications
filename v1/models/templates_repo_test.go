package models_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesRepo", func() {
	var repo models.TemplatesRepo
	var conn db.ConnectionInterface
	var template models.Template
	var createdAt time.Time

	BeforeEach(func() {
		repo = models.NewTemplatesRepo()
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)
		conn = database.Connection()
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
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Template with ID \"silly_template\" could not be found")}))
			})
		})
	})

	Describe("#Create", func() {
		It("inserts a template into the database", func() {
			newTemplate := models.Template{
				Name:    "A Nice Template",
				Subject: "Kind Words",
				Text:    "Some kind of compliment.",
				HTML:    "<h1>Genuine Smile</h1>",
			}

			createdTemplate, err := repo.Create(conn, newTemplate)
			Expect(err).ToNot(HaveOccurred())

			foundTemplate, err := repo.FindByID(conn, createdTemplate.ID)
			if err != nil {
				panic(err)
			}

			regularExpression := `^[[:xdigit:]]{8}\-[[:xdigit:]]{4}\-[[:xdigit:]]{4}\-[[:xdigit:]]{4}\-[[:xdigit:]]{12}$`
			Expect(foundTemplate.ID).To(MatchRegexp(regularExpression))
			Expect(foundTemplate.ID).To(Equal(createdTemplate.ID))
			Expect(foundTemplate.Name).To(Equal(newTemplate.Name))
			Expect(foundTemplate.Subject).To(Equal(newTemplate.Subject))
			Expect(foundTemplate.Text).To(Equal(newTemplate.Text))
			Expect(foundTemplate.HTML).To(Equal(newTemplate.HTML))
			Expect(foundTemplate.CreatedAt).To(BeTemporally("~", time.Now().UTC(), 2*time.Second))
			Expect(foundTemplate.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), 2*time.Second))
		})
	})

	Describe("Update", func() {
		var aNewTemplate models.Template

		BeforeEach(func() {
			aNewTemplate = models.Template{
				Name:     "a brand new name",
				Subject:  "Some new subject",
				Text:     "some newer text",
				HTML:     "<p>new HTML</p>",
				Metadata: "{\"cloudy\": true}",
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
				Expect(foundTemplate.Metadata).To(Equal(aNewTemplate.Metadata))
				Expect(foundTemplate.CreatedAt).To(Equal(createdAt))
				Expect(foundTemplate.UpdatedAt).ToNot(Equal(createdAt))
				Expect(foundTemplate.UpdatedAt).To(BeTemporally(">", createdAt))
				Expect(foundTemplate.Overridden).To(BeTrue())
			})
		})

		Context("the template does not exist in the database", func() {
			It("bubbles up the error", func() {
				_, err := repo.Update(conn, "a-bad-id", aNewTemplate)
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Template with ID \"a-bad-id\" could not be found")}))
			})
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
					{
						ID:   "raptor_template",
						Name: "Raptors On The Run",
					},
					{
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
				Expect(err).To(MatchError(models.NotFoundError{Err: fmt.Errorf("Template with ID %q could not be found", template.ID)}))
			})
		})

		Context("the template does not exist in the database", func() {
			It("returns an RecordNotFoundError", func() {
				err := repo.Destroy(conn, "knockknock")
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Template with ID \"knockknock\" could not be found")}))
			})
		})
	})
})
