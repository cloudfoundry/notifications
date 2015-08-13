package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesCollection", func() {
	var (
		templatesCollection collections.TemplatesCollection
		templatesRepository *fakes.TemplatesRepository
		conn                *fakes.Connection
	)

	BeforeEach(func() {
		templatesRepository = fakes.NewTemplatesRepository()

		templatesCollection = collections.NewTemplatesCollection(templatesRepository)
		conn = fakes.NewConnection()
	})

	Describe("Set", func() {
		BeforeEach(func() {
			templatesRepository.InsertCall.Returns.Template = models.Template{
				ID:       "some-template-id",
				Name:     "some-template",
				HTML:     "<h1>My Cool Template</h1>",
				Subject:  "{{.Subject}}",
				ClientID: "some-client-id",
			}
		})

		It("will insert a template into the collection", func() {
			template, err := templatesCollection.Set(conn, collections.Template{
				Name:     "some-template",
				HTML:     "<h1>My Cool Template</h1>",
				Subject:  "{{.Subject}}",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(collections.Template{
				ID:       "some-template-id",
				Name:     "some-template",
				HTML:     "<h1>My Cool Template</h1>",
				Subject:  "{{.Subject}}",
				ClientID: "some-client-id",
			}))
			Expect(templatesRepository.InsertCall.Receives.Template.Name).To(Equal("some-template"))
		})

		Context("failure cases", func() {
			It("returns a DuplicateRecordError if the repo returns it", func() {
				templatesRepository.InsertCall.Returns.Err = models.DuplicateRecordError{}

				_, err := templatesCollection.Set(conn, collections.Template{
					Name:     "some-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				})

				Expect(err).To(BeAssignableToTypeOf(collections.DuplicateRecordError{}))
			})

			It("returns a persistence error for anything else", func() {
				templatesRepository.InsertCall.Returns.Err = errors.New("failed to save")

				_, err := templatesCollection.Set(conn, collections.Template{
					Name:     "some-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				})

				Expect(err).To(Equal(collections.PersistenceError{
					Err: errors.New("failed to save"),
				}))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			templatesRepository.GetCall.Returns.Template = models.Template{
				ID:       "some-template-id",
				Name:     "some-template",
				HTML:     "<h1>My Cool Template</h1>",
				Subject:  "{{.Subject}}",
				ClientID: "some-client-id",
			}
		})

		It("will retrieve a template from the collection", func() {
			template, err := templatesCollection.Get(conn, "some-template-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(collections.Template{
				ID:       "some-template-id",
				Name:     "some-template",
				HTML:     "<h1>My Cool Template</h1>",
				Subject:  "{{.Subject}}",
				ClientID: "some-client-id",
			}))

			Expect(templatesRepository.GetCall.Receives.Conn).To(Equal(conn))
			Expect(templatesRepository.GetCall.Receives.TemplateID).To(Equal("some-template-id"))
		})

		Context("failure cases", func() {
			It("returns a not found error if the template does not exist", func() {
				templatesRepository.GetCall.Returns.Err = models.RecordNotFoundError("")
				_, err := templatesCollection.Get(conn, "missing-template-id", "some-client-id")
				Expect(err).To(BeAssignableToTypeOf(collections.NotFoundError{}))
			})

			It("returns a not found error if the template belongs to a different client ID", func() {
				templatesRepository.GetCall.Returns.Template = models.Template{
					ID:       "some-template-id",
					Name:     "some-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "other-client-id",
				}
				_, err := templatesCollection.Get(conn, "some-template-id", "some-client-id")
				Expect(err).To(BeAssignableToTypeOf(collections.NotFoundError{}))
			})

			It("returns a persistence error if one occurs", func() {
				templatesRepository.GetCall.Returns.Err = errors.New("failed to retrieve")
				_, err := templatesCollection.Get(conn, "some-template-id", "some-client-id")
				Expect(err).To(BeAssignableToTypeOf(collections.PersistenceError{}))
			})
		})
	})
})
