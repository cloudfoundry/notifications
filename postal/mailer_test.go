package postal_test

import (
    "bytes"
    "log"

    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Mailer", func() {
    var mailClient FakeMailClient
    var mailer postal.Mailer
    var logger *log.Logger
    var buffer *bytes.Buffer
    var queue *FakeQueue
    var worker postal.DeliveryWorker

    BeforeEach(func() {
        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        mailClient = FakeMailClient{}
        queue = NewFakeQueue()
        mailer = postal.NewMailer(queue, FakeGuidGenerator)

        worker = postal.NewDeliveryWorker(1234, logger, &mailClient, queue)
        go worker.Work()
    })

    AfterEach(func() {
        worker.Halt()
    })

    Describe("Deliver", func() {
        It("returns the correct types of responses for users", func() {
            users := map[string]uaa.User{
                "user-1": {ID: "user-1", Emails: []string{"user-1@example.com"}},
                "user-2": {},
                "user-3": {ID: "user-3"},
                "user-4": {ID: "user-4", Emails: []string{"user-4"}},
            }
            responses := mailer.Deliver(postal.Templates{}, users, postal.Options{}, "the-space", "the-org", "the-client")

            Expect(len(responses)).To(Equal(4))
            Expect(responses).To(ContainElement(postal.Response{
                Status:         "queued",
                Recipient:      "user-1",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            Expect(responses).To(ContainElement(postal.Response{
                Status:         "queued",
                Recipient:      "user-2",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            Expect(responses).To(ContainElement(postal.Response{
                Status:         "queued",
                Recipient:      "user-3",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            Expect(responses).To(ContainElement(postal.Response{
                Status:         "queued",
                Recipient:      "user-4",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))
        })
    })
})
