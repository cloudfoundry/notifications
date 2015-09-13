package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateAssigner", func() {
	var (
		assigner      services.TemplateAssigner
		kindsRepo     *mocks.KindsRepo
		clientsRepo   *mocks.ClientsRepository
		templatesRepo *mocks.TemplatesRepo
		conn          *mocks.Connection
		database      *mocks.Database
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		clientsRepo = mocks.NewClientsRepository()
		kindsRepo = mocks.NewKindsRepo()
		templatesRepo = mocks.NewTemplatesRepo()

		assigner = services.NewTemplateAssigner(clientsRepo, kindsRepo, templatesRepo)
	})

	Describe("AssignToClient", func() {
		BeforeEach(func() {
			var err error

			clientsRepo.FindCall.Returns.Client = models.Client{
				ID: "my-client",
			}

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "default",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "my-template",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("assigns the template to the given client", func() {
			err := assigner.AssignToClient(database, "my-client", "my-template")
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
				clientsRepo.FindCall.Returns.Error = models.NotFoundError{errors.New("not found")}

				err := assigner.AssignToClient(database, "missing-client", "my-template")
				Expect(err).To(MatchError(models.NotFoundError{errors.New("not found")}))
			})

			It("reports that the template cannot be found", func() {
				templatesRepo.FindByIDCall.Returns.Error = models.NotFoundError{errors.New("not found")}

				err := assigner.AssignToClient(database, "my-client", "non-existant-template")
				Expect(err).To(MatchError(services.TemplateAssignmentError{errors.New("No template with id \"non-existant-template\"")}))
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
				err := assigner.AssignToClient(database, "my-client", "")
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
				err := assigner.AssignToClient(database, "my-client", models.DefaultTemplateID)
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

					err := assigner.AssignToClient(database, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindByIDCall.Returns.Error = errors.New("database failure")

					err := assigner.AssignToClient(database, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					clientsRepo.UpdateCall.Returns.Error = errors.New("database fail")

					err := assigner.AssignToClient(database, "my-client", "my-template")
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

			_, err := templatesRepo.Create(conn, models.Template{
				ID: "default",
			})
			Expect(err).NotTo(HaveOccurred())

			kind = models.Kind{
				ID:       "my-kind",
				ClientID: "my-client",
			}

			kindsRepo.FindCall.Returns.Kinds = []models.Kind{kind}

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "my-template",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("assigns the template to the given kind", func() {
			err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
			Expect(err).NotTo(HaveOccurred())

			Expect(kindsRepo.UpdateCall.Receives.Kind).To(Equal(models.Kind{
				ID:         "my-kind",
				ClientID:   "my-client",
				TemplateID: "my-template",
			}))
		})

		Context("when the request includes a non-existant id", func() {
			It("reports that the client cannot be found", func() {
				kindsRepo.FindCall.Returns.Error = models.NotFoundError{errors.New("not found")}

				err := assigner.AssignToNotification(database, "bad-client", "my-kind", "my-template")
				Expect(err).To(MatchError(models.NotFoundError{errors.New("not found")}))
			})

			It("reports that the kind cannot be found", func() {
				kindsRepo.FindCall.Returns.Error = models.NotFoundError{errors.New("not found")}

				err := assigner.AssignToNotification(database, "my-client", "bad-kind", "my-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.NotFoundError{errors.New("not found")}))
			})

			It("reports that the template cannot be found", func() {
				templatesRepo.FindByIDCall.Returns.Error = models.NotFoundError{errors.New("not found")}

				err := assigner.AssignToNotification(database, "my-client", "my-kind", "non-existant-template")
				Expect(err).To(MatchError(services.TemplateAssignmentError{errors.New("No template with id \"non-existant-template\"")}))
			})
		})

		Context("when the request should reset the template assignment", func() {
			BeforeEach(func() {
				var err error
				kind.TemplateID = "some-random-template"
				kind, err = kindsRepo.Update(conn, kind)
				Expect(err).NotTo(HaveOccurred())
			})

			It("allows template id of empty string to reset the assignment", func() {
				err := assigner.AssignToNotification(database, "my-client", "my-kind", "")
				Expect(err).NotTo(HaveOccurred())

				Expect(kindsRepo.UpdateCall.Receives.Kind).To(Equal(models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: models.DefaultTemplateID,
				}))
			})

			It("allows template id of default template id to reset the assignment", func() {
				err := assigner.AssignToNotification(database, "my-client", "my-kind", models.DefaultTemplateID)
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

					err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindByIDCall.Returns.Error = errors.New("database failure")

					err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					kindsRepo.UpdateCall.Returns.Error = errors.New("database fail")

					err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
