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
	var templatesRepo *fakes.TemplatesRepo

	BeforeEach(func() {
		templatesRepo = fakes.NewTemplatesRepo()
		deleter = services.NewTemplateDeleter(templatesRepo, fakes.NewDatabase())
	})

	Describe("#Delete", func() {
		It("calls destroy on its repo", func() {
			err := deleter.Delete("templateName")
			if err != nil {
				panic(err)
			}

			Expect(templatesRepo.DestroyArgument).To(Equal("templateName"))
		})

		It("returns an error if repo destroy returns an error", func() {
			templatesRepo.DestroyError = errors.New("Boom!!")
			err := deleter.Delete("templateName")
			Expect(err).To(Equal(templatesRepo.DestroyError))
		})
	})
})
