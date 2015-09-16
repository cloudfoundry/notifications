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

	Describe("DefaultTemplate", func() {
		It("defines a default template", func() {
			Expect(collections.DefaultTemplate).To(Equal(collections.Template{
				ID:       "default",
				Name:     "The Default Template",
				Subject:  "{{.Subject}}",
				Text:     "{{.Text}}",
				HTML:     "{{.HTML}}",
				Metadata: "{}",
			}))
		})
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
				templatesRepository.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("not found")}
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
					templatesRepository.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("not found")}
					templatesRepository.InsertCall.Returns.Error = models.DuplicateRecordError{}

					_, err := templatesCollection.Set(conn, collections.Template{
						Name:     "some-template",
						HTML:     "<h1>My Cool Template</h1>",
						Subject:  "{{.Subject}}",
						ClientID: "some-client-id",
					})

					Expect(err).To(BeAssignableToTypeOf(collections.DuplicateRecordError{}))
				})
			})
		})

		Context("when an existing ID is supplied", func() {
			BeforeEach(func() {
				templatesRepository.GetCall.Returns.Template = models.Template{
					ID:       "existing-id",
					Name:     "old-template",
					HTML:     "<h1>My Cool Template</h1>",
					Subject:  "{{.Subject}}",
					ClientID: "some-client-id",
				}
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

				Expect(templatesRepository.GetCall.Receives.Connection).To(Equal(conn))
				Expect(templatesRepository.GetCall.Receives.TemplateID).To(Equal("existing-id"))

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
					templatesRepository.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("not found")}

					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "some-client-id",
					})
					Expect(err).NotTo(HaveOccurred())

					Expect(templatesRepository.GetCall.Receives.Connection).To(Equal(conn))
					Expect(templatesRepository.GetCall.Receives.TemplateID).To(Equal("default"))

					Expect(templatesRepository.InsertCall.Receives.Connection).To(Equal(conn))
					Expect(templatesRepository.InsertCall.Receives.Template).To(Equal(models.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "some-client-id",
					}))
				})

				It("will update the saved template if it already exists", func() {
					templatesRepository.GetCall.Returns.Template = models.Template{
						ID:       "default",
						Name:     "some-template",
						HTML:     "<h1>My Cool Template</h1>",
						Subject:  "{{.Subject}}",
						ClientID: "some-client-id",
					}

					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "some-client-id",
					})
					Expect(err).NotTo(HaveOccurred())

					Expect(templatesRepository.GetCall.Receives.Connection).To(Equal(conn))
					Expect(templatesRepository.GetCall.Receives.TemplateID).To(Equal("default"))

					Expect(templatesRepository.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(templatesRepository.UpdateCall.Receives.Template).To(Equal(models.Template{
						ID:       "default",
						Name:     "updated default",
						Text:     "new default text",
						HTML:     "new default html",
						Subject:  "New Default Subject",
						ClientID: "some-client-id",
					}))
				})
			})

			Context("failure cases", func() {
				It("returns a NotFoundError when the ID supplied does not exist", func() {
					repoError := models.RecordNotFoundError{errors.New("whatever")}
					templatesRepository.GetCall.Returns.Error = repoError

					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "not-existing-id",
						Name:     "new-template",
						HTML:     "<h1>My Cool Template</h1>",
						Subject:  "{{.Subject}}",
						ClientID: "some-client-id",
					})
					Expect(err).To(MatchError(collections.NotFoundError{repoError}))
				})

				It("returns a PersistenceError when the template repo returns an error from Get", func() {
					repoError := errors.New("whoops!")
					templatesRepository.GetCall.Returns.Error = repoError

					_, err := templatesCollection.Set(conn, collections.Template{
						ID:       "not-existing-id",
						Name:     "new-template",
						HTML:     "<h1>My Cool Template</h1>",
						Subject:  "{{.Subject}}",
						ClientID: "some-client-id",
					})
					Expect(err).To(MatchError(collections.PersistenceError{repoError}))
				})

				It("returns a PersistenceError when the template repo returns an error from Insert", func() {
					templatesRepository.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("not found")}
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
				templatesRepository.GetCall.Returns.Error = models.NewRecordNotFoundError("")

				template, err := templatesCollection.Get(conn, "default", "some-client-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(template).To(Equal(collections.DefaultTemplate))
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
