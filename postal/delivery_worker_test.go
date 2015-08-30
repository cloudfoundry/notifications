package postal_test

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	v2Queue "github.com/cloudfoundry-incubator/notifications/v2/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type logLine struct {
	Source   string                 `json:"source"`
	Message  string                 `json:"message"`
	LogLevel int                    `json:"log_level"`
	Data     map[string]interface{} `json:"data"`
}

func parseLogLines(b []byte) ([]logLine, error) {
	var lines []logLine
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		var ll logLine
		err := json.Unmarshal(line, &ll)
		if err != nil {
			return lines, err
		}

		lines = append(lines, ll)
	}

	return lines, nil
}

var _ = Describe("DeliveryWorker", func() {
	var (
		mailClient             *mocks.MailClient
		worker                 postal.DeliveryWorker
		id                     int
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
		id = 1234
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		mailClient = mocks.NewMailClient()
		queue = mocks.NewQueue()
		unsubscribesRepo = mocks.NewUnsubscribesRepo()
		globalUnsubscribesRepo = mocks.NewGlobalUnsubscribesRepo()
		kindsRepo = mocks.NewKindsRepo()

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

		worker = postal.NewDeliveryWorker(postal.DeliveryWorkerConfig{
			ID:            id,
			Sender:        "from@example.com",
			Domain:        "example.com",
			UAAHost:       "https://uaa.example.com",
			EncryptionKey: encryptionKey,
			Logger:        logger,
			Queue:         queue,

			Database:               database,
			DBTrace:                false,
			GlobalUnsubscribesRepo: globalUnsubscribesRepo,
			UnsubscribesRepo:       unsubscribesRepo,
			KindsRepo:              kindsRepo,
			UserLoader:             userLoader,
			TemplatesLoader:        templateLoader,
			ReceiptsRepo:           receiptsRepo,
			TokenLoader:            tokenLoader,
			MailClient:             mailClient,
			StrategyDeterminer:     strategyDeterminer,
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

	Describe("Work", func() {
		It("pops Deliveries off the queue, sending emails for each", func() {
			queue.Enqueue(gobble.NewJob(delivery))

			worker.Work()

			delivery2 := postal.Delivery{
				UserGUID: "user-456",
			}
			queue.Enqueue(gobble.NewJob(delivery2))

			<-time.After(10 * time.Millisecond)
			worker.Halt()

			Expect(mailClient.Messages).To(HaveLen(2))
		})

		It("can be halted", func() {
			go func() {
				worker.Halt()
			}()

			Eventually(func() bool {
				worker.Work()
				return true
			}).Should(BeTrue())
		})
	})

	Describe("Deliver to zone", func() {
		It("makes a call to getNewClientToken for a zone during a delivery", func() {
			job := gobble.NewJob(delivery)
			worker.Deliver(&job)

			Expect(tokenLoader.LoadCall.Receives.UAAHost).To(Equal("https://uaa.example.com"))
		})
	})

	Describe("Deliver", func() {
		var job *gobble.Job

		BeforeEach(func() {
			j := gobble.NewJob(delivery)
			job = &j
		})

		Context("when Deliver receives a campaign", func() {
			BeforeEach(func() {
				campaignJob := v2Queue.CampaignJob{JobType: "campaign", Campaign: collections.Campaign{}}
				j := gobble.NewJob(campaignJob)
				job = &j
			})

			It("sends the job to the strategyDeterminer", func() {
				worker.Deliver(job)
				Expect(strategyDeterminer.DetermineCall.Receives.Job).To(Equal(*job))
				Expect(receiptsRepo.WasCalled).To(BeFalse())
			})

			It("logs that it is determining the strategy", func() {
				worker.Deliver(job)
				lines, err := parseLogLines(buffer.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(lines).To(ContainElement(logLine{
					Source:   "notifications",
					Message:  "notifications.worker.determining-strategy",
					LogLevel: int(lager.INFO),
					Data: map[string]interface{}{
						"session":   "1",
						"worker_id": float64(1234),
					},
				}))
			})

			It("does not retry the job", func() {
				worker.Deliver(job)

				Expect(deliveryFailureHandler.HandleCall.WasCalled).To(BeFalse())
			})

			Context("when the campaign cannot be determined", func() {
				BeforeEach(func() {
					strategyDeterminer.DetermineCall.Returns.Error = errors.New("looks like that worked out pretty well")
					campaignJob := v2Queue.CampaignJob{JobType: "campaign", Campaign: collections.Campaign{}}
					j := gobble.NewJob(campaignJob)
					job = &j

					worker.Deliver(job)
				})

				It("retries the job", func() {
					Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
					Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})

				It("logs the error", func() {
					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.determining-strategy-failed",
						LogLevel: int(lager.ERROR),
						Data: map[string]interface{}{
							"session":   "1",
							"worker_id": float64(1234),
							"error":     "looks like that worked out pretty well",
						},
					}))
				})
			})
		})

		Context("when Deliver does not receive a campaign", func() {
			It("does not send the job to the strategyDeterminer", func() {
				worker.Deliver(job)
				Expect(strategyDeterminer.DetermineCall.WasCalled).To(BeFalse())
			})
		})

		It("logs the email address of the recipient", func() {
			worker.Deliver(job)

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
			worker.Deliver(job)

			Expect(templateLoader.LoadTemplatesCall.Receives.ClientID).To(Equal("some-client"))
			Expect(templateLoader.LoadTemplatesCall.Receives.KindID).To(Equal("some-kind"))
			Expect(templateLoader.LoadTemplatesCall.Receives.TemplateID).To(Equal("some-template-id"))

		})

		It("logs successful delivery", func() {
			worker.Deliver(job)

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
			worker = postal.NewDeliveryWorker(postal.DeliveryWorkerConfig{
				ID:         id,
				Logger:     logger,
				MailClient: mailClient,
				Queue:      queue,
				GlobalUnsubscribesRepo: globalUnsubscribesRepo,
				UnsubscribesRepo:       unsubscribesRepo,
				KindsRepo:              kindsRepo,
				Database:               database,
				DBTrace:                true,
				Sender:                 "from@example.com",
				EncryptionKey:          encryptionKey,
				UserLoader:             userLoader,
				TemplatesLoader:        templateLoader,
				ReceiptsRepo:           receiptsRepo,
				TokenLoader:            tokenLoader,
				MessageStatusUpdater:   messageStatusUpdater,
			})
			worker.Deliver(job)

			Expect(database.TraceOnCall.Receives.Prefix).To(BeEmpty())
			Expect(database.TraceOnCall.Receives.Logger).NotTo(BeNil())
		})

		It("does not log database operations when database traces are disabled", func() {
			worker.Deliver(job)
			Expect(database.TraceOnCall.Receives.Prefix).To(BeEmpty())
			Expect(database.TraceOnCall.Receives.Logger).To(BeNil())
		})

		It("updates the message status as delivered", func() {
			worker.Deliver(job)

			Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
			Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusDelivered))
			Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
		})

		It("creates a reciept for the delivery", func() {
			worker.Deliver(job)

			Expect(receiptsRepo.ClientID).To(Equal("some-client"))
			Expect(receiptsRepo.KindID).To(Equal("some-kind"))
			Expect(receiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-123"}))
		})

		Context("when the receipt fails to be created", func() {
			It("retries the job", func() {
				receiptsRepo.CreateReceiptsError = true
				worker.Deliver(job)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})

		Context("when loading a zoned token fails", func() {
			It("retries the job", func() {
				j := gobble.NewJob(delivery)
				job = &j

				tokenLoader.LoadCall.Returns.Error = errors.New("failed to load a zoned UAA token")
				worker.Deliver(job)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})

		It("ensures message delivery", func() {
			worker.Deliver(job)

			Expect(mailClient.Messages).To(HaveLen(1))
			msg := mailClient.Messages[0]
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
			worker.Deliver(job)
			Expect(mailClient.ConnectLogger.SessionName()).To(Equal("notifications.worker"))
			Expect(mailClient.SendLogger.SessionName()).To(Equal("notifications.worker"))
		})

		Context("when the delivery fails to be sent", func() {
			Context("because of a send error", func() {
				BeforeEach(func() {
					mailClient.SendError = errors.New("Error sending message!!!")
				})

				It("marks the job for retry", func() {
					worker.Deliver(job)

					Expect(mailClient.Messages).To(HaveLen(0))
					Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
					Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})

				It("logs an SMTP send error", func() {
					worker.Deliver(job)

					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.delivery-failed-smtp-error",
						LogLevel: int(lager.ERROR),
						Data: map[string]interface{}{
							"session":         "1",
							"error":           mailClient.SendError.Error(),
							"recipient":       "user-123@example.com",
							"worker_id":       float64(1234),
							"message_id":      "randomly-generated-guid",
							"vcap_request_id": "some-request-id",
						},
					}))
				})

				It("updates the message status as failed", func() {
					worker.Deliver(job)

					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusFailed))
					Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})
			})

			Context("and the error is a connect error", func() {
				It("logs an SMTP connection error", func() {
					mailClient.ConnectError = errors.New("server timeout")
					worker.Deliver(job)

					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.smtp-connection-error",
						LogLevel: int(lager.ERROR),
						Data: map[string]interface{}{
							"session":         "1",
							"error":           mailClient.ConnectError.Error(),
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

					mailClient.ConnectError = errors.New("BOOM!")
					messageID := jobDelivery.MessageID
					worker.Deliver(job)

					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUnavailable))
					Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
				})
			})
		})

		Context("when recipient has globally unsubscribed", func() {
			BeforeEach(func() {
				err := globalUnsubscribesRepo.Set(conn, userGUID, true)
				if err != nil {
					panic(err)
				}
				worker.Deliver(job)
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
				Expect(mailClient.Messages).To(HaveLen(0))
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

					worker.Deliver(job)
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

					worker.Deliver(job)
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
				err := unsubscribesRepo.Set(conn, userGUID, "some-client", "some-kind", true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("logs that the user has unsubscribed from this notification", func() {
				worker.Deliver(job)

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
				worker.Deliver(job)

				Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(conn))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal(messageID))
				Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal(postal.StatusUndeliverable))
				Expect(messageStatusUpdater.UpdateCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})

			Context("and the notification is not registered", func() {
				It("does not send the email", func() {
					worker.Deliver(job)

					Expect(len(mailClient.Messages)).To(Equal(0))
				})
			})

			Context("and the notification is registered as not critical", func() {
				BeforeEach(func() {
					_, err := kindsRepo.Create(conn, models.Kind{
						ID:       "some-kind",
						ClientID: "some-client",
						Critical: false,
					})

					if err != nil {
						panic(err)
					}
				})
				It("does not send the email", func() {
					worker.Deliver(job)

					Expect(len(mailClient.Messages)).To(Equal(0))
				})
			})

			Context("and the notification is registered as critical", func() {
				BeforeEach(func() {
					_, err := kindsRepo.Create(conn, models.Kind{
						ID:       "some-kind",
						ClientID: "some-client",
						Critical: true,
					})

					if err != nil {
						panic(err)
					}
				})

				It("does send the email", func() {
					worker.Deliver(job)

					Expect(len(mailClient.Messages)).To(Equal(1))
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
					worker.Deliver(job)
				}).ToNot(Panic())
			})

			It("marks the job for retry later", func() {
				worker.Deliver(job)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})

			It("logs that the packer errored", func() {
				worker.Deliver(job)

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
				worker.Deliver(job)

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
					worker.Deliver(job)
				}).ToNot(Panic())
			})

			It("marks the job for retry later", func() {
				worker.Deliver(job)

				Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger.SessionName()).To(Equal("notifications.worker"))
			})
		})
	})
})
