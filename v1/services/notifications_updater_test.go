package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationUpdater", func() {
	var (
		notificationsUpdater services.NotificationsUpdater
		kindsRepo            *mocks.KindsRepo
		database             *mocks.Database
		conn                 *mocks.Connection
	)

	BeforeEach(func() {
		kindsRepo = mocks.NewKindsRepo()
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		notificationsUpdater = services.NewNotificationsUpdater(kindsRepo)
	})

	Describe("Update", func() {
		It("updates the correct model with the new fields provided", func() {
			kindsRepo.UpdateCall.Returns.Kind = models.Kind{
				ID:          "my-current-kind-id",
				ClientID:    "my-current-client-id",
				Description: "What a beautiful description",
				TemplateID:  "my-current-template-id",
				Critical:    false,
			}

			err := notificationsUpdater.Update(database, models.Kind{
				ID:          "my-current-kind-id",
				Description: "some-description",
				Critical:    true,
				TemplateID:  "a-brand-new-template",
				ClientID:    "my-current-client-id",
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(kindsRepo.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(kindsRepo.UpdateCall.Receives.Kind).To(Equal(models.Kind{
				ID:          "my-current-kind-id",
				Description: "some-description",
				Critical:    true,
				TemplateID:  "a-brand-new-template",
				ClientID:    "my-current-client-id",
			}))
		})

		It("propagates errors returned by the repo", func() {
			kindsRepo.UpdateCall.Returns.Error = errors.New("Boom")

			err := notificationsUpdater.Update(database, models.Kind{})
			Expect(err).To(MatchError(errors.New("Boom")))
		})
	})
})
