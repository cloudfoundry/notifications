package mail_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/mail"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Mail", func() {
    var mailServer *SMTPServer
    var client mail.Client

    Context("when the SMTP server is properly configured", func() {
        BeforeEach(func() {
            var err error

            mailServer = NewSMTPServer("user", "pass")
            serverURL := mailServer.URL.String()
            client, err = mail.NewClient("user", "pass", serverURL)
            if err != nil {
                panic(err)
            }
        })

        It("can send mail", func() {
            msg := mail.Message{
                From:    "me@example.com",
                To:      "you@example.com",
                Subject: "Urgent! Read now!",
                Body:    "This email is the most important thing you will read all day!",
            }

            err := client.Send(msg)
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

        It("can make multiple requests", func() {
            firstMsg := mail.Message{
                From:    "me@example.com",
                To:      "you@example.com",
                Subject: "Urgent! Read now!",
                Body:    "This email is the most important thing you will read all day!",
            }

            err := client.Send(firstMsg)
            if err != nil {
                panic(err)
            }

            Eventually(func() int {
                return len(mailServer.Deliveries)
            }).Should(Equal(1))
            delivery := mailServer.Deliveries[0]

            Expect(delivery.Sender).To(Equal("me@example.com"))
            Expect(delivery.Recipient).To(Equal("you@example.com"))
            Expect(delivery.Data).To(Equal(strings.Split(firstMsg.Data(), "\n")))

            secondMsg := mail.Message{
                From:    "first@example.com",
                To:      "second@example.com",
                Subject: "Boring. Do not read.",
                Body:    "This email is the least important thing you will read all day. Sorry.",
            }

            err = client.Send(secondMsg)
            if err != nil {
                panic(err)
            }

            Eventually(func() int {
                return len(mailServer.Deliveries)
            }).Should(Equal(2))
            delivery = mailServer.Deliveries[1]

            Expect(delivery.Sender).To(Equal("first@example.com"))
            Expect(delivery.Recipient).To(Equal("second@example.com"))
            Expect(delivery.Data).To(Equal(strings.Split(secondMsg.Data(), "\n")))

        })
    })
})
