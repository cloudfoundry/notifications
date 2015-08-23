package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesCollection", func() {
	var (
		templatesCollection collections.TemplatesCollection
		templatesRepository *mocks.TemplatesRepository
		conn                *mocks.Connection
	)

	BeforeEach(func() {
		templatesRepository = mocks.NewTemplatesRepository()

		templatesCollection = collections.NewTemplatesCollection(templatesRepository)
		conn = mocks.NewConnection()
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

			Expect(templatesRepository.InsertCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepository.InsertCall.Receives.Template).To(Equal(models.Template{
				Name:     "some-template",
				HTML:     "<h1>My Cool Template</h1>",
				Subject:  "{{.Subject}}",
				ClientID: "some-client-id",
			}))
		})

		Context("failure cases", func() {
			It("returns a DuplicateRecordError if the repo returns it", func() {
				templatesRepository.InsertCall.Returns.Error = models.DuplicateRecordError{}

				_, err := templatesCollection.Set(conn, collections.Template{
					Name:     "some-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				})

				Expect(err).To(BeAssignableToTypeOf(collections.DuplicateRecordError{}))
			})

			It("returns a persistence error for anything else", func() {
				templatesRepository.InsertCall.Returns.Error = errors.New("failed to save")

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

			Expect(templatesRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepository.GetCall.Receives.TemplateID).To(Equal("some-template-id"))
		})

		Context("failure cases", func() {
			It("returns a not found error if the template does not exist", func() {
				templatesRepository.GetCall.Returns.Error = models.NewRecordNotFoundError("")
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
				templatesRepository.GetCall.Returns.Error = errors.New("failed to retrieve")
				_, err := templatesCollection.Get(conn, "some-template-id", "some-client-id")
				Expect(err).To(BeAssignableToTypeOf(collections.PersistenceError{}))
			})
		})
	})

	Describe("Delete", func() {
		It("deletes a template from the collection", func() {
			err := templatesCollection.Delete(conn, "some-template-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(templatesRepository.DeleteCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepository.DeleteCall.Receives.TemplateID).To(Equal("some-template-id"))
		})

		Context("failure cases", func() {
			It("returns a not found error if the template does not exist", func() {
				templatesRepository.DeleteCall.Returns.Error = models.NewRecordNotFoundError("")
				err := templatesCollection.Delete(conn, "missing-template-id")
				Expect(err).To(BeAssignableToTypeOf(collections.NotFoundError{}))
			})

			It("returns a persistence error if one occurs", func() {
				templatesRepository.DeleteCall.Returns.Error = errors.New("failed to delete")
				err := templatesCollection.Delete(conn, "some-template-id")
				Expect(err).To(MatchError(collections.PersistenceError{errors.New("failed to delete")}))
			})
		})
	})
})
