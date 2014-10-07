package postal_test

import (
    "bytes"
    "log"
    "os"
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
    var conn *fakes.FakeDBConn
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
        conn = &fakes.FakeDBConn{}
        userGUID = "user-123"

        worker = postal.NewDeliveryWorker(id, logger, &mailClient, queue, globalUnsubscribesRepo, unsubscribesRepo, kindsRepo)

        os.Setenv("SENDER", "from@email.com")

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

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
        })

        It("logs the message envelope", func() {
            worker.Deliver(&job)

            data := []string{
                "From: from@email.com",
                "To: fake-user@example.com",
                "Subject: the subject",
                `body content`,
            }
            results := strings.Split(buffer.String(), "\n")
            for _, item := range data {
                Expect(results).To(ContainElement(item))
            }
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
                mailClient.ErrorOnSend = true

                worker.Deliver(&job)

                Expect(len(mailClient.Messages)).To(Equal(0))
                Expect(job.ShouldRetry).To(BeTrue())
            })

            It("sets the retry duration using an exponential backoff algorithm", func() {
                mailClient.ErrorOnConnect = true

                worker.Deliver(&job)
                Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(1*time.Minute), 10*time.Second))
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

        Context("when recipient has globally unsubscribed", func() {
            BeforeEach(func() {
                err := globalUnsubscribesRepo.Set(conn, userGUID, true)
                if err != nil {
                    panic(err)
                }
            })

            It("does not send any non-critical notifications", func() {
                worker.Deliver(&job)
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
