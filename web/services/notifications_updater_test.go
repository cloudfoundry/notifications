package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationUpdater", func() {
	var notificationsUpdater services.NotificationsUpdater
	var kindsRepo *fakes.KindsRepo
	var clientID string
	var notificationID string

	BeforeEach(func() {
		kindsRepo = fakes.NewKindsRepo()
		notificationsUpdater = services.NewNotificationsUpdater(kindsRepo, fakes.NewDatabase())
		clientID = "my-current-client-id"
		notificationID = "my-current-kind-id"
	})

	Describe("Update", func() {
		BeforeEach(func() {
			kindsRepo.Kinds[notificationID+clientID] = models.Kind{
				ID:          notificationID,
				ClientID:    clientID,
				Description: "What a beautiful description",
				TemplateID:  "my-current-template-id",
				Critical:    false,
			}
		})

		It("updates the correct model with the new fields provided", func() {
			err := notificationsUpdater.Update(clientID, notificationID, models.Kind{
				Description: "some-description",
				Critical:    true,
				TemplateID:  "a-brand-new-template"})

			Expect(err).ToNot(HaveOccurred())
			updatedKind := kindsRepo.Kinds[notificationID+clientID]

			Expect(updatedKind.Description).To(Equal("some-description"))
			Expect(updatedKind.TemplateID).To(Equal("a-brand-new-template"))
			Expect(updatedKind.Critical).To(BeTrue())
			Expect(updatedKind.ID).To(Equal(notificationID))
			Expect(updatedKind.ClientID).To(Equal(clientID))
		})

		It("propagates errors returned by the repo", func() {
			boomError := errors.New("Boom")
			kindsRepo.UpdateError = boomError
			err := notificationsUpdater.Update(clientID, notificationID, models.Kind{})

			Expect(err).To(Equal(boomError))
		})
	})
})
