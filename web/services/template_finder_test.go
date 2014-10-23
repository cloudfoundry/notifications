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
            finder = services.NewTemplateFinder(fakeTemplatesRepo, env.RootPath)
        })

        Context("when the finder returns a template", func() {
            It("returns the default value for an unknown space template", func() {
                fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                template, err := finder.Find(services.SpaceBody, "missing_template_file")
                Expect(err).ToNot(HaveOccurred())
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

            It("returns the default value for an unknown user template", func() {
                fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                template, err := finder.Find(services.UserBody, "missing_template_file")
                Expect(err).ToNot(HaveOccurred())

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

            It("returns a template for a known template", func() {
                expectedTemplate := models.Template{
                    Text:       "hungry raptors template",
                    HTML:       "<p>hungry raptors template</p>",
                    Overridden: true,
                }

                fakeTemplatesRepo.Templates["raptors.hungry.user_body"] = expectedTemplate
                template, err := finder.Find(services.UserBody, "raptors.hungry.user_body")
                Expect(err).ToNot(HaveOccurred())
                Expect(template).To(Equal(expectedTemplate))
            })
        })

        Context("when the finder returns an error", func() {
            It("propagates the error", func() {
                fakeTemplatesRepo.FindError = errors.New("some-error")
                _, err := finder.Find(services.SpaceBody, "missing_template_file")
                Expect(err.Error()).To(Equal("some-error"))
            })
        })
    })
})
