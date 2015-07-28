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
		notificationTypesCollection     collections.CampaignTypesCollection
		fakeNotificationTypesRepository *fakes.NotificationTypesRepository
		fakeSendersRepository           *fakes.SendersRepository
		fakeDatabaseConnection          *fakes.Connection
	)

	BeforeEach(func() {
		fakeNotificationTypesRepository = fakes.NewNotificationTypesRepository()
		fakeSendersRepository = fakes.NewSendersRepository()

		notificationTypesCollection = collections.NewCampaignTypesCollection(fakeNotificationTypesRepository, fakeSendersRepository)
		fakeDatabaseConnection = fakes.NewConnection()
	})

	Describe("Add", func() {
		var (
			notificationType collections.CampaignType
		)

		BeforeEach(func() {
			notificationType = collections.CampaignType{
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
			notificationType = collections.CampaignType{
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
			Expect(err).To(MatchError(collections.NewValidationError("missing campaign type name")))
		})

		It("requires a description to be specified", func() {
			notificationType = collections.CampaignType{
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
			Expect(err).To(MatchError(collections.NewValidationError("missing campaign type description")))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")

				_, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "sender not found",
					Err:     models.RecordNotFoundError("sender not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := notificationTypesCollection.Add(fakeDatabaseConnection, notificationType, "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender not found")))
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
			It("validates that a sender id was specified", func() {
				_, err := notificationTypesCollection.List(fakeDatabaseConnection, "", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing sender id")))
			})

			It("validates that a client id was specified", func() {
				_, err := notificationTypesCollection.List(fakeDatabaseConnection, "some-sender-id", "")
				Expect(err).To(MatchError(collections.NewValidationError("missing client id")))
			})

			It("generates a not found error when the sender does not exist", func() {
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")

				_, err := notificationTypesCollection.List(fakeDatabaseConnection, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "sender not found",
					Err:     models.RecordNotFoundError("sender not found"),
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
				Expect(err).To(MatchError(collections.NewNotFoundError("sender not found")))
			})

			It("handles unexpected database errors", func() {
				fakeNotificationTypesRepository.ListCall.ReturnNotificationTypeList = []models.NotificationType{}
				fakeNotificationTypesRepository.ListCall.Err = errors.New("BOOM!")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}

				_, err := notificationTypesCollection.List(fakeDatabaseConnection, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})

	Describe("Get", func() {
		It("returns the ID if it is found", func() {
			fakeNotificationTypesRepository.GetReturn.NotificationType = models.NotificationType{
				ID:       "a-notification-type-id",
				Name:     "typename",
				SenderID: "senderID",
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "senderID",
				Name:     "I dont matter",
				ClientID: "some-client-id",
			}

			notificationType, err := notificationTypesCollection.Get(fakeDatabaseConnection, "a-notification-type-id", "senderID", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(notificationType.Name).To(Equal("typename"))
		})

		Context("failure cases", func() {
			It("validates that a notification type id was specified", func() {
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing campaign type id")))
			})

			It("validates that a sender id was specified", func() {
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing sender id")))
			})

			It("validates that a client id was specified", func() {
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "some-sender-id", "")
				Expect(err).To(MatchError(collections.NewValidationError("missing client id")))
			})

			It("returns a not found error if the notification type does not exist", func() {
				fakeNotificationTypesRepository.GetReturn.Err = models.RecordNotFoundError("campaign type not found")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "I dont matter",
					ClientID: "some-client-id",
				}
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "missing-notification-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "campaign type missing-notification-type-id not found",
					Err:     models.RecordNotFoundError("campaign type not found"),
				}))
			})

			It("returns a not found error if the sender does not exist", func() {
				fakeNotificationTypesRepository.GetReturn.NotificationType = models.NotificationType{
					ID:       "some-notification-type-id",
					Name:     "typename",
					SenderID: "some-sender-id",
				}
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "missing-sender-id", "some-client-id")
				Expect(err.(collections.NotFoundError).Message).To(Equal("sender some-notification-type-id not found"))
			})

			It("returns a not found error if the notification type does not belong to the sender", func() {
				fakeNotificationTypesRepository.GetReturn.NotificationType = models.NotificationType{
					ID:       "some-notification-type-id",
					Name:     "typename",
					SenderID: "my-sender-id",
				}
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "someone-elses-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "someone-elses-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("campaign type some-notification-type-id not found")))
			})

			It("returns a not found error if the sender does not belong to the client", func() {
				fakeNotificationTypesRepository.GetReturn.NotificationType = models.NotificationType{
					ID:       "some-notification-type-id",
					Name:     "typename",
					SenderID: "my-sender-id",
				}
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "my-sender-id",
					Name:     "some-sender",
					ClientID: "client_id",
				}
				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "my-sender-id", "someone-elses-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender my-sender-id not found")))
			})

			It("handles unexpected database errors from the senders repository", func() {
				fakeSendersRepository.GetCall.Err = errors.New("BOOM!")

				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})

			It("handles unexpected database errors from the notification types repository", func() {
				fakeNotificationTypesRepository.GetReturn.Err = errors.New("BOOM!")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}

				_, err := notificationTypesCollection.Get(fakeDatabaseConnection, "some-notification-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})
})
