package postal_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {
    Describe("NewMessageContext", func() {
        var templates postal.Templates
        var user uaa.User
        var env config.Environment
        var options postal.Options

        BeforeEach(func() {
            user = uaa.User{
                ID:     "user-456",
                Emails: []string{"bounce@example.com"},
            }

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
            messageContext := postal.NewMessageContext(user, options, env, "the-space", "the-org",
                "the-client-ID", FakeGuidGenerator, templates)

            guid, err := FakeGuidGenerator()
            if err != nil {
                panic(err)
            }

            Expect(messageContext.From).To(Equal(os.Getenv("SENDER")))
            Expect(messageContext.ReplyTo).To(Equal(options.ReplyTo))
            Expect(messageContext.To).To(Equal(user.Emails[0]))
            Expect(messageContext.Subject).To(Equal(options.Subject))
            Expect(messageContext.Text).To(Equal(options.Text))
            Expect(messageContext.HTML).To(Equal(options.HTML))
            Expect(messageContext.TextTemplate).To(Equal(templates.Text))
            Expect(messageContext.HTMLTemplate).To(Equal(templates.HTML))
            Expect(messageContext.SubjectTemplate).To(Equal(templates.Subject))
            Expect(messageContext.KindDescription).To(Equal(options.KindDescription))
            Expect(messageContext.SourceDescription).To(Equal(options.SourceDescription))
            Expect(messageContext.ClientID).To(Equal("the-client-ID"))
            Expect(messageContext.MessageID).To(Equal(guid.String()))
            Expect(messageContext.Space).To(Equal("the-space"))
            Expect(messageContext.Organization).To(Equal("the-org"))
        })

        It("falls back to Kind if KindDescription is missing", func() {
            options.KindDescription = ""

            messageContext := postal.NewMessageContext(user, options, env, "the-space",
                "the-org", "the-client-ID", FakeGuidGenerator, templates)

            Expect(messageContext.KindDescription).To(Equal("the-kind"))
        })

        It("falls back to clientID when SourceDescription is missing", func() {
            options.SourceDescription = ""

            messageContext := postal.NewMessageContext(user, options, env, "the-space",
                "the-org", "the-client-ID", FakeGuidGenerator, templates)

            Expect(messageContext.SourceDescription).To(Equal("the-client-ID"))
        })
    })
})
