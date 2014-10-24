package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Finder", func() {
    var finder services.TemplateFinder
    var fakeTemplatesRepo *fakes.FakeTemplatesRepo

    Describe("#Find", func() {
        BeforeEach(func() {
            env := config.NewEnvironment()
            fakeTemplatesRepo = fakes.NewFakeTemplatesRepo()
            finder = services.NewTemplateFinder(fakeTemplatesRepo, env.RootPath, fakes.NewDatabase())
        })

        Context("when the finder returns a template", func() {
            Context("when the override does not exist", func() {
                It("returns the default template space template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp.space_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal(`Hello {{.To}},

The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}"
component of Cloud Foundry because you are a member of the "{{.Space}}" space
in the "{{.Organization}}" organization:

{{.Text}}

This message was sent from {{.From}} and can be replied to at {{.ReplyTo}}. The
notification can be identified with the {{.MessageID}} identifier and was sent
with the {{.ClientID}} UAA client. The notification can be unsubscribed from
using the "{{.UnsubscribeID}}" unsubscribe token.
`))
                    Expect(template.HTML).To(Equal(`<p>Hello {{.To}},</p>

<p>The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}"
    component of Cloud Foundry because you are a member of the "{{.Space}}" space
    in the "{{.Organization}}" organization:</p>

{{.HTML}}

<p>This message was sent from {{.From}} and can be replied to at {{.ReplyTo}}. The
    notification can be identified with the {{.MessageID}} identifier and was sent
    with the {{.ClientID}} UAA client. The notification can be unsubscribed from
    using the "{{.UnsubscribeID}}" unsubscribe token.</p>
`))
                })

                It("returns the default user template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp.user_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal(`Hello {{.To}},

The following "{{.KindDescription}}" notification was sent to you directly by the
"{{.SourceDescription}}" component of Cloud Foundry:

{{.Text}}

This message was sent from {{.From}} and can be replied to at {{.ReplyTo}}. The
notification can be identified with the {{.MessageID}} identifier and was sent
with the {{.ClientID}} UAA client. The notification can be unsubscribed from
using the "{{.UnsubscribeID}}" unsubscribe token.
`))

                    Expect(template.HTML).To(Equal(`<p>Hello {{.To}},</p>

<p>The following "{{.KindDescription}}" notification was sent to you directly by the
    "{{.SourceDescription}}" component of Cloud Foundry:</p>

{{.HTML}}

<p>This message was sent from {{.From}} and can be replied to at {{.ReplyTo}}. The
    notification can be identified with the {{.MessageID}} identifier and was sent
    with the {{.ClientID}} UAA client. The notification can be unsubscribed from
    using the "{{.UnsubscribeID}}" unsubscribe token.</p>
`))
                })
            })

            Context("when the override exists in the database", func() {
                var expectedTemplate models.Template

                BeforeEach(func() {
                    expectedTemplate = models.Template{
                        Text:       "authenticate new hungry raptors template",
                        HTML:       "<p>hungry raptors are newly authenticated template</p>",
                        Overridden: true,
                    }
                    fakeTemplatesRepo.Templates["authentication.new.user_body"] = expectedTemplate
                })

                It("returns the requested override template", func() {
                    template, err := finder.Find("authentication.new.user_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeTrue())
                    Expect(template).To(Equal(expectedTemplate))
                })

            })

            Context("when the requested client/kind override does not exist in db", func() {
                Context("but the client override does", func() {
                    var expectedTemplate models.Template

                    BeforeEach(func() {
                        expectedTemplate = models.Template{
                            Text:       "authentication template for hungry raptors",
                            HTML:       "<h1>Wow you are authentic!</h1>",
                            Overridden: true,
                        }
                        fakeTemplatesRepo.Templates["authentication.user_body"] = expectedTemplate
                    })

                    It("returns the fallback override that exists", func() {
                        template, err := finder.Find("authentication.new.user_body")
                        Expect(err).ToNot(HaveOccurred())
                        Expect(template.Overridden).To(BeTrue())
                        Expect(template).To(Equal(expectedTemplate))
                    })
                })

                Context("when the client override does not exist, but the notification type does", func() {
                    var expectedTemplate models.Template

                    BeforeEach(func() {
                        expectedTemplate = models.Template{
                            Text:       "special user template",
                            HTML:       "<h1>Wow you are a special user!</h1>",
                            Overridden: true,
                        }
                        fakeTemplatesRepo.Templates["user_body"] = expectedTemplate
                    })

                    It("returns the fallback override that exists", func() {
                        template, err := finder.Find("authentication.new.user_body")
                        Expect(err).ToNot(HaveOccurred())
                        Expect(template.Overridden).To(BeTrue())
                        Expect(template).To(Equal(expectedTemplate))
                    })
                })
            })
        })

        Context("when the finder returns an error", func() {
            It("propagates the error", func() {
                fakeTemplatesRepo.FindError = errors.New("some-error")
                _, err := finder.Find("missing_template_file")
                Expect(err.Error()).To(Equal("some-error"))
            })
        })

        Context("when something wacky is requested", func() {
            It("returns an error", func() {
                fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}
                _, err := finder.Find("...")
                Expect(err).To(HaveOccurred())
                Expect(err).To(MatchError(models.ErrRecordNotFound{}))
            })
        })
    })
})
