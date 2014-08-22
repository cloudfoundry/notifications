package postal_test

import (
    "bytes"
    "log"
    "os"
    "strings"
    "time"

    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("DeliveryWorker", func() {
    var mailClient FakeMailClient
    var worker postal.DeliveryWorker
    var logger *log.Logger
    var buffer *bytes.Buffer
    var delivery postal.Delivery
    var queue *postal.DeliveryQueue

    BeforeEach(func() {
        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        mailClient = FakeMailClient{}
        queue = postal.NewDeliveryQueue()

        worker = postal.NewDeliveryWorker(logger, &mailClient, queue)

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
        It("pops Deliveries off the queue, passing them to Deliver, and sending their responses on the Delivery.Response chan", func() {
            queue.Enqueue(delivery)

            worker.Run()

            delivery2 := postal.Delivery{
                User: uaa.User{
                    Emails: []string{"fake-user@example.com"},
                },
                UserGUID: "user-456",
            }
            queue.Enqueue(delivery2)

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
        It("logs the email address of the recipient", func() {
            worker.Deliver(delivery)

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
        })

        It("logs the message envelope", func() {
            worker.Deliver(delivery)

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
            worker.Deliver(delivery)

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
    })
})
