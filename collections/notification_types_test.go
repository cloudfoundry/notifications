package collections_test

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationTypesCollection", func() {
	var (
		notificationTypesCollection     collections.NotificationTypesCollection
		fakeNotificationTypesRepository *fakes.NotificationTypesRepository
		fakeDatabaseConnection          *fakes.Connection
	)

	BeforeEach(func() {
		fakeNotificationTypesRepository = fakes.NewNotificationTypesRepository()
		fakeNotificationTypesRepository.InsertCall.ReturnNotificationType.ID = "generated-id"

		notificationTypesCollection = collections.NewNotificationTypesCollection(fakeNotificationTypesRepository)
		fakeDatabaseConnection = fakes.NewConnection()
	})

	Describe("Add", func() {
		It("adds a notification type to the collection", func() {
			notificationType := collections.NotificationType{
				Name:        "My cool notification type",
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}

			returnedNotificationType, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedNotificationType.ID).To(Equal("generated-id"))
			Expect(fakeNotificationTypesRepository.InsertCall.Connection).To(Equal(fakeDatabaseConnection))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.Name).To(Equal(notificationType.Name))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.Description).To(Equal(notificationType.Description))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.Critical).To(Equal(notificationType.Critical))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.TemplateID).To(Equal(notificationType.TemplateID))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.SenderID).To(Equal(notificationType.SenderID))
		})
	})
})
