package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deleter", func() {
	var (
		deleter       services.TemplateDeleter
		templatesRepo *mocks.TemplatesRepo
		database      *mocks.Database
		conn          *mocks.Connection
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		templatesRepo = mocks.NewTemplatesRepo()
		_, err := templatesRepo.Create(conn, models.Template{
			ID: "templateID",
		})
		Expect(err).NotTo(HaveOccurred())

		deleter = services.NewTemplateDeleter(templatesRepo)
	})

	Describe("#Delete", func() {
		It("calls destroy on its repo", func() {
			err := deleter.Delete(database, "templateID")
			Expect(err).NotTo(HaveOccurred())

			Expect(templatesRepo.DestroyCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepo.DestroyCall.Receives.TemplateID).To(Equal("templateID"))

			_, err = templatesRepo.FindByID(database.Connection(), "templateID")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
		})

		It("returns an error if repo destroy returns an error", func() {
			templatesRepo.DestroyError = errors.New("Boom!!")

			err := deleter.Delete(database, "templateID")
			Expect(err).To(Equal(templatesRepo.DestroyError))
		})
	})
})
