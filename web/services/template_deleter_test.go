package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Deleter", func() {
    var deleter services.TemplateDeleter
    var fakeTemplatesRepo *fakes.FakeTemplatesRepo

    BeforeEach(func() {
        fakeTemplatesRepo = fakes.NewFakeTemplatesRepo()
        deleter = services.NewTemplateDeleter(fakeTemplatesRepo, fakes.NewDatabase())
    })

    Describe("#Delete", func() {
        It("calls destroy on its repo", func() {
            err := deleter.Delete("templateName")
            if err != nil {
                panic(err)
            }

            Expect(fakeTemplatesRepo.DestroyArgument).To(Equal("templateName"))
        })

        It("returns an error if repo destroy returns an error", func() {
            fakeTemplatesRepo.DestroyError = errors.New("Boom!!")
            err := deleter.Delete("templateName")
            Expect(err).To(Equal(fakeTemplatesRepo.DestroyError))
        })
    })
})
