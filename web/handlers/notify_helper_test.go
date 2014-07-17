package handlers_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyHelper", func() {
    var helper handlers.NotifyHelper
    var emailTemplate string
    var user uaa.User
    var env config.Environment
    var params handlers.NotifyParams

    BeforeEach(func() {
        user = uaa.User{
            ID:     "user-456",
            Emails: []string{"bounce@example.com"},
        }

        env = config.NewEnvironment()

        emailTemplate = "the email template"

        params = handlers.NotifyParams{
            Subject:           "the subject",
            KindDescription:   "the kind description",
            SourceDescription: "the source description",
            Text:              "the email text",
            Kind:              "the-kind",
        }
    })

    Describe("BuildContext", func() {

        It("returns the appropriate MessageContext when all params are specified", func() {
            messageContext := helper.BuildSpaceContext(user, params, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, emailTemplate)

            guid, err := FakeGuidGenerator()
            if err != nil {
                panic(err)
            }

            Expect(messageContext.From).To(Equal(os.Getenv("SENDER")))
            Expect(messageContext.To).To(Equal(user.Emails[0]))
            Expect(messageContext.Subject).To(Equal(params.Subject))
            Expect(messageContext.Text).To(Equal(params.Text))
            Expect(messageContext.Template).To(Equal(emailTemplate))
            Expect(messageContext.KindDescription).To(Equal(params.KindDescription))
            Expect(messageContext.SourceDescription).To(Equal(params.SourceDescription))
            Expect(messageContext.ClientID).To(Equal("the-client-ID"))
            Expect(messageContext.MessageID).To(Equal(guid.String()))
            Expect(messageContext.Space).To(Equal("the-space"))
            Expect(messageContext.Organization).To(Equal("the-org"))
        })

        It("falls back to Kind if KindDescription is missing", func() {
            params.KindDescription = ""

            messageContext := helper.BuildUserContext(user, params, env, "the-client-ID", FakeGuidGenerator, emailTemplate)

            Expect(messageContext.KindDescription).To(Equal("the-kind"))
        })

        It("falls back to clientID when SourceDescription is missing", func() {
            params.SourceDescription = ""

            messageContext := helper.BuildSpaceContext(user, params, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, emailTemplate)

            Expect(messageContext.SourceDescription).To(Equal("the-client-ID"))
        })
    })
})
