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
			err := deleter.Delete("templateID")
			if err != nil {
				panic(err)
			}

			Expect(templatesRepo.DestroyArgument).To(Equal("templateID"))
		})

		It("returns an error if repo destroy returns an error", func() {
			templatesRepo.DestroyError = errors.New("Boom!!")
			err := deleter.Delete("templateID")
			Expect(err).To(Equal(templatesRepo.DestroyError))
		})
	})

	Describe("#DeprecatedDelete", func() {
		It("calls destroy by template name on its repo", func() {
			err := deleter.DeprecatedDelete("templateName")
			if err != nil {
				panic(err)
			}

			Expect(templatesRepo.DeprecatedDestroyArgument).To(Equal("templateName"))
		})

		It("returns an error if repo destroy returns an error", func() {
			templatesRepo.DeprecatedDestroyError = errors.New("Boom!!")
			err := deleter.DeprecatedDelete("templateName")
			Expect(err).To(Equal(templatesRepo.DeprecatedDestroyError))
		})
	})
})
