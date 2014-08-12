package postal_test

import (
    "html"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/postal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {

    Describe("NewMessageContext", func() {
        var templates postal.Templates
        var email string
        var env config.Environment
        var options postal.Options

        BeforeEach(func() {
            email = "bounce@example.com"

            env = config.NewEnvironment()

            templates = postal.Templates{
                Text:    "the plainText email template",
                HTML:    "the html email template",
                Subject: "the subject template",
            }

            options = postal.Options{
                ReplyTo:           "awesomeness",
                Subject:           "the subject",
                KindDescription:   "the kind description",
                SourceDescription: "the source description",
                Text:              "user supplied email text",
                HTML:              "user supplied html",
                Kind:              "the-kind",
            }
        })

        It("returns the appropriate MessageContext when all options are specified", func() {
            context := postal.NewMessageContext(email, options, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, templates)

            guid, err := FakeGuidGenerator()
            if err != nil {
                panic(err)
            }

            Expect(context.From).To(Equal(env.Sender))
            Expect(context.ReplyTo).To(Equal(options.ReplyTo))
            Expect(context.To).To(Equal(email))
            Expect(context.Subject).To(Equal(options.Subject))
            Expect(context.Text).To(Equal(options.Text))
            Expect(context.HTML).To(Equal(options.HTML))
            Expect(context.TextTemplate).To(Equal(templates.Text))
            Expect(context.HTMLTemplate).To(Equal(templates.HTML))
            Expect(context.SubjectTemplate).To(Equal(templates.Subject))
            Expect(context.KindDescription).To(Equal(options.KindDescription))
            Expect(context.SourceDescription).To(Equal(options.SourceDescription))
            Expect(context.ClientID).To(Equal("the-client-ID"))
            Expect(context.MessageID).To(Equal(guid.String()))
            Expect(context.Space).To(Equal("the-space"))
            Expect(context.Organization).To(Equal("the-org"))
        })

        It("falls back to Kind if KindDescription is missing", func() {
            options.KindDescription = ""
            context := postal.NewMessageContext(email, options, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, templates)

            Expect(context.KindDescription).To(Equal("the-kind"))
        })

        It("falls back to clientID when SourceDescription is missing", func() {
            options.SourceDescription = ""
            context := postal.NewMessageContext(email, options, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, templates)

            Expect(context.SourceDescription).To(Equal("the-client-ID"))
        })
    })

    Describe("Escape", func() {
        var templates postal.Templates
        var email string
        var env config.Environment
        var options postal.Options

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
                HTML:              "user & supplied html",
                Kind:              "the & kind",
            }
        })

        It("html escapes various fields on the message context", func() {
            context := postal.NewMessageContext(email, options, env, "the<space", "the>org", "the\"client ID", FakeGuidGenerator, templates)

            guid, err := FakeGuidGenerator()
            if err != nil {
                panic(err)
            }

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
            Expect(context.ClientID).To(Equal("the&#34;client ID"))
            Expect(context.MessageID).To(Equal(guid.String()))
            Expect(context.Space).To(Equal("the&lt;space"))
            Expect(context.Organization).To(Equal("the&gt;org"))
        })
    })
})
