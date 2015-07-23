package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationTypesCollection", func() {
	var (
		notificationTypesCollection     collections.NotificationTypesCollection
		fakeNotificationTypesRepository *fakes.NotificationTypesRepository
		fakeSendersRepository           *fakes.SendersRepository
		fakeDatabaseConnection          *fakes.Connection
	)

	BeforeEach(func() {
		fakeNotificationTypesRepository = fakes.NewNotificationTypesRepository()
		fakeSendersRepository = fakes.NewSendersRepository()

		notificationTypesCollection = collections.NewNotificationTypesCollection(fakeNotificationTypesRepository, fakeSendersRepository)
		fakeDatabaseConnection = fakes.NewConnection()
	})

	Describe("Add", func() {
		var (
			notificationType collections.NotificationType
		)

		BeforeEach(func() {
			notificationType = collections.NotificationType{
				Name:        "My cool notification type",
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}

			fakeNotificationTypesRepository.InsertCall.ReturnNotificationType.ID = "generated-id"
		})

		It("adds a notification type to the collection", func() {
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client_id",
			}

			returnedNotificationType, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "client_id")
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedNotificationType.ID).To(Equal("generated-id"))
			Expect(fakeNotificationTypesRepository.InsertCall.Connection).To(Equal(fakeDatabaseConnection))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.Name).To(Equal(notificationType.Name))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.Description).To(Equal(notificationType.Description))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.Critical).To(Equal(notificationType.Critical))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.TemplateID).To(Equal(notificationType.TemplateID))
			Expect(fakeNotificationTypesRepository.InsertCall.NotificationType.SenderID).To(Equal(notificationType.SenderID))
		})

		It("requires a name to be specified", func() {
			notificationType = collections.NotificationType{
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client_id",
			}

			_, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "client_id")
			Expect(err).To(MatchError(collections.ValidationError{
				Err: errors.New("missing notification type name"),
			}))
		})

		It("requires a description to be specified", func() {
			notificationType = collections.NotificationType{
				Name:       "some-notification-type",
				Critical:   false,
				TemplateID: "",
				SenderID:   "mysender",
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client_id",
			}

			_, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "client_id")
			Expect(err).To(MatchError(collections.ValidationError{
				Err: errors.New("missing notification type description"),
			}))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				fakeNotificationTypesRepository.InsertCall.Err = models.RecordNotFoundError("sender not found")
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")

				_, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Err: models.RecordNotFoundError("sender not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeNotificationTypesRepository.InsertCall.Err = models.RecordNotFoundError("sender not found")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Err: errors.New("sender not found"),
				}))
			})
		})
	})

	Describe("List", func() {
		It("retrieves a list of notification types from the collection", func() {
			fakeNotificationTypesRepository.ListCall.ReturnNotificationTypeList = []models.NotificationType{
				{
					ID:          "notification-type-id-one",
					Name:        "first-notification-type",
					Description: "first-notification-type-description",
					Critical:    false,
					TemplateID:  "",
					SenderID:    "some-sender-id",
				},
				{
					ID:          "notification-type-id-two",
					Name:        "second-notification-type",
					Description: "second-notification-type-description",
					Critical:    true,
					TemplateID:  "",
					SenderID:    "some-sender-id",
				},
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			returnedNotificationTypeList, err := notificationTypesCollection.List(fakeDatabaseConnection, "some-sender-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(returnedNotificationTypeList)).To(Equal(2))

			Expect(returnedNotificationTypeList[0].ID).To(Equal("notification-type-id-one"))
			Expect(returnedNotificationTypeList[0].SenderID).To(Equal("some-sender-id"))

			Expect(returnedNotificationTypeList[1].ID).To(Equal("notification-type-id-two"))
			Expect(returnedNotificationTypeList[1].SenderID).To(Equal("some-sender-id"))
		})

		It("retrieves an empty list of notification types from the collection if no records have been added", func() {
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			returnedNotificationTypeList, err := notificationTypesCollection.List(fakeDatabaseConnection, "some-senderid", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(returnedNotificationTypeList)).To(Equal(0))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				fakeNotificationTypesRepository.ListCall.Err = models.RecordNotFoundError("sender not found")
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")

				_, err := notificationTypesCollection.List(fakeDatabaseConnection, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Err: models.RecordNotFoundError("sender not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeNotificationTypesRepository.ListCall.Err = models.RecordNotFoundError("sender not found")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := notificationTypesCollection.List(fakeDatabaseConnection, "mismatch-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Err: errors.New("sender not found"),
				}))
			})
		})
	})
})
