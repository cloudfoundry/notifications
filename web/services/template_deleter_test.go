package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deleter", func() {
	var (
		deleter       services.TemplateDeleter
		templatesRepo *fakes.TemplatesRepo
		database      *fakes.Database
	)

	BeforeEach(func() {
		database = fakes.NewDatabase()

		templatesRepo = fakes.NewTemplatesRepo()
		_, err := templatesRepo.Create(database.Connection(), models.Template{
			ID: "templateID",
		})
		Expect(err).NotTo(HaveOccurred())

		deleter = services.NewTemplateDeleter(templatesRepo)
	})

	Describe("#Delete", func() {
		It("calls destroy on its repo", func() {
			err := deleter.Delete(database, "templateID")
			Expect(err).NotTo(HaveOccurred())
			Expect(database.ConnectionWasCalled).To(BeTrue())

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
