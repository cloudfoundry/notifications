package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateAssociationLister", func() {
	var lister services.TemplateAssociationLister
	var expectedAssociations []services.TemplateAssociation
	var templateID string
	var clientsRepo *fakes.ClientsRepo
	var kindsRepo *fakes.KindsRepo
	var templatesRepo *fakes.TemplatesRepo
	var database *fakes.Database

	Describe("List", func() {
		BeforeEach(func() {
			clientsRepo = fakes.NewClientsRepo()
			kindsRepo = fakes.NewKindsRepo()
			templatesRepo = fakes.NewTemplatesRepo()
			database = fakes.NewDatabase()

			templateID = "a-template-id"
			_, err := templatesRepo.Create(database.Connection(), models.Template{
				ID: templateID,
			})
			if err != nil {
				panic(err)
			}

			lister = services.NewTemplateAssociationLister(clientsRepo, kindsRepo, templatesRepo, database)
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

				_, err := clientsRepo.Create(database.Connection(), models.Client{
					ID:         "some-client",
					TemplateID: templateID,
				})
				if err != nil {
					panic(err)
				}

				_, err = kindsRepo.Create(database.Connection(), models.Kind{
					ID:         "some-notification",
					ClientID:   "some-client",
					TemplateID: templateID,
				})
				if err != nil {
					panic(err)
				}

				_, err = kindsRepo.Create(database.Connection(), models.Kind{
					ID:         "another-notification",
					ClientID:   "another-client",
					TemplateID: templateID,
				})
				if err != nil {
					panic(err)
				}
			})

			It("returns the full list of associations", func() {
				associations, err := lister.List(templateID)
				Expect(err).ToNot(HaveOccurred())
				Expect(associations).To(ConsistOf(expectedAssociations))
			})
		})

		Context("when errors occur", func() {
			Context("when the clients repo returns an error", func() {
				It("returns the underlying error", func() {
					clientsRepo.FindAllByTemplateIDError = errors.New("something bad happened")

					_, err := lister.List(templateID)
					Expect(err).To(MatchError(errors.New("something bad happened")))
				})
			})

			Context("when the kinds repo returns an error", func() {
				It("returns the underlying error", func() {
					kindsRepo.FindAllByTemplateIDError = errors.New("more bad happened")

					_, err := lister.List(templateID)
					Expect(err).To(MatchError(errors.New("more bad happened")))
				})
			})

			Context("when the template repo returns an error", func() {
				It("returns the underlying error", func() {
					templatesRepo.FindError = errors.New("something terrible happened")

					_, err := lister.List(templateID)
					Expect(err).To(MatchError(errors.New("something terrible happened")))
				})
			})
		})
	})
})
