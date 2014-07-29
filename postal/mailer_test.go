package postal_test

import (
    "bytes"
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/postal"

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
