package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesCollection", func() {
	var (
		kindsRepo     *mocks.KindsRepo
		clientsRepo   *mocks.ClientsRepository
		templatesRepo *mocks.TemplatesRepo
		conn          *mocks.Connection

		collection collections.TemplatesCollection
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()

		clientsRepo = mocks.NewClientsRepository()
		kindsRepo = mocks.NewKindsRepo()
		templatesRepo = mocks.NewTemplatesRepo()

		collection = collections.NewTemplatesCollection(clientsRepo, kindsRepo, templatesRepo)
	})

	Describe("AssignToClient", func() {
		BeforeEach(func() {
			clientsRepo.FindCall.Returns.Client = models.Client{
				ID: "my-client",
			}
		})

		It("assigns the template to the given client", func() {
			err := collection.AssignToClient(conn, "my-client", "my-template")
			Expect(err).NotTo(HaveOccurred())

			Expect(clientsRepo.FindCall.Receives.Connection).To(Equal(conn))
			Expect(clientsRepo.FindCall.Receives.ClientID).To(Equal("my-client"))

			Expect(clientsRepo.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(clientsRepo.UpdateCall.Receives.Client).To(Equal(models.Client{
				ID:         "my-client",
				TemplateID: "my-template",
			}))
		})

		Context("when the request includes a non-existant id", func() {
			It("reports that the client cannot be found", func() {
				clientsRepo.FindCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				err := collection.AssignToClient(conn, "missing-client", "my-template")
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("not found")}))
			})

			It("reports that the template cannot be found", func() {
				templatesRepo.FindByIDCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				err := collection.AssignToClient(conn, "my-client", "non-existant-template")
				Expect(err).To(MatchError(collections.TemplateAssignmentError{Err: errors.New("No template with id \"non-existant-template\"")}))
			})
		})

		Context("when the request should reset the template assignment", func() {
			BeforeEach(func() {
				clientsRepo.FindCall.Returns.Client = models.Client{
					ID:         "my-client",
					TemplateID: "some-random-template",
				}
			})

			It("allows template id of empty string to reset the assignment", func() {
				err := collection.AssignToClient(conn, "my-client", "")
				Expect(err).NotTo(HaveOccurred())

				Expect(clientsRepo.FindCall.Receives.Connection).To(Equal(conn))
				Expect(clientsRepo.FindCall.Receives.ClientID).To(Equal("my-client"))

				Expect(clientsRepo.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(clientsRepo.UpdateCall.Receives.Client).To(Equal(models.Client{
					ID:         "my-client",
					TemplateID: models.DefaultTemplateID,
				}))
			})

			It("allows template id of default template id to reset the assignment", func() {
				err := collection.AssignToClient(conn, "my-client", models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())

				Expect(clientsRepo.FindCall.Receives.Connection).To(Equal(conn))
				Expect(clientsRepo.FindCall.Receives.ClientID).To(Equal("my-client"))

				Expect(clientsRepo.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(clientsRepo.UpdateCall.Receives.Client).To(Equal(models.Client{
					ID:         "my-client",
					TemplateID: models.DefaultTemplateID,
				}))
			})
		})

		Context("when it gets an error it doesn't understand", func() {
			Context("on finding the client", func() {
				It("returns any errors it doesn't understand", func() {
					clientsRepo.FindCall.Returns.Error = errors.New("database connection failure")

					err := collection.AssignToClient(conn, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindByIDCall.Returns.Error = errors.New("database failure")

					err := collection.AssignToClient(conn, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					clientsRepo.UpdateCall.Returns.Error = errors.New("database fail")

					err := collection.AssignToClient(conn, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
				})
			})

		})
	})

	Describe("AssignToNotification", func() {
		var kind models.Kind

		BeforeEach(func() {
			clientsRepo.FindCall.Returns.Client = models.Client{
				ID: "my-client",
			}

			kind = models.Kind{
				ID:       "my-kind",
				ClientID: "my-client",
			}

			kindsRepo.FindCall.Returns.Kinds = []models.Kind{kind}
		})

		It("assigns the template to the given kind", func() {
			err := collection.AssignToNotification(conn, "my-client", "my-kind", "my-template")
			Expect(err).NotTo(HaveOccurred())

			Expect(kindsRepo.UpdateCall.Receives.Kind).To(Equal(models.Kind{
				ID:         "my-kind",
				ClientID:   "my-client",
				TemplateID: "my-template",
			}))
		})

		Context("when the request includes a non-existant id", func() {
			It("reports that the client cannot be found", func() {
				kindsRepo.FindCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				err := collection.AssignToNotification(conn, "bad-client", "my-kind", "my-template")
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("not found")}))
			})

			It("reports that the kind cannot be found", func() {
				kindsRepo.FindCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				err := collection.AssignToNotification(conn, "my-client", "bad-kind", "my-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.NotFoundError{Err: errors.New("not found")}))
			})

			It("reports that the template cannot be found", func() {
				templatesRepo.FindByIDCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				err := collection.AssignToNotification(conn, "my-client", "my-kind", "non-existant-template")
				Expect(err).To(MatchError(collections.TemplateAssignmentError{Err: errors.New("No template with id \"non-existant-template\"")}))
			})
		})

		Context("when the request should reset the template assignment", func() {
			BeforeEach(func() {
				kind.TemplateID = "some-random-template"
			})

			It("allows template id of empty string to reset the assignment", func() {
				err := collection.AssignToNotification(conn, "my-client", "my-kind", "")
				Expect(err).NotTo(HaveOccurred())

				Expect(kindsRepo.UpdateCall.Receives.Kind).To(Equal(models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: models.DefaultTemplateID,
				}))
			})

			It("allows template id of default template id to reset the assignment", func() {
				err := collection.AssignToNotification(conn, "my-client", "my-kind", models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())

				Expect(kindsRepo.UpdateCall.Receives.Kind).To(Equal(models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: models.DefaultTemplateID,
				}))
			})
		})

		Context("when it gets an error it doesn't understand", func() {
			Context("on finding the client", func() {
				It("returns any errors it doesn't understand", func() {
					clientsRepo.FindCall.Returns.Error = errors.New("database connection failure")

					err := collection.AssignToNotification(conn, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindByIDCall.Returns.Error = errors.New("database failure")

					err := collection.AssignToNotification(conn, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					kindsRepo.UpdateCall.Returns.Error = errors.New("database fail")

					err := collection.AssignToNotification(conn, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("ListAssociations", func() {
		Context("when a template has been associated to some clients and notifications", func() {
			BeforeEach(func() {
				clientsRepo.FindAllByTemplateIDCall.Returns.Clients = []models.Client{
					{
						ID:         "some-client",
						TemplateID: "some-template-id",
					},
				}

				kindsRepo.FindAllByTemplateIDCall.Returns.Kinds = []models.Kind{
					{
						ID:         "some-notification",
						ClientID:   "some-client",
						TemplateID: "some-template-id",
					},
					{
						ID:         "another-notification",
						ClientID:   "another-client",
						TemplateID: "some-template-id",
					},
				}
			})

			It("returns the full list of associations", func() {
				associations, err := collection.ListAssociations(conn, "some-template-id")
				Expect(err).ToNot(HaveOccurred())

				Expect(associations).To(Equal([]collections.TemplateAssociation{
					{
						ClientID: "some-client",
					},
					{
						ClientID:       "some-client",
						NotificationID: "some-notification",
					},
					{
						ClientID:       "another-client",
						NotificationID: "another-notification",
					},
				}))
				Expect(templatesRepo.FindByIDCall.Receives.Connection).To(Equal(conn))
				Expect(templatesRepo.FindByIDCall.Receives.TemplateID).To(Equal("some-template-id"))
			})
		})

		Context("when errors occur", func() {
			Context("when the clients repo returns an error", func() {
				It("returns the underlying error", func() {
					clientsRepo.FindAllByTemplateIDCall.Returns.Error = errors.New("something bad happened")

					_, err := collection.ListAssociations(conn, "some-template-id")
					Expect(err).To(MatchError(errors.New("something bad happened")))
				})
			})

			Context("when the kinds repo returns an error", func() {
				It("returns the underlying error", func() {
					kindsRepo.FindAllByTemplateIDCall.Returns.Error = errors.New("more bad happened")

					_, err := collection.ListAssociations(conn, "some-template-id")
					Expect(err).To(MatchError(errors.New("more bad happened")))
				})
			})

			Context("when the template repo returns an error", func() {
				It("returns the underlying error", func() {
					templatesRepo.FindByIDCall.Returns.Error = errors.New("something terrible happened")

					_, err := collection.ListAssociations(conn, "some-template-id")
					Expect(err).To(MatchError(errors.New("something terrible happened")))
				})
			})
		})
	})

	Describe("Create", func() {
		It("creates a new template via the templates repo", func() {
			templatesRepo.CreateCall.Returns.Template = models.Template{
				ID:       "some-template-guid",
				Name:     "some-template-name",
				Text:     "some-text",
				HTML:     "some-html",
				Subject:  "some-subject",
				Metadata: "some-metadata",
			}

			template, err := collection.Create(conn, collections.Template{
				Name:     "some-template-name",
				Text:     "some-text",
				HTML:     "some-html",
				Subject:  "some-subject",
				Metadata: "some-metadata",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(template).To(Equal(collections.Template{
				ID:       "some-template-guid",
				Name:     "some-template-name",
				Text:     "some-text",
				HTML:     "some-html",
				Subject:  "some-subject",
				Metadata: "some-metadata",
			}))

			Expect(templatesRepo.CreateCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepo.CreateCall.Receives.Template).To(Equal(models.Template{
				Name:     "some-template-name",
				Text:     "some-text",
				HTML:     "some-html",
				Subject:  "some-subject",
				Metadata: "some-metadata",
			}))
		})

		It("propagates errors from repo", func() {
			templatesRepo.CreateCall.Returns.Error = errors.New("Boom!")

			_, err := collection.Create(conn, collections.Template{})
			Expect(err).To(Equal(errors.New("Boom!")))
		})
	})

	Describe("Delete", func() {
		It("calls destroy on its repo", func() {
			err := collection.Delete(conn, "templateID")
			Expect(err).NotTo(HaveOccurred())

			Expect(templatesRepo.DestroyCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepo.DestroyCall.Receives.TemplateID).To(Equal("templateID"))
		})

		It("returns an error if repo destroy returns an error", func() {
			templatesRepo.DestroyCall.Returns.Error = errors.New("Boom!!")

			err := collection.Delete(conn, "templateID")
			Expect(err).To(MatchError(errors.New("Boom!!")))
		})
	})
})
