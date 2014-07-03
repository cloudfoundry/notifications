package mail_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/mail"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Mail", func() {
    var mailServer *SMTPServer

    Context("when the SMTP server is properly configured", func() {
        BeforeEach(func() {
            mailServer = NewSMTPServer("user", "pass")
        })

        It("can send mail", func() {
            serverURL := mailServer.URL.String()
            client, err := mail.NewClient("user", "pass", serverURL)
            if err != nil {
                panic(err)
            }

            msg := mail.Message{
                From:    "me@example.com",
                To:      "you@example.com",
                Subject: "Urgent! Read now!",
                Body:    "This email is the most important thing you will read all day!",
            }

            err = client.Send(msg)
            if err != nil {
                panic(err)
            }

            Eventually(func() int {
                return len(mailServer.Deliveries)
            }).Should(Equal(1))
            delivery := mailServer.Deliveries[0]

            Expect(delivery.Sender).To(Equal("me@example.com"))
            Expect(delivery.Recipient).To(Equal("you@example.com"))
            Expect(delivery.Data).To(Equal(strings.Split(msg.Data(), "\n")))
        })
    })
})
