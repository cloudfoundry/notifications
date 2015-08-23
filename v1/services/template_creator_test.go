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
		)

		BeforeEach(func() {
			templatesRepo = mocks.NewTemplatesRepo()
			template = models.Template{
				Name:    "Big Hero 6 Template",
				Text:    "Adorable robot.",
				HTML:    "<p>Many heroes.</p>",
				Subject: "Robots and Heroes",
			}
			database = mocks.NewDatabase()
			creator = services.NewTemplateCreator(templatesRepo)
		})

		It("Creates a new template via the templates repo", func() {
			Expect(templatesRepo.Templates).ToNot(ContainElement(template))

			_, err := creator.Create(database, template)
			Expect(err).ToNot(HaveOccurred())
			Expect(templatesRepo.Templates).To(ContainElement(template))
		})

		It("propagates errors from repo", func() {
			expectedErr := errors.New("Boom!")
			templatesRepo.CreateError = expectedErr

			_, err := creator.Create(database, template)

			Expect(err).To(Equal(expectedErr))
		})
	})
})
