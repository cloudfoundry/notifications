package postal_test

import (
    "bytes"
    "log"
    "os"
    "strings"

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

        worker = postal.NewDeliveryWorker(FakeGuidGenerator, logger, &mailClient, queue)

        os.Setenv("SENDER", "from@email.com")

        delivery = postal.Delivery{
            User: uaa.User{
                Emails: []string{"fake-user@example.com"},
            },
            UserGUID: "user-123",
            Options: postal.Options{
                Subject: "the subject",
                Text:    "body content",
            },
            Templates: postal.Templates{
                Text:    "{{.Text}}",
                Subject: "{{.Subject}}",
            },
            Response: make(chan postal.Response),
        }
    })

    Describe("Work", func() {
        It("pops Deliveries off the queue, passing them to Deliver, and sending their responses on the Delivery.Response chan", func() {
            queue.Enqueue(delivery)

            worker.Run()

            Expect(<-delivery.Response).To(Equal(postal.Response{
                Status:         "delivered",
                Recipient:      "user-123",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            delivery2 := postal.Delivery{
                User: uaa.User{
                    Emails: []string{"fake-user@example.com"},
                },
                UserGUID: "user-456",
                Response: make(chan postal.Response),
            }
            queue.Enqueue(delivery2)

            Expect(<-delivery2.Response).To(Equal(postal.Response{
                Status:         "delivered",
                Recipient:      "user-456",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            worker.Halt()
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
        It("logs the email address of the recipient and returns the response object", func() {
            response := worker.Deliver(delivery)

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
            Expect(response).To(Equal(postal.Response{
                Status:         "delivered",
                Recipient:      "user-123",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))
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
    })
})
