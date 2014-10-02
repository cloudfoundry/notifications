package postal_test

import (
    "html"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {

    Describe("NewMessageContext", func() {
        var templates postal.Templates
        var email string
        var env config.Environment
        var options postal.Options
        var html postal.HTML
        var cryptographer *fakes.FakeCryptoClient

        BeforeEach(func() {
            email = "bounce@example.com"
            env = config.NewEnvironment()

            templates = postal.Templates{
                Text:    "the plainText email template",
                HTML:    "the html email template",
                Subject: "the subject template",
            }

            html = postal.HTML{
                BodyContent: "user supplied html",
            }

            options = postal.Options{
                ReplyTo:           "awesomeness",
                Subject:           "the subject",
                KindDescription:   "the kind description",
                SourceDescription: "the source description",
                Text:              "user supplied email text",
                HTML:              html,
                KindID:            "the-kind-id",
            }

            cryptographer = &fakes.FakeCryptoClient{
                EncryptedResult: "the-encoded-result",
            }
        })

        It("returns the appropriate MessageContext when all options are specified", func() {
            context := postal.NewMessageContext(email, options, env, "the-space", "the-org", "the-client-id",
                "message-id", "the-user", templates, cryptographer)

            Expect(context.From).To(Equal(env.Sender))
            Expect(context.ReplyTo).To(Equal(options.ReplyTo))
            Expect(context.To).To(Equal(email))
            Expect(context.Subject).To(Equal(options.Subject))
            Expect(context.Text).To(Equal(options.Text))
            Expect(context.HTML).To(Equal(options.HTML.BodyContent))
            Expect(context.HTMLComponents).To(Equal(options.HTML))
            Expect(context.TextTemplate).To(Equal(templates.Text))
            Expect(context.HTMLTemplate).To(Equal(templates.HTML))
            Expect(context.SubjectTemplate).To(Equal(templates.Subject))
            Expect(context.KindDescription).To(Equal(options.KindDescription))
            Expect(context.SourceDescription).To(Equal(options.SourceDescription))
            Expect(context.UserGUID).To(Equal("the-user"))
            Expect(context.ClientID).To(Equal("the-client-id"))
            Expect(context.MessageID).To(Equal("message-id"))
            Expect(context.Space).To(Equal("the-space"))
            Expect(context.Organization).To(Equal("the-org"))
            Expect(context.UnsubscribeID).To(Equal("the-encoded-result"))
            Expect(cryptographer.EncryptArgument).To(Equal("the-user|the-client-id|the-kind-id"))
        })

        It("falls back to Kind if KindDescription is missing", func() {
            options.KindDescription = ""
            context := postal.NewMessageContext(email, options, env, "the-space", "the-org", "the-client-id", "random-id", "the-user", templates, cryptographer)

            Expect(context.KindDescription).To(Equal("the-kind-id"))
        })

        It("falls back to clientID when SourceDescription is missing", func() {
            options.SourceDescription = ""
            context := postal.NewMessageContext(email, options, env, "the-space", "the-org", "the-client-id", "my-message-id", "the-user", templates, cryptographer)

            Expect(context.SourceDescription).To(Equal("the-client-id"))
        })
    })

    Describe("Escape", func() {
        var templates postal.Templates
        var email string
        var env config.Environment
        var options postal.Options
        var cryptographer *fakes.FakeCryptoClient

        BeforeEach(func() {
            email = "bounce@example.com"

            env = config.NewEnvironment()

            templates = postal.Templates{
                Text:    "the plainText email < template",
                HTML:    "the html <h1> email < template</h1>",
                Subject: "the subject < template",
            }

            options = postal.Options{
                ReplyTo:           "awesomeness",
                Subject:           "the & subject",
                KindDescription:   "the & kind description",
                SourceDescription: "the & source description",
                Text:              "user & supplied email text",
                HTML:              postal.HTML{BodyContent: "user & supplied html"},
                KindID:            "the & kind",
            }

            cryptographer = &fakes.FakeCryptoClient{
                EncryptedResult: "the-encrypted-result",
            }
        })

        It("html escapes various fields on the message context", func() {
            context := postal.NewMessageContext(email, options, env, "the<space", "the>org", "the\"client id", "some>id", "the-user", templates, cryptographer)
            context.Escape()

            Expect(context.From).To(Equal(html.EscapeString(env.Sender)))
            Expect(context.ReplyTo).To(Equal("awesomeness"))
            Expect(context.To).To(Equal("bounce@example.com"))
            Expect(context.Subject).To(Equal("the &amp; subject"))
            Expect(context.Text).To(Equal("user &amp; supplied email text"))
            Expect(context.HTML).To(Equal("user & supplied html"))
            Expect(context.TextTemplate).To(Equal("the plainText email < template"))
            Expect(context.HTMLTemplate).To(Equal("the html <h1> email < template</h1>"))
            Expect(context.SubjectTemplate).To(Equal("the subject < template"))
            Expect(context.KindDescription).To(Equal("the &amp; kind description"))
            Expect(context.SourceDescription).To(Equal("the &amp; source description"))
            Expect(context.UserGUID).To(Equal("the-user"))
            Expect(context.ClientID).To(Equal("the&#34;client id"))
            Expect(context.MessageID).To(Equal("some&gt;id"))
            Expect(context.Space).To(Equal("the&lt;space"))
            Expect(context.Organization).To(Equal("the&gt;org"))
        })
    })
})
