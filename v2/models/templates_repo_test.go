package models_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesRepo", func() {
	var (
		repo          models.TemplatesRepository
		conn          db.ConnectionInterface
		guidGenerator *mocks.IDGenerator
	)

	BeforeEach(func() {
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)

		guidGenerator = mocks.NewIDGenerator()
		guidGenerator.GenerateCall.Returns.IDs = []string{"first-random-guid", "second-random-guid"}

		repo = models.NewTemplatesRepository(guidGenerator.Generate)
		conn = database.Connection()
	})

	Describe("DefaultTemplate", func() {
		It("defines a default template", func() {
			Expect(models.DefaultTemplate).To(Equal(models.Template{
				ID:       "default",
				Name:     "The Default Template",
				Subject:  "{{.Subject}}",
				Text:     "{{.Text}}",
				HTML:     "{{.HTML}}",
				Metadata: "{}",
			}))
		})
	})

	Describe("Insert", func() {
		It("returns the data", func() {
			createdTemplate, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createdTemplate.ID).To(Equal("first-random-guid"))
		})

		Context("when the 'default' template ID is supplied", func() {
			It("does not generate a random guid", func() {
				createdTemplate, err := repo.Insert(conn, models.Template{
					ID:       "default",
					Name:     "some-template",
					ClientID: "some-client-id",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(guidGenerator.GenerateCall.CallCount).To(Equal(0))
				Expect(createdTemplate.ID).To(Equal("default"))
			})
		})

		Context("failure cases", func() {
			It("returns an error if the database returns one", func() {
				connection := mocks.NewConnection()
				connection.InsertCall.Returns.Error = errors.New("some error")

				_, err := repo.Insert(connection, models.Template{
					ID:       "default",
					Name:     "some-template",
					ClientID: "some-client-id",
				})
				Expect(err).To(MatchError(errors.New("some error")))
			})

			It("returns an error if the guid generator returns one", func() {
				guidGenerator.GenerateCall.Returns.Error = errors.New("some error")

				_, err := repo.Insert(conn, models.Template{})
				Expect(err).To(MatchError(errors.New("some error")))
			})
		})
	})

	Describe("List", func() {
		It("returns the templates", func() {
			_, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-other-client-id",
			})

			_, err = repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			templates, err := repo.List(conn, "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(HaveLen(1))
			Expect(templates[0].ClientID).To(Equal("some-client-id"))
			Expect(templates[0].Name).To(Equal("some-template"))
		})

		Context("failure cases", func() {
			It("returns not found error if it happens", func() {
				connection := mocks.NewConnection()
				connection.SelectCall.Returns.Error = errors.New("an error")
				_, err := repo.List(connection, "client-id")
				Expect(err).To(MatchError(errors.New("an error")))
			})
		})
	})

	Describe("Get", func() {
		It("fetches the template given a template_id", func() {
			createdTemplate, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			template, err := repo.Get(conn, createdTemplate.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(createdTemplate))
		})

		It("fetches the default template if it does not exist", func() {
			template, err := repo.Get(conn, "default")
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(models.DefaultTemplate))
		})

		Context("failure cases", func() {
			It("returns not found error if it happens", func() {
				_, err := repo.Get(conn, "missing-template-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
			})
		})
	})

	Describe("Delete", func() {
		It("deletes the template given a template_id", func() {
			template, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			err = repo.Delete(conn, template.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Get(conn, template.ID)
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
		})

		Context("failure cases", func() {
			It("returns not found error if it happens", func() {
				err := repo.Delete(conn, "missing-template-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
			})
		})
	})

	Describe("Update", func() {
		It("updates the template", func() {
			createdTemplate, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			createdTemplate.Name = "new-template"

			updatedTemplate, err := repo.Update(conn, createdTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedTemplate.ID).To(Equal(createdTemplate.ID))
			Expect(updatedTemplate.Name).To(Equal("new-template"))

			template, err := repo.Get(conn, updatedTemplate.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(updatedTemplate))
		})

		Context("failure cases", func() {
			It("returns an error when the update call fails", func() {
				fakeConn := mocks.NewConnection()
				fakeConn.UpdateCall.Returns.Error = errors.New("some db error")
				conn = fakeConn

				_, err := repo.Update(conn, models.Template{})
				Expect(err).To(MatchError(errors.New("some db error")))
			})
		})
	})
})
