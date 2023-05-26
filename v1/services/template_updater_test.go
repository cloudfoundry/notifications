package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updater", func() {
	Describe("Update", func() {
		var (
			conn          *mocks.Connection
			database      *mocks.Database
			templatesRepo *mocks.TemplatesRepo
			updater       services.TemplateUpdater
		)

		BeforeEach(func() {
			conn = mocks.NewConnection()
			database = mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = conn
			templatesRepo = mocks.NewTemplatesRepo()

			updater = services.NewTemplateUpdater(templatesRepo)
		})

		It("Inserts templates into the templates repo", func() {
			err := updater.Update(database, "my-awesome-id", models.Template{
				Name: "gobble template",
				Text: "gobble",
				HTML: "<p>gobble</p>",
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(templatesRepo.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(templatesRepo.UpdateCall.Receives.TemplateID).To(Equal("my-awesome-id"))
			Expect(templatesRepo.UpdateCall.Receives.Template).To(Equal(models.Template{
				Name: "gobble template",
				Text: "gobble",
				HTML: "<p>gobble</p>",
			}))
		})

		It("propagates errors from repo", func() {
			templatesRepo.UpdateCall.Returns.Error = errors.New("Boom!")

			err := updater.Update(database, "unimportant", models.Template{})
			Expect(err).To(MatchError(errors.New("Boom!")))
		})
	})
})
