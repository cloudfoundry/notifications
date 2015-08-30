package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updater", func() {
	Describe("Update", func() {
		var (
			templatesRepo *mocks.TemplatesRepo
			template      models.Template
			updater       services.TemplateUpdater
			database      *mocks.Database
			conn          *mocks.Connection
		)

		BeforeEach(func() {
			templatesRepo = mocks.NewTemplatesRepo()
			template = models.Template{
				Name: "gobble template",
				Text: "gobble",
				HTML: "<p>gobble</p>",
			}
			conn = mocks.NewConnection()
			database = mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = conn

			updater = services.NewTemplateUpdater(templatesRepo)
		})

		It("Inserts templates into the templates repo", func() {
			Expect(templatesRepo.Templates).ToNot(ContainElement(template))

			err := updater.Update(database, "my-awesome-id", template)
			Expect(err).ToNot(HaveOccurred())

			Expect(templatesRepo.Templates).To(ContainElement(template))
			Expect(templatesRepo.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepo.UpdateCall.Receives.TemplateID).To(Equal("my-awesome-id"))
			Expect(templatesRepo.UpdateCall.Receives.Template).To(Equal(models.Template{
				Name: "gobble template",
				Text: "gobble",
				HTML: "<p>gobble</p>",
			}))
		})

		It("propagates errors from repo", func() {
			expectedErr := errors.New("Boom!")

			templatesRepo.UpdateError = expectedErr
			err := updater.Update(database, "unimportant", template)

			Expect(err).To(Equal(expectedErr))
		})
	})
})
