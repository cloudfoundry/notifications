package mail_test

import (
    "os"
    "strings"
    "time"

    "github.com/cloudfoundry-incubator/notifications/mail"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Mail", func() {
    var mailServer *SMTPServer
    var client *mail.Client

    Context("NewClient", func() {
        It("defaults the ConnectTimeout to 15 seconds", func() {
            client, err := mail.NewClient("user", "pass", "0.0.0.0:3000")
            if err != nil {
                panic(err)
            }

            Expect(client.ConnectTimeout).To(Equal(15 * time.Second))
        })
    })

    Describe("Send", func() {
        Context("when in Testmode", func() {
            BeforeEach(func() {
                var err error

                mailServer = NewSMTPServer("user", "pass")
                mailServer.SupportsTLS = true
                serverURL := mailServer.URL.String()
                client, err = mail.NewClient("user", "pass", serverURL)
                if err != nil {
                    panic(err)
                }

                os.Setenv("TEST_MODE", "true")
            })

            AfterEach(func() {
                os.Setenv("TEST_MODE", "false")
            })

            It("does not connect to the smtp server", func() {
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
                }).Should(Equal(0))
            })
        })
    })

    Context("when the SMTP server is properly configured", func() {
        BeforeEach(func() {
            var err error

            mailServer = NewSMTPServer("user", "pass")
            mailServer.SupportsTLS = true
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
            Expect(delivery.UsedTLS).To(BeFalse())
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

        Context("when configured to use TLS", func() {
            var smtpTLS string

            BeforeEach(func() {
                smtpTLS = os.Getenv("SMTP_TLS")
                os.Setenv("SMTP_TLS", "true")
                client.Insecure = true
            })

            AfterEach(func() {
                os.Setenv("SMTP_TLS", smtpTLS)
            })

            It("communicates over TLS", func() {

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
                Expect(delivery.UsedTLS).To(BeTrue())
            })
        })
    })

    Describe("Connect", func() {
        Context("when in test mode", func() {
            BeforeEach(func() {
                os.Setenv("TEST_MODE", "TRUE")
            })

            AfterEach(func() {
                os.Setenv("TEST_MODE", "FALSE")
            })

            It("does not connect to the smtp server", func() {

                serverURL := "fakewebsiteoninternet.com:587"
                client, err := mail.NewClient("user", "pass", serverURL)
                if err != nil {
                    panic(err)
                }
                err = client.Connect()

                Expect(err).To(BeNil())
            })
        })

        It("returns an error if it cannot connect within the given timeout duration", func() {
            mailServer = NewSMTPServer("user", "pass")
            mailServer.ConnectWait = 5 * time.Second

            serverURL := mailServer.URL.String()
            client, err := mail.NewClient("user", "pass", serverURL)
            if err != nil {
                panic(err)
            }

            client.ConnectTimeout = 100 * time.Millisecond
            err = client.Connect()

            Expect(err).ToNot(BeNil())
            Expect(err.Error()).To(ContainSubstring("timeout"))
        })
    })

    Context("Extension", func() {
        BeforeEach(func() {
            var err error

            mailServer = NewSMTPServer("user", "pass")
            serverURL := mailServer.URL.String()
            client, err = mail.NewClient("user", "pass", serverURL)
            if err != nil {
                panic(err)
            }
        })

        It("returns a bool, representing presence of, and parameters for a given SMTP extension", func() {
            err := client.Connect()
            if err != nil {
                panic(err)
            }

            err = client.Hello()
            if err != nil {
                panic(err)
            }

            ok, params := client.Extension("AUTH")
            Expect(ok).To(BeTrue())
            Expect(params).To(Equal("PLAIN LOGIN"))

            ok, params = client.Extension("STARTTLS")
            Expect(ok).To(BeFalse())
            Expect(params).To(Equal(""))
        })
    })
})
