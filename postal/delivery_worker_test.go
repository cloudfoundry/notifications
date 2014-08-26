package postal_test

import (
    "bytes"
    "log"
    "os"
    "strings"
    "time"

    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("DeliveryWorker", func() {
    var mailClient FakeMailClient
    var worker postal.DeliveryWorker
    var id int
    var logger *log.Logger
    var buffer *bytes.Buffer
    var delivery postal.Delivery
    var queue *FakeQueue

    BeforeEach(func() {
        buffer = bytes.NewBuffer([]byte{})
        id = 1234
        logger = log.New(buffer, "", 0)
        mailClient = FakeMailClient{}
        queue = NewFakeQueue()

        worker = postal.NewDeliveryWorker(id, logger, &mailClient, queue)

        os.Setenv("SENDER", "from@email.com")

        delivery = postal.Delivery{
            User: uaa.User{
                Emails: []string{"fake-user@example.com"},
            },
            ClientID: "some-client",
            UserGUID: "user-123",
            Options: postal.Options{
                Subject: "the subject",
                Text:    "body content",
                ReplyTo: "thesender@example.com",
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

            Expect(len(mailClient.messages)).To(Equal(2))
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

            Expect(mailClient.messages).To(ContainElement(mail.Message{
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
                mailClient.errorOnSend = true

                worker.Deliver(&job)

                Expect(len(mailClient.messages)).To(Equal(0))
                Expect(job.ShouldRetry).To(BeTrue())
            })

            It("sets the retry duration using an exponential backoff algorithm", func() {
                mailClient.errorOnConnect = true

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
    })
})
