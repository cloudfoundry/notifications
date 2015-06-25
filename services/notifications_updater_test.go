package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationUpdater", func() {
	var (
		notificationsUpdater services.NotificationsUpdater
		kindsRepo            *fakes.KindsRepo
		database             *fakes.Database
		clientID             string
		notificationID       string
	)

	BeforeEach(func() {
		kindsRepo = fakes.NewKindsRepo()
		database = fakes.NewDatabase()
		clientID = "my-current-client-id"
		notificationID = "my-current-kind-id"

		notificationsUpdater = services.NewNotificationsUpdater(kindsRepo)
	})

	Describe("Update", func() {
		It("updates the correct model with the new fields provided", func() {
			kindsRepo.Kinds[notificationID+clientID] = models.Kind{
				ID:          notificationID,
				ClientID:    clientID,
				Description: "What a beautiful description",
				TemplateID:  "my-current-template-id",
				Critical:    false,
			}

			err := notificationsUpdater.Update(database, models.Kind{
				Description: "some-description",
				Critical:    true,
				TemplateID:  "a-brand-new-template",
				ID:          notificationID,
				ClientID:    clientID,
			})

			Expect(err).ToNot(HaveOccurred())
			updatedKind := kindsRepo.Kinds[notificationID+clientID]

			Expect(updatedKind.Description).To(Equal("some-description"))
			Expect(updatedKind.TemplateID).To(Equal("a-brand-new-template"))
			Expect(updatedKind.Critical).To(BeTrue())
			Expect(updatedKind.ID).To(Equal(notificationID))
			Expect(updatedKind.ClientID).To(Equal(clientID))

			Expect(database.ConnectionWasCalled).To(BeTrue())
		})

		It("propagates errors returned by the repo", func() {
			boomError := errors.New("Boom")
			kindsRepo.UpdateError = boomError
			err := notificationsUpdater.Update(database, models.Kind{})

			Expect(err).To(Equal(boomError))
		})
	})
})
