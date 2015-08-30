package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateAssociationLister", func() {
	var (
		lister               services.TemplateAssociationLister
		expectedAssociations []services.TemplateAssociation
		templateID           string
		clientsRepo          *mocks.ClientsRepository
		kindsRepo            *mocks.KindsRepo
		templatesRepo        *mocks.TemplatesRepo
		database             *mocks.Database
		conn                 *mocks.Connection
	)

	Describe("List", func() {
		BeforeEach(func() {
			clientsRepo = mocks.NewClientsRepository()
			kindsRepo = mocks.NewKindsRepo()
			templatesRepo = mocks.NewTemplatesRepo()
			conn = mocks.NewConnection()
			database = mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = conn

			templateID = "a-template-id"
			_, err := templatesRepo.Create(conn, models.Template{
				ID: templateID,
			})
			Expect(err).NotTo(HaveOccurred())
			lister = services.NewTemplateAssociationLister(clientsRepo, kindsRepo, templatesRepo)
		})

		Context("when a template has been associated to some clients and notifications", func() {
			BeforeEach(func() {
				expectedAssociations = []services.TemplateAssociation{
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
				}

				clientsRepo.FindAllByTemplateIDCall.Returns.Clients = []models.Client{
					{
						ID:         "some-client",
						TemplateID: templateID,
					},
				}

				_, err := kindsRepo.Create(database.Connection(), models.Kind{
					ID:         "some-notification",
					ClientID:   "some-client",
					TemplateID: templateID,
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = kindsRepo.Create(database.Connection(), models.Kind{
					ID:         "another-notification",
					ClientID:   "another-client",
					TemplateID: templateID,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the full list of associations", func() {
				associations, err := lister.List(database, templateID)
				Expect(err).ToNot(HaveOccurred())

				Expect(associations).To(ConsistOf(expectedAssociations))
				Expect(templatesRepo.FindByIDCall.Receives.Connection).To(Equal(conn))
				Expect(templatesRepo.FindByIDCall.Receives.TemplateID).To(Equal("a-template-id"))
			})
		})

		Context("when errors occur", func() {
			Context("when the clients repo returns an error", func() {
				It("returns the underlying error", func() {
					clientsRepo.FindAllByTemplateIDCall.Returns.Error = errors.New("something bad happened")

					_, err := lister.List(database, templateID)
					Expect(err).To(MatchError(errors.New("something bad happened")))
				})
			})

			Context("when the kinds repo returns an error", func() {
				It("returns the underlying error", func() {
					kindsRepo.FindAllByTemplateIDError = errors.New("more bad happened")

					_, err := lister.List(database, templateID)
					Expect(err).To(MatchError(errors.New("more bad happened")))
				})
			})

			Context("when the template repo returns an error", func() {
				It("returns the underlying error", func() {
					templatesRepo.FindError = errors.New("something terrible happened")

					_, err := lister.List(database, templateID)
					Expect(err).To(MatchError(errors.New("something terrible happened")))
				})
			})
		})
	})
})
