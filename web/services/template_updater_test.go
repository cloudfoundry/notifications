package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Updater", func() {
    Describe("#Update", func() {
        var fakeTemplatesRepo *fakes.FakeTemplatesRepo
        var template models.Template
        var updater services.TemplateUpdater

        BeforeEach(func() {
            fakeTemplatesRepo = fakes.NewFakeTemplatesRepo()
            template = models.Template{
                Name:       "gobble.user_body",
                Text:       "gobble",
                HTML:       "<p>gobble</p>",
                Overridden: true,
            }

            updater = services.NewTemplateUpdater(fakeTemplatesRepo, fakes.NewDatabase())
        })

        It("Inserts templates into the templates repo", func() {
            Expect(fakeTemplatesRepo.Templates).ToNot(ContainElement(template))
            err := updater.Update(template)
            Expect(err).ToNot(HaveOccurred())
            Expect(fakeTemplatesRepo.Templates).To(ContainElement(template))
        })

        It("propagates errors from repo", func() {
            expectedErr := errors.New("Boom!")

            fakeTemplatesRepo.UpsertError = expectedErr
            err := updater.Update(template)

            Expect(err).To(Equal(expectedErr))
        })
    })
})
