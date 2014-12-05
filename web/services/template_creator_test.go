package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Creator", func() {
	Describe("Create", func() {
		var templatesRepo *fakes.TemplatesRepo
		var template models.Template
		var creator services.TemplateCreator

		BeforeEach(func() {
			templatesRepo = fakes.NewTemplatesRepo()
			template = models.Template{
				Name:    "Big Hero 6 Template",
				Text:    "Adorable robot.",
				HTML:    "<p>Many heroes.</p>",
				Subject: "Robots and Heroes",
			}

			creator = services.NewTemplateCreator(templatesRepo, fakes.NewDatabase())
		})

		It("Creates a new template via the templates repo", func() {
			Expect(templatesRepo.Templates).ToNot(ContainElement(template))
			_, err := creator.Create(template)
			if err != nil {
				panic(err)
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(templatesRepo.Templates).To(ContainElement(template))
		})

		It("propagates errors from repo", func() {
			expectedErr := errors.New("Boom!")

			templatesRepo.CreateError = expectedErr
			_, err := creator.Create(template)

			Expect(err).To(Equal(expectedErr))
		})
	})
})
