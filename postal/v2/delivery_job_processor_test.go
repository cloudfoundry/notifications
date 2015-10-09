package v2_test

import (
	"bytes"
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/postal/v2"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeliveryJobProcessor", func() {
	var (
		buffer                  *bytes.Buffer
		processor               v2.DeliveryJobProcessor
		logger                  lager.Logger
		mailClient              *mocks.MailClient
		userLoader              *mocks.UserLoader
		tokenLoader             *mocks.TokenLoader
		messageStatusUpdater    *mocks.MessageStatusUpdater
		packager                *mocks.Packager
		conn                    *mocks.Connection
		database                *mocks.Database
		delivery                common.Delivery
		campaignsRepository     *mocks.CampaignsRepository
		unsubscribersRepository *mocks.UnsubscribersRepository
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		logger = logger.Session("worker", lager.Data{"worker_id": 1234})

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		messageStatusUpdater = mocks.NewMessageStatusUpdater()

		tokenLoader = mocks.NewTokenLoader()
		tokenLoader.LoadCall.Returns.Token = "some-token"

		userLoader = mocks.NewUserLoader()
		userLoader.LoadCall.Returns.Users = map[string]uaa.User{
			"user-123": {
				Emails: []string{"user-123@example.com"},
			},
		}

		campaignsRepository = mocks.NewCampaignsRepository()
		unsubscribersRepository = mocks.NewUnsubscribersRepository()
		unsubscribersRepository.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("not unsubscribed == will be delivered!")}

		packager = mocks.NewPackager()
		packager.PrepareContextCall.Returns.MessageContext = common.MessageContext{
			From:    "from@example.com",
			ReplyTo: "thesender@example.com",
			To:      "user-123@example.com",
			Subject: "the subject",
			Text:    "body content",
			HTML:    "",
			HTMLComponents: common.HTML{
				BodyContent:    "",
				BodyAttributes: "",
				Head:           "",
				Doctype:        "",
			},
			TextTemplate:      "{{.Text}} {{.Domain}}",
			HTMLTemplate:      "<p>{{.HTML}}</p>",
			SubjectTemplate:   "{{.Subject}}",
			KindDescription:   "some-kind",
			SourceDescription: "some-client",
			UserGUID:          "user-123",
			ClientID:          "some-client",
			MessageID:         "randomly-generated-guid",
			Space:             "",
			SpaceGUID:         "",
			Organization:      "",
			OrganizationGUID:  "",
			UnsubscribeID:     "eFGlsyNvaxtJ_lbV6KcY9BCb6O7H78pEPcLIARVkbTQt4dDrf2sqFjd9pfOOi439mVtNrTZJwhM=",
			Scope:             "",
			Endorsement:       "",
			OrganizationRole:  "",
			RequestReceived:   time.Time{},
			Domain:            "example.com",
		}
		packager.PackCall.Returns.Message = mail.Message{
			Date:                    "",
			MimeVersion:             "",
			ContentType:             "",
			ContentTransferEncoding: "",
			From:    "from@example.com",
			ReplyTo: "thesender@example.com",
			To:      "user-123@example.com",
			Subject: "the subject",
			Body: []mail.Part{
				{
					ContentType: "text/plain",
					Content:     "body content example.com",
				},
			},
			Headers: []string{
				"X-CF-Client-ID: some-client",
				"X-CF-Notification-ID: randomly-generated-guid",
				"X-CF-Notification-Timestamp: 2015-09-10T12:27:11.675885866-07:00",
				"X-CF-Notification-Request-Received: 0001-01-01T00:00:00Z",
			},
			CompiledBody: "",
		}

		mailClient = mocks.NewMailClient()

		delivery = common.Delivery{
			ClientID: "some-client",
			UserGUID: "user-123",
			Options: common.Options{
				Subject:    "the subject",
				Text:       "body content",
				ReplyTo:    "thesender@example.com",
				KindID:     "some-kind",
				TemplateID: "some-template-id",
			},
			MessageID:     "randomly-generated-guid",
			VCAPRequestID: "some-request-id",
			CampaignID:    "some-campaign-id",
		}

		processor = v2.NewDeliveryJobProcessor(mailClient, packager, userLoader, tokenLoader, messageStatusUpdater, database, unsubscribersRepository, campaignsRepository, "from@example.com", "example.com", "uaa-host")
	})

	It("ensures message delivery", func() {
		err := processor.Process(delivery, logger)
		Expect(err).NotTo(HaveOccurred())

		Expect(tokenLoader.LoadCall.Receives.UAAHost).To(Equal("uaa-host"))

		Expect(userLoader.LoadCall.Receives.UserGUIDs).To(Equal([]string{"user-123"}))
		Expect(userLoader.LoadCall.Receives.Token).To(Equal("some-token"))

		delivery.Email = "user-123@example.com"
		Expect(packager.PrepareContextCall.Receives.Delivery).To(Equal(delivery))
		Expect(packager.PrepareContextCall.Receives.Sender).To(Equal("from@example.com"))
		Expect(packager.PrepareContextCall.Receives.Domain).To(Equal("example.com"))

		Expect(packager.PackCall.Receives.MessageContext).To(Equal(packager.PrepareContextCall.Returns.MessageContext))

		Expect(mailClient.SendCall.CallCount).To(Equal(1))
		msg := mailClient.SendCall.Receives.Message
		Expect(msg.From).To(Equal("from@example.com"))
		Expect(msg.ReplyTo).To(Equal("thesender@example.com"))
		Expect(msg.To).To(Equal("user-123@example.com"))
		Expect(msg.Subject).To(Equal("the subject"))
		Expect(msg.Body).To(ConsistOf([]mail.Part{
			{
				ContentType: "text/plain",
				Content:     "body content example.com",
			},
		}))
		Expect(msg.Headers).To(Equal([]string{
			"X-CF-Client-ID: some-client",
			"X-CF-Notification-ID: randomly-generated-guid",
			"X-CF-Notification-Timestamp: 2015-09-10T12:27:11.675885866-07:00",
			"X-CF-Notification-Request-Received: 0001-01-01T00:00:00Z",
		}))
	})

	It("updates the message status as delivered", func() {
		err := processor.Process(delivery, logger)
		Expect(err).NotTo(HaveOccurred())

		Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
		Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal("randomly-generated-guid"))
		Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(common.StatusDelivered))
		Expect(messageStatusUpdater.UpdateCall.Receives.CampaignID).To(Equal("some-campaign-id"))
		Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
	})

	Context("when the delivery does not have a user GUID", func() {
		BeforeEach(func() {
			delivery.Email = "user-123@example.com"
			delivery.UserGUID = ""
		})

		It("should not call the userLoader", func() {
			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(tokenLoader.LoadCall.Receives.UAAHost).To(BeEmpty())

			Expect(userLoader.LoadCall.Receives.UserGUIDs).To(BeNil())
			Expect(userLoader.LoadCall.Receives.Token).To(Equal(""))
		})

		It("should call PrepareContext with the correct arguments", func() {
			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(packager.PrepareContextCall.Receives.Delivery).To(Equal(delivery))
			Expect(packager.PrepareContextCall.Receives.Sender).To(Equal("from@example.com"))
			Expect(packager.PrepareContextCall.Receives.Domain).To(Equal("example.com"))
		})

		It("should call Pack with the correct arguments", func() {
			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(packager.PackCall.Receives.MessageContext).To(Equal(packager.PrepareContextCall.Returns.MessageContext))
		})

		It("should still deliver the email (assuming there is an email address)", func() {
			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(mailClient.SendCall.CallCount).To(Equal(1))
			msg := mailClient.SendCall.Receives.Message
			Expect(msg.From).To(Equal("from@example.com"))
			Expect(msg.ReplyTo).To(Equal("thesender@example.com"))
			Expect(msg.To).To(Equal("user-123@example.com"))
			Expect(msg.Subject).To(Equal("the subject"))
			Expect(msg.Body).To(ConsistOf([]mail.Part{
				{
					ContentType: "text/plain",
					Content:     "body content example.com",
				},
			}))
			Expect(msg.Headers).To(Equal([]string{
				"X-CF-Client-ID: some-client",
				"X-CF-Notification-ID: randomly-generated-guid",
				"X-CF-Notification-Timestamp: 2015-09-10T12:27:11.675885866-07:00",
				"X-CF-Notification-Request-Received: 0001-01-01T00:00:00Z",
			}))
		})

		Context("when the delivery does not have an email address", func() {
			BeforeEach(func() {
				delivery.Email = ""
			})

			It("should not call the userLoader", func() {
				err := processor.Process(delivery, logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(tokenLoader.LoadCall.Receives.UAAHost).To(BeEmpty())

				Expect(userLoader.LoadCall.Receives.UserGUIDs).To(BeNil())
				Expect(userLoader.LoadCall.Receives.Token).To(Equal(""))
			})

			It("should mark the status as undeliverable", func() {
				err := processor.Process(delivery, logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(common.StatusUndeliverable))
			})
		})
	})

	Context("when the delivery has both an email and a userGUID", func() {
		It("should be marked as undeliverable", func() {
			delivery.Email = "some-email@example.com"
			delivery.UserGUID = "some-user-guid"

			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(common.StatusUndeliverable))
		})
	})

	Context("when the user email address is malformed", func() {
		It("marks the message status as undeliverable", func() {
			userLoader.LoadCall.Returns.Users = map[string]uaa.User{
				"user-123": {
					Emails: []string{"something"},
				},
			}

			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(common.StatusUndeliverable))
		})
	})

	Context("when the user is unsubscribed from the campaign type", func() {
		BeforeEach(func() {
			unsubscribersRepository.GetCall.Returns.Unsubscriber = models.Unsubscriber{
				ID:             "some-id",
				CampaignTypeID: "some-campaign-type-id",
				UserGUID:       "user-123",
			}
			campaignsRepository.GetCall.Returns.Campaign = models.Campaign{
				CampaignTypeID: "some-campaign-type-id",
			}

			err := processor.Process(delivery, logger)
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not send the notification", func() {
			Expect(campaignsRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(campaignsRepository.GetCall.Receives.CampaignID).To(Equal("some-campaign-id"))

			Expect(unsubscribersRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(unsubscribersRepository.GetCall.Receives.UserGUID).To(Equal("user-123"))
			Expect(unsubscribersRepository.GetCall.Receives.CampaignTypeID).To(Equal("some-campaign-type-id"))

			Expect(mailClient.SendCall.CallCount).To(Equal(0))
		})

		It("marks the message as delivered", func() {
			Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal("randomly-generated-guid"))
			Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(common.StatusDelivered))
			Expect(messageStatusUpdater.UpdateCall.Receives.CampaignID).To(Equal("some-campaign-id"))
			Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
		})
	})

	Context("failure cases", func() {
		Context("when the campaigns repository has an error", func() {
			It("returns the error", func() {
				campaignsRepository.GetCall.Returns.Error = errors.New("some-campaigns-repository-error")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("some-campaigns-repository-error")))
			})
		})

		Context("when the unsubscriber has an unknown error", func() {
			It("returns the error", func() {
				unsubscribersRepository.GetCall.Returns.Error = errors.New("some-unsubscriber-error")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("some-unsubscriber-error")))
			})
		})

		Context("when the token cannot be loaded", func() {
			It("returns the error", func() {
				tokenLoader.LoadCall.Returns.Error = errors.New("some-token-error")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("some-token-error")))
			})
		})

		Context("when the user cannot be loaded", func() {
			It("returns the error", func() {
				userLoader.LoadCall.Returns.Error = errors.New("something happened")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("something happened")))
			})
		})

		Context("when the packager fails to prepare the context", func() {
			It("returns the error", func() {
				packager.PrepareContextCall.Returns.Error = errors.New("some-packaging-error")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("some-packaging-error")))
			})
		})

		Context("when the packager fails to pack the message", func() {
			It("returns the error", func() {
				packager.PackCall.Returns.Error = errors.New("some-packaging-error")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("some-packaging-error")))
			})
		})

		Context("when the mail client fails to send the message", func() {
			It("returns the error", func() {
				mailClient.SendCall.Returns.Error = errors.New("smtp error")

				err := processor.Process(delivery, logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("smtp error")))
			})
		})
	})
})
