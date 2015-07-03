package postal_test

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/pivotal-golang/lager"

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
		mailClient             *fakes.MailClient
		worker                 postal.DeliveryWorker
		id                     int
		logger                 lager.Logger
		buffer                 *bytes.Buffer
		delivery               postal.Delivery
		queue                  *fakes.Queue
		unsubscribesRepo       *fakes.UnsubscribesRepo
		globalUnsubscribesRepo *fakes.GlobalUnsubscribesRepo
		kindsRepo              *fakes.KindsRepo
		messagesRepo           *fakes.MessagesRepo
		database               *fakes.Database
		conn                   models.ConnectionInterface
		userLoader             *fakes.UserLoader
		userGUID               string
		fakeUserEmail          string
		templateLoader         *fakes.TemplatesLoader
		receiptsRepo           *fakes.ReceiptsRepo
		tokenLoader            *fakes.TokenLoader
		zonedTokenLoader       *fakes.ZonedTokenLoader
		messageID              string
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		id = 1234
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		mailClient = fakes.NewMailClient()
		queue = fakes.NewQueue()
		unsubscribesRepo = fakes.NewUnsubscribesRepo()
		globalUnsubscribesRepo = fakes.NewGlobalUnsubscribesRepo()
		kindsRepo = fakes.NewKindsRepo()
		messagesRepo = fakes.NewMessagesRepo()
		database = fakes.NewDatabase()
		conn = database.Connection()
		userGUID = "user-123"
		sender := "from@email.com"
		sum := md5.Sum([]byte("banana's are so very tasty"))
		encryptionKey := sum[:]
		fakeUserEmail = "user-123@example.com"
		userLoader = fakes.NewUserLoader()
		userLoader.Users["user-123"] = uaa.User{Emails: []string{fakeUserEmail}}
		userLoader.Users["user-456"] = uaa.User{Emails: []string{"user-456@example.com"}}
		tokenLoader = fakes.NewTokenLoader()
		zonedTokenLoader = fakes.NewZonedTokenLoader()
		templateLoader = fakes.NewTemplatesLoader()
		templateLoader.Templates = postal.Templates{
			Text:    "{{.Text}}",
			HTML:    "<p>{{.HTML}}</p>",
			Subject: "{{.Subject}}",
		}
		receiptsRepo = fakes.NewReceiptsRepo()

		worker = postal.NewDeliveryWorker(id, logger, mailClient, queue, globalUnsubscribesRepo, unsubscribesRepo, kindsRepo,
			messagesRepo, database, false, sender, encryptionKey, "canonical-uaa-host", userLoader, templateLoader, receiptsRepo, tokenLoader, zonedTokenLoader)

		messageID = "randomly-generated-guid"
		delivery = postal.Delivery{
			ClientID: "some-client",
			UserGUID: userGUID,
			UAAHost:  "canonical-uaa-host",
			Options: postal.Options{
				Subject: "the subject",
				Text:    "body content",
				ReplyTo: "thesender@example.com",
				KindID:  "some-kind",
			},
			MessageID:     messageID,
			VCAPRequestID: "some-request-id",
		}

		_, err := messagesRepo.Upsert(database.Connection(), models.Message{ID: messageID, Status: postal.StatusQueued})
		Expect(err).NotTo(HaveOccurred())
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
			delivery.UAAHost = "zoned-uaa-host"

			job := gobble.NewJob(delivery)
			worker.Deliver(&job)

			Expect(zonedTokenLoader.LoadArgument).To(Equal("zoned-uaa-host"))
			Expect(tokenLoader.LoadWasCalled).To(BeFalse())
		})
	})

	Describe("Deliver", func() {
		var job gobble.Job

		BeforeEach(func() {
			job = gobble.NewJob(delivery)
		})

		It("logs the email address of the recipient", func() {
			worker.Deliver(&job)

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

		It("logs successful delivery", func() {
			worker.Deliver(&job)

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
			worker = postal.NewDeliveryWorker(id, logger, mailClient, queue, globalUnsubscribesRepo, unsubscribesRepo, kindsRepo,
				messagesRepo, database, true, "from@email.com", encryptionKey, "canonical-uaa-host", userLoader, templateLoader, receiptsRepo, tokenLoader, zonedTokenLoader)
			worker.Deliver(&job)
			database.TraceLogger.Printf("some statement")

			Expect(database.TracePrefix).To(BeEmpty())
			lines, err := parseLogLines(buffer.Bytes())
			Expect(err).NotTo(HaveOccurred())

			Expect(lines).To(ContainElement(logLine{
				Source:   "notifications",
				Message:  "notifications.worker.db",
				LogLevel: int(lager.INFO),
				Data: map[string]interface{}{
					"session":         "2",
					"statement":       "some statement",
					"worker_id":       float64(1234),
					"message_id":      "randomly-generated-guid",
					"vcap_request_id": "some-request-id",
				},
			}))
		})

		It("does not log database operations when database traces are disabled", func() {
			worker.Deliver(&job)
			Expect(database.TraceLogger).To(BeNil())
			Expect(database.TracePrefix).To(BeEmpty())
		})

		It("upserts the StatusDelivered to the database", func() {
			worker.Deliver(&job)

			message, err := messagesRepo.FindByID(conn, messageID)
			if err != nil {
				panic(err)
			}

			Expect(message.Status).To(Equal(postal.StatusDelivered))
		})

		It("creates a reciept for the delivery", func() {
			worker.Deliver(&job)

			Expect(receiptsRepo.ClientID).To(Equal("some-client"))
			Expect(receiptsRepo.KindID).To(Equal("some-kind"))
			Expect(receiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-123"}))
		})

		Context("when the receipt fails to be created", func() {
			It("retries the job", func() {
				receiptsRepo.CreateReceiptsError = true
				worker.Deliver(&job)

				Expect(job.RetryCount).To(Equal(1))
			})
		})

		Context("when loading a zoned token fails", func() {
			It("retries the job", func() {
				delivery.UAAHost = "zoned-uaa-host"
				job = gobble.NewJob(delivery)

				zonedTokenLoader.LoadError = errors.New("failed to load a zoned UAA token")
				worker.Deliver(&job)

				Expect(job.RetryCount).To(Equal(1))
			})
		})

		Context("when the StatusDelivered failed to be upserted to the database", func() {
			It("logs the error", func() {
				messagesRepo.UpsertError = errors.New("An unforseen error in upserting to our db")
				worker.Deliver(&job)

				lines, err := parseLogLines(buffer.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(lines).To(ContainElement(logLine{
					Source:   "notifications",
					Message:  "notifications.worker.failed-message-status-upsert",
					LogLevel: int(lager.ERROR),
					Data: map[string]interface{}{
						"session":         "1",
						"error":           messagesRepo.UpsertError.Error(),
						"recipient":       "user-123@example.com",
						"worker_id":       float64(1234),
						"message_id":      "randomly-generated-guid",
						"status":          postal.StatusDelivered,
						"vcap_request_id": "some-request-id",
					},
				}))
			})
		})

		It("ensures message delivery", func() {
			worker.Deliver(&job)

			Expect(mailClient.Messages).To(HaveLen(1))
			msg := mailClient.Messages[0]
			Expect(msg.From).To(Equal("from@email.com"))
			Expect(msg.ReplyTo).To(Equal("thesender@example.com"))
			Expect(msg.To).To(Equal(fakeUserEmail))
			Expect(msg.Subject).To(Equal("the subject"))
			Expect(msg.Body).To(ConsistOf([]mail.Part{
				{
					ContentType: "text/plain",
					Content:     "body content",
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
			worker.Deliver(&job)
			Expect(mailClient.ConnectLogger.SessionName()).To(Equal("notifications.worker"))
			Expect(mailClient.SendLogger.SessionName()).To(Equal("notifications.worker"))
		})

		Context("when the delivery fails to be sent", func() {
			Context("because of a send error", func() {
				BeforeEach(func() {
					mailClient.SendError = errors.New("Error sending message!!!")
				})

				It("marks the job for retry", func() {
					worker.Deliver(&job)
					Expect(len(mailClient.Messages)).To(Equal(0))
					Expect(job.ShouldRetry).To(BeTrue())
				})

				It("logs an SMTP send error", func() {
					worker.Deliver(&job)

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

				It("upserts the StatusFailed to the database", func() {
					worker.Deliver(&job)

					message, err := messagesRepo.FindByID(conn, messageID)
					if err != nil {
						panic(err)
					}

					Expect(message.Status).To(Equal(postal.StatusFailed))
				})

				Context("when the StatusFailed fails to be upserted into the db", func() {
					It("logs the failure", func() {
						messagesRepo.UpsertError = errors.New("An unforseen error in upserting to our db")
						worker.Deliver(&job)

						lines, err := parseLogLines(buffer.Bytes())
						Expect(err).NotTo(HaveOccurred())

						Expect(lines).To(ContainElement(logLine{
							Source:   "notifications",
							Message:  "notifications.worker.failed-message-status-upsert",
							LogLevel: int(lager.ERROR),
							Data: map[string]interface{}{
								"session":         "1",
								"error":           messagesRepo.UpsertError.Error(),
								"recipient":       "user-123@example.com",
								"worker_id":       float64(1234),
								"message_id":      "randomly-generated-guid",
								"status":          postal.StatusFailed,
								"vcap_request_id": "some-request-id",
							},
						}))
					})
				})
			})

			Context("and the error is a connect error", func() {
				It("logs an SMTP connection error", func() {
					mailClient.ConnectError = errors.New("server timeout")
					worker.Deliver(&job)

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

				It("upserts the StatusUnavailable to the database", func() {
					var jobDelivery postal.Delivery
					err := job.Unmarshal(&jobDelivery)
					if err != nil {
						panic(err)
					}

					mailClient.ConnectError = errors.New("BOOM!")
					messageID := jobDelivery.MessageID
					worker.Deliver(&job)

					message, err := messagesRepo.FindByID(conn, messageID)
					if err != nil {
						panic(err)
					}

					Expect(message.Status).To(Equal(postal.StatusUnavailable))
				})

				It("sets the retry duration using an exponential backoff algorithm", func() {
					mailClient.ConnectError = errors.New("BOOM!")
					worker.Deliver(&job)

					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(1*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(1))

					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					line := lines[1]
					Expect(line.Source).To(Equal("notifications"))
					Expect(line.Message).To(Equal("notifications.worker.delivery-failed-retrying"))
					Expect(line.LogLevel).To(Equal(int(lager.INFO)))
					Expect(line.Data).To(HaveKeyWithValue("session", "1"))
					Expect(line.Data).To(HaveKeyWithValue("recipient", "user-123@example.com"))
					Expect(line.Data).To(HaveKeyWithValue("worker_id", float64(1234)))
					Expect(line.Data).To(HaveKeyWithValue("message_id", "randomly-generated-guid"))
					Expect(line.Data).To(HaveKeyWithValue("retry_count", float64(1)))
					Expect(line.Data).To(HaveKeyWithValue("vcap_request_id", "some-request-id"))

					Expect(line.Data).To(HaveKey("active_at"))
					activeAt, err := time.Parse(time.RFC3339, line.Data["active_at"].(string))
					Expect(err).NotTo(HaveOccurred())
					Expect(activeAt).To(BeTemporally("~", time.Now().Add(1*time.Minute), 10*time.Second))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(2*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(2))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(4*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(3))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(8*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(4))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(16*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(5))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(32*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(6))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(64*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(7))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(128*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(8))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(256*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(9))

					worker.Deliver(&job)
					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(512*time.Minute), 10*time.Second))
					Expect(job.RetryCount).To(Equal(10))

					job.ShouldRetry = false
					worker.Deliver(&job)
					Expect(job.ShouldRetry).To(BeFalse())
				})
			})
		})

		Context("when recipient has globally unsubscribed", func() {
			BeforeEach(func() {
				err := globalUnsubscribesRepo.Set(conn, userGUID, true)
				if err != nil {
					panic(err)
				}
				worker.Deliver(&job)
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

			It("upserts the StatusUndeliverable to the database", func() {
				message, err := messagesRepo.FindByID(conn, messageID)
				if err != nil {
					panic(err)
				}

				Expect(message.Status).To(Equal(postal.StatusUndeliverable))
			})
		})

		Context("when the recipient hasn't unsubscribed, but doesn't have a valid email address", func() {
			Context("when the recipient has no emails", func() {
				BeforeEach(func() {
					delivery.Email = ""
					userLoader.Users["user-123"] = uaa.User{}
					job = gobble.NewJob(delivery)

					worker.Deliver(&job)
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

				It("upserts the StatusUndeliverable to the database", func() {
					message, err := messagesRepo.FindByID(conn, messageID)
					if err != nil {
						panic(err)
					}

					Expect(message.Status).To(Equal(postal.StatusUndeliverable))
				})
			})

			Context("when the recipient's first email address is missing an @ symbol", func() {
				BeforeEach(func() {
					delivery.Email = "nope"
					job = gobble.NewJob(delivery)

					worker.Deliver(&job)
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

				It("upserts the StatusUndeliverable to the database", func() {
					message, err := messagesRepo.FindByID(conn, messageID)
					if err != nil {
						panic(err)
					}

					Expect(message.Status).To(Equal(postal.StatusUndeliverable))
				})
			})
		})

		Context("when recipient has unsubscribed", func() {
			BeforeEach(func() {
				err := unsubscribesRepo.Set(conn, userGUID, "some-client", "some-kind", true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("logs that the user has unsubscribed from this notification", func() {
				worker.Deliver(&job)

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

			It("upserts the StatusUndeliverable to the database", func() {
				worker.Deliver(&job)
				message, err := messagesRepo.FindByID(conn, messageID)
				if err != nil {
					panic(err)
				}

				Expect(message.Status).To(Equal(postal.StatusUndeliverable))
			})

			Context("and the notification is not registered", func() {
				It("does not send the email", func() {
					worker.Deliver(&job)

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
					worker.Deliver(&job)

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
					worker.Deliver(&job)

					Expect(len(mailClient.Messages)).To(Equal(1))
				})
			})
		})

		Context("when the template contains syntax errors", func() {
			BeforeEach(func() {
				templateLoader.Templates = postal.Templates{
					Text:    "This message is a test of the endorsement broadcast system. \n\n {{.Text}} \n\n ==Endorsement== \n {{.Endorsement} \n ==End Endorsement==",
					HTML:    "<h3>This message is a test of the Endorsement Broadcast System</h3><p>{{.HTML}}</p><h3>Endorsement:</h3><p>{.Endorsement}</p>",
					Subject: "Endorsement Test: {{.Subject}}",
				}
				job = gobble.NewJob(delivery)
			})

			It("does not panic", func() {
				Expect(func() {
					worker.Deliver(&job)
				}).ToNot(Panic())
			})

			It("marks the job for retry later", func() {
				worker.Deliver(&job)
				Expect(job.RetryCount).To(Equal(1))
			})

			It("logs that the packer errored", func() {
				worker.Deliver(&job)

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

			It("upserts the StatusFailed to the database", func() {
				worker.Deliver(&job)

				message, err := messagesRepo.FindByID(conn, messageID)
				Expect(err).ToNot(HaveOccurred())
				Expect(message.Status).To(Equal(postal.StatusFailed))
			})

			Context("when the StatusFailed fails to be upserted into the db", func() {
				It("logs the failure", func() {
					messagesRepo.UpsertError = errors.New("An unforseen error in upserting to our db")
					worker.Deliver(&job)

					lines, err := parseLogLines(buffer.Bytes())
					Expect(err).NotTo(HaveOccurred())

					Expect(lines).To(ContainElement(logLine{
						Source:   "notifications",
						Message:  "notifications.worker.failed-message-status-upsert",
						LogLevel: int(lager.ERROR),
						Data: map[string]interface{}{
							"session":         "1",
							"error":           messagesRepo.UpsertError.Error(),
							"recipient":       "user-123@example.com",
							"worker_id":       float64(1234),
							"message_id":      "randomly-generated-guid",
							"status":          postal.StatusFailed,
							"vcap_request_id": "some-request-id",
						},
					}))
				})
			})
		})

		Context("when the job contains malformed JSON", func() {
			BeforeEach(func() {
				job.Payload = `{"Space":"my-space","Options":{"HTML":"<p>some text that just abruptly ends`
			})

			It("does not crash the process", func() {
				Expect(func() {
					worker.Deliver(&job)
				}).ToNot(Panic())
			})

			It("marks the job for retry later", func() {
				worker.Deliver(&job)
				Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(1*time.Minute), 10*time.Second))
				Expect(job.RetryCount).To(Equal(1))
			})
		})

		Context("when the message status fails to be upserted into the db", func() {
			It("logs the failure", func() {
				messagesRepo.UpsertError = errors.New("An unforseen error in upserting to our db")
				worker.Deliver(&job)

				lines, err := parseLogLines(buffer.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(lines).To(ContainElement(logLine{
					Source:   "notifications",
					Message:  "notifications.worker.failed-message-status-upsert",
					LogLevel: int(lager.ERROR),
					Data: map[string]interface{}{
						"session":         "1",
						"error":           messagesRepo.UpsertError.Error(),
						"recipient":       "user-123@example.com",
						"worker_id":       float64(1234),
						"message_id":      "randomly-generated-guid",
						"status":          postal.StatusDelivered,
						"vcap_request_id": "some-request-id",
					},
				}))
			})

			It("still delivers the message", func() {
				worker.Deliver(&job)

				Expect(mailClient.Messages).To(HaveLen(1))
				Expect(mailClient.Messages[0].To).To(Equal(fakeUserEmail))
			})
		})
	})
})
