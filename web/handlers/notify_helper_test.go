package handlers_test

import (
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyHelper", func() {
    var helper handlers.NotifyHelper
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

    Describe("BuildContext", func() {
        It("returns the appropriate MessageContext when all params are specified", func() {
            messageContext := helper.BuildSpaceContext(user, params, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, plainTextEmailTemplate, htmlEmailTemplate)

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

            messageContext := helper.BuildUserContext(user, params, env, "the-client-ID", FakeGuidGenerator, plainTextEmailTemplate, htmlEmailTemplate)

            Expect(messageContext.KindDescription).To(Equal("the-kind"))
        })

        It("falls back to clientID when SourceDescription is missing", func() {
            params.SourceDescription = ""

            messageContext := helper.BuildSpaceContext(user, params, env, "the-space", "the-org", "the-client-ID", FakeGuidGenerator, plainTextEmailTemplate, htmlEmailTemplate)

            Expect(messageContext.SourceDescription).To(Equal("the-client-ID"))
        })
    })

    Describe("LoadTemplates", func() {
        Context("when there are no template overrides", func() {
            It("loads the templates from the default location", func() {
                helper.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "html") {
                        return "the fake html", nil
                    } else {
                        return "the fake text", nil
                    }
                }

                helper.FileExists = func(path string) bool {
                    return false
                }

                plainTextTemplate, htmlTemplate, err := helper.LoadTemplates(false, env.RootPath)
                if err != nil {
                    panic(err)
                }
                Expect(htmlTemplate).To(Equal("the fake html"))
                Expect(plainTextTemplate).To(Equal("the fake text"))
            })
        })

        Context("when a template has an override set", func() {
            It("replaces the default template with the user provided one", func() {
                var loadedPaths []string

                helper.ReadFile = func(path string) (string, error) {
                    loadedPaths = append(loadedPaths, path)
                    return "the fake template", nil
                }

                helper.FileExists = func(path string) bool {
                    return true
                }

                _, _, err := helper.LoadTemplates(false, env.RootPath)
                if err != nil {
                    panic(err)
                }

                Expect(loadedPaths).To(ContainElement(env.RootPath + "/templates/overrides/user_body.text"))
                Expect(loadedPaths).To(ContainElement(env.RootPath + "/templates/overrides/user_body.html"))
            })
        })
    })
})
