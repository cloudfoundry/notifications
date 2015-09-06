package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Creator", func() {
	Describe("Create", func() {
		var (
			templatesRepo *mocks.TemplatesRepo
			template      models.Template
			creator       services.TemplateCreator
			database      *mocks.Database
			conn          *mocks.Connection
		)

		BeforeEach(func() {
			templatesRepo = mocks.NewTemplatesRepo()
			template = models.Template{
				Name:    "Big Hero 6 Template",
				Text:    "Adorable robot.",
				HTML:    "<p>Many heroes.</p>",
				Subject: "Robots and Heroes",
			}

			conn = mocks.NewConnection()
			database = mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = conn

			creator = services.NewTemplateCreator(templatesRepo)
		})

		It("Creates a new template via the templates repo", func() {
			_, err := creator.Create(database, template)
			Expect(err).ToNot(HaveOccurred())

			Expect(templatesRepo.CreateCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepo.CreateCall.Receives.Template).To(Equal(template))
		})

		It("propagates errors from repo", func() {
			templatesRepo.CreateCall.Returns.Error = errors.New("Boom!")

			_, err := creator.Create(database, template)
			Expect(err).To(Equal(errors.New("Boom!")))
		})
	})
})
