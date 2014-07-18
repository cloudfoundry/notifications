package handlers_test

import (
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("MailSender", func() {

    var mailSender handlers.MailSender
    var context handlers.MessageContext
    var client mail.Client

    BeforeEach(func() {
        client = mail.Client{}
        context = handlers.MessageContext{
            From:      "banana man",
            To:        "endless monkeys",
            Subject:   "we will be eaten",
            ClientID:  "333",
            MessageID: "4444",
            Text:      "User supplied banana text",
            HTML:      "<p>user supplied banana html</p>",
            PlainTextEmailTemplate: "Banana preamble {{.Text}}",
            HTMLEmailTemplate:      "Banana preamble {{.HTML}}",
        }
        mailSender = handlers.NewMailSender(&client, context)
    })

    Describe("CompileBody", func() {
        It("returns the compiled email containing both the plaintext and html portions", func() {
            body, err := mailSender.CompileBody()
            if err != nil {
                panic(err)
            }

            emailBody := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-type: text/plain

Banana preamble User supplied banana text
--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        Banana preamble <p>user supplied banana html</p>
    </body>
</html>
--our-content-boundary--`

            Expect(body).To(Equal(emailBody))
        })

        Context("when no html is set", func() {
            It("only sends a plaintext of the email", func() {
                context.HTML = ""
                mailSender = handlers.NewMailSender(&client, context)

                body, err := mailSender.CompileBody()
                if err != nil {
                    panic(err)
                }

                emailBody := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-type: text/plain

Banana preamble User supplied banana text
--our-content-boundary--`
                Expect(body).To(Equal(emailBody))
            })
        })

        Context("when no text is set", func() {
            It("omits the plaintext portion of the email", func() {
                context.Text = ""
                mailSender = handlers.NewMailSender(&client, context)

                body, err := mailSender.CompileBody()
                if err != nil {
                    panic(err)
                }

                emailBody := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        Banana preamble <p>user supplied banana html</p>
    </body>
</html>
--our-content-boundary--`
                Expect(body).To(Equal(emailBody))
            })
        })
    })

    Describe("CompileMessage", func() {
        It("returns a mail message with all fields", func() {
            message := mailSender.CompileMessage("New Body")
            Expect(message.From).To(Equal("banana man"))
            Expect(message.To).To(Equal("endless monkeys"))
            Expect(message.Subject).To(Equal("CF Notification: we will be eaten"))
            Expect(message.Body).To(Equal("New Body"))
            Expect(message.Headers).To(Equal([]string{"X-CF-Client-ID: 333", "X-CF-Notification-ID: 4444"}))
        })
    })
})
