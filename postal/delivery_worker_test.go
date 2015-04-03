package postal_test

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getMessageIDFromJob(job gobble.Job) string {
	var jobDelivery postal.Delivery
	err := job.Unmarshal(&jobDelivery)
	if err != nil {
		panic(err)
	}
	return jobDelivery.MessageID
}

var _ = Describe("DeliveryWorker", func() {
	var (
		mailClient             fakes.MailClient
		worker                 postal.DeliveryWorker
		id                     int
		logger                 *log.Logger
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
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		id = 1234
		logger = log.New(buffer, "", 0)
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
		templateLoader = fakes.NewTemplatesLoader()
		templateLoader.Templates = postal.Templates{
			Text:    "{{.Text}}",
			HTML:    "<p>{{.HTML}}</p>",
			Subject: "{{.Subject}}",
		}
		receiptsRepo = fakes.NewReceiptsRepo()

		worker = postal.NewDeliveryWorker(id, logger, &mailClient, queue, globalUnsubscribesRepo, unsubscribesRepo, kindsRepo,
			messagesRepo, database, sender, encryptionKey, userLoader, templateLoader, receiptsRepo, tokenLoader)

		delivery = postal.Delivery{
			ClientID: "some-client",
			UserGUID: userGUID,
			Options: postal.Options{
				Subject: "the subject",
				Text:    "body content",
				ReplyTo: "thesender@example.com",
				KindID:  "some-kind",
			},
			MessageID: "randomly-generated-guid",
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

	Describe("Deliver", func() {
		var job gobble.Job

		BeforeEach(func() {
			job = gobble.NewJob(delivery)
		})

		It("logs the email address of the recipient", func() {
			worker.Deliver(&job)
			Expect(buffer.String()).To(ContainSubstring("Worker 1234: Attempting to deliver message to user-123@example.com"))
		})

		It("logs successful delivery", func() {
			worker.Deliver(&job)

			results := strings.Split(buffer.String(), "\n")
			Expect(results).To(ContainElement("Worker 1234: Message was successfully sent to user-123@example.com"))
		})

		It("upserts the StatusDelivered to the database", func() {
			messageID := getMessageIDFromJob(job)
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

		It("makes a call to getNewClientToken during a delivery", func() {
			worker.Deliver(&job)

			Expect(tokenLoader.LoadWasCalled).To(BeTrue())
		})

		Context("when loading a token fails", func() {
			It("retries the job", func() {
				tokenLoader.LoadError = errors.New("failed to load a UAA token")
				worker.Deliver(&job)

				Expect(job.RetryCount).To(Equal(1))
			})
		})

		Context("when the StatusDelivered failed to be upserted to the database", func() {
			It("Logs the error", func() {
				messagesRepo.UpsertError = errors.New("An unforseen error in upserting to our db")
				worker.Deliver(&job)
				Expect(buffer.String()).To(ContainSubstring(
					fmt.Sprintf("Worker 1234: Failed to upsert status '%s' of notification %s. Error: %s",
						postal.StatusDelivered,
						getMessageIDFromJob(job),
						messagesRepo.UpsertError.Error(),
					)))
			})
		})

		It("ensures message delivery", func() {
			worker.Deliver(&job)

			Expect(mailClient.Messages).To(ContainElement(mail.Message{
				From:    "from@email.com",
				ReplyTo: "thesender@example.com",
				To:      fakeUserEmail,
				Subject: "the subject",
				Body: []mail.Part{
					{
						ContentType: "text/plain",
						Content:     "body content",
					},
				},
				Headers: []string{
					"X-CF-Client-ID: some-client",
					"X-CF-Notification-ID: randomly-generated-guid",
				},
			}))
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
					Expect(buffer.String()).To(ContainSubstring("Worker 1234: Failed to deliver message due to SMTP error: " + mailClient.SendError.Error()))
				})

				It("upserts the StatusFailed to the database", func() {
					messageID := getMessageIDFromJob(job)
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
						Expect(buffer.String()).To(ContainSubstring("Worker 1234: Failed to upsert status '%s' of notification %s. Error: %s",
							postal.StatusFailed,
							getMessageIDFromJob(job),
							messagesRepo.UpsertError.Error()))
					})
				})
			})

			Context("and the error is a connect error", func() {
				It("logs an SMTP timeout error", func() {
					mailClient.ConnectError = errors.New("server timeout")
					worker.Deliver(&job)
					Expect(buffer.String()).To(ContainSubstring("Worker 1234: Error Establishing SMTP Connection: server timeout"))
				})

				It("logs other connect errors", func() {
					mailClient.ConnectError = errors.New("BOOM!")
					worker.Deliver(&job)
					Expect(buffer.String()).ToNot(ContainSubstring("server timeout"))
					Expect(buffer.String()).To(ContainSubstring("Worker 1234: Error Establishing SMTP Connection: BOOM!"))
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
					layout := "Jan 2, 2006 at 3:04pm (MST)"
					retryString := fmt.Sprintf("Worker 1234: Message failed to send, retrying at: %s", job.ActiveAt.Format(layout))

					Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(1*time.Minute), 10*time.Second))
					Expect(buffer.String()).To(ContainSubstring(retryString))
					Expect(job.RetryCount).To(Equal(1))

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
				Expect(buffer.String()).To(ContainSubstring("Worker 1234: Not delivering because user-123@example.com has unsubscribed"))
			})

			It("does not send any non-critical notifications", func() {
				Expect(mailClient.Messages).To(HaveLen(0))
			})
		})

		Context("when the recipient hasn't unsubscribed, but doesn't have a valid email address", func() {
			Context("when the recipient has no emails", func() {
				It("logs the error", func() {
					delivery.Email = ""
					userLoader.Users["user-123"] = uaa.User{}
					job = gobble.NewJob(delivery)

					worker.Deliver(&job)

					Expect(buffer.String()).To(ContainSubstring("Worker 1234: Not delivering because recipient has no email addresses"))
				})
			})

			Context("when the recipient's first email address is missing an @ symbol", func() {
				It("logs the error", func() {
					delivery.Email = "nope"
					job = gobble.NewJob(delivery)

					worker.Deliver(&job)

					Expect(buffer.String()).To(ContainSubstring("Worker 1234: Not delivering because recipient's email address is invalid"))
				})
			})
		})

		Context("when recipient has unsubscribed", func() {
			BeforeEach(func() {
				_, err := unsubscribesRepo.Create(conn, models.Unsubscribe{
					UserID:   userGUID,
					ClientID: "some-client",
					KindID:   "some-kind",
				})
				if err != nil {
					panic(err)
				}
			})

			It("logs that the user has unsubscribed from this notification", func() {
				worker.Deliver(&job)
				Expect(buffer.String()).To(ContainSubstring("Worker 1234: Not delivering because user-123@example.com has unsubscribed"))
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
				Expect(buffer.String()).To(ContainSubstring("Worker 1234: Not delivering because template failed to pack"))
			})

			It("upserts the StatusFailed to the database", func() {
				worker.Deliver(&job)
				messageID := getMessageIDFromJob(job)

				message, err := messagesRepo.FindByID(conn, messageID)
				Expect(err).ToNot(HaveOccurred())
				Expect(message.Status).To(Equal(postal.StatusFailed))
			})

			Context("when the StatusFailed fails to be upserted into the db", func() {
				It("logs the failure", func() {
					messagesRepo.UpsertError = errors.New("An unforseen error in upserting to our db")
					worker.Deliver(&job)
					Expect(buffer.String()).To(ContainSubstring("Worker 1234: Failed to upsert status '%s' of notification %s. Error: %s",
						postal.StatusFailed,
						getMessageIDFromJob(job),
						messagesRepo.UpsertError.Error()))
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
				Expect(buffer.String()).To(ContainSubstring(
					fmt.Sprintf("Worker 1234: Failed to upsert status '%s' of notification %s. Error: %s",
						postal.StatusDelivered,
						getMessageIDFromJob(job),
						messagesRepo.UpsertError.Error(),
					)))
			})

			It("still delivers the message", func() {
				worker.Deliver(&job)

				Expect(mailClient.Messages).To(ContainElement(mail.Message{
					From:    "from@email.com",
					ReplyTo: "thesender@example.com",
					To:      fakeUserEmail,
					Subject: "the subject",
					Body: []mail.Part{
						{
							ContentType: "text/plain",
							Content:     "body content",
						},
					},
					Headers: []string{
						"X-CF-Client-ID: some-client",
						"X-CF-Notification-ID: randomly-generated-guid",
					},
				}))
			})
		})
	})
})
