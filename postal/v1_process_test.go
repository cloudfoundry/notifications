package postal_test

import (
	"bytes"
	"crypto/md5"
	"errors"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/conceal"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("V1Process", func() {
	var (
		mailClient             *mocks.MailClient
		v1Process              postal.V1Process
		logger                 lager.Logger
		buffer                 *bytes.Buffer
		delivery               postal.Delivery
		queue                  *mocks.Queue
		unsubscribesRepo       *mocks.UnsubscribesRepo
		globalUnsubscribesRepo *mocks.GlobalUnsubscribesRepo
		kindsRepo              *mocks.KindsRepo
		database               *mocks.Database
		strategyDeterminer     *mocks.StrategyDeterminer
		conn                   *mocks.Connection
		userLoader             *mocks.UserLoader
		userGUID               string
		fakeUserEmail          string
		templateLoader         *mocks.TemplatesLoader
		receiptsRepo           *mocks.ReceiptsRepo
		tokenLoader            *mocks.TokenLoader
		messageID              string
		messageStatusUpdater   *mocks.MessageStatusUpdater
		deliveryFailureHandler *mocks.DeliveryFailureHandler
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})

		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		logger = logger.Session("worker", lager.Data{"worker_id": 1234})

		mailClient = mocks.NewMailClient()
		queue = mocks.NewQueue()
		unsubscribesRepo = mocks.NewUnsubscribesRepo()
		globalUnsubscribesRepo = mocks.NewGlobalUnsubscribesRepo()

		kindsRepo = mocks.NewKindsRepo()
		kindsRepo.FindCall.Returns.Kinds = []models.Kind{
			{
				ID:       "some-kind",
				ClientID: "some-client",
				Critical: false,
			},
		}

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		strategyDeterminer = mocks.NewStrategyDeterminer()
		userGUID = "user-123"
		sum := md5.Sum([]byte("banana's are so very tasty"))
		encryptionKey := sum[:]
		fakeUserEmail = "user-123@example.com"
		userLoader = mocks.NewUserLoader()
		userLoader.LoadCall.Returns.Users = map[string]uaa.User{
			"user-123": {Emails: []string{fakeUserEmail}},
			"user-456": {Emails: []string{"user-456@example.com"}},
		}
		tokenLoader = mocks.NewTokenLoader()
		templateLoader = mocks.NewTemplatesLoader()
		templateLoader.LoadTemplatesCall.Returns.Templates = postal.Templates{
			Text:    "{{.Text}} {{.Domain}}",
			HTML:    "<p>{{.HTML}}</p>",
			Subject: "{{.Subject}}",
		}
		receiptsRepo = mocks.NewReceiptsRepo()
		messageStatusUpdater = mocks.NewMessageStatusUpdater()
		deliveryFailureHandler = mocks.NewDeliveryFailureHandler()

		cloak, err := conceal.NewCloak(encryptionKey)
		Expect(err).NotTo(HaveOccurred())

		v1Process = postal.NewV1Process(postal.V1ProcessConfig{
			DBTrace: false,
			UAAHost: "https://uaa.example.com",
			Sender:  "from@example.com",
			Domain:  "example.com",

			Packager:    postal.NewPackager(templateLoader, cloak),
			MailClient:  mailClient,
			Database:    database,
			TokenLoader: tokenLoader,
			UserLoader:  userLoader,

			KindsRepo:              kindsRepo,
			ReceiptsRepo:           receiptsRepo,
			UnsubscribesRepo:       unsubscribesRepo,
			GlobalUnsubscribesRepo: globalUnsubscribesRepo,
			MessageStatusUpdater:   messageStatusUpdater,
			DeliveryFailureHandler: deliveryFailureHandler,
		})

		messageID = "randomly-generated-guid"
		delivery = postal.Delivery{
			ClientID: "some-client",
			UserGUID: userGUID,
			Options: postal.Options{
				Subject:    "the subject",
				Text:       "body content",
				ReplyTo:    "thesender@example.com",
				KindID:     "some-kind",
				TemplateID: "some-template-id",
			},
			MessageID:     messageID,
			VCAPRequestID: "some-request-id",
		}
	})

	Describe("Deliver", func() {
		var job *gobble.Job

		BeforeEach(func() {
			j := gobble.NewJob(delivery)
			job = &j
		})

		It("logs the email address of the recipient", func() {
			v1Process.Deliver(job, logger)

			lines, err := parseLogLines(buffer.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(lines).To(ContainElement(logLine{
				Source:   "notifications",
				Message:  "notifications.worker.delivery-start",
				LogLevel: int(lager.INFO),
				Data: map[string]interface{}{
					"session":         "1",
					"recipient":       "user-123@example.com",
					"worker_id":       float64(1234),
					"message_id":      "randomly-generated-guid",
					"vcap_request_id": "some-request-id",
				},
			}))
		})

		It("loads the correct template", func() {
			v1Process.Deliver(job, logger)

			Expect(templateLoader.LoadTemplatesCall.Receives.ClientID).To(Equal("some-client"))
			Expect(templateLoader.LoadTemplatesCall.Receives.KindID).To(Equal("some-kind"))
			Expect(templateLoader.LoadTemplatesCall.Receives.TemplateID).To(Equal("some-template-id"))
		})

		It("logs successful delivery", func() {
			v1Process.Deliver(job, logger)

			lines, err := parseLogLines(buffer.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(lines).To(ContainElement(logLine{
				Source:   "notifications",
				Message:  "notifications.worker.message-sent",
				LogLevel: int(lager.INFO),
				Data: map[string]interface{}{
					"session":         "1",
					"recipient":       "user-123@example.com",
					"worker_id":       float64(1234),
					"message_id":      "randomly-generated-guid",
					"vcap_request_id": "some-request-id",
				},
			}))
		})

		It("logs database operations when database traces are enabled", func() {
			sum := md5.Sum([]byte("banana's are so very tasty"))
			encryptionKey := sum[:]
			cloak, err := conceal.NewCloak(encryptionKey)
			Expect(err).NotTo(HaveOccurred())
			v1Process = postal.NewV1Process(postal.V1ProcessConfig{
				DBTrace: true,
				UAAHost: "https://uaa.example.com",
				Sender:  "from@example.com",
				Domain:  "example.com",

				Packager:    postal.NewPackager(templateLoader, cloak),
				MailClient:  mailClient,
				Database:    database,
				TokenLoader: tokenLoader,
				UserLoader:  userLoader,

				KindsRepo:              kindsRepo,
				ReceiptsRepo:           receiptsRepo,
				UnsubscribesRepo:       unsubscribesRepo,
				GlobalUnsubscribesRepo: globalUnsubscribesRepo,
				MessageStatusUpdater:   messageStatusUpdater,
				DeliveryFailureHandler: deliveryFailureHandler,
			})
			v1Process.Deliver(job, logger)

			Expect(database.TraceOnCall.Receives.Prefix).To(BeEmpty())
			Expect(database.TraceOnCall.Receives.Logger).NotTo(BeNil())
		})

		It("does not log database operations when database traces are disabled", func() {
			v1Process.Deliver(job, logger)
			Expect(database.TraceOnCall.Receives.Prefix).To(BeEmpty())
			Expect(database.TraceOnCall.Receives.Logger).To(BeNil())
		})

		It("updates the message status as delivered", func() {
			v1Process.Deliver(job, logger)

			Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
			Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusDelivered))
			Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
		})

		It("creates a reciept for the delivery", func() {
			v1Process.Deliver(job, logger)

			Expect(receiptsRepo.CreateReceiptsCall.Receives.Connection).To(Equal(conn))
			Expect(receiptsRepo.CreateReceiptsCall.Receives.ClientID).To(Equal("some-client"))
			Expect(receiptsRepo.CreateReceiptsCall.Receives.KindID).To(Equal("some-kind"))
			Expect(receiptsRepo.CreateReceiptsCall.Receives.UserGUIDs).To(Equal([]string{"user-123"}))
		})

		Context("when the receipt fails to be created", func() {
			It("retries the job", func() {
				receiptsRepo.CreateReceiptsCall.Returns.Error = errors.New("something happened")
				v1Process.Deliver(job, logger)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})

		Context("when loading a zoned token fails", func() {
			It("retries the job", func() {
				j := gobble.NewJob(delivery)
				job = &j

				tokenLoader.LoadCall.Returns.Error = errors.New("failed to load a zoned UAA token")
				v1Process.Deliver(job, logger)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})

		It("ensures message delivery", func() {
			v1Process.Deliver(job, logger)

			Expect(mailClient.SendCall.CallCount).To(Equal(1))
			msg := mailClient.SendCall.Receives.Message
			Expect(msg.From).To(Equal("from@example.com"))
			Expect(msg.ReplyTo).To(Equal("thesender@example.com"))
			Expect(msg.To).To(Equal(fakeUserEmail))
			Expect(msg.Subject).To(Equal("the subject"))
			Expect(msg.Body).To(ConsistOf([]mail.Part{
				{
					ContentType: "text/plain",
					Content:     "body content example.com",
				},
			}))
			Expect(msg.Headers).To(ContainElement("X-CF-Client-ID: some-client"))
			Expect(msg.Headers).To(ContainElement("X-CF-Notification-ID: randomly-generated-guid"))

			var formattedTimestamp string
			prefix := "X-CF-Notification-Timestamp: "
			for _, header := range msg.Headers {
				if strings.Contains(header, prefix) {
					formattedTimestamp = strings.TrimPrefix(header, prefix)
					break
				}
			}
			Expect(formattedTimestamp).NotTo(BeEmpty())

			timestamp, err := time.Parse(time.RFC3339, formattedTimestamp)
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(BeTemporally("~", time.Now(), 2*time.Second))
		})

		It("should connect and send the message with the worker's logger session", func() {
			v1Process.Deliver(job, logger)
			Expect(mailClient.ConnectCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			Expect(mailClient.SendCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
		})

		Context("when the delivery fails to be sent", func() {
			Context("because of a send error", func() {
				BeforeEach(func() {
					mailClient.SendCall.Returns.Error = errors.New("Error sending message!!!")
				})

				It("marks the job for retry", func() {
					v1Process.Deliver(job, logger)

					Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
					Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})

				It("logs an SMTP send error", func() {
					v1Process.Deliver(job, logger)

					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.delivery-failed-smtp-error",
						LogLevel: int(lager.ERROR),
						Data: map[string]interface{}{
							"session":         "1",
							"error":           "Error sending message!!!",
							"recipient":       "user-123@example.com",
							"worker_id":       float64(1234),
							"message_id":      "randomly-generated-guid",
							"vcap_request_id": "some-request-id",
						},
					}))
				})

				It("updates the message status as failed", func() {
					v1Process.Deliver(job, logger)

					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusFailed))
					Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})
			})

			Context("and the error is a connect error", func() {
				It("logs an SMTP connection error", func() {
					mailClient.ConnectCall.Returns.Error = errors.New("server timeout")
					v1Process.Deliver(job, logger)

					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.smtp-connection-error",
						LogLevel: int(lager.ERROR),
						Data: map[string]interface{}{
							"session":         "1",
							"error":           "server timeout",
							"recipient":       "user-123@example.com",
							"worker_id":       float64(1234),
							"message_id":      "randomly-generated-guid",
							"vcap_request_id": "some-request-id",
						},
					}))
				})

				It("updates the message status as unavailable", func() {
					var jobDelivery postal.Delivery
					err := job.Unmarshal(&jobDelivery)
					if err != nil {
						panic(err)
					}

					mailClient.ConnectCall.Returns.Error = errors.New("BOOM!")
					messageID := jobDelivery.MessageID
					v1Process.Deliver(job, logger)

					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUnavailable))
					Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})
			})
		})

		Context("when recipient has globally unsubscribed", func() {
			BeforeEach(func() {
				globalUnsubscribesRepo.GetCall.Returns.Unsubscribed = true

				v1Process.Deliver(job, logger)
			})

			It("logs that the user has unsubscribed from this notification", func() {
				lines, err := parseLogLines(buffer.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(lines).To(ContainElement(logLine{
					Source:   "notifications",
					Message:  "notifications.worker.user-unsubscribed",
					LogLevel: int(lager.INFO),
					Data: map[string]interface{}{
						"session":         "1",
						"recipient":       "user-123@example.com",
						"worker_id":       float64(1234),
						"message_id":      "randomly-generated-guid",
						"vcap_request_id": "some-request-id",
					},
				}))
			})

			It("does not send any non-critical notifications", func() {
				Expect(mailClient.SendCall.CallCount).To(Equal(0))
			})

			It("updates the message status as undeliverable", func() {
				Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUndeliverable))
				Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})

		Context("when the recipient hasn't unsubscribed, but doesn't have a valid email address", func() {
			Context("when the recipient has no emails", func() {
				BeforeEach(func() {
					delivery.Email = ""
					userLoader.LoadCall.Returns.Users = map[string]uaa.User{
						"user-123": {},
					}
					j := gobble.NewJob(delivery)
					job = &j

					v1Process.Deliver(job, logger)
				})

				It("logs the info", func() {
					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.no-email-address-for-user",
						LogLevel: int(lager.INFO),
						Data: map[string]interface{}{
							"session":         "1",
							"recipient":       "",
							"worker_id":       float64(1234),
							"message_id":      "randomly-generated-guid",
							"vcap_request_id": "some-request-id",
						},
					}))
				})

				It("updates the message status as undeliverable", func() {
					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUndeliverable))
					Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})
			})

			Context("when the recipient's first email address is missing an @ symbol", func() {
				BeforeEach(func() {
					delivery.Email = "nope"
					j := gobble.NewJob(delivery)
					job = &j

					v1Process.Deliver(job, logger)
				})

				It("logs the info", func() {
					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.malformatted-email-address",
						LogLevel: int(lager.INFO),
						Data: map[string]interface{}{
							"session":         "1",
							"recipient":       "nope",
							"worker_id":       float64(1234),
							"message_id":      "randomly-generated-guid",
							"vcap_request_id": "some-request-id",
						},
					}))
				})

				It("updates the message status as undeliverable", func() {
					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUndeliverable))
					Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})
			})
		})

		Context("when recipient has unsubscribed", func() {
			BeforeEach(func() {
				unsubscribesRepo.GetCall.Returns.Unsubscribed = true
			})

			It("logs that the user has unsubscribed from this notification", func() {
				v1Process.Deliver(job, logger)

				lines, err := parseLogLines(buffer.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(lines).To(ContainElement(logLine{
					Source:   "notifications",
					Message:  "notifications.worker.user-unsubscribed",
					LogLevel: int(lager.INFO),
					Data: map[string]interface{}{
						"session":         "1",
						"recipient":       "user-123@example.com",
						"worker_id":       float64(1234),
						"message_id":      "randomly-generated-guid",
						"vcap_request_id": "some-request-id",
					},
				}))
			})

			It("updates the message status as undeliverable", func() {
				v1Process.Deliver(job, logger)

				Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUndeliverable))
				Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})

			Context("and the notification is not registered", func() {
				It("does not send the email", func() {
					v1Process.Deliver(job, logger)

					Expect(mailClient.SendCall.CallCount).To(Equal(0))
				})
			})

			Context("and the notification is registered as not critical", func() {
				BeforeEach(func() {
					kindsRepo.FindCall.Returns.Kinds = []models.Kind{
						{
							ID:       "some-kind",
							ClientID: "some-client",
							Critical: false,
						},
					}
				})

				It("does not send the email", func() {
					v1Process.Deliver(job, logger)

					Expect(mailClient.SendCall.CallCount).To(Equal(0))
				})
			})

			Context("and the notification is registered as critical", func() {
				BeforeEach(func() {
					kindsRepo.FindCall.Returns.Kinds = []models.Kind{
						{
							ID:       "some-kind",
							ClientID: "some-client",
							Critical: true,
						},
					}
				})

				It("does send the email", func() {
					v1Process.Deliver(job, logger)

					Expect(mailClient.SendCall.CallCount).To(Equal(1))
				})
			})
		})

		Context("when the template contains syntax errors", func() {
			BeforeEach(func() {
				templateLoader.LoadTemplatesCall.Returns.Templates = postal.Templates{
					Text:    "This message is a test of the endorsement broadcast system. \n\n {{.Text}} \n\n ==Endorsement== \n {{.Endorsement} \n ==End Endorsement==",
					HTML:    "<h3>This message is a test of the Endorsement Broadcast System</h3><p>{{.HTML}}</p><h3>Endorsement:</h3><p>{.Endorsement}</p>",
					Subject: "Endorsement Test: {{.Subject}}",
				}
				j := gobble.NewJob(delivery)
				job = &j
			})

			It("does not panic", func() {
				Expect(func() {
					v1Process.Deliver(job, logger)
				}).ToNot(Panic())
			})

			It("marks the job for retry later", func() {
				v1Process.Deliver(job, logger)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})

			It("logs that the packer errored", func() {
				v1Process.Deliver(job, logger)

				lines, err := parseLogLines(buffer.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(lines).To(ContainElement(logLine{
					Source:   "notifications",
					Message:  "notifications.worker.template-pack-failed",
					LogLevel: int(lager.INFO),
					Data: map[string]interface{}{
						"session":         "1",
						"recipient":       "user-123@example.com",
						"worker_id":       float64(1234),
						"message_id":      "randomly-generated-guid",
						"vcap_request_id": "some-request-id",
					},
				}))
			})

			It("updates the message status as failed", func() {
				v1Process.Deliver(job, logger)

				Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusFailed))
				Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})

		Context("when the job contains malformed JSON", func() {
			BeforeEach(func() {
				job.Payload = `{"Space":"my-space","Options":{"HTML":"<p>some text that just abruptly ends`
			})

			It("does not crash the process", func() {
				Expect(func() {
					v1Process.Deliver(job, logger)
				}).ToNot(Panic())
			})

			It("marks the job for retry later", func() {
				v1Process.Deliver(job, logger)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})
	})
})
