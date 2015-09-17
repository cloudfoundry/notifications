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
		Context("when no ID is supplied", func() {
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
		})

		Context("when an existing ID is supplied", func() {
			BeforeEach(func() {
				templatesRepository.UpdateCall.Returns.Template = models.Template{
					ID:       "existing-id",
					Name:     "new-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				}
			})

			It("will update a template if it already exists", func() {
				template, err := templatesCollection.Set(conn, collections.Template{
					ID:       "existing-id",
					Name:     "new-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(template).To(Equal(collections.Template{
					ID:       "existing-id",
					Name:     "new-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				}))

				Expect(templatesRepository.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(templatesRepository.UpdateCall.Receives.Template).To(Equal(models.Template{
					ID:       "existing-id",
					Name:     "new-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				}))
			})

			Context("when the default template ID is supplied", func() {
				It("will create a new record if it does not already exist", func() {
					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "",
					})
					Expect(err).NotTo(HaveOccurred())

					Expect(templatesRepository.InsertCall.Receives.Connection).To(Equal(conn))
					Expect(templatesRepository.InsertCall.Receives.Template).To(Equal(models.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "",
					}))
				})

				It("will update the saved template if it already exists", func() {
					templatesRepository.InsertCall.Returns.Error = models.DuplicateRecordError{errors.New("dup")}
					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "",
					})
					Expect(err).NotTo(HaveOccurred())

					Expect(templatesRepository.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(templatesRepository.UpdateCall.Receives.Template).To(Equal(models.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "",
					}))
				})
			})

			Context("failure cases", func() {
				It("returns a PersistenceError when the template repo returns an error from Insert", func() {
					repoError := errors.New("failed to save")
					templatesRepository.InsertCall.Returns.Error = repoError

					_, err := templatesCollection.Set(conn, collections.Template{
						Name:     "some-template",
						HTML:     "<h1>My Cool Template</h1>",
						Subject:  "{{.Subject}}",
						ClientID: "some-client-id",
					})
					Expect(err).To(MatchError(collections.PersistenceError{repoError}))
				})

				It("returns a PersistenceError when the template repo returns an error from Update", func() {
					repoError := errors.New("fail!")
					templatesRepository.UpdateCall.Returns.Error = repoError

					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "not-existing-id",
						Name:     "new-template",
						HTML:     "<h1>My Cool Template</h1>",
						Subject:  "{{.Subject}}",
						ClientID: "some-client-id",
					})
					Expect(err).To(MatchError(collections.PersistenceError{repoError}))
				})
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

		Context("when the default template does not exist in the repo", func() {
			It("returns the 'stock' default template", func() {
				templatesRepository.GetCall.Returns.Template = models.DefaultTemplate

				template, err := templatesCollection.Get(conn, "default", "some-client-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(template).To(Equal(collections.Template{
					ID:       models.DefaultTemplate.ID,
					Name:     models.DefaultTemplate.Name,
					Text:     models.DefaultTemplate.Text,
					HTML:     models.DefaultTemplate.HTML,
					Subject:  models.DefaultTemplate.Subject,
					Metadata: models.DefaultTemplate.Metadata,
				}))
			})
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

	Describe("List", func() {
		BeforeEach(func() {
			templatesRepository.ListCall.Returns.Templates = []models.Template{
				{
					ID:       "my-template-id",
					ClientID: "some-client-id",
				},
			}
		})
		It("returns templates for a client id", func() {
			templates, err := templatesCollection.List(conn, "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(HaveLen(1))
			Expect(templates[0].ID).To(Equal("my-template-id"))

			Expect(templatesRepository.ListCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepository.ListCall.Receives.ClientID).To(Equal("some-client-id"))
		})

		Context("failure cases", func() {
			It("returns an unknown error if the repo returns an error", func() {
				templatesRepository.ListCall.Returns.Error = errors.New("failed to list")

				_, err := templatesCollection.List(conn, "some-client-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("failed to list")}))
			})
		})
	})
})
