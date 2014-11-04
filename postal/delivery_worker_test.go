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

var _ = Describe("DeliveryWorker", func() {
    var mailClient fakes.FakeMailClient
    var worker postal.DeliveryWorker
    var id int
    var logger *log.Logger
    var buffer *bytes.Buffer
    var delivery postal.Delivery
    var queue *fakes.FakeQueue
    var unsubscribesRepo *fakes.FakeUnsubscribesRepo
    var globalUnsubscribesRepo *fakes.GlobalUnsubscribesRepo
    var kindsRepo *fakes.FakeKindsRepo
    var database *fakes.Database
    var conn models.ConnectionInterface
    var userGUID string

    BeforeEach(func() {
        buffer = bytes.NewBuffer([]byte{})
        id = 1234
        logger = log.New(buffer, "", 0)
        mailClient = fakes.FakeMailClient{}
        queue = fakes.NewFakeQueue()
        unsubscribesRepo = fakes.NewFakeUnsubscribesRepo()
        globalUnsubscribesRepo = fakes.NewGlobalUnsubscribesRepo()
        kindsRepo = fakes.NewFakeKindsRepo()
        database = fakes.NewDatabase()
        conn = database.Connection()
        userGUID = "user-123"
        sender := "from@email.com"
        sum := md5.Sum([]byte("banana's are so very tasty"))
        encryptionKey := string(sum[:])

        worker = postal.NewDeliveryWorker(id, logger, &mailClient, queue, globalUnsubscribesRepo, unsubscribesRepo, kindsRepo, database, sender, encryptionKey)

        delivery = postal.Delivery{
            User: uaa.User{
                Emails: []string{"fake-user@example.com"},
            },
            ClientID: "some-client",
            UserGUID: userGUID,
            Options: postal.Options{
                Subject: "the subject",
                Text:    "body content",
                ReplyTo: "thesender@example.com",
                KindID:  "some-kind",
            },
            Templates: postal.Templates{
                Text:    "{{.Text}}",
                Subject: "{{.Subject}}",
            },
            MessageID: "randomly-generated-guid",
        }
    })

    Describe("Work", func() {
        It("pops Deliveries off the queue, sending emails for each", func() {
            queue.Enqueue(gobble.NewJob(delivery))

            worker.Work()

            delivery2 := postal.Delivery{
                User: uaa.User{
                    Emails: []string{"fake-user@example.com"},
                },
                UserGUID: "user-456",
            }
            queue.Enqueue(gobble.NewJob(delivery2))

            <-time.After(10 * time.Millisecond)
            worker.Halt()

            Expect(len(mailClient.Messages)).To(Equal(2))
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
            Expect(buffer.String()).To(ContainSubstring("Attempting to deliver message to fake-user@example.com"))
        })

        It("logs successful delivery", func() {
            worker.Deliver(&job)

            results := strings.Split(buffer.String(), "\n")
            Expect(results).To(ContainElement("Message was successfully sent to fake-user@example.com"))
        })

        It("ensures message delivery", func() {
            worker.Deliver(&job)

            Expect(mailClient.Messages).To(ContainElement(mail.Message{
                From:    "from@email.com",
                ReplyTo: "thesender@example.com",
                To:      "fake-user@example.com",
                Subject: "the subject",
                Body:    "\nThis is a multi-part message in MIME format...\n\n--our-content-boundary\nContent-type: text/plain\n\nbody content\n--our-content-boundary--",
                Headers: []string{
                    "X-CF-Client-ID: some-client",
                    "X-CF-Notification-ID: randomly-generated-guid",
                },
            }))
        })

        Context("when the delivery fails to be sent", func() {
            It("marks the job for retry", func() {
                mailClient.SendError = errors.New("my awesome error")
                worker.Deliver(&job)
                Expect(len(mailClient.Messages)).To(Equal(0))
                Expect(job.ShouldRetry).To(BeTrue())
            })

            It("logs an SMTP send error", func() {
                mailClient.SendError = errors.New("BOOM!")
                worker.Deliver(&job)
                Expect(buffer.String()).To(ContainSubstring("Failed to deliver message due to SMTP error: BOOM!"))
            })

            Context("and the error is a connect error", func() {
                It("logs an SMTP timeout error", func() {
                    mailClient.ConnectError = errors.New("server timeout")
                    worker.Deliver(&job)
                    Expect(buffer.String()).To(ContainSubstring("Error Establishing SMTP Connection: server timeout"))
                })

                It("logs other connect errors", func() {
                    mailClient.ConnectError = errors.New("BOOM!")
                    worker.Deliver(&job)
                    Expect(buffer.String()).ToNot(ContainSubstring("server timeout"))
                    Expect(buffer.String()).To(ContainSubstring("Error Establishing SMTP Connection: BOOM!"))
                })

                It("sets the retry duration using an exponential backoff algorithm", func() {
                    mailClient.ConnectError = errors.New("BOOM!")
                    worker.Deliver(&job)
                    layout := "Jan 2, 2006 at 3:04pm (MST)"
                    retryString := fmt.Sprintf("Message failed to send, retrying at: %s", job.ActiveAt.Format(layout))

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
                Expect(buffer.String()).To(ContainSubstring("Not delivering because fake-user@example.com has unsubscribed"))
            })

            It("does not send any non-critical notifications", func() {
                Expect(len(mailClient.Messages)).To(Equal(0))
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
                Expect(buffer.String()).To(ContainSubstring("Not delivering because fake-user@example.com has unsubscribed"))
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
    })
})
