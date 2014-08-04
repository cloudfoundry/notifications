package postal_test

import (
    "bytes"
    "log"
    "strings"

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

    BeforeEach(func() {
        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        mailClient = FakeMailClient{}
        mailer = postal.NewMailer(FakeGuidGenerator, logger, &mailClient)
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
                Status:         "delivered",
                Recipient:      "user-1",
                NotificationID: "deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            Expect(responses).To(ContainElement(postal.Response{
                Status:         "notfound",
                Recipient:      "user-2",
                NotificationID: "",
            }))

            Expect(responses).To(ContainElement(postal.Response{
                Status:         "noaddress",
                Recipient:      "user-3",
                NotificationID: "",
            }))

            Expect(responses).To(ContainElement(postal.Response{
                Status:         "noaddress",
                Recipient:      "user-4",
                NotificationID: "",
            }))
        })
    })

    Describe("SendMailToUser", func() {
        It("logs the email address of the recipient and returns the status", func() {
            messageContext := postal.MessageContext{
                To: "fake-user@example.com",
            }

            mailClient = FakeMailClient{}

            status := mailer.SendMailToUser(messageContext, logger, &mailClient)

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
            Expect(status).To(Equal("delivered"))
        })

        It("logs the message envelope", func() {
            messageContext := postal.MessageContext{
                To:              "fake-user@example.com",
                From:            "from@email.com",
                Subject:         "the subject",
                Text:            "body content",
                KindDescription: "the kind description",
                TextTemplate:    "{{.Text}}",
                SubjectTemplate: "{{.Subject}}",
            }

            mailClient = FakeMailClient{}

            mailer.SendMailToUser(messageContext, logger, &mailClient)

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
