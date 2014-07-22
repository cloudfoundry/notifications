package handlers_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

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

    Describe("NewMessageContext", func() {
        var plainTextEmailTemplate string
        var htmlEmailTemplate string
        var user uaa.User
        var env config.Environment
        var params handlers.NotifyParams

        BeforeEach(func() {
            user = uaa.User{
                ID:     "user-456",
                Emails: []string{"bounce@example.com"},
            }

            env = config.NewEnvironment()

            plainTextEmailTemplate = "the plainText email template"
            htmlEmailTemplate = "the html email template"

            params = handlers.NotifyParams{
                Subject:           "the subject",
                KindDescription:   "the kind description",
                SourceDescription: "the source description",
                Text:              "user supplied email text",
                HTML:              "user supplied html",
                Kind:              "the-kind",
            }
        })

        It("returns the appropriate MessageContext when all params are specified", func() {
            messageContext := handlers.NewMessageContext(user, params, env, "the-space", "the-org",
                "the-client-ID", FakeGuidGenerator, plainTextEmailTemplate, htmlEmailTemplate)

            guid, err := FakeGuidGenerator()
            if err != nil {
                panic(err)
            }

            Expect(messageContext.From).To(Equal(os.Getenv("SENDER")))
            Expect(messageContext.To).To(Equal(user.Emails[0]))
            Expect(messageContext.Subject).To(Equal(params.Subject))
            Expect(messageContext.Text).To(Equal(params.Text))
            Expect(messageContext.HTML).To(Equal(params.HTML))
            Expect(messageContext.PlainTextEmailTemplate).To(Equal(plainTextEmailTemplate))
            Expect(messageContext.HTMLEmailTemplate).To(Equal(htmlEmailTemplate))
            Expect(messageContext.KindDescription).To(Equal(params.KindDescription))
            Expect(messageContext.SourceDescription).To(Equal(params.SourceDescription))
            Expect(messageContext.ClientID).To(Equal("the-client-ID"))
            Expect(messageContext.MessageID).To(Equal(guid.String()))
            Expect(messageContext.Space).To(Equal("the-space"))
            Expect(messageContext.Organization).To(Equal("the-org"))
        })

        It("falls back to Kind if KindDescription is missing", func() {
            params.KindDescription = ""

            messageContext := handlers.NewMessageContext(user, params, env, "the-space",
                "the-org", "the-client-ID", FakeGuidGenerator, plainTextEmailTemplate, htmlEmailTemplate)

            Expect(messageContext.KindDescription).To(Equal("the-kind"))
        })

        It("falls back to clientID when SourceDescription is missing", func() {
            params.SourceDescription = ""

            messageContext := handlers.NewMessageContext(user, params, env, "the-space",
                "the-org", "the-client-ID", FakeGuidGenerator, plainTextEmailTemplate, htmlEmailTemplate)

            Expect(messageContext.SourceDescription).To(Equal("the-client-ID"))
        })
    })
})
